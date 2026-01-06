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
	repoURL    = "https://github.com/yejune/do.git"
	releaseURL = "https://github.com/yejune/do/releases/latest/download/do-release.tar.gz"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		runInit()
	case "update":
		runUpdate()
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
  godo init      Install Do in current directory
  godo update    Update existing Do installation
  godo version   Show version
  godo help      Show this help

Examples:
  cd my-project
  godo init      # Install Do
  godo update    # Update to latest version`)
}

func runInit() {
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
	if !isInstalled() {
		fmt.Println("Do가 설치되지 않음. 'godo init'을 먼저 실행하세요.")
		os.Exit(1)
	}

	fmt.Println("업데이트 중...")
	install(true)
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
	fmt.Println("설정: .do/config/sections/")
	fmt.Println("  - language.yaml (대화 언어)")
	fmt.Println("  - user.yaml (사용자 이름)")
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
