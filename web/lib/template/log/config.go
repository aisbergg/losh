package log

import "losh/lib/log"

type AppLogConfig struct {
	log.CommonLogConfig
	Level string
}

type AccessLogConfig struct {
	log.CommonLogConfig
	Enabled bool
	Fields  []string
}
