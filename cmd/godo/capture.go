package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Terminal type detection
type TerminalType int

const (
	Unknown TerminalType = iota
	ITerm2
	TerminalApp
	VSCode
	Cursor
	Antigravity
	Tmux
	Screen
)

func detectTerminal() TerminalType {
	// Check tmux first
	if os.Getenv("TMUX") != "" {
		return Tmux
	}

	// Check screen
	if os.Getenv("STY") != "" {
		return Screen
	}

	// Check TERM_PROGRAM env var
	termProgram := os.Getenv("TERM_PROGRAM")
	switch termProgram {
	case "iTerm.app":
		return ITerm2
	case "Apple_Terminal":
		return TerminalApp
	case "vscode":
		return VSCode
	case "cursor":
		return Cursor
	case "antigravity":
		return Antigravity
	}

	return Unknown
}

func captureITerm2(lines int) (string, error) {
	// Use AppleScript to get iTerm2 entire session contents
	// Note: "contents" returns all visible and scrollback buffer text
	script := `tell application "iTerm2"
		tell current session of current window
			set sessionContents to contents
			return sessionContents
		end tell
	end tell`

	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("iTerm2 capture failed: %v", err)
	}

	return string(output), nil
}

func captureTerminalApp(lines int) (string, error) {
	// Use AppleScript to get Terminal.app contents
	script := `tell application "Terminal"
		set sessionContents to contents of selected tab of front window
		return sessionContents
	end tell`

	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Terminal.app capture failed: %v", err)
	}

	return string(output), nil
}

func captureTmux(lines int) (string, error) {
	// Capture entire scrollback buffer (from beginning to end)
	// -S - means start from the beginning of history
	// -E - means end at the last line
	cmd := exec.Command("tmux", "capture-pane", "-p", "-S", "-", "-E", "-")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("tmux capture failed: %v", err)
	}

	return string(output), nil
}

func captureScreen(lines int) (string, error) {
	// Capture entire scrollback buffer including history
	tmpFile := "/tmp/screen-capture-" + strconv.Itoa(os.Getpid()) + ".txt"
	defer os.Remove(tmpFile)

	// -h flag includes scrollback history
	cmd := exec.Command("screen", "-X", "hardcopy", "-h", tmpFile)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("screen capture failed: %v", err)
	}

	// Read the file
	data, err := os.ReadFile(tmpFile)
	if err != nil {
		return "", fmt.Errorf("failed to read screen capture: %v", err)
	}

	return string(data), nil
}

func isTmuxAvailable() bool {
	_, err := exec.LookPath("tmux")
	return err == nil
}

func inTmuxSession() bool {
	// Check TMUX environment variable
	if os.Getenv("TMUX") != "" {
		return true
	}

	// Fallback: check parent process for tmux
	// This can be implemented by checking ps, but TMUX env var is more reliable
	return false
}

func captureIDE(lines int) (string, error) {
	// 1. Check if tmux is available
	if !isTmuxAvailable() {
		return "", fmt.Errorf(`IDE 터미널은 tmux 필요. 설치: brew install tmux

해결책:
1. tmux 설치 후 다시 시도
   brew install tmux

2. tmux 실행 후 claude 시작:
   tmux
   claude

3. 또는 ~/.do/claude-session.log 확인
   (claude 실행 시 자동 로깅)`)
	}

	// 2. Check if currently in tmux session
	if !inTmuxSession() {
		return "", fmt.Errorf(`IDE 터미널에서는 tmux 필요합니다.

해결책:
1. tmux 실행 후 claude 시작:
   tmux
   claude

2. 또는 ~/.do/claude-session.log 확인
   (claude 실행 시 자동 로깅)

현재 상태: tmux 세션이 감지되지 않습니다.`)
	}

	// 3. Capture from tmux session
	return captureTmux(0)
}

func limitLines(content string, lines int) string {
	allLines := strings.Split(content, "\n")
	if len(allLines) <= lines {
		return content
	}

	// Return last N lines
	start := len(allLines) - lines
	return strings.Join(allLines[start:], "\n")
}

func runCapture() {
	// Parse flags
	var outputPath string
	var lines int = 0 // 0 = capture entire scrollback buffer

	args := os.Args[2:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--output", "-o":
			if i+1 < len(args) {
				outputPath = args[i+1]
				i++
			} else {
				fmt.Println("Error: --output requires a path")
				os.Exit(1)
			}
		case "--lines", "-n":
			// Note: --lines is ignored; always captures entire scrollback
			if i+1 < len(args) {
				i++ // skip the value
			} else {
				fmt.Println("Error: --lines requires a number")
				os.Exit(1)
			}
		default:
			fmt.Printf("Error: unknown flag: %s\n", args[i])
			printCaptureUsage()
			os.Exit(1)
		}
	}

	if outputPath == "" {
		fmt.Println("Error: --output is required")
		printCaptureUsage()
		os.Exit(1)
	}

	// Detect terminal type
	termType := detectTerminal()
	var content string
	var err error

	// Capture based on terminal type
	switch termType {
	case ITerm2:
		fmt.Println("Detected: iTerm2")
		content, err = captureITerm2(lines)
	case TerminalApp:
		fmt.Println("Detected: Terminal.app")
		content, err = captureTerminalApp(lines)
	case Tmux:
		fmt.Println("Detected: tmux")
		content, err = captureTmux(lines)
	case Screen:
		fmt.Println("Detected: screen")
		content, err = captureScreen(lines)
	case VSCode:
		fmt.Println("Detected: IDE Terminal (VSCode)")
		content, err = captureIDE(lines)
	case Cursor:
		fmt.Println("Detected: IDE Terminal (Cursor)")
		content, err = captureIDE(lines)
	case Antigravity:
		fmt.Println("Detected: IDE Terminal (Antigravity)")
		content, err = captureIDE(lines)
	default:
		fmt.Println("Error: Unable to detect terminal type")
		fmt.Println("Supported: iTerm2, Terminal.app, tmux, screen, VSCode, Cursor, Antigravity")
		fmt.Println("Current TERM_PROGRAM:", os.Getenv("TERM_PROGRAM"))
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Write entire buffer to file
	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		fmt.Printf("Error: failed to write file: %v\n", err)
		os.Exit(1)
	}

	lineCount := len(strings.Split(strings.TrimSpace(content), "\n"))
	fmt.Printf("✓ 터미널 캡처 완료: %s (%d줄)\n", outputPath, lineCount)
}

func printCaptureUsage() {
	fmt.Println(`Usage: godo capture --output <file> [--lines <number>]

Options:
  --output, -o <file>    Output file path (required)
  --lines, -n <number>   Number of lines to capture (default: 500)

Examples:
  godo capture --output terminal.txt
  godo capture --output debug.txt --lines 1000
  godo capture -o output.txt -n 200

Supported terminals:
  - iTerm2
  - Terminal.app (macOS)
  - tmux
  - screen
  - VSCode (requires tmux)
  - Cursor (requires tmux)
  - Antigravity (requires tmux)

Note: IDE terminals (VSCode, Cursor, Antigravity) require tmux to be installed and running.`)
}
