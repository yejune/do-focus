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
	}

	return Unknown
}

func captureITerm2(lines int) (string, error) {
	// Use AppleScript to get iTerm2 contents
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
	// Use tmux capture-pane command
	cmd := exec.Command("tmux", "capture-pane", "-p", "-S", fmt.Sprintf("-%d", lines))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("tmux capture failed: %v", err)
	}

	return string(output), nil
}

func captureScreen(lines int) (string, error) {
	// Use screen hardcopy command
	tmpFile := "/tmp/screen-capture-" + strconv.Itoa(os.Getpid()) + ".txt"
	defer os.Remove(tmpFile)

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
	var lines int = 500 // default

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
			if i+1 < len(args) {
				var err error
				lines, err = strconv.Atoi(args[i+1])
				if err != nil {
					fmt.Printf("Error: invalid lines value: %v\n", err)
					os.Exit(1)
				}
				i++
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
		fmt.Println("Error: VSCode terminal capture not supported")
		fmt.Println("Hint: Use tmux or screen for terminal capture in VSCode")
		os.Exit(1)
	default:
		fmt.Println("Error: Unable to detect terminal type")
		fmt.Println("Supported: iTerm2, Terminal.app, tmux, screen")
		fmt.Println("Current TERM_PROGRAM:", os.Getenv("TERM_PROGRAM"))
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Limit to requested number of lines
	content = limitLines(content, lines)

	// Write to file
	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		fmt.Printf("Error: failed to write file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ 터미널 %d줄 캡처 완료: %s\n", lines, outputPath)
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
  - screen`)
}
