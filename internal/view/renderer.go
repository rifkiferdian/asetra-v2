package view

import (
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin/render"
)

func formatDate(value string) string {
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return value
	}
	return formatIndonesianDate(parsed)
}

func formatDateTime(value string) string {
	parsed, err := time.Parse("2006-01-02 15:04", value)
	if err != nil {
		return value
	}
	return formatIndonesianDate(parsed) + parsed.Format(", 15:04")
}

func formatIndonesianDate(value time.Time) string {
	months := [...]string{"Jan", "Feb", "Mar", "Apr", "Mei", "Jun", "Jul", "Agu", "Sep", "Okt", "Nov", "Des"}
	return fmt.Sprintf("%02d %s %04d", value.Day(), months[value.Month()-1], value.Year())
}

func formatRupiah(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	sign := ""
	if strings.HasPrefix(value, "-") {
		sign = "-"
		value = strings.TrimPrefix(value, "-")
	}
	parts := strings.SplitN(value, ".", 2)
	integer := parts[0]
	if integer == "" {
		integer = "0"
	}
	for i := len(integer) - 3; i > 0; i -= 3 {
		integer = integer[:i] + "." + integer[i:]
	}
	fraction := ""
	if len(parts) == 2 && strings.TrimRight(parts[1], "0") != "" {
		fraction = "," + strings.TrimRight(parts[1], "0")
	}
	return "Rp " + sign + integer + fraction
}

func dict(values ...any) map[string]any {
	result := make(map[string]any, len(values)/2)
	for i := 0; i+1 < len(values); i += 2 {
		key, ok := values[i].(string)
		if ok {
			result[key] = values[i+1]
		}
	}
	return result
}

// Each page has its own template set, so content blocks cannot overwrite other pages.
type Renderer struct{ pages map[string]*template.Template }

func New(root string) (*Renderer, error) {
	pages, err := filepath.Glob(filepath.Join(root, "*", "*.html"))
	if err != nil {
		return nil, err
	}
	r := &Renderer{pages: make(map[string]*template.Template)}
	for _, page := range pages {
		folder := filepath.Base(filepath.Dir(page))
		if folder == "layouts" || folder == "partials" {
			continue
		}
		t, err := template.New("").Funcs(template.FuncMap{
			"dict":           dict,
			"addOne":         func(value int) int { return value + 1 },
			"sub":            func(a, b float64) float64 { return a - b },
			"formatDate":     formatDate,
			"formatDateTime": formatDateTime,
			"formatRupiah":   formatRupiah,
		}).ParseGlob(filepath.Join(root, "layouts", "*.html"))
		if err != nil {
			return nil, fmt.Errorf("parse layouts: %w", err)
		}
		t, err = t.ParseGlob(filepath.Join(root, "partials", "*.html"))
		if err != nil {
			return nil, fmt.Errorf("parse partials: %w", err)
		}
		t, err = t.ParseFiles(page)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", page, err)
		}
		r.pages[folder+"/"+filepath.Base(page)] = t
	}
	if len(r.pages) == 0 {
		return nil, fmt.Errorf("no page templates found in %s", root)
	}
	return r, nil
}

func (r *Renderer) Instance(name string, data any) render.Render {
	t, ok := r.pages[name]
	if !ok {
		panic("unknown page template: " + name)
	}
	return render.HTML{Template: t, Name: "base", Data: data}
}
