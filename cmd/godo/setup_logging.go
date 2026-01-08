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

        # Run Claude with logging
        "$claude_original" "$@" 2>&1 | tee -a "$log_file"
        local exit_code=${PIPESTATUS[0]}

        # Log end to file only
        echo "🏁 Claude session ended at $(date) (exit code: $exit_code)" >> "$log_file"
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
	fmt.Println("이제 claude 실행 시 프로젝트 디렉토리의 .do/claude-session/YYYY/MM/DD/에 기록됩니다")
}
