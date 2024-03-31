package main

import (
	"log"
	"time"

	sentry "github.com/getsentry/sentry-go"
	"github.com/ryomak/serrs"
)

const DefaultCodeUnexpected serrs.DefaultCode = "unexpected"

func main() {

	err := sentry.Init(sentry.ClientOptions{
		Dsn: "",
		// Enable printing of SDK debug messages.
		// Useful when getting started or trying to figure something out.
		Debug: true,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)

	if err := Do4(); err != nil {
		serrs.ReportSentry(
			err,
			serrs.WithSentryContexts(map[string]sentry.Context{
				"custom": map[string]any{
					"key": "value",
				},
			}),
			serrs.WithSentryTags(map[string]string{
				"code": serrs.GetErrorCodeString(err),
			}),
			serrs.WithSentryLevel(sentry.LevelInfo),
		)
	}
}

func Do() error {
	return serrs.New(serrs.DefaultCode("do_err"), "unexpected do error")
}

func Do2() error {
	if err := Do(); err != nil {
		return serrs.Wrap(err)
	}
	return nil
}

func Do3() error {
	if err := Do2(); err != nil {
		return serrs.Wrap(
			err,
			serrs.WithCode(DefaultCodeUnexpected),
			serrs.WithMessage("do2 error"),
			serrs.WithData(serrs.DefaultCustomData{
				"key": "value",
			}),
		)
	}
	return nil
}

func Do4() error {
	if err := Do3(); err != nil {
		return serrs.Wrap(
			err,
			serrs.WithCode(DefaultCodeUnexpected),
			serrs.WithMessage("Do4 error"),
			serrs.WithData(serrs.DefaultCustomData{
				"userName": "hoge",
			}),
		)
	}
	return nil
}
