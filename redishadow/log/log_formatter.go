package log

// log formatter
type logFormatter struct {
}

func newFormatter(format string) *logFormatter {
	return &logFormatter{}
}

func (f *logFormatter) format(v interface{}) []byte {
	return nil
}
