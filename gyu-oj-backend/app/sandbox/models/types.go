package models

type ExecResult struct {
	StdOut string
	StdErr string
	Time   int64
	Memory int64
}
