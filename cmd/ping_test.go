package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestPingCommand(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		count    int
		wantPass bool
	}{
		{
			name:     "Valid localhost ping",
			host:     "127.0.0.1",
			count:    1,
			wantPass: true,
		},
		{
			name:     "Invalid host ping",
			host:     "invalid.host.name",
			count:    1,
			wantPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a pipe to capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Run the command
			err := runPingCommand(tt.host, tt.count)

			// Close writer and restore stdout
			w.Close()
			os.Stdout = oldStdout

			// Read captured output
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			// Verify results
			if tt.wantPass {
				if err != nil {
					t.Errorf("expected success, got error: %v", err)
				}
				if !strings.Contains(output, "ping statistics") {
					t.Errorf("expected ping statistics in output, got: %s", output)
				}
			} else {
				if err == nil {
					t.Error("expected error, got success")
				}
			}
		})
	}
}
