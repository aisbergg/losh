package log

// CommonLogConfig contains the common log configuration options.
type CommonLogConfig struct {
	FileLogConfig

	// Format defines the encoding format of the logs.
	Format string
}

// FileLogConfig contains the configuration options for the file output.
type FileLogConfig struct {
	// Filename is the file to write logs to. Backup log files will be retained
	// in the same directory. If empty logs will be written to stderr instead of
	// a file.
	Filename string

	// Rotate determines if the file should be automatically rotated with the
	// given parameters.
	Rotate bool

	// MaxSize is the maximum size in megabytes of the log file before it gets
	// rotated.
	MaxSize int

	// MaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	MaxAge int

	// MaxBackups is the maximum number of old log files to retain.
	MaxBackups int

	// LocalTime determines if the time used for formatting the timestamps in
	// backup files is the computer's local time.
	LocalTime bool

	// Compress determines if the rotated log files should be compressed
	// using gzip.
	Compress bool

	// Permission defines the mode which is used to create the file.
	Permissions uint32
}
