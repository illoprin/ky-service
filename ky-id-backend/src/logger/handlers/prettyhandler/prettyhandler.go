package prettyhandler

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"

	"github.com/fatih/color"
)

/*
	type Handler interface {
		Enabled(context.Context, Level) bool
		Handle(context.Context, r Record) error
		WithAttrs(attrs []Attr) Handler
		WithGroup(name string) Handler
	}
*/

type PrettyHandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

type PrettyHandler struct {
	slog.Handler
	l *log.Logger
}

func formatLevel(lvl slog.Level) string {
	var level string = lvl.String()
	switch lvl {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}
	return level
}

func (h *PrettyHandler) Handle(
	_ context.Context, r slog.Record) error {

	// Level
	level := formatLevel(r.Level)

	// Attributes
	fields := make([]string, 0, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		field := fmt.Sprintf("%s: %v", a.Key, a.Value.Any())
		fields = append(fields, field)

		return true
	})

	// Time
	time := r.Time.Format("[15:10:05.000]")

	// Message
	msg := color.CyanString(r.Message)

	h.l.Println()

	// Print message
	h.l.Printf("%s %s %s: %s",
		color.GreenString("[ky-id]"), time, level, msg,
	)

	// Print attributes
	for _, attr := range fields {
		h.l.Println("    " + color.WhiteString(attr))
	}

	return nil
}

func NewPrettyHandler(
	out io.Writer,
	opts *PrettyHandlerOptions,
) *PrettyHandler {
	h := &PrettyHandler{
		l:       log.New(out, "", 0),
		Handler: slog.NewTextHandler(out, &opts.SlogOpts),
	}
	return h
}
