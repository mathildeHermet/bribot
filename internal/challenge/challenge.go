package challenge

import (
	"fmt"
	"time"

	"github.com/mathildeHermet/bribot/internal/discord"
)

type Challenge struct {
	name      string
	url       string
	teamName  string
	teamPwd   string
	beginning time.Time
	deadline  time.Time
}

func NewChallenge(name, url, teamName, teamPwd, beginning, deadline, layout string) (*Challenge, error) {
	//  "Mon, 02 Jan 2006 15:04:05 MST"
	b, err := time.Parse(layout, beginning)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse start date")
	}
	d, err := time.Parse(layout, deadline)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse end date")
	}

	// Load Paris location (Europe/Paris is always UTC+2)
	loc, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		return nil, fmt.Errorf("Cannot load Paris location: %v", err)
	}

	// Convert to Paris local time
	localBegin := b.In(loc)
	localEnd := d.In(loc)

	duration := d.Sub(b)
	if duration < 0 {
		return nil, fmt.Errorf("End date should be after start date")
	}

	return &Challenge{
		name:      name,
		url:       url,
		teamName:  teamName,
		teamPwd:   teamPwd,
		beginning: localBegin,
		deadline:  localEnd,
	}, nil
}

func (c *Challenge) remainingTime() string {
	// Calculate time remaining until 23:00 today
	now := time.Now().Local()
	end := c.deadline
	remaining := end.Sub(now)
	second := time.Second * 1
	remainingTime := fmt.Sprintf("%+v", remaining.Round(second))

	return remainingTime
}

func (c *Challenge) MakeReminderMessage() *discord.Message {

	content := fmt.Sprintf("%s incoming. Register to %s team", c.name, c.teamName)
	registerMsg := fmt.Sprintf("Team Name: %s\nPassword: %s", c.teamName, c.teamPwd)
	timeLine := fmt.Sprintf("Start at: %s\nEnd at: %s\n", c.beginning.Format(time.RFC1123), c.deadline.Format(time.RFC1123))
	untill := c.remainingTime()

	fieldRemaining := &discord.Field{
		Name:   "Time Remaining",
		Value:  untill,
		Inline: true,
	}
	fieldRegister := &discord.Field{
		Name:   "Register",
		Value:  registerMsg,
		Inline: false,
	}
	fieldCtfTimeline := &discord.Field{
		Name:   "CTF Challenge Time Line",
		Value:  timeLine,
		Inline: false,
	}
	fields := []discord.Field{*fieldRemaining, *fieldRegister, *fieldCtfTimeline}
	msgEmbeds := &discord.Embed{
		Title:       "CTF challenge is starting soon !!",
		Description: "Make sure to prepare and join in time!",
		URL:         c.url,
		Fields:      fields,
		Color:       15258703,
	}
	return &discord.Message{
		Content: content,
		Embeds:  []discord.Embed{*msgEmbeds},
	}
}

// parseTimeInLocation parses the time using the provided layout and timezone.
func parseTimeInLocation(timeStr, layout, timezone string) (time.Time, error) {
	// Load the desired location
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to load location: %w", err)
	}

	// Parse the time with the specified layout and location
	t, err := time.ParseInLocation(layout, timeStr, location)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse time: %w", err)
	}

	return t, nil
}
