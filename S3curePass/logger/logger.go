package logger

import (
	"fmt"
	"time"
)

//logs an info
func Log(msg string) {
	fmt.Println(getTime() + "(INF) -> " + msg)
}

//logs a warning
func Warning(msg string) {
	fmt.Println(getTime() + "(WAR) -> " + msg)
}

//logs en error message
func Error(msg string) {
	fmt.Println(getTime() + "(ERR) -> " + msg)
}

//logs en error message
func Fatal(msg string) {
	fmt.Println(getTime() + "(FAT) -> " + msg)
}

//returns the prefix for the logger (the current date + time)
func getTime() string {
	time := time.Now()
	return "[" + time.Format("Mon, 02 Jan 2006 15:04:05 MST") + "] "
}
