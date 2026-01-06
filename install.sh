#!/bin/sh
# Do - Claude Code 프로젝트 환경 설치
# curl -fsSL https://raw.githubusercontent.com/yejune/do/main/install.sh | sh

set -e

echo ""
echo "Do - Claude Code 프로젝트 환경"
echo "================================"
echo ""

# 이미 설치 확인
if [ -d ".claude/agents/do" ] && [ "$1" != "--force" ]; then
    echo "이미 설치됨. 업데이트: curl ... | sh -s -- --force"
    exit 0
fi

# 임시 디렉토리
TMP=$(mktemp -d)
trap "rm -rf $TMP" EXIT

# 다운로드
echo "다운로드 중..."
if ! git clone --depth 1 --quiet https://github.com/yejune/do.git "$TMP/do"; then
    echo "오류: git clone 실패"
    exit 1
fi

# 디렉토리 생성
echo "설치 중..."
mkdir -p .claude
mkdir -p .do/config/sections

# .claude 복사 (기존 settings.json 보존)
if [ -f ".claude/settings.json" ]; then
    cp .claude/settings.json "$TMP/settings.backup.json"
fi

cp -r "$TMP/do/.claude/agents" .claude/
cp -r "$TMP/do/.claude/skills" .claude/
cp -r "$TMP/do/.claude/hooks" .claude/
cp -r "$TMP/do/.claude/commands" .claude/
cp -r "$TMP/do/.claude/styles" .claude/

if [ -f "$TMP/settings.backup.json" ]; then
    echo "settings.json 보존됨 (새 버전: .claude/settings.json.new)"
    cp "$TMP/do/.claude/settings.json" .claude/settings.json.new
    mv "$TMP/settings.backup.json" .claude/settings.json
else
    cp "$TMP/do/.claude/settings.json" .claude/
fi

# .do/config 복사 (기존 설정 보존)
if [ ! -f ".do/config/sections/language.yaml" ]; then
    cp -r "$TMP/do/.do/config/"* .do/config/ 2>/dev/null || true
fi

# CLAUDE.md
if [ ! -f "CLAUDE.md" ]; then
    cp "$TMP/do/CLAUDE.md" .
else
    echo "CLAUDE.md 보존됨 (새 버전: CLAUDE.md.new)"
    cp "$TMP/do/CLAUDE.md" CLAUDE.md.new
fi

# 권한 설정
chmod +x .claude/hooks/do/*.py 2>/dev/null || true

echo ""
echo "설치 완료!"
echo ""
echo "포함:"
echo "  - 27개 에이전트 (expert/manager/builder/mcp)"
echo "  - 47개 스킬 (lang/platform/workflow/domain)"
echo ""
echo "설정: .do/config/sections/"
echo "  - language.yaml (대화 언어)"
echo "  - user.yaml (사용자 이름)"
echo ""
