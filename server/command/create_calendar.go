package command

import (
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/config"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/remote"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/utils"
)

func (c *Command) createCalendar(parameters ...string) (string, bool, error) {
	if len(parameters) != 1 {
		return "Please provide the name of one calendar to create", false, nil
	}

	calIn := &remote.Calendar{
		Name: parameters[0],
	}

	resp, err := c.MSCalendar.CreateCalendar(c.user(), calIn)
	if err != nil {
		return "", false, err
	}
	return utils.JSONHeading(config.CreateCalendarHeading) + utils.JSONBlock(resp), false, nil
}
