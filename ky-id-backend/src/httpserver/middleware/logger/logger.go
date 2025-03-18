package logger

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/fatih/color"
	mw "github.com/go-chi/chi/v5/middleware"
)

func formatStatus(status int) string {
	if status >= 100 && status < 200 {
		return color.Set(color.Bold, color.FgYellow).Sprint(status)
	} else if status >= 200 && status < 300 {
		return color.Set(color.Bold, color.FgGreen).Sprint(status)
	} else if status >= 300 && status < 400 {
		return color.Set(color.Bold, color.FgBlue).Sprint(status)
	} else if status >= 400 && status < 500 {
		return color.Set(color.Bold, color.FgRed).Sprint(status)
	} else if status >= 500 && status <= 511 {
		return color.Set(color.Bold, color.FgMagenta).Sprint(status)
	}
	return strconv.Itoa(status)
}

// Log request some info
func LoggerMW(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		slog.Info("new request",
			slog.String("method", color.CyanString(r.Method)),
			slog.String("path", color.CyanString(r.URL.Path)),
			slog.String("remote_addr", color.BlueString(r.RemoteAddr)),
		)

		wrapper := mw.NewWrapResponseWriter(w, r.ProtoMajor)

		t1 := time.Now()

		defer func() {
			slog.Info("request completed",
				slog.String("bytes", color.Set(color.Bold, color.FgWhite).Sprint(wrapper.BytesWritten())),
				slog.String("status", formatStatus(wrapper.Status())),
				slog.String("duration", color.Set(color.Bold, color.FgCyan).Sprint(time.Since(t1).String())),
			)
		}()

		next.ServeHTTP(wrapper, r)
	}
	return http.HandlerFunc(fn)

}
