package merge

type IMergeUtil interface {
	Merge(destination interface{}, source interface{}) error
}
