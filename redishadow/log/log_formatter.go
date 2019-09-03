package log

import "strings"

// log formatter
type logFormatter struct {
	providers []logPartProvider
}

type logPartProvider interface {
	consume(v interface{}) (string, bool)
}

func newFormatter(format string) *logFormatter {
	l := &logFormatter{
		providers: make([]logPartProvider, 0),
	}
	var preRune rune
	builder := strings.Builder{}
	for i, r := range format {
		if preRune == '%' {
			switch r {
			case 'T':
				l.providers = append(l.providers, nil)
			case 'y':
				l.providers = append(l.providers, nil)
			case 'M':
				l.providers = append(l.providers, nil)
			case 'd':
				l.providers = append(l.providers, nil)
			case 'h':
				l.providers = append(l.providers, nil)
			case 'm':
				l.providers = append(l.providers, nil)
			case 's':
				l.providers = append(l.providers, nil)
			case 'L':
				l.providers = append(l.providers, nil)
			case 'S':
				l.providers = append(l.providers, nil)
			default:
				fatalF("Invalid control character at pos %d of `%s`", i, format)
			}
		} else if r == '%' {
			l.providers = append(l.providers, nil)
			builder.Reset()
		} else {
			builder.WriteRune(r)
		}
		preRune = r
	}
	return l
}

/*
ctrl literal:
%T timestamp
%y years
%M months
%d days
%h hours
%m minutes
%s seconds
%L log level
%S placeholder
*/
func (l *logFormatter) format(v ...interface{}) []byte {
	builder := strings.Builder{}
	i := 0
	tmp := v[i]
	for _, provider := range l.providers {
		part, ok := provider.consume(tmp)
		builder.WriteString(part)
		if ok {
			if i < len(v) {
				i++
				tmp = v[i]
			} else {
				tmp = nil
			}
		}
	}
	return []byte(builder.String())
}
