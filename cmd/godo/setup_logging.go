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
	aliasLine := `
# Do - Claude logging (added by godo)
if command -v claude &> /dev/null; then
    claude_original="$(which claude)"
    claude() {
        mkdir -p ~/.do
        local log_file=~/.do/claude-session-$(date +%Y%m%d_%H%M%S).log
        echo "🎬 Claude session started at $(date)" | tee -a "$log_file"
        "$claude_original" "$@" 2>&1 | tee -a "$log_file"
        local exit_code=${PIPESTATUS[0]}
        echo "🏁 Claude session ended at $(date) (exit code: $exit_code)" | tee -a "$log_file"
        return $exit_code
    }
fi
`

	// 4. 이미 설정되어 있는지 확인
	content, err := os.ReadFile(rcFile)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		// 파일이 없으면 새로 생성
		content = []byte("")
	}

	if strings.Contains(string(content), "Do - Claude logging") {
		fmt.Println("✓ 이미 설정되어 있습니다")
		return
	}

	// 5. rc 파일에 추가
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
	fmt.Println("이제 claude 실행 시 자동으로 ~/.do/claude-session.log에 기록됩니다")
}
