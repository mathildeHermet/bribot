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
	"github.com/spf13/cobra"
)

var (
	verbose = false
	dryrun  = false
	version = "0.0.1"
	rootCmd = &cobra.Command{
		Use:     "bribot",
		Version: version,
		Short:   "Manage discord bot for CTF annoncements.",
	}
	cmds = []func() (*cobra.Command, error){
		NewReminderCmd,
	}
)

func Run() error {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Make the operation more talkative.")
	rootCmd.PersistentFlags().BoolVarP(&dryrun, "dry-run", "d", false, "Do not perform the operation, only show what would be sent.")
	for i := 0; i < len(cmds); i++ {
		cmd, err := cmds[i]()
		if err != nil {
			return err
		}
		rootCmd.AddCommand(cmd)
	}
	return rootCmd.Execute()
}
