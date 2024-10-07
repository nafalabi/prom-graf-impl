// credit prometheus alertmanager default template
package template

import (
	"fmt"
	tmplhtml "html/template"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func ConvertToFloat(i interface{}) (float64, error) {
	switch v := i.(type) {
	case float64:
		return v, nil
	case string:
		return strconv.ParseFloat(v, 64)
	case int:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case time.Duration:
		return v.Seconds(), nil
	default:
		return 0, fmt.Errorf("can't convert %T to float", v)
	}
}

var DefaultFuncs = map[string]interface{}{
	"toUpper": strings.ToUpper,
	"toLower": strings.ToLower,
	"title": func(text string) string {
		// Casers should not be shared between goroutines, instead
		// create a new caser each time this function is called.
		return cases.Title(language.AmericanEnglish).String(text)
	},
	"trimSpace": strings.TrimSpace,
	// join is equal to strings.Join but inverts the argument order
	// for easier pipelining in templates.
	"join": func(sep string, s []string) string {
		return strings.Join(s, sep)
	},
	"match": regexp.MatchString,
	"safeHtml": func(text string) tmplhtml.HTML {
		return tmplhtml.HTML(text)
	},
	"reReplaceAll": func(pattern, repl, text string) string {
		re := regexp.MustCompile(pattern)
		return re.ReplaceAllString(text, repl)
	},
	"stringSlice": func(s ...string) []string {
		return s
	},
	// date returns the text representation of the time in the specified format.
	"date": func(fmt string, t time.Time) string {
		return t.Format(fmt)
	},
	// tz returns the time in the timezone.
	"tz": func(name string, t time.Time) (time.Time, error) {
		loc, err := time.LoadLocation(name)
		if err != nil {
			return time.Time{}, err
		}
		return t.In(loc), nil
	},
	"since": time.Since,
	"humanizeDuration": func(i interface{}) (string, error) {
		v, err := ConvertToFloat(i)
		if err != nil {
			return "", err
		}

		if math.IsNaN(v) || math.IsInf(v, 0) {
			return fmt.Sprintf("%.4g", v), nil
		}
		if v == 0 {
			return fmt.Sprintf("%.4gs", v), nil
		}
		if math.Abs(v) >= 1 {
			sign := ""
			if v < 0 {
				sign = "-"
				v = -v
			}
			duration := int64(v)
			seconds := duration % 60
			minutes := (duration / 60) % 60
			hours := (duration / 60 / 60) % 24
			days := duration / 60 / 60 / 24
			// For days to minutes, we display seconds as an integer.
			if days != 0 {
				return fmt.Sprintf("%s%dd %dh %dm %ds", sign, days, hours, minutes, seconds), nil
			}
			if hours != 0 {
				return fmt.Sprintf("%s%dh %dm %ds", sign, hours, minutes, seconds), nil
			}
			if minutes != 0 {
				return fmt.Sprintf("%s%dm %ds", sign, minutes, seconds), nil
			}
			// For seconds, we display 4 significant digits.
			return fmt.Sprintf("%s%.4gs", sign, v), nil
		}
		prefix := ""
		for _, p := range []string{"m", "u", "n", "p", "f", "a", "z", "y"} {
			if math.Abs(v) >= 1 {
				break
			}
			prefix = p
			v *= 1000
		}
		return fmt.Sprintf("%.4g%ss", v, prefix), nil
	},
}

func GetTemplate() *tmplhtml.Template {
	tmpl := tmplhtml.Must(
		tmplhtml.New("email.html").Funcs(DefaultFuncs).ParseFiles(
			"./template/email.html",
      // "./template/email.tmpl",
			"./template/default.tmpl",
		),
	)

	return tmpl
}
