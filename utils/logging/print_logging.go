package logging

import "log"

type printLogUtil struct {
}

func (p printLogUtil) Error(ref string, err error) {
	log.Println(ref, err)
}

// PrintLogUtil is an implementation of ILoggingUtil that prints to the console of the Go executable
var PrintLogUtil = printLogUtil{}
