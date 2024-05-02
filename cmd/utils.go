package cmd

import (
	"strings"

	"github.com/mathildeHermet/bribot/internal/log"
	"github.com/spf13/cobra"
)

func newLogger(cmd *cobra.Command) log.Logger {
	var loggerOptLogLevel log.Option
	switch strings.ToLower(cmd.Flag(flagLogLevel).Value.String()) {
	case "debug":
		loggerOptLogLevel = log.WithLevelDebug()
	case "info":
		loggerOptLogLevel = log.WithLevelInfo()
	case "warn":
		loggerOptLogLevel = log.WithLevelWarn()
	case "error":
		loggerOptLogLevel = log.WithLevelError()
	default:
	}
	return log.NewLogger(
		log.WithAppName("bribot"),
		log.WithServiceName(cmd.Use),
		loggerOptLogLevel,
	)
}
