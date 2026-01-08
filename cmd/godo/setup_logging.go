package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func runSetupLogging() {
	fmt.Println("Claude 로깅 설정 중...")

	// 1. 홈 디렉토리 찾기
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// 2. Shell 감지
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

	fmt.Printf("Shell: %s\n", shell)
	fmt.Printf("RC 파일: %s\n", rcFile)

	// 3. alias 설정
	aliasLine := "\n" + getExpectedAliasContent() + "\n"

	// 4. 기존 내용 읽기
	content, err := os.ReadFile(rcFile)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		// 파일이 없으면 새로 생성
		content = []byte("")
	}

	// 5. 기존 alias 제거 (있으면)
	if strings.Contains(string(content), "Do - Claude logging") {
		lines := strings.Split(string(content), "\n")
		newLines := []string{}
		skipUntilFi := false

		for _, line := range lines {
			if strings.Contains(line, "# Do - Claude logging") {
				skipUntilFi = true
				continue
			}
			if skipUntilFi {
				if strings.TrimSpace(line) == "fi" {
					skipUntilFi = false
				}
				continue
			}
			newLines = append(newLines, line)
		}

		content = []byte(strings.Join(newLines, "\n"))
		err := os.WriteFile(rcFile, content, 0644)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}

	// 6. rc 파일에 새 alias 추가
	f, err := os.OpenFile(rcFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	if _, err := f.WriteString(aliasLine); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Claude 로깅 설정 완료: %s\n", rcFile)
	fmt.Println()
	fmt.Println("다음 명령 실행:")
	fmt.Printf("  source %s\n", rcFile)
	fmt.Println()
	fmt.Println("이제 claude 실행 시 프로젝트 디렉토리의 .do/claude-session/YYYY/MM/DD/에 기록됩니다")
}
