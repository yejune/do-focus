package main

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var version = "dev" // 빌드 시 -ldflags로 주입

const (
	repoURL    = "https://github.com/yejune/do-focus.git"
	releaseURL = "https://github.com/yejune/do-focus/releases/latest/download/do-release.tar.gz"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "sync":
		runSync()
	case "selfupdate", "self-update":
		runSelfUpdate()
	case "capture":
		runCapture()
	case "setup-logging":
		runSetupLogging()
	case "version", "-v", "--version":
		fmt.Printf("godo version %s\n", version)
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`godo - Do CLI installer

Usage:
  godo sync           Install or update Do
  godo selfupdate     Update godo itself
  godo capture        Capture terminal buffer to file
  godo setup-logging  Setup Claude session logging
  godo version        Show version
  godo help           Show this help

Examples:
  cd my-project
  godo sync                              # Install or update Do
  godo selfupdate                        # Update godo CLI
  godo capture --output terminal.txt     # Capture terminal buffer
  godo setup-logging                     # Setup Claude logging to ~/.do/claude-session.log`)
}

func runSelfUpdate() {
	fmt.Println("godo 업데이트 중...")
	fmt.Printf("현재 버전: %s\n", version)

	// Try brew first
	cmd := exec.Command("brew", "upgrade", "yejune/tap/godo")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		// Fallback: direct download
		fmt.Println("brew 업그레이드 실패. 직접 다운로드 시도...")
		selfUpdateDirect()
		return
	}

	fmt.Println("✓ godo 업데이트 완료!")
}

func selfUpdateDirect() {
	// Detect OS and arch
	goos := os.Getenv("GOOS")
	goarch := os.Getenv("GOARCH")

	if goos == "" {
		switch {
		case strings.Contains(strings.ToLower(os.Getenv("OS")), "windows"):
			goos = "windows"
		default:
			// Use uname
			out, _ := exec.Command("uname", "-s").Output()
			switch strings.TrimSpace(strings.ToLower(string(out))) {
			case "darwin":
				goos = "darwin"
			default:
				goos = "linux"
			}
		}
	}

	if goarch == "" {
		out, _ := exec.Command("uname", "-m").Output()
		arch := strings.TrimSpace(strings.ToLower(string(out)))
		switch arch {
		case "arm64", "aarch64":
			goarch = "arm64"
		default:
			goarch = "amd64"
		}
	}

	// Download URL
	binaryName := fmt.Sprintf("godo-%s-%s", goos, goarch)
	if goos == "windows" {
		binaryName += ".exe"
	}
	url := fmt.Sprintf("https://github.com/yejune/do/releases/latest/download/%s", binaryName)

	fmt.Printf("다운로드: %s\n", url)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("오류: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("오류: HTTP %d\n", resp.StatusCode)
		os.Exit(1)
	}

	// Get current executable path
	exePath, err := os.Executable()
	if err != nil {
		fmt.Printf("오류: %v\n", err)
		os.Exit(1)
	}

	// Write to temp file
	tmpFile := exePath + ".new"
	f, err := os.OpenFile(tmpFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Printf("오류: %v\n", err)
		os.Exit(1)
	}

	io.Copy(f, resp.Body)
	f.Close()

	// Replace old binary
	oldFile := exePath + ".old"
	os.Remove(oldFile)
	os.Rename(exePath, oldFile)

	if err := os.Rename(tmpFile, exePath); err != nil {
		// Rollback
		os.Rename(oldFile, exePath)
		fmt.Printf("오류: %v\n", err)
		os.Exit(1)
	}

	os.Remove(oldFile)
	fmt.Println("✓ godo 업데이트 완료!")
}

func isDevFolder() bool {
	// Check if this is the Do development folder
	if fileExists("tobrew.yaml") && fileExists("cmd/godo/main.go") {
		return true
	}
	return false
}

func runInit() {
	// Prevent running in development folder
	if isDevFolder() {
		fmt.Println("오류: Do 개발 폴더에서는 godo init/update를 실행할 수 없습니다.")
		fmt.Println("      다른 프로젝트 폴더에서 실행하세요.")
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("Do - Claude Code 프로젝트 환경")
	fmt.Println("================================")
	fmt.Println()

	// Check if already installed
	if isInstalled() {
		if !confirm("이미 설치됨. 업데이트할까요?") {
			fmt.Println("취소됨.")
			return
		}
		runUpdate()
		return
	}

	install(false)
}

func runUpdate() {
	// Prevent running in development folder
	if isDevFolder() {
		fmt.Println("오류: Do 개발 폴더에서는 godo init/update를 실행할 수 없습니다.")
		fmt.Println("      다른 프로젝트 폴더에서 실행하세요.")
		os.Exit(1)
	}

	if !isInstalled() {
		fmt.Println("Do가 설치되지 않음. 'godo init'을 먼저 실행하세요.")
		os.Exit(1)
	}

	fmt.Println("업데이트 중...")
	install(true)
}

func runSync() {
	// Prevent running in development folder
	if isDevFolder() {
		fmt.Println("오류: Do 개발 폴더에서는 godo init/update/sync를 실행할 수 없습니다.")
		fmt.Println("      다른 프로젝트 폴더에서 실행하세요.")
		os.Exit(1)
	}

	if isInstalled() {
		// Already installed - run update
		fmt.Println("업데이트 중...")
		install(true)
		fmt.Println()
		fmt.Println("✓ Do 업데이트 완료!")

		// Auto-setup logging alias
		aliasStatus := checkAliasStatus()

		switch aliasStatus {
		case "none":
			fmt.Println()
			fmt.Println("⚙️  Claude 로깅 설정 중...")
			runSetupLogging()
			printSourceInstructions()

		case "outdated":
			fmt.Println()
			fmt.Println("⚙️  Claude 로깅 설정 업데이트 중...")
			updateAliasInRcFile()
			fmt.Println("✓ Claude 로깅 설정 업데이트 완료!")
			printSourceInstructions()

		case "ok":
			// Already correct, do nothing
		}

		fmt.Println("💡 Claude Code를 시작하세요")
	} else {
		// New installation - run init
		fmt.Println()
		fmt.Println("Do - Claude Code 프로젝트 환경")
		fmt.Println("================================")
		fmt.Println()
		install(false)
	}
}

func isInstalled() bool {
	_, err := os.Stat(".claude/agents/do")
	return err == nil
}

func confirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s [y/N]: ", prompt)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

func install(force bool) {
	fmt.Println("다운로드 중...")

	// Try release download first, fallback to git clone
	if err := installFromRelease(force); err != nil {
		fmt.Printf("릴리즈 다운로드 실패, git clone 시도: %v\n", err)
		installFromGit(force)
	}

	printSuccess()
}

func installFromRelease(force bool) error {
	// Download tarball
	resp, err := http.Get(releaseURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	// Create directories
	os.MkdirAll(".claude", 0755)
	os.MkdirAll(".do/config/sections", 0755)

	// Extract tarball
	gzr, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Skip files that shouldn't be overwritten
		target := header.Name
		if !force {
			if target == "CLAUDE.md" && fileExists("CLAUDE.md") {
				continue
			}
			if target == ".claude/settings.json" && fileExists(".claude/settings.json") {
				continue
			}
			if strings.HasPrefix(target, ".do/config/") && fileExists(".do/config/sections/language.yaml") {
				continue
			}
		}

		switch header.Typeflag {
		case tar.TypeDir:
			os.MkdirAll(target, 0755)
		case tar.TypeReg:
			os.MkdirAll(filepath.Dir(target), 0755)
			f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			io.Copy(f, tr)
			f.Close()
		}
	}

	// Copy user.yaml.example to user.yaml if user.yaml doesn't exist
	userYamlExample := ".do/config/sections/user.yaml.example"
	userYaml := ".do/config/sections/user.yaml"
	if fileExists(userYamlExample) && !fileExists(userYaml) {
		copyFile(userYamlExample, userYaml)
	}

	return nil
}

func installFromGit(force bool) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "do-install-*")
	if err != nil {
		fmt.Printf("오류: 임시 디렉토리 생성 실패: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmpDir)

	// Clone repository
	cloneDir := filepath.Join(tmpDir, "do")
	cmd := exec.Command("git", "clone", "--depth", "1", "--quiet", repoURL, cloneDir)
	if err := cmd.Run(); err != nil {
		fmt.Printf("오류: git clone 실패: %v\n", err)
		os.Exit(1)
	}

	// Create directories
	fmt.Println("설치 중...")
	os.MkdirAll(".claude", 0755)
	os.MkdirAll(".do/config/sections", 0755)

	// Copy .claude directories
	dirs := []string{"agents", "skills", "hooks", "commands", "styles", "lib"}
	for _, dir := range dirs {
		src := filepath.Join(cloneDir, ".claude", dir)
		dst := filepath.Join(".claude", dir)
		if _, err := os.Stat(src); err == nil {
			os.RemoveAll(dst)
			copyDir(src, dst)
		}
	}

	// settings.json
	settingsSrc := filepath.Join(cloneDir, ".claude", "settings.json")
	settingsDst := filepath.Join(".claude", "settings.json")
	if force || !fileExists(settingsDst) {
		copyFile(settingsSrc, settingsDst)
	}

	// .do/config
	configSrc := filepath.Join(cloneDir, ".do", "config")
	if force || !fileExists(".do/config/sections/language.yaml") {
		if entries, err := os.ReadDir(configSrc); err == nil {
			for _, entry := range entries {
				src := filepath.Join(configSrc, entry.Name())
				dst := filepath.Join(".do/config", entry.Name())
				if entry.IsDir() {
					copyDir(src, dst)
				} else {
					copyFile(src, dst)
				}
			}
		}
	}

	// Copy user.yaml.example to user.yaml if user.yaml doesn't exist
	userYamlExample := ".do/config/sections/user.yaml.example"
	userYaml := ".do/config/sections/user.yaml"
	if fileExists(userYamlExample) && !fileExists(userYaml) {
		copyFile(userYamlExample, userYaml)
	}

	// CLAUDE.md
	claudeMdSrc := filepath.Join(cloneDir, "CLAUDE.md")
	if force || !fileExists("CLAUDE.md") {
		copyFile(claudeMdSrc, "CLAUDE.md")
	}

	// Set permissions
	filepath.Walk(".claude/hooks", func(path string, info os.FileInfo, err error) error {
		if err == nil && strings.HasSuffix(path, ".py") {
			os.Chmod(path, 0755)
		}
		return nil
	})
}

func printSuccess() {
	fmt.Println()
	fmt.Println("설치 완료!")
	fmt.Println()
	fmt.Println("포함:")
	fmt.Println("  - 27개 에이전트 (expert/manager/builder/mcp)")
	fmt.Println("  - 47개 스킬 (lang/platform/workflow/domain)")
	fmt.Println()
	fmt.Println("설정: /do:setup 실행하여 개인 설정")
	fmt.Println("  - .claude/settings.local.json에 저장됨")
	fmt.Println()
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(src, path)
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, 0755)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(dstPath, data, 0644)
	})
}

func checkAliasInstalled() bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	shell := os.Getenv("SHELL")
	var rcFile string

	if strings.Contains(shell, "zsh") {
		rcFile = filepath.Join(home, ".zshrc")
	} else if strings.Contains(shell, "bash") {
		rcFile = filepath.Join(home, ".bashrc")
	} else {
		return false
	}

	content, err := os.ReadFile(rcFile)
	if err != nil {
		return false
	}

	// Check for "Do - Claude logging" comment
	return strings.Contains(string(content), "Do - Claude logging")
}

// checkAliasStatus returns "none", "outdated", or "ok"
func checkAliasStatus() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "none"
	}

	shell := os.Getenv("SHELL")
	var rcFile string

	if strings.Contains(shell, "zsh") {
		rcFile = filepath.Join(home, ".zshrc")
	} else if strings.Contains(shell, "bash") {
		rcFile = filepath.Join(home, ".bashrc")
	} else {
		return "none"
	}

	content, err := os.ReadFile(rcFile)
	if err != nil {
		return "none"
	}

	contentStr := string(content)

	// Not installed
	if !strings.Contains(contentStr, "Do - Claude logging") {
		return "none"
	}

	// Check if outdated (compare with expected alias content)
	expectedAlias := getExpectedAliasContent()
	if !strings.Contains(contentStr, expectedAlias) {
		return "outdated"
	}

	return "ok"
}

func getExpectedAliasContent() string {
	return `# Do - Claude logging (added by godo)
if command -v claude &> /dev/null; then
    claude_original="$(which claude)"
    claude() {
        # Find Git root or use current directory
        local git_root=$(git rev-parse --show-toplevel 2>/dev/null || echo "$PWD")

        # Create hierarchical directory structure: <project>/.do/claude-session/YYYY/MM/DD/
        local session_date=$(date +%Y/%m/%d)
        local session_dir="${git_root}/.do/claude-session/${session_date}"
        mkdir -p "$session_dir"

        # Session ID format: YYYYMMDD-HHmmss (dash instead of underscore)
        export CLAUDE_SESSION_ID=$(date +%Y%m%d-%H%M%S)
        local log_file=${session_dir}/${CLAUDE_SESSION_ID}.session

        # Show session ID to user (stderr - not sent to Claude)
        echo "🔗 Session: $CLAUDE_SESSION_ID" >&2

        # Log detailed info to file only
        echo "🎬 Claude session started at $(date)" >> "$log_file"

        # Run Claude with logging (use process substitution to preserve stdin)
        "$claude_original" "$@" > >(tee -a "$log_file") 2>&1
        local exit_code=$?

        # Log end to file only
        echo "🏁 Claude session ended at $(date) (exit code: $exit_code)" >> "$log_file"
        return $exit_code
    }
fi`
}

func printSourceInstructions() {
	shell := os.Getenv("SHELL")
	var rcFile string

	if strings.Contains(shell, "zsh") {
		rcFile = "~/.zshrc"
	} else if strings.Contains(shell, "bash") {
		rcFile = "~/.bashrc"
	} else {
		rcFile = "~/.bashrc or ~/.zshrc"
	}

	fmt.Println()
	fmt.Println("📌 다음 명령을 실행하여 현재 shell에 적용하세요:")
	fmt.Printf("   source %s\n", rcFile)
	fmt.Println()
	fmt.Println("또는 새 터미널 창을 여세요")
	fmt.Println()
}

func updateAliasInRcFile() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	shell := os.Getenv("SHELL")
	var rcFile string

	if strings.Contains(shell, "zsh") {
		rcFile = filepath.Join(home, ".zshrc")
	} else if strings.Contains(shell, "bash") {
		rcFile = filepath.Join(home, ".bashrc")
	} else {
		fmt.Println("지원하지 않는 shell:", shell)
		os.Exit(1)
	}

	content, err := os.ReadFile(rcFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Remove old alias section
	lines := strings.Split(string(content), "\n")
	var newLines []string
	skipMode := false

	for _, line := range lines {
		if strings.Contains(line, "# Do - Claude logging") {
			skipMode = true
			continue
		}
		if skipMode && strings.TrimSpace(line) == "fi" {
			skipMode = false
			continue
		}
		if !skipMode {
			newLines = append(newLines, line)
		}
	}

	// Write new content
	newContent := strings.Join(newLines, "\n")
	if err := os.WriteFile(rcFile, []byte(newContent), 0644); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Append new alias
	f, err := os.OpenFile(rcFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	aliasContent := "\n" + getExpectedAliasContent() + "\n"
	if _, err := f.WriteString(aliasContent); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
