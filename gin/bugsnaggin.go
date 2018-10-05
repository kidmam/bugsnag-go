package bugsnaggin

import (
	"github.com/bugsnag/bugsnag-go"
	"github.com/gin-gonic/gin"
)

const FrameworkName string = "Gin"

// AutoNotify sends any panics to bugsnag, and then re-raises them.
// You should use this after another middleware that
// returns an error page to the client, for example gin.Recovery().
// The arguments can be any RawData to pass to Bugsnag, most usually
// you'll pass a bugsnag.Configuration object.
func AutoNotify(rawData ...interface{}) gin.HandlerFunc {
	// Configure bugsnag with the passed in configuration (for manual notifications)
	for _, datum := range rawData {
		if c, ok := datum.(bugsnag.Configuration); ok {
			bugsnag.Configure(c)
		}
	}

	state := bugsnag.HandledState{
		SeverityReason:   bugsnag.SeverityReasonUnhandledMiddlewareError,
		OriginalSeverity: bugsnag.SeverityError,
		Unhandled:        true,
		Framework:        FrameworkName,
	}
	rawData = append(rawData, state)
	return func(c *gin.Context) {
		r := c.Copy().Request
		ctx := r.Context()
		ctx = bugsnag.StartSession(ctx)
		ctx = bugsnag.AttachRequestData(ctx, r)
		c.Request = r.WithContext(ctx)

		// create a notifier that has the current request bound to it
		notifier := bugsnag.New(append(rawData, r)...)
		notifier.FlushSessionsOnRepanic(false)
		defer notifier.AutoNotify(ctx, r)
		c.Next()
	}
}
