package log

import (
	"context"
	"io"
	"log/slog"
	"maps"
	"slices"
	"sync"
	"text/template"

	"github.com/fatih/color"
)

var logTemplate = `{{- .time}} {{.level}} {{if .prefix}}[{{.prefix}}]{{end}} {{.message}} {{range .attrs}} {{.Key}}={{.Value.String}} {{end -}}
{{range $i, $group := .groups}} {{range $j, $attr := index $.gAttrs $i}} {{"\n"}}---  {{- $group }}  {{$attr.Key -}} = {{- $attr.Value.String}} {{end}} {{end}}`

const timeFormat = "2006-01-02 15:04:05"

var levelStringColorsMap = map[slog.Level]string{
	slog.LevelDebug: color.New(color.FgGreen, color.Bold).Sprint("DEBUG"),
	slog.LevelInfo:  color.New(color.FgCyan, color.Bold).Sprint("INFO"),
	slog.LevelWarn:  color.New(color.FgYellow, color.Bold).Sprint("WARN"),
	slog.LevelError: color.New(color.FgRed, color.Bold, color.Underline).Sprint("ERROR"),
}

type Handler struct {
	stdout io.Writer
	stderr io.Writer

	prefix string
	groups []string
	gAttrs map[int][]slog.Attr
	attrs  []slog.Attr
	tmplt  *template.Template
	mu     *sync.Mutex
}

func NewHandler(stdout, stderr io.Writer, prefix ...string) *Handler {
	var p string
	if len(prefix) > 0 {
		p = prefix[0]
	}
	return &Handler{
		stdout: stdout,
		stderr: stderr,
		prefix: p,
		gAttrs: make(map[int][]slog.Attr),
		tmplt:  template.Must(template.New("log").Parse(logTemplate)),
		mu:     new(sync.Mutex),
	}
}

func (*Handler) Enabled(_ context.Context, l slog.Level) bool {
	return l >= slog.LevelInfo
}

func (h *Handler) WithPrefix(prefix string) slog.Handler {
	h2 := h.clone()
	h2.prefix = prefix
	return h2
}

func (h *Handler) WithGroup(name string) slog.Handler {
	h2 := h.clone()
	h2.groups = append(h2.groups, name)
	h2.gAttrs[len(h2.groups)-1] = make([]slog.Attr, 0)
	return h2
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h2 := h.clone()
	if len(h2.groups) == 0 {
		h2.attrs = append(h2.attrs, attrs...)

	} else {
		h2.gAttrs[len(h2.groups)-1] = append(h2.gAttrs[len(h2.groups)-1], attrs...)
	}
	return h2
}

func (h *Handler) Handle(_ context.Context, r slog.Record) error {
	attrs := slices.Clone(h.attrs)
	groups := slices.Clone(h.groups)
	gAttrs := maps.Clone(h.gAttrs)
	r.Attrs(func(attr slog.Attr) bool {
		if attr.Value.Kind() == slog.KindGroup {
			groups = append(groups, attr.Key)
			gAttrs[len(groups)-1] = slices.Clone(attr.Value.Group())
			return true
		}
		if len(h.groups) > 0 {
			gAttrs[len(h.groups)-1] = append(gAttrs[len(h.groups)-1], attr)
			return true
		}
		attrs = append(attrs, attr)
		return true
	})

	if r.Level >= slog.LevelError {
		return h.tmplt.Execute(h.stderr, map[string]interface{}{
			"time":    r.Time.Format(timeFormat),
			"level":   levelStringColorsMap[r.Level],
			"prefix":  h.prefix,
			"message": r.Message,
			"attrs":   attrs,
			"groups":  groups,
			"gAttrs":  gAttrs,
		})
	}

	return h.tmplt.Execute(h.stdout, map[string]interface{}{
		"time":    r.Time.Format(timeFormat),
		"level":   levelStringColorsMap[r.Level],
		"prefix":  h.prefix,
		"message": r.Message,
		"attrs":   attrs,
		"groups":  groups,
		"gAttrs":  gAttrs,
	})
}
func (h *Handler) clone() *Handler {
	return &Handler{
		stdout: h.stdout,
		stderr: h.stderr,
		prefix: h.prefix,
		groups: slices.Clone(h.groups),
		gAttrs: maps.Clone(h.gAttrs),
		attrs:  slices.Clone(h.attrs),
		tmplt:  h.tmplt,
		mu:     h.mu, // mutex shared among all clones of this handler
	}
}
