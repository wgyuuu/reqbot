package util

import (
	"runtime"

	"github.com/wgyuuu/reqbot/common/dlog"
)

func SafeFunc(f func()) {
	defer func() {
		if err := recover(); err != nil {
			dlog.Error("panic:%v.\n", err)

			i := 0
			funcName, file, line, ok := runtime.Caller(i)
			for ok {
				dlog.Error("trace:", i, "func:", runtime.FuncForPC(funcName).Name(), "file:", file, "line:", line)
				i++
				funcName, file, line, ok = runtime.Caller(i)
			}
		}
	}()

	f()
}
