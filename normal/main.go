package main

import (
	"net/http"

	"github.com/cockroachdb/errors"
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

		errors.ReportError(errors.New("This is an error"))

		return ctx.String(http.StatusNotFound, "an error")
	})

	// And run it
	// curl -v localhost:3000
	app.Logger.Fatal(app.Start(":3000"))

}
