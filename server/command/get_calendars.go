package command

import (
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/utils"
)

func (c *Command) showCalendars(parameters ...string) (string, bool, error) {
	resp, err := c.MSCalendar.GetCalendars(c.user())
	if err != nil {
		return "", false, err
	}
	return utils.JSONBlock(resp), false, nil
}
