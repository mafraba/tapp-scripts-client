package main

import "time"


type Execution struct {
	Order  int    `json:"execution_order"`
	UUID   string `json:"uuid"`
	Script Script `json:"script"`
}

type Script struct {
	Code string `json:"code"`
	UUID string `json:"uuid"`
}


func ExecCode(code string) (output string, exitCode int, startedAt time.Time, finishedAt time.Time) {
	output = ""
	exitCode = 0
	startedAt = time.Now()
	finishedAt = time.Now()
	return
}
