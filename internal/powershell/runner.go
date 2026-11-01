package powershell

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const DefaultTimeout = 30 * time.Second

type Runner struct {
	Executable string
	Timeout    time.Duration
}

type ProcessError struct {
	Err    error
	Stderr string
}

func (e *ProcessError) Error() string {
	if e.Stderr == "" {
		return e.Err.Error()
	}
	return fmt.Sprintf("%v: %s", e.Err, strings.TrimSpace(e.Stderr))
}

func (e *ProcessError) Unwrap() error {
	return e.Err
}

func NewRunner() Runner {
	return Runner{Executable: "powershell.exe", Timeout: DefaultTimeout}
}

func (r Runner) Run(ctx context.Context, script []byte, args []string) ([]byte, error) {
	if r.Executable == "" {
		r.Executable = "powershell.exe"
	}
	if r.Timeout <= 0 {
		r.Timeout = DefaultTimeout
	}

	tempDir, err := os.MkdirTemp("", "cortex-powershell-")
	if err != nil {
		return nil, fmt.Errorf("create temporary script directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	scriptPath := filepath.Join(tempDir, "action.ps1")
	if err := os.WriteFile(scriptPath, script, 0o600); err != nil {
		return nil, fmt.Errorf("write temporary script: %w", err)
	}

	commandContext, cancel := context.WithTimeout(ctx, r.Timeout)
	defer cancel()
	commandArgs := []string{"-NoLogo", "-NoProfile", "-NonInteractive", "-ExecutionPolicy", "RemoteSigned", "-File", scriptPath}
	commandArgs = append(commandArgs, args...)
	command := exec.CommandContext(commandContext, r.Executable, commandArgs...)
	var stdout, stderr strings.Builder
	command.Stdout = &stdout
	command.Stderr = &stderr
	if err := command.Run(); err != nil {
		if commandContext.Err() != nil {
			return nil, &ProcessError{Err: commandContext.Err(), Stderr: stderr.String()}
		}
		return nil, &ProcessError{Err: err, Stderr: stderr.String()}
	}
	if stdout.Len() == 0 {
		return nil, &ProcessError{Err: io.ErrUnexpectedEOF, Stderr: stderr.String()}
	}
	return []byte(stdout.String()), nil
}
