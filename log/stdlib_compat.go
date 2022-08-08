package log

import (
	"fmt"
	"os"
)

func Printf(format string, v ...interface{}) {
	DefaultLogger.Info(fmt.Sprintf(format, v...))
}

func Println(v ...interface{}) {
	DefaultLogger.Info(fmt.Sprint(v...))
}

func Fatalf(format string, v ...interface{}) {
	DefaultLogger.Error(nil, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func Fatalln(v ...interface{}) {
	DefaultLogger.Error(nil, fmt.Sprint(v...))
	os.Exit(1)
}
