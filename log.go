package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"

	"github.com/fatih/color"
)

type PrettyLogHandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

type PrettyLogHandler struct {
	slog.Handler
	l *log.Logger
}

const TIMESTAMP_FORMAT = "15:04:05.000"

func (h *PrettyLogHandler) Handle(ctx context.Context, r slog.Record) error {
	timeStr := r.Time.Format(TIMESTAMP_FORMAT)

	// Level
	level := r.Level.String()
	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	// Message
	msg := color.CyanString(r.Message)

	// Additional fields
	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()

		return true
	})
	addlData := ""
	if len(fields) > 0 {
		data, err := json.Marshal(fields)
		if err != nil {
			return err
		}
		addlData = color.WhiteString(string(data))
	}

	h.l.Println(timeStr, level, msg, addlData)

	return nil
}

func NewPrettyLogHandler(out io.Writer, opts PrettyLogHandlerOptions) *PrettyLogHandler {
	return &PrettyLogHandler{
		Handler: slog.NewJSONHandler(out, &opts.SlogOpts),
		l:       log.New(out, "", 0),
	}
}
