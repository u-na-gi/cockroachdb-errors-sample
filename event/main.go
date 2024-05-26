package main

import (
	"net/http"

	sentrygo "github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/u-na-gi/cockroachdb-errors-sample/sentry"
)

func main() {
	sentry.Init()

	// パニックをキャプチャ
	defer sentrygo.Recover()

	app := echo.New()

	app.Use(middleware.Logger())
	app.Use(middleware.Recover())

	// Once it's done, you can attach the handler as one of your middleware
	app.Use(sentryecho.New(sentryecho.Options{}))

	// Set up routes
	app.GET("/", func(ctx echo.Context) error {

		// errors.ReportError(errors.New("This is an error"))

		event := sentrygo.NewEvent()
		uid := uuid.NewString()
		event.Message = "Custom Event with differences" + uid
		event.Level = sentrygo.LevelInfo

		event.Extra = map[string]interface{}{
			"uid": uid,
		}

		// カスタムイベントの送信
		// sentrygo.CaptureEvent(event)

		// 追加の情報を含める
		event.Extra["expect_row_json"] = `{"id": 1, "name": "John Doe", "age": 20, "email": "", "address": "fsdfsfs"}`
		event.Extra["actual_row_json"] = `{"id": 1, "name": "alice", "age": 20}`

		// カスタムイベントの送信
		// sentrygo.CaptureEvent(event)

		// さらに追加の情報を含める
		event.Extra["additional"] = "This is additional informationnnnnnn"

		// カスタムイベントの送信
		sentrygo.CaptureEvent(event)

		return ctx.String(http.StatusNotFound, "an error")
	})

	// And run it
	// curl -v localhost:3000
	app.Logger.Fatal(app.Start(":3000"))

}
