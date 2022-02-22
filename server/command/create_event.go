package command

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"

	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/remote"
	"github.com/Brightscout/mattermost-plugin-exchange-mscalendar/server/utils"
)

func getCreateEventFlagSet() *flag.FlagSet {
	flagSet := flag.NewFlagSet("create", flag.ContinueOnError)
	flagSet.Bool("help", false, "show help")
	flagSet.String("test-subject", "", "Subject of the event (no spaces for now)")
	flagSet.String("test-body", "", "Body of the event (no spaces for now)")
	flagSet.String("test-location", "", "Location of the event (no spaces for now)")
	flagSet.String("starttime", time.Now().Format(time.RFC3339), "Start time for the event")
	flagSet.Bool("allday", false, "Set as all day event (starttime/endtime must be set to midnight on different days - 2019-12-19T00:00:00-00:00)")
	flagSet.Int("reminder", 15, "Reminder (in minutes)")
	flagSet.String("endtime", time.Now().Add(time.Hour).Format(time.RFC3339), "End time for the event")
	flagSet.StringSlice("attendees", nil, "A comma separated list of Mattermost User Emails")

	return flagSet
}

func (c *Command) createEvent(parameters ...string) (string, bool, error) {
	if len(parameters) == 0 {
		return getCreateEventFlagSet().FlagUsages(), false, nil
	}

	event, err := parseCreateArgs(parameters)
	if err != nil {
		return err.Error(), false, nil
	}

	createFlagSet := getCreateEventFlagSet()
	err = createFlagSet.Parse(parameters)
	if err != nil {
		return "", false, err
	}

	calEvent, err := c.MSCalendar.CreateEvent(c.user(), event)
	if err != nil {
		return "", false, err
	}
	resp := "Event Created\n" + utils.JSONBlock(&calEvent)

	return resp, false, nil
}

func parseCreateArgs(args []string) (*remote.Event, error) {
	event := &remote.Event{}

	createFlagSet := getCreateEventFlagSet()
	err := createFlagSet.Parse(args)
	if err != nil {
		return nil, err
	}

	// check for required flags
	requiredFlags := []string{"test-subject"}
	flags := make(map[string]bool)
	createFlagSet.Visit(
		func(f *flag.Flag) {
			flags[f.Name] = true
		})
	for _, req := range requiredFlags {
		if !flags[req] {
			return nil, fmt.Errorf("missing required flag: `--%s` ", req)
		}
	}

	help, err := createFlagSet.GetBool("help")
	if err != nil {
		return nil, err
	}

	if help {
		return nil, errors.New(getCreateEventFlagSet().FlagUsages())
	}

	subject, err := createFlagSet.GetString("test-subject")
	if err != nil {
		return nil, err
	}
	// check that next arg is not a flag "--"
	if strings.HasPrefix(subject, "--") {
		return nil, errors.New("test-subject flag requires an argument")
	}
	event.Subject = subject

	body, err := createFlagSet.GetString("test-body")
	if err != nil {
		return nil, err
	}
	// check that next arg is not a flag "--"
	if strings.HasPrefix(body, "--") {
		return nil, errors.New("body flag requires an argument")
	}
	event.Body = &remote.ItemBody{
		Content: body,
	}

	startTime, err := createFlagSet.GetString("starttime")
	if err != nil {
		return nil, err
	}
	if strings.HasPrefix(startTime, "--") {
		return nil, errors.New("starttime flag requires an argument")
	}
	event.Start = &remote.DateTime{
		DateTime: startTime,
	}

	endTime, err := createFlagSet.GetString("endtime")
	if err != nil {
		return nil, err
	}
	if strings.HasPrefix(endTime, "--") {
		return nil, errors.New("endtime flag requires an argument")
	}
	event.End = &remote.DateTime{
		DateTime: endTime,
	}

	mattermostUserEmails, err := createFlagSet.GetStringSlice("attendees")
	if err != nil {
		return nil, err
	}
	if len(mattermostUserEmails) != 0 {
		attendees := make([]*remote.Attendee, len(mattermostUserEmails))
		for idx, email := range mattermostUserEmails {
			attendees[idx] = &remote.Attendee{
				EmailAddress: &remote.EmailAddress{
					Address: email,
				},
			}
		}
		event.Attendees = attendees
	}

	allday, err := createFlagSet.GetBool("allday")
	if err != nil {
		return nil, err
	}
	event.IsAllDay = allday

	reminder, err := createFlagSet.GetInt("reminder")
	if err != nil {
		return nil, err
	}
	if strings.HasPrefix(strconv.Itoa(reminder), "--") {
		return nil, errors.New("reminder flag requires an argument")
	}
	event.ReminderMinutesBeforeStart = reminder

	location, err := createFlagSet.GetString("test-location")
	if err != nil {
		return nil, err
	}
	if len(location) != 0 {
		event.Location = location
	}

	return event, nil
}
