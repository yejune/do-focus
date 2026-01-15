package main

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
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
	case "worker":
		runWorker()
	case "selfupdate", "self-update":
		runSelfUpdate()
	case "capture":
		runCapture()
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
  godo sync             Install or update Do
  godo worker start     Start the memory worker
  godo worker stop      Stop the memory worker
  godo worker status    Show worker status
  godo selfupdate       Update godo itself
  godo capture          Capture terminal buffer to file
  godo version          Show version
  godo help             Show this help

Examples:
  cd my-project
  godo sync                              # Install or update Do
  godo worker start                      # Start memory worker
  godo selfupdate                        # Update godo CLI
  godo capture --output terminal.txt     # Capture terminal buffer`)
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

	// 1. 글로벌 디렉토리 초기화
	if err := initializeGlobalDir(); err != nil {
		fmt.Printf("Warning: Failed to initialize global dir: %v\n", err)
	}

	// 2. 현재 프로젝트 등록
	cwd, _ := os.Getwd()
	if err := registerProject(cwd); err != nil {
		fmt.Printf("Warning: Failed to register project: %v\n", err)
	}

	// 3. 기존 sync 로직 (hooks 복사 등)
	if isInstalled() {
		// Already installed - run update
		fmt.Println("업데이트 중...")
		install(true)
		fmt.Println()
		fmt.Println("✓ Do 업데이트 완료!")
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

		// Skip files that shouldn't be overwritten or extracted
		target := header.Name

		// Always skip install.sh (should not be in project folders)
		if target == "install.sh" {
			continue
		}

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

	// Copy .claude directories (always overwrite)
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

	// CLAUDE.md (always overwrite)
	claudeMdSrc := filepath.Join(cloneDir, "CLAUDE.md")
	copyFile(claudeMdSrc, "CLAUDE.md")

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

// initializeGlobalDir은 ~/.do/ 디렉토리를 초기화합니다
func initializeGlobalDir() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	baseDir := filepath.Join(homeDir, ".do")

	// 디렉토리 생성
	dirs := []string{
		baseDir,
		filepath.Join(baseDir, "bin"),
		filepath.Join(baseDir, "config"),
		filepath.Join(baseDir, "logs"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// config.json 초기화 (없으면)
	configPath := filepath.Join(baseDir, "config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaultConfig := `{
  "version": "1.0.0",
  "worker": {
    "port": 3778,
    "auto_start": true
  },
  "viewer": {
    "port": 3777
  },
  "database": {
    "type": "sqlite",
    "path": "~/.do/memory.db"
  }
}`
		if err := os.WriteFile(configPath, []byte(defaultConfig), 0644); err != nil {
			return err
		}
	}

	// projects.json 초기화 (없으면)
	projectsPath := filepath.Join(baseDir, "projects.json")
	if _, err := os.Stat(projectsPath); os.IsNotExist(err) {
		if err := os.WriteFile(projectsPath, []byte(`{"projects":[]}`), 0644); err != nil {
			return err
		}
	}

	return nil
}

// registerProject는 현재 프로젝트를 projects.json에 등록합니다
func registerProject(projectPath string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	projectsPath := filepath.Join(homeDir, ".do", "projects.json")

	// 기존 데이터 읽기
	data, err := os.ReadFile(projectsPath)
	if err != nil {
		return err
	}

	var projects struct {
		Projects []struct {
			Path         string `json:"path"`
			Name         string `json:"name"`
			RegisteredAt int64  `json:"registered_at"`
		} `json:"projects"`
	}

	if err := json.Unmarshal(data, &projects); err != nil {
		return err
	}

	// 중복 체크
	for _, p := range projects.Projects {
		if p.Path == projectPath {
			return nil // 이미 등록됨
		}
	}

	// 새 프로젝트 추가
	projects.Projects = append(projects.Projects, struct {
		Path         string `json:"path"`
		Name         string `json:"name"`
		RegisteredAt int64  `json:"registered_at"`
	}{
		Path:         projectPath,
		Name:         filepath.Base(projectPath),
		RegisteredAt: time.Now().Unix(),
	})

	// 저장
	newData, err := json.MarshalIndent(projects, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(projectsPath, newData, 0644)
}

// runWorker handles worker subcommands: start, stop, status
func runWorker() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: godo worker [start|stop|status]")
		os.Exit(1)
	}

	switch os.Args[2] {
	case "start", "restart":
		workerStart()
	case "stop":
		workerStop()
	case "status":
		workerStatus()
	default:
		fmt.Printf("Unknown worker command: %s\n", os.Args[2])
		fmt.Println("Usage: godo worker [start|stop|status]")
		os.Exit(1)
	}
}

func getWorkerPath() string {
	// First check if godo-worker is in PATH (installed via homebrew)
	if path, err := exec.LookPath("godo-worker"); err == nil {
		return path
	}
	// Fallback to ~/.do/bin/godo-worker
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".do", "bin", "godo-worker")
}

func isWorkerRunning() bool {
	resp, err := http.Get("http://127.0.0.1:3778/health")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

func getWorkerPID() int {
	// Find worker process by port
	cmd := exec.Command("lsof", "-ti", ":3778")
	out, err := cmd.Output()
	if err != nil {
		return 0
	}
	pid := strings.TrimSpace(string(out))
	if pid == "" {
		return 0
	}
	var p int
	fmt.Sscanf(pid, "%d", &p)
	return p
}

func workerStart() {
	workerPath := getWorkerPath()

	// Check if worker binary exists
	if _, err := os.Stat(workerPath); os.IsNotExist(err) {
		// Check if it's in PATH (homebrew install)
		if _, err := exec.LookPath("godo-worker"); err != nil {
			fmt.Println("Error: godo-worker not found")
			fmt.Println("       Install with: brew upgrade godo")
			os.Exit(1)
		}
	}

	// Kill existing if running (idempotent restart)
	if pid := getWorkerPID(); pid > 0 {
		// Try SIGTERM first
		exec.Command("kill", fmt.Sprintf("%d", pid)).Run()
		time.Sleep(500 * time.Millisecond)

		// If still running, SIGKILL
		if getWorkerPID() > 0 {
			exec.Command("kill", "-9", fmt.Sprintf("%d", pid)).Run()
			time.Sleep(300 * time.Millisecond)
		}
	}

	// Start worker
	cmd := exec.Command(workerPath)
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error: Failed to start worker: %v\n", err)
		os.Exit(1)
	}

	// Wait for startup and verify
	for i := 0; i < 10; i++ {
		time.Sleep(200 * time.Millisecond)
		if isWorkerRunning() {
			fmt.Printf("✓ Worker started (PID: %d)\n", cmd.Process.Pid)
			fmt.Println("  http://127.0.0.1:3778")
			return
		}
	}

	fmt.Println("Warning: Worker may not have started correctly")
}

func workerStop() {
	pid := getWorkerPID()
	if pid == 0 {
		fmt.Println("Worker is not running")
		return
	}

	if err := exec.Command("kill", fmt.Sprintf("%d", pid)).Run(); err != nil {
		fmt.Printf("Error: Failed to stop worker: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Worker stopped (PID: %d)\n", pid)
}

func workerStatus() {
	if isWorkerRunning() {
		pid := getWorkerPID()
		fmt.Printf("✓ Worker is running (PID: %d)\n", pid)
		fmt.Println("  http://127.0.0.1:3778")

		// Get health details
		resp, err := http.Get("http://127.0.0.1:3778/health")
		if err == nil {
			defer resp.Body.Close()
			var health map[string]interface{}
			if json.NewDecoder(resp.Body).Decode(&health) == nil {
				if v, ok := health["version"]; ok {
					fmt.Printf("  Version: %v\n", v)
				}
				if v, ok := health["db_type"]; ok {
					fmt.Printf("  DB: %v\n", v)
				}
			}
		}
	} else {
		fmt.Println("✗ Worker is not running")
		fmt.Println("  Start with: godo worker start")
	}
}

