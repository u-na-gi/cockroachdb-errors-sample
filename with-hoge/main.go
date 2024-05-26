package main

import (
	"net/http"

	"github.com/cockroachdb/errors"
	sentrygo "github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/u-na-gi/cockroachdb-errors-sample/sentry"
)

func main() {
	sentry.Init()

	// Then create your app
	app := echo.New()

	app.Use(middleware.Logger())
	app.Use(middleware.Recover())

	// Once it's done, you can attach the handler as one of your middleware
	app.Use(sentryecho.New(sentryecho.Options{}))

	// Set up routes
	app.GET("/", func(ctx echo.Context) error {

		err := errors.New("This is an error with sentry by with SetTag and SetExtras")
		// err = errors.WithMessage(err, "error report with sentry")
		// err = errors.WithHint(err, "this error has a hint")
		// err = errors.WithDetail(err, "this error has a detail")
		// errors.ReportError(err)

		sentrygo.WithScope(func(scope *sentrygo.Scope) {

			scope.SetTag("kind", "Youtube Comment")
			scope.SetExtras(map[string]interface{}{
				"expected_comment_struct": struct {
					ID   string
					Mail string
				}{},
				"actual_comment_struct": struct {
					ID       string
					UserMail string
				}{ID: "54321",
					UserMail: "user_mail@ggg.com"},
			})

			sentrygo.CaptureException(err)

		})

		return ctx.String(http.StatusNotFound, "an error")
	})

	// And run it
	// curl -v localhost:3000
	app.Logger.Fatal(app.Start(":3000"))

}
