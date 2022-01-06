package command

import (
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/config"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/utils"
)

func (c *Command) deleteCalendar(parameters ...string) (string, bool, error) {
	if len(parameters) != 1 {
		return "Please provide the ID of only one calendar ", false, nil
	}

	resp, err := c.MSCalendar.DeleteCalendar(c.user(), parameters[0])
	if err != nil {
		return "", false, err
	}
	return utils.JSONHeading(config.DeleteCalendarHeading) + utils.JSONBlock(resp), false, nil
}
