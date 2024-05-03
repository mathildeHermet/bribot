/*
Copyright © 2024 Mathilde Hermet

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/mathildeHermet/bribot/internal/challenge"
	"github.com/mathildeHermet/bribot/internal/discord"
	"github.com/spf13/cobra"
)

const layout = "2006-01-02 15:04:05"

func NewReminderCmd() (*cobra.Command, error) {
	reminderCmd := &cobra.Command{
		Use:   "reminder",
		Short: "Remind CTF url, timelines adn credentials.",
		Long:  "Enroll in a specific CTF challenge.",
		Example: fmt.Sprintf("bribot reminder --%s debug --%s https://discord.com/api/webhooks/id/pwd --%s punkctf --%s https://ctf.example.com --%s 2024-01-01 00:00:00 UTC --%s 2024-02-01 00:00:00 UTC --%s UTC --%s teamName --%s teamPassword",
			flagLogLevel,
			flagWebhookUrl,
			flagCTFName,
			flagCTFURL,
			flagStartDate,
			flagEndDate,
			flagTz,
			flagTeamName,
			flagTeamPassword,
		),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := newLogger(cmd)
			logger.Info("Reminder command called")

			dryrun, err := cmd.Flags().GetBool("dry-run")
			if err != nil {
				logger.Error("Failed to get dry-run flag", err)
				return err
			}
			challengeName := cmd.Flag(flagCTFName).Value.String()
			challengeURL := cmd.Flag(flagCTFURL).Value.String()
			teamName := cmd.Flag(flagTeamName).Value.String()
			teamPassword := cmd.Flag(flagTeamPassword).Value.String()
			startDate := cmd.Flag(flagStartDate).Value.String()
			endDate := cmd.Flag(flagEndDate).Value.String()
			tz := cmd.Flag(flagTz).Value.String()

			challenge, err := challenge.NewChallenge(logger, challengeName, challengeURL, teamName, teamPassword, startDate, endDate, tz, layout)
			if err != nil {
				logger.Error("Failed to create a new challenge", err)
				return err
			}
			reminderMsg := challenge.MakeReminderMessage()

			sender := &discord.Sender{
				WebhookURL: cmd.Flag(flagWebhookUrl).Value.String(),
				Message:    reminderMsg,
			}

			if dryrun {
				logger.Info("Dry-run enabled. Skipping sending message")
				logger.Info(fmt.Sprintf("Challenge Message %+v", reminderMsg))
				logger.Info(fmt.Sprintf("Webhook URL: %s", sender.WebhookURL))
			} else {
				logger.Info("Sending reminder message")
				err := sender.Send()
				if err != nil {
					logger.Error("Failed to send message", err)
					return err
				}
			}
			return nil
		},
	}
	reminderCmd.Flags().String(flagLogLevel, "info", "Specify the log level (debug, info, warn, error)")
	reminderCmd.Flags().String(flagWebhookUrl, "", "Specify the webhook url")
	reminderCmd.Flags().String(flagCTFName, "", "Specify the CTF name")
	reminderCmd.Flags().String(flagCTFURL, "", "Specify the CTF url")
	reminderCmd.Flags().String(flagStartDate, "", "Specify the start date. Format should be 'YYYY-MM-DD HH:MM:SS'. ")
	reminderCmd.Flags().String(flagEndDate, "", "Specify the end date. Format should be 'YYYY-MM-DD HH:MM:SS'.")
	reminderCmd.Flags().String(flagTz, "Europe/Paris", "Specify the timezone of date in input. EX: 'Europe/London'.")
	reminderCmd.Flags().String(flagTeamName, "", "Specify the team name for this CTF.")
	reminderCmd.Flags().String(flagTeamPassword, "", "Specify the team password for this CTF.")
	return reminderCmd, nil
}
