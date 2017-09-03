package util

import "github.com/transcovo/go-chpr-logger"

func PanicOnError(err error, message string) {
	if err != nil {
		logger.GetLogger().WithError(err).Fatal(message)
	}
}

