package log_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/vanillazen/stl/backend/internal/sys/log"
)

// MockLogger is a mock implementation of the Logger interface for testing.
type MockLogger struct {
	debug *bytes.Buffer
	info  *bytes.Buffer
	error *bytes.Buffer
}

func NewMockLogger() *MockLogger {
	return &MockLogger{
		debug: new(bytes.Buffer),
		info:  new(bytes.Buffer),
		error: new(bytes.Buffer),
	}
}

func TestSimpleLoggerDebug(t *testing.T) {
	output := NewMockLogger()

	sl := log.NewLogger("debug")
	sl.SetDebugOutput(output.debug)

	sl.Debug("debug message")
	expectedOutput := "debug message\n"
	actualOutput := output.debug.String()
	if actualOutput != expectedOutput {
		t.Errorf("Expected debug output:\n%s\nBut got:\n%s", expectedOutput, actualOutput)
	}
}

func TestSimpleLoggerDebugf(t *testing.T) {
	output := NewMockLogger()

	sl := log.NewLogger("debug")
	sl.SetDebugOutput(output.debug)

	sl.Debugf("debug message with value: %d", 42)
	expectedOutput := "debug message with value: 42\n"
	actualOutput := output.debug.String()
	if actualOutput != expectedOutput {
		t.Errorf("Expected debugf output:\n%s\nBut got:\n%s", expectedOutput, actualOutput)
	}
}

func TestSimpleLoggerInfo(t *testing.T) {
	output := NewMockLogger()

	sl := log.NewLogger("info")
	sl.SetInfoOutput(output.info)

	sl.Info("info message")
	expectedOutput := "info message\n"
	actualOutput := output.info.String()
	if actualOutput != expectedOutput {
		t.Errorf("Expected info output:\n%s\nBut got:\n%s", expectedOutput, actualOutput)
	}
}

func TestSimpleLoggerInfof(t *testing.T) {
	output := NewMockLogger()

	sl := log.NewLogger("info")
	sl.SetInfoOutput(output.info)

	sl.Infof("info message with value: %d", 42)
	expectedOutput := "info message with value: 42\n"
	actualOutput := output.info.String()
	if actualOutput != expectedOutput {
		t.Errorf("Expected infof output:\n%s\nBut got:\n%s", expectedOutput, actualOutput)
	}
}

func TestSimpleLoggerError(t *testing.T) {
	output := NewMockLogger()

	sl := log.NewLogger("error")
	sl.SetErrorOutput(output.error)

	sl.Error("error message")
	expectedOutput := "error message\n"
	actualOutput := output.error.String()
	if actualOutput != expectedOutput {
		t.Errorf("Expected error output:\n%s\nBut got:\n%s", expectedOutput, actualOutput)
	}
}

func TestSimpleLoggerErrorf(t *testing.T) {
	output := NewMockLogger()

	sl := log.NewLogger("error")
	sl.SetErrorOutput(output.error)

	sl.Errorf("error message with value: %d", 42)
	expectedOutput := "error message with value: 42\n"
	actualOutput := output.error.String()
	if actualOutput != expectedOutput {
		t.Errorf("Expected errorf output:\n%s\nBut got:\n%s", expectedOutput, actualOutput)
	}
}

func (m *MockLogger) Write(p []byte) (n int, err error) {
	// Not needed for testing
	return 0, nil
}

func (m *MockLogger) SetOutput(out io.Writer) {
	// Redirect the output to the respective buffers
	m.debug.Reset()
	m.info.Reset()
	m.error.Reset()

	if output, ok := out.(*bytes.Buffer); ok {
		// If the provided output is a bytes.Buffer, copy its contents to the respective buffers
		m.debug.Write(output.Bytes())
		m.info.Write(output.Bytes())
		m.error.Write(output.Bytes())
	}
}
