package models

import (
	"log"
	"math"
	"net/url"
	"regexp"
	"strings"

	"github.com/nerdyworm/sourcemap"
)

const (
	NumberOfContextLines = 10
)

type Resource struct {
	URL       string
	Source    string
	SourceMap *sourcemap.Consumer
}

func (r *Resource) baseURL() string {
	return r.URL[0 : strings.LastIndex(r.URL, "/")+1]
}

func (r *Resource) SourceMapURL() string {
	reg := regexp.MustCompile("//# sourceMappingURL=(.*)")
	if reg.MatchString(r.Source) {
		matches := reg.FindStringSubmatch(r.Source)
		return r.baseURL() + matches[1]
	}

	return ""
}

func (r *Resource) Context(lineno int, column int) SourceContext {
	if r.SourceMap == nil {
		return r.contextFromSource(lineno, column)
	} else {
		return r.contextFromSourceMap(lineno, column)
	}
}

func (r *Resource) contextFromSourceMap(lineno, column int) SourceContext {
	source, _, line, col, ok := r.SourceMap.Source(lineno, column)
	if !ok {
		log.Println("not sure why not ok...")
	}

	context := SourceContext{}
	context.OrigLine = line
	context.OrigColumn = col
	context.OrigFile = source

	content := r.SourceMap.SourcesContent(source)

	lines := strings.Split(content, "\n")
	context.PreContext, context.ContextLine, context.PostContext =
		getSourceContext(lines, context.OrigLine, NumberOfContextLines)

	return context
}

func (r *Resource) contextFromSource(lineno, column int) SourceContext {
	lines := strings.Split(r.Source, "\n")
	context := SourceContext{}
	context.PreContext, context.ContextLine, context.PostContext = getSourceContext(lines, lineno, NumberOfContextLines)
	context.OrigLine = lineno
	context.OrigColumn = column

	u, err := url.Parse(r.URL)
	if err != nil {
		log.Println(err)
	}

	context.OrigFile = u.Path
	return context
}

func getSourceContext(lines []string, lineno int, linesOfContext int) ([]string, string, []string) {
	// JavaScript line numbers start from 1
	if lineno > 0 {
		lineno -= 1
	}

	lenLines := len(lines)
	if lenLines == 0 {
		return []string{}, "", []string{}
	}

	if lineno > lenLines {
		lineno = lenLines - 1
	}

	lower := int(math.Max(0, float64(lineno-linesOfContext)))
	upper := int(math.Min(float64(lineno+1+linesOfContext), float64(lenLines)))

	return lines[lower:lineno], lines[lineno], lines[lineno+1 : upper]
}