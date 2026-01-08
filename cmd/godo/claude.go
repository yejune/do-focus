package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func runClaude() {
	// Find Git root or use current directory
	gitRoot, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	var baseDir string
	if err != nil {
		baseDir, _ = os.Getwd()
	} else {
		baseDir = string(gitRoot[:len(gitRoot)-1]) // remove newline
	}

	// Create session directory
	now := time.Now()
	sessionID := now.Format("20060102-150405")
	sessionDir := filepath.Join(baseDir, ".do", "claude-session", now.Format("2006/01/02"))
	os.MkdirAll(sessionDir, 0755)

	logFile := filepath.Join(sessionDir, sessionID+".session")
	f, err := os.Create(logFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not create log file: %v\n", err)
		// Continue without logging
		runClaudeNoLog()
		return
	}
	defer f.Close()

	// Write session start
	fmt.Fprintf(f, "🎬 Claude session started at %s\n", now.Format(time.RFC3339))
	fmt.Fprintf(os.Stderr, "🔗 Session: %s\n", sessionID)

	// Find claude binary
	claudePath, err := exec.LookPath("claude")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: claude not found in PATH\n")
		os.Exit(1)
	}

	// Create command with all args
	args := os.Args[2:] // skip "godo" and "claude"
	cmd := exec.Command(claudePath, args...)

	// Set environment variables including session ID
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("CLAUDE_SESSION_ID=%s", sessionID))

	// Connect stdin directly
	cmd.Stdin = os.Stdin

	// Use MultiWriter to write to both stdout and log file
	cmd.Stdout = io.MultiWriter(os.Stdout, f)
	cmd.Stderr = io.MultiWriter(os.Stderr, f)

	// Run claude
	err = cmd.Run()

	// Write session end
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}

	fmt.Fprintf(f, "🏁 Claude session ended at %s (exit code: %d)\n", time.Now().Format(time.RFC3339), exitCode)

	os.Exit(exitCode)
}

func runClaudeNoLog() {
	claudePath, err := exec.LookPath("claude")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: claude not found in PATH\n")
		os.Exit(1)
	}

	// Create session ID even when no logging
	sessionID := time.Now().Format("20060102-150405")

	args := os.Args[2:]
	cmd := exec.Command(claudePath, args...)

	// Set environment variables including session ID
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("CLAUDE_SESSION_ID=%s", sessionID))

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		os.Exit(1)
	}
}
