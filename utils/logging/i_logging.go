package logging

// ILoggingUtil is an abstraction for a utility that logs server events like errors for analysis
type ILoggingUtil interface {
	Error(ref string, err error)
}
