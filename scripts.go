package main

type Execution struct {
	Order int    `json:"execution_order"`
	UUID  string `json:"uuid"`
	Script Script `json:"script"`
}

type Script struct {
	Code string `json:"code"`
	UUID string `json:"uuid"`
}