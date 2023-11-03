package httplog

import (
	"log/slog"
	"strings"
	"time"
)

var defaultOptions = Options{
	LevelFieldName:     "level",
	Concise:            true,
	Tags:               nil,
	RequestHeaders:     true,
	HideRequestHeaders: nil,
	QuietDownRoutes:    nil,
	QuietDownPeriod:    0,
	MessageFieldName:   "message",
}

type Options struct {
	// LevelFieldName sets the field name for the log level or severity.
	// Some providers parse and search for different field names.
	LevelFieldName string

	// MessageFieldName sets the field name for the message.
	// Default is "msg".
	MessageFieldName string

	// Pretty enables pretty printing of the stacktraces.
	Pretty bool

	// Concise mode includes fewer log details during the request flow. For example
	// excluding details like request content length, user-agent and other details.
	// This is useful if during development your console is too noisy.
	Concise bool

	// Tags are additional fields included at the root level of all logs.
	// These can be useful for example the commit hash of a build, or an environment
	// name like prod/stg/dev
	Tags map[string]string

	// RequestHeaders enables logging of all request headers, however sensitive
	// headers like authorization, cookie and set-cookie are hidden.
	RequestHeaders bool

	// HideRequestHeaders are additional requests headers which are redacted from the logs
	HideRequestHeaders []string

	// ResponseHeaders enables logging of all response headers.
	ResponseHeaders bool

	// QuietDownRoutes are routes which are temporarily excluded from logging for a QuietDownPeriod after it occurs
	// for the first time
	// to cancel noise from logging for routes that are known to be noisy.
	QuietDownRoutes []string

	// QuietDownPeriod is the duration for which a route is excluded from logging after it occurs for the first time
	// if the route is in QuietDownRoutes
	QuietDownPeriod time.Duration

	// TimeFieldName sets the field name for the time field.
	// Some providers parse and search for different field names.
	TimeFieldName string

	// SourceFieldName sets the field name for the source field which logs
	// the location in the program source code where the logger was called.
	// If set to "" then it'll be disabled.
	SourceFieldName string
}

// Configure will set new options for the httplog instance and behaviour
// of underlying slog pkg and its global logger.
func (l *Logger) Configure(opts Options) {
	// if opts.LogLevel is not set
	// it would be 0 which is LevelInfo

	if opts.LevelFieldName == "" {
		opts.LevelFieldName = "level"
	}

	if len(opts.QuietDownRoutes) > 0 {
		if opts.QuietDownPeriod == 0 {
			opts.QuietDownPeriod = 5 * time.Minute
		}
	}

	// Pre-downcase all SkipHeaders
	for i, header := range opts.HideRequestHeaders {
		opts.HideRequestHeaders[i] = strings.ToLower(header)
	}

	l.Options = opts
}

func LevelByName(name string) slog.Level {
	switch strings.ToUpper(name) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return 0
	}
}
