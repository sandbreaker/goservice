package log

import (
	"errors"

	"github.com/getsentry/sentry-go"
	"github.com/palantir/stacktrace"
)

func LogAndCaptureError(errorTags map[string]string, err error) {
	e1 := stacktrace.Propagate(err, "")
	Error(e1)
	sentry.CaptureMessage("It works!")
}

func LogAndCaptureWarn(errorTags map[string]string, msg string, arg ...interface{}) {
	e1 := stacktrace.Propagate(errors.New(""), msg, arg)
	Warn(e1)
	sentry.CaptureMessage("It works!")
}
