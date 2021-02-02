package merge

// IMergeUtil is an abstraction for a utility that parses the keys of a source interface to a destination interface
type IMergeUtil interface {
	Merge(destination interface{}, source interface{}) error
}
