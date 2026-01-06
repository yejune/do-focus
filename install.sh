#!/bin/bash
# Do - Claude Code Project Template Installer
#
# 사용법:
#   curl -fsSL https://raw.githubusercontent.com/USER/do/main/install.sh | bash
#
# 이 스크립트는 현재 프로젝트에 Do 설정을 설치합니다.
# 프로젝트 레벨 설정은 사용자 레벨(~/.claude/)보다 우선합니다.

set -e

# GitHub 저장소 설정 (배포 시 수정)
REPO="USER/do"
BRANCH="main"
BASE_URL="https://raw.githubusercontent.com/${REPO}/${BRANCH}"

# 색상
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo ""
echo -e "${GREEN}Do${NC} - 말하면 한다"
echo "================================"
echo ""

# 기존 설정 확인
if [ -f "CLAUDE.md" ]; then
  echo -e "${YELLOW}경고: 기존 CLAUDE.md가 있습니다.${NC}"
  read -p "덮어쓸까요? (y/N) " -n 1 -r
  echo
  if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "설치 취소됨"
    exit 1
  fi
fi

# 디렉토리 생성
echo "디렉토리 생성..."
mkdir -p .claude/agents/do
mkdir -p .claude/commands/do
mkdir -p .claude/hooks/do
mkdir -p .claude/styles

# CLAUDE.md 다운로드
echo "CLAUDE.md 다운로드..."
curl -fsSL "${BASE_URL}/CLAUDE.md" -o CLAUDE.md

# settings.json 다운로드
echo "settings.json 다운로드..."
curl -fsSL "${BASE_URL}/.claude/settings.json" -o .claude/settings.json

# 스타일 다운로드
echo "스타일 다운로드..."
for style in sprint pair direct; do
  curl -fsSL "${BASE_URL}/.claude/styles/${style}.md" -o ".claude/styles/${style}.md" 2>/dev/null || true
done

# 커맨드 다운로드
echo "커맨드 다운로드..."
curl -fsSL "${BASE_URL}/.claude/commands/do/style.md" -o ".claude/commands/do/style.md" 2>/dev/null || true

# 에이전트 다운로드 (핵심만)
echo "에이전트 다운로드..."
CORE_AGENTS=(
  "expert-backend" "expert-frontend" "expert-database"
  "manager-tdd" "manager-git" "manager-docs"
)

for agent in "${CORE_AGENTS[@]}"; do
  curl -fsSL "${BASE_URL}/.claude/agents/do/${agent}.md" -o ".claude/agents/do/${agent}.md" 2>/dev/null || true
done

echo ""
echo -e "${GREEN}설치 완료!${NC}"
echo ""
echo "================================"
echo ""
echo "사용법:"
echo "  \"로그인 기능 만들어줘\""
echo "  \"테스트 작성해줘\""
echo "  \"버그 고쳐줘\""
echo ""
echo "스타일 변경:"
echo "  /do style"
echo ""
echo "================================"
