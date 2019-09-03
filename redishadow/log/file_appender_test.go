package log

import (
	"io/ioutil"
	"testing"
)

const (
	logPath           = "./test-output"
	logFileNamePrefix = "test-log-"
	maxSingleFileSize = "0.5kb"
	testLog           = "test log"
)

func TestAll(t *testing.T) {
	f := newFileAppender(logPath, logFileNamePrefix, maxSingleFileSize)
	defer func() {

	}()
	if f.logPath != logPath {
		t.Errorf("logPath should be `%s`, got `%s`", logPath, f.logPath)
	}
	if f.logFileNamePrefix != logFileNamePrefix {
		t.Errorf("logFileNamePrefix should be `%s`, got `%s`", logPath, f.logPath)
	}
	if f.maxSingleFileSize != 0.5*1024 {
		t.Errorf("maxSingleFileSize should be 2*1024*1024, got %d", f.maxSingleFileSize)
	}
	if f.current == nil {
		t.Error("failed to open log file")
	}

	if err := f.Write([]byte(testLog)); err != nil {
		t.Error(err)
	}

	f.Close()
	if f.current != nil {
		t.Error("log file should be closed")
	}

	if data, err := ioutil.ReadFile(f.getLogFilePath()); err != nil {
		t.Error(err)
	} else if string(data) != testLog {
		t.Errorf("log output should be `%s`, got `%s`", testLog, string(data))
	}

	// check if file_appender creates the log path correctlly
}
