// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package plugin

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	pluginapiclient "github.com/mattermost/mattermost-plugin-api"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin"
	"github.com/pkg/errors"

	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/api"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/command"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/config"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/jobs"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/mscalendar"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/remote"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/remote/msgraph"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/store"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/tracker"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/utils/bot"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/utils/flow"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/utils/httputils"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/utils/pluginapi"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/utils/settingspanel"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/utils/telemetry"
)

type Env struct {
	mscalendar.Env
	bot                   bot.Bot
	jobManager            *jobs.JobManager
	notificationProcessor mscalendar.NotificationProcessor
	httpHandler           *httputils.Handler
	configError           error
}

type Plugin struct {
	plugin.MattermostPlugin

	envLock         *sync.RWMutex
	env             Env
	Templates       map[string]*template.Template
	telemetryClient telemetry.Client
}

func NewWithEnv(env mscalendar.Env) *Plugin {
	return &Plugin{
		env: Env{
			Env: env,
		},
		envLock: &sync.RWMutex{},
	}
}

func (p *Plugin) OnActivate() error {
	pluginAPIClient := pluginapiclient.NewClient(p.API, p.Driver)
	stored := config.StoredConfig{}
	err := p.API.LoadPluginConfiguration(&stored)
	if err != nil {
		return errors.WithMessage(err, "failed to load plugin configuration")
	}

	_ = p.initEnv(&p.env, pluginAPIClient, "")
	if stored.EWSProxyServerBaseURL == "" || stored.EWSProxyServerAuthKey == "" {
		return errors.New("failed to configure: Exchange Server credentials to be set in the config")
	}

	bundlePath, err := p.API.GetBundlePath()
	if err != nil {
		return errors.Wrap(err, "couldn't get bundle path")
	}
	err = p.loadTemplates(bundlePath)
	if err != nil {
		return err
	}

	err = command.Register(pluginAPIClient, stored.AutoConnectUsers)
	if err != nil {
		return errors.Wrap(err, "failed to register command")
	}

	p.telemetryClient, err = telemetry.NewRudderClient()
	if err != nil {
		p.env.bot.Errorf("Cannot create telemetry client. err=%v", err)
	}

	// Auto-connect all the users present on server
	if stored.AutoConnectUsers {
		go p.ConnectUsers()
	}

	return nil
}

func (p *Plugin) ConnectUsers() {
	teams, err := p.API.GetTeams()
	if err != nil {
		p.API.LogError(fmt.Sprintf("error occurred while fetching teams. Error: %s", err.Error()))
		return
	}
	m := mscalendar.New(p.getEnv().Env, "")
	for _, team := range teams {
		teamStats, err := p.API.GetTeamStats(team.Id)
		if err != nil {
			p.API.LogError(fmt.Sprintf("error occurred while fetching team stats of team %s. Error: %s", team.Name, err.Error()))
			continue
		}
		totalPages := int(teamStats.TotalMemberCount) / config.UsersCountPerPage
		if int(teamStats.TotalMemberCount)%config.UsersCountPerPage != 0 {
			totalPages++
		}
		for page := 0; page < totalPages; page++ {
			users, err := p.API.GetUsersInTeam(team.Id, page, config.UsersCountPerPage)
			if err != nil {
				p.API.LogError(fmt.Sprintf("error occurred while fetching users of team %s. Error: %s", team.Name, err.Error()))
				continue
			}
			connectErr := m.CompleteOAuth2ForUsers(users)
			if connectErr != nil {
				p.API.LogError(fmt.Sprintf("error occurred while connecting users of team %s. Error: %s", team.Name, connectErr.Error()))
				continue
			}
		}
	}
}

func (p *Plugin) OnDeactivate() error {
	if p.telemetryClient != nil {
		err := p.telemetryClient.Close()
		if err != nil {
			p.env.Logger.Warnf("OnDeactivate: Failed to close telemetryClient. err=%v", err)
		}
	}

	e := p.getEnv()
	if e.jobManager != nil {
		if err := e.jobManager.Close(); err != nil {
			p.env.Logger.Warnf("OnDeactivate: Failed to close job manager. err=%v", err)
			return err
		}
	}
	return nil
}

func (p *Plugin) OnConfigurationChange() (err error) {
	pluginAPIClient := pluginapiclient.NewClient(p.API, p.Driver)
	defer func() {
		p.updateEnv(func(e *Env) {
			e.configError = err
		})
	}()

	env := p.getEnv()
	stored := config.StoredConfig{}
	err = p.API.LoadPluginConfiguration(&stored)
	if err != nil {
		return errors.WithMessage(err, "failed to load plugin configuration")
	}
	if stored.StatusSyncJobInterval <= 0 {
		return errors.New("status sync job interval should be greater than 0")
	}

	if p.getEnv().Config != nil && p.getEnv().Config.AutoConnectUsers != stored.AutoConnectUsers {
		if err = command.UnRegister(pluginAPIClient); err != nil {
			return errors.Wrap(err, "failed to unregister slash command")
		}

		if err = command.Register(pluginAPIClient, stored.AutoConnectUsers); err != nil {
			return errors.Wrap(err, "failed to register slash command")
		}
	}

	mattermostSiteURL := p.API.GetConfig().ServiceSettings.SiteURL
	if mattermostSiteURL == nil {
		return errors.New("plugin requires Mattermost Site URL to be set")
	}
	mattermostURL, err := url.Parse(*mattermostSiteURL)
	if err != nil {
		return err
	}
	pluginURLPath := "/plugins/" + env.Config.PluginID
	pluginURL := strings.TrimRight(*mattermostSiteURL, "/") + pluginURLPath

	p.updateEnv(func(e *Env) {
		_ = p.initEnv(e, pluginAPIClient, pluginURL)
		previousConfig := e.Config
		if previousConfig != nil && previousConfig.StatusSyncJobInterval != stored.StatusSyncJobInterval && e.jobManager != nil {
			if err = env.jobManager.Close(); err != nil {
				e.Logger.Warnf("Failed to close job manager. err=%v", err)
				e.configError = err
				return
			}
			e.jobManager = nil
		}
		e.StoredConfig = stored
		e.Config.MattermostSiteURL = *mattermostSiteURL
		e.Config.MattermostSiteHostname = mattermostURL.Hostname()
		e.Config.PluginURL = pluginURL
		e.Config.PluginURLPath = pluginURLPath

		e.bot = e.bot.WithConfig(stored.Config)
		e.Dependencies.Remote = remote.Makers[msgraph.Kind](e.Config, e.bot)

		mscalendarBot := mscalendar.NewMSCalendarBot(e.bot, e.Env, pluginURL)

		e.Dependencies.Logger = e.bot

		diagnostics := false
		if p.API.GetConfig() != nil && p.API.GetConfig().LogSettings.EnableDiagnostics != nil {
			diagnostics = *p.API.GetConfig().LogSettings.EnableDiagnostics
		}
		e.Dependencies.Tracker = tracker.New(telemetry.NewTracker(p.telemetryClient, p.API.GetDiagnosticId(), p.API.GetServerVersion(), e.PluginID, e.PluginVersion, config.TelemetryShortName, diagnostics, e.Logger))

		e.Dependencies.Poster = e.bot
		e.Dependencies.Welcomer = mscalendarBot
		e.Dependencies.Store = store.NewPluginStore(p.API, e.bot, e.Dependencies.Tracker)
		e.Dependencies.SettingsPanel = mscalendar.NewSettingsPanel(
			e.bot,
			e.Dependencies.Store,
			e.Dependencies.Store,
			"/settings",
			pluginURL,
			func(userID string) mscalendar.MSCalendar {
				return mscalendar.New(e.Env, userID)
			},
		)

		welcomeFlow := mscalendar.NewWelcomeFlow(e.bot, e.Dependencies.Welcomer)
		e.bot.RegisterFlow(welcomeFlow, mscalendarBot)

		if e.notificationProcessor == nil {
			e.notificationProcessor = mscalendar.NewNotificationProcessor(e.Env)
		} else {
			e.notificationProcessor.Configure(e.Env)
		}

		e.httpHandler = httputils.NewHandler()
		// oauth2connect.Init(e.httpHandler, mscalendar.NewOAuth2App(e.Env))
		flow.Init(e.httpHandler, welcomeFlow, mscalendarBot)
		settingspanel.Init(e.httpHandler, e.Dependencies.SettingsPanel)
		api.Init(e.httpHandler, e.Env, e.notificationProcessor)

		if e.jobManager == nil {
			e.jobManager = jobs.NewJobManager(p.API, e.Env)
			e.jobManager.AddJob(jobs.NewStatusSyncJob(e.Env))
			e.jobManager.AddJob(jobs.NewDailySummaryJob())
		}
	})

	return nil
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	env := p.getEnv()
	if env.configError != nil {
		p.API.LogError(env.configError.Error())
		return nil, model.NewAppError("mscalendarplugin.ExecuteCommand", "Unable to execute command.", nil, env.configError.Error(), http.StatusInternalServerError)
	}

	command := command.Command{
		Context:    c,
		Args:       args,
		ChannelID:  args.ChannelId,
		Config:     env.Config,
		MSCalendar: mscalendar.New(env.Env, args.UserId),
		API:        p.API,
	}
	out, mustRedirectToDM, err := command.Handle()
	if err != nil {
		p.API.LogError(err.Error())
		return nil, model.NewAppError("mscalendarplugin.ExecuteCommand", "Unable to execute command.", nil, err.Error(), http.StatusInternalServerError)
	}

	if out != "" {
		env.Poster.Ephemeral(args.UserId, args.ChannelId, "%s", out)
	}

	response := &model.CommandResponse{}
	if mustRedirectToDM {
		t, appErr := p.API.GetTeam(args.TeamId)
		if appErr != nil {
			return nil, model.NewAppError("mscalendarplugin.ExecuteCommand", "Unable to execute command.", nil, appErr.Error(), http.StatusInternalServerError)
		}
		dmURL := fmt.Sprintf("%s/%s/messages/@%s", env.MattermostSiteURL, t.Name, config.BotUserName)
		response.GotoLocation = dmURL
	}

	return response, nil
}

func (p *Plugin) ServeHTTP(pc *plugin.Context, w http.ResponseWriter, req *http.Request) {
	env := p.getEnv()
	if env.configError != nil {
		p.API.LogError(env.configError.Error())
		http.Error(w, env.configError.Error(), http.StatusInternalServerError)
		return
	}

	env.httpHandler.ServeHTTP(w, req)
}

func (p *Plugin) getEnv() Env {
	p.envLock.RLock()
	defer p.envLock.RUnlock()
	return p.env
}

func (p *Plugin) updateEnv(f func(*Env)) Env {
	p.envLock.Lock()
	defer p.envLock.Unlock()

	f(&p.env)
	return p.env
}

func (p *Plugin) loadTemplates(bundlePath string) error {
	if p.Templates != nil {
		return nil
	}

	templatesPath := filepath.Join(bundlePath, "assets", "templates")
	templates := make(map[string]*template.Template)
	err := filepath.Walk(templatesPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		template, err := template.ParseFiles(path)
		if err != nil {
			return nil
		}
		key := path[len(templatesPath):]
		templates[key] = template
		return nil
	})
	if err != nil {
		return errors.WithMessage(err, "OnActivate/loadTemplates failed")
	}
	p.Templates = templates
	return nil
}

func (p *Plugin) initEnv(e *Env, client *pluginapiclient.Client, pluginURL string) error {
	e.Dependencies.PluginAPI = pluginapi.New(p.API)
	if e.bot == nil {
		e.bot = bot.New(p.API, *client, pluginURL)
		err := e.bot.Ensure(
			&model.Bot{
				Username:    config.BotUserName,
				DisplayName: config.BotDisplayName,
				Description: config.BotDescription,
			},
			"assets/profile.png")
		if err != nil {
			return errors.Wrap(err, "failed to ensure bot account")
		}
	}

	return nil
}

func (p *Plugin) UserHasLoggedIn(c *plugin.Context, user *model.User) {
	if p.getEnv().Config.AutoConnectUsers {
		// Auto-connect the user to ms-calendar after he logs in
		m := mscalendar.New(p.getEnv().Env, "")
		if _, err := m.GetRemoteUser(user.Id); err == nil {
			// If user is already connected do nothing
			return
		}

		if err := m.CompleteOAuth2(user.Id); err != nil {
			p.API.LogError("Error occurred while connecting user to mscalendar", "Email", user.Email, "Error", err.Error())
		}
	}
}
