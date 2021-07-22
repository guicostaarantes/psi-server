package logging

import "log"

type PrintLoggingUtil struct {
}

func (p PrintLoggingUtil) Error(ref string, err error) {
	log.Println(ref, err)
}
