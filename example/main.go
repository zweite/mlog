package main

import (
	"github.com/zweite/mlog"
)

func main() {
	mlog.SetLogger(mlog.NewLogger(mlog.LogConfig{
		LogLevel:   "info",
		Buffer:     100,
		Dir:        "/data/logs/mimi",
		SubRelPath: "2006-01/2006-01-02/2006-01-02-15.log",
		ServerIP:   "111.111.111.111",
	}))

	defer mlog.Close()
	defer mlog.Flush()

	for index := 0; index < 1000000; index++ {
		mlog.GetMlog().Record("mq_product", "topic"+"\t"+"msg"+"\t"+"error_msg")
	}
}
