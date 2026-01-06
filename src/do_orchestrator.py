#!/usr/bin/env python3
from __future__ import annotations

"""Do Orchestrator - 병렬 에이전트 실행 엔진

Agent SDK를 사용하여 여러 전문 에이전트를 병렬로 실행하고 결과를 조율함.
"""

import asyncio
import json
from dataclasses import dataclass
from enum import Enum
from pathlib import Path
from typing import Any


class AgentType(Enum):
    """사용 가능한 전문 에이전트"""
    BACKEND = "expert-backend"
    FRONTEND = "expert-frontend"
    SECURITY = "expert-security"
    TESTING = "expert-testing"
    DEBUG = "expert-debug"
    PERFORMANCE = "expert-performance"
    DATABASE = "expert-database"
    QUALITY = "manager-quality"
    GIT = "manager-git"
    DOCS = "manager-docs"


@dataclass
class TaskResult:
    """에이전트 작업 결과"""
    agent: str
    success: bool
    result: str
    error: str | None = None


class DoOrchestrator:
    """Do 총괄 에이전트 - 병렬 실행 엔진"""

    def __init__(self, project_dir: str | Path | None = None):
        self.project_dir = Path(project_dir) if project_dir else Path.cwd()
        self.results: list[TaskResult] = []

    async def run_agent(
        self,
        agent_type: AgentType,
        prompt: str,
        tools: list[str] | None = None,
    ) -> TaskResult:
        """단일 에이전트 실행

        Agent SDK가 설치되어 있으면 사용, 아니면 subprocess로 claude 호출
        """
        agent_name = agent_type.value
        allowed_tools = tools or ["Read", "Grep", "Glob", "Bash"]

        try:
            # Agent SDK 사용 시도
            result = await self._run_with_sdk(agent_name, prompt, allowed_tools)
            return TaskResult(agent=agent_name, success=True, result=result)
        except ImportError:
            # SDK 없으면 CLI 사용
            result = await self._run_with_cli(agent_name, prompt)
            return TaskResult(agent=agent_name, success=True, result=result)
        except Exception as e:
            return TaskResult(
                agent=agent_name,
                success=False,
                result="",
                error=str(e)
            )

    async def _run_with_sdk(
        self,
        agent_name: str,
        prompt: str,
        allowed_tools: list[str],
    ) -> str:
        """Agent SDK를 사용한 실행"""
        # claude_agent_sdk import 시도
        from claude_agent_sdk import ClaudeSDKClient, ClaudeAgentOptions

        options = ClaudeAgentOptions(
            allowed_tools=allowed_tools,
            permission_mode="acceptEdits",
        )

        async with ClaudeSDKClient(options=options) as client:
            await client.query(f"[{agent_name}] {prompt}")

            result_text = ""
            async for message in client.receive_response():
                if hasattr(message, 'content'):
                    for block in message.content:
                        if hasattr(block, 'text'):
                            result_text = block.text

            return result_text

    async def _run_with_cli(self, agent_name: str, prompt: str) -> str:
        """Claude CLI를 사용한 실행 (fallback)"""
        import subprocess

        cmd = [
            "claude",
            "--print",
            f"Use {agent_name} agent: {prompt}"
        ]

        proc = await asyncio.create_subprocess_exec(
            *cmd,
            stdout=asyncio.subprocess.PIPE,
            stderr=asyncio.subprocess.PIPE,
            cwd=str(self.project_dir),
        )

        stdout, stderr = await proc.communicate()

        if proc.returncode != 0:
            raise RuntimeError(f"CLI error: {stderr.decode()}")

        return stdout.decode()

    async def run_parallel(
        self,
        tasks: list[tuple[AgentType, str]],
    ) -> list[TaskResult]:
        """여러 에이전트를 병렬로 실행

        Args:
            tasks: (에이전트 타입, 프롬프트) 튜플 리스트

        Returns:
            각 에이전트의 실행 결과 리스트
        """
        coroutines = [
            self.run_agent(agent_type, prompt)
            for agent_type, prompt in tasks
        ]

        self.results = await asyncio.gather(*coroutines, return_exceptions=True)

        # 예외를 TaskResult로 변환
        processed_results = []
        for i, result in enumerate(self.results):
            if isinstance(result, Exception):
                agent_name = tasks[i][0].value
                processed_results.append(TaskResult(
                    agent=agent_name,
                    success=False,
                    result="",
                    error=str(result)
                ))
            else:
                processed_results.append(result)

        return processed_results

    async def orchestrate(self, request: str) -> dict[str, Any]:
        """요청을 분석하고 적절한 에이전트들에게 병렬로 위임

        Args:
            request: 사용자 요청

        Returns:
            종합된 결과
        """
        # 요청 분석하여 필요한 에이전트 결정
        agents_needed = self._analyze_request(request)

        if not agents_needed:
            return {"error": "적절한 에이전트를 찾을 수 없음"}

        # 병렬 실행
        tasks = [(agent, request) for agent in agents_needed]
        results = await self.run_parallel(tasks)

        # 결과 종합
        return self._synthesize_results(results)

    def _analyze_request(self, request: str) -> list[AgentType]:
        """요청을 분석하여 필요한 에이전트 목록 반환"""
        request_lower = request.lower()
        agents = []

        # 키워드 기반 에이전트 매칭
        keywords = {
            AgentType.BACKEND: ["api", "백엔드", "서버", "endpoint", "rest", "graphql"],
            AgentType.FRONTEND: ["ui", "프론트", "컴포넌트", "react", "화면", "페이지"],
            AgentType.SECURITY: ["보안", "security", "취약점", "인증", "권한"],
            AgentType.TESTING: ["테스트", "test", "tdd", "검증"],
            AgentType.DEBUG: ["버그", "오류", "에러", "debug", "fix"],
            AgentType.PERFORMANCE: ["성능", "최적화", "performance", "속도"],
            AgentType.DATABASE: ["db", "데이터베이스", "스키마", "쿼리", "sql"],
            AgentType.QUALITY: ["품질", "리뷰", "코드 검토", "quality"],
            AgentType.GIT: ["git", "커밋", "브랜치", "pr", "merge"],
            AgentType.DOCS: ["문서", "doc", "readme", "주석"],
        }

        for agent_type, words in keywords.items():
            if any(word in request_lower for word in words):
                agents.append(agent_type)

        # 기본값: 백엔드 + 테스팅
        if not agents:
            agents = [AgentType.BACKEND, AgentType.TESTING]

        return agents

    def _synthesize_results(self, results: list[TaskResult]) -> dict[str, Any]:
        """여러 에이전트의 결과를 종합"""
        successful = [r for r in results if r.success]
        failed = [r for r in results if not r.success]

        return {
            "success": len(failed) == 0,
            "total_agents": len(results),
            "successful": len(successful),
            "failed": len(failed),
            "results": {
                r.agent: r.result for r in successful
            },
            "errors": {
                r.agent: r.error for r in failed
            } if failed else None,
        }


# CLI 진입점
async def main():
    import sys

    if len(sys.argv) < 2:
        print("Usage: do_orchestrator.py <request>")
        print("Example: do_orchestrator.py '새로운 사용자 인증 API 구현'")
        sys.exit(1)

    request = " ".join(sys.argv[1:])

    orchestrator = DoOrchestrator()
    result = await orchestrator.orchestrate(request)

    print(json.dumps(result, ensure_ascii=False, indent=2))


if __name__ == "__main__":
    asyncio.run(main())
