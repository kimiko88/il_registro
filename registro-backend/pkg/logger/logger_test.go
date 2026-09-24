package logger

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestInit_LogLevels(t *testing.T) {
	tests := []struct {
		name          string
		inputLevel    string
		expectedLevel logrus.Level
	}{
		{"debug level", "debug", logrus.DebugLevel},
		{"info level", "info", logrus.InfoLevel},
		{"warn level", "warn", logrus.WarnLevel},
		{"warning level", "warning", logrus.WarnLevel},
		{"error level", "error", logrus.ErrorLevel},
		{"invalid level defaults to info", "nonexistent_level", logrus.InfoLevel},
		{"empty level defaults to info", "", logrus.InfoLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Init(tt.inputLevel)
			assert.Equal(t, tt.expectedLevel, Log.GetLevel())
		})
	}
}

func TestInit_JSONFormatting(t *testing.T) {
	Init("info")
	var buf bytes.Buffer
	Log.SetOutput(&buf)

	Log.WithField("school_id", "school-test-123").Info("test log message")

	output := buf.String()
	assert.NotEmpty(t, output)

	var parsed map[string]interface{}
	err := json.Unmarshal([]byte(output), &parsed)
	assert.NoError(t, err, "log output must be valid JSON")
	assert.Equal(t, "test log message", parsed["msg"])
	assert.Equal(t, "info", parsed["level"])
	assert.Equal(t, "school-test-123", parsed["school_id"])
}
