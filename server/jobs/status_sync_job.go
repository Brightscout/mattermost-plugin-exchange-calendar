// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package jobs

import (
	"time"

	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/mscalendar"
)

// Unique id for the status sync job
const statusSyncJobID = "status_sync"

// NewStatusSyncJob creates a RegisteredJob with the parameters specific to the StatusSyncJob
func NewStatusSyncJob(env mscalendar.Env) RegisteredJob {
	return RegisteredJob{
		id:       statusSyncJobID,
		interval: time.Minute * time.Duration(env.StatusSyncJobInterval),
		work:     runSyncJob,
	}
}

// runSyncJob synchronizes all users' statuses between mscalendar and Mattermost.
func runSyncJob(env mscalendar.Env) {
	env.Logger.Debugf("User status sync job beginning")

	_, err := mscalendar.New(env, "").SyncAll()
	if err != nil {
		env.Logger.Errorf("Error during user status sync job. err=%v", err)
	}

	env.Logger.Debugf("User status sync job finished")
}
