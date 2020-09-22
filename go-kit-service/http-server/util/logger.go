package util

import (
	kitlog "github.com/go-kit/kit/log"
	"os"
)

func GetLogger () kitlog.Logger {
	var logger kitlog.Logger
	{
		logger = kitlog.NewLogfmtLogger(os.Stdout)
		logger = kitlog.WithPrefix(logger, "micro-srv", "1.0")
		logger = kitlog.With(logger, "time", kitlog.DefaultTimestampUTC)
		logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)
	}
	return logger
}

