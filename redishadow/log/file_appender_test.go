package log

import "testing"

const (
	logPath = "./"
)

func TestAll(t *testing.T) {
	f := newFileAppender("./", "test-log-", "2MB")
	if f.logPath != "./"
	if f.maxSingleFileSize != 2 * 1024 * 1024 {
		t.Errorf("maxSingleFileSize should be 2 * 1024 * 1024, got %d", f.maxSingleFileSize)
	}
}