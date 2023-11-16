package logging

import (
	"context"
	"fmt"
	"os"
	"strings"

	"log/slog"
)

func NewHandler(group string, level slog.Level) slog.Handler {
	return &Handler{
		group: group,
		level: level,
		attrs: make(map[string]slog.Value),
	}
}

type Handler struct {
	level slog.Level
	group string
	attrs map[string]slog.Value
}

func (h *Handler) Enabled(ctx context.Context, lvl slog.Level) bool {
	return lvl >= h.level
}

// Handle handles the Record.
// It will only be called Enabled returns true.
// The Context argument is as for Enabled.
// It is present solely to provide Handlers access to the context's values.
// Canceling the context should not affect record processing.
// (Among other things, log messages may be necessary to debug a
// cancellation-related problem.)
//
// Handle methods that produce output should observe the following rules:
//   - If r.Time is the zero time, ignore the time.
//   - If r.PC is zero, ignore it.
//   - Attr's values should be resolved.
//   - If an Attr's key and value are both the zero value, ignore the Attr.
//     This can be tested with attr.Equal(Attr{}).
//   - If a group's key is empty, inline the group's Attrs.
//   - If a group has no Attrs (even if it has a non-empty key),
//     ignore it.
func (h *Handler) Handle(ctx context.Context, record slog.Record) error {
	sb := new(strings.Builder)
	indent := strings.Count(h.group, ".")
	if h.group != "" && indent == 0 {
		indent = 1
	}
	sb.WriteString(strings.Repeat("  ", indent))
	sb.WriteString(fmt.Sprintf("%-40s", record.Message))
	for k, v := range h.attrs {
		h.writeAttr(sb, slog.Any(k, v))
	}
	record.Attrs(func(a slog.Attr) bool {
		h.writeAttr(sb, a)
		return true
	})
	sb.WriteRune('\n')
	os.Stdout.WriteString(sb.String())
	return nil
}
func (h *Handler) writeAttr(sb *strings.Builder, a slog.Attr) {
	sb.WriteString(" ")
	if h.group != "" {
		sb.WriteString(h.group)
		sb.WriteRune('.')
	}
	sb.WriteString(a.Key)
	sb.WriteRune('=')
	if val := a.Value.String(); strings.ContainsAny(val, " \t\n") {
		sb.WriteRune('"')
		sb.WriteString(val)
		sb.WriteRune('"')
	} else {
		sb.WriteString(val)
	}
}
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	nattrs := make(map[string]slog.Value)
	for k, v := range h.attrs {
		nattrs[k] = v
	}
	for _, attr := range attrs {
		nattrs[attr.Key] = attr.Value
	}
	return &Handler{
		level: h.level,
		group: h.group,
		attrs: nattrs,
	}
}

// WithGroup returns a new Handler with the given group appended to
// the receiver's existing groups.
// The keys of all subsequent attributes, whether added by With or in a
// Record, should be qualified by the sequence of group names.
//
// How this qualification happens is up to the Handler, so long as
// this Handler's attribute keys differ from those of another Handler
// with a different sequence of group names.
//
// A Handler should treat WithGroup as starting a Group of Attrs that ends
// at the end of the log event. That is,
//
//	logger.WithGroup("s").LogAttrs(level, msg, slog.Int("a", 1), slog.Int("b", 2))
//
// should behave like
//
//	logger.LogAttrs(level, msg, slog.Group("s", slog.Int("a", 1), slog.Int("b", 2)))
//
// If the name is empty, WithGroup returns the receiver.
func (h *Handler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	var group string
	if h.group != "" {
		group = h.group + "."
	}
	return &Handler{
		level: h.level,
		group: group + name,
		attrs: h.attrs,
	}
}
