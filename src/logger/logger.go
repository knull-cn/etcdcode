package logger

import (
	"fmt"
	"time"
)

func timestr() string {
	return time.Now().Format("2006-01-02 15:04:05.000000000")
}

func LOGGING(hdr, fmat string, arg []interface{}) {
	fmt.Println(timestr(), "]", hdr, fmt.Sprintf(fmat, arg...))
}

func LogDbg(fmat string, arg ...interface{}) {
	LOGGING("debug : ", fmat, arg)
}

func LogInfo(fmat string, arg ...interface{}) {
	LOGGING(" info : ", fmat, arg)
}

func LogErr(fmat string, arg ...interface{}) {
	LOGGING(" err : ", fmat, arg)
}
