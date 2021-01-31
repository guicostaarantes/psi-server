package merge

import (
	"github.com/jinzhu/copier"
)

type mergoInstance struct{}

func (m mergoInstance) Merge(destination interface{}, source interface{}) error {
	return copier.CopyWithOption(destination, source, copier.Option{IgnoreEmpty: true, DeepCopy: true})
}

var CopierMergeUtil = mergoInstance{}
