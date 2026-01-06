---
name: do-ai-nano-banana
description: 콘텐츠 생성, 이미지 생성 및 AI 기반 워크플로우를 위한 Nano-Banana AI 서비스 통합. AI 서비스를 콘텐츠 생성에 통합할 때 사용
version: 1.1.0
category: integration
allowed-tools: Read, Write, Edit, Bash, Grep, Glob
tags:
  - ai
  - content-generation
  - image-generation
  - nano-banana
  - ai-service
  - gemini-3-pro
related-skills:
  - do-docs-generation
  - do-domain-uiux
updated: 2025-12-23
status: active
author: Do Team
---

# Nano-Banana AI 서비스 통합

## 빠른 참조

Nano Banana Pro (gemini-3-pro-image-preview) 통합으로 고품질 이미지 생성, 1K/2K/4K 해상도 지원, 스타일 프리픽스, 배치 처리 제공

핵심 기능:
- 이미지 생성: 1K, 2K, 4K 해상도 및 10가지 종횡비 지원
- 스타일 제어: 일관된 미학을 위한 스타일 프리픽스 지원
- 배치 처리: JSON/YAML 설정을 통한 동시 생성
- 스마트 재시도: API 할당량 처리를 위한 지수 백오프
- Google 검색: 사실 기반 콘텐츠를 위한 그라운딩 생성

CLI 빠른 시작:
- 단일 이미지: generate_image.py -p "Dashboard UI design" -o ui.png -r 4K
- 스타일 포함: generate_image.py -p "Logo" -o logo.png --style "minimalist, vector"
- 배치 생성: batch_generate.py -c prompts.json -d output/ --concurrency 2

사용 시점:
- 문서화를 위한 고품질 이미지 생성
- 일관된 스타일링을 가진 시각적 자산 생성
- 다중 이미지 프롬프트 배치 처리
- 자동화된 이미지 생성 파이프라인 구축

---

## 구현 가이드

### CLI 스크립트 (권장)

단일 이미지 생성 (scripts/generate_image.py):
- 기본 사용: generate_image.py -p "A fluffy cat" -o cat.png
- 고해상도 및 종횡비: generate_image.py -p "Dashboard UI" -o ui.png -r 4K -a 16:9
- 일관된 미학을 위한 스타일 프리픽스: generate_image.py -p "Mountain landscape" -o landscape.png --style "photorealistic, dramatic lighting"
- Google 검색 그라운딩: generate_image.py -p "Mount Fuji at sunset" -o fuji.png -g
- 상세 모드: generate_image.py -p "Tech logo" -o logo.png -r 4K -a 1:1 --style "minimalist, vector art" -v

배치 생성 (scripts/batch_generate.py):
- JSON 설정에서: batch_generate.py -c prompts.json -d output/
- YAML 설정에서 동시성 포함: batch_generate.py -c prompts.yaml -d output/ --concurrency 3
- 명령줄 프롬프트에서: batch_generate.py --prompts "Cat" "Dog" "Bird" -d output/ --style "watercolor painting"
- 보고서 생성 포함: batch_generate.py -c prompts.json -d output/ --report report.json

### 배치 설정 파일 예제

prompts.json 형식:
- defaults 객체: style, resolution, aspect_ratio 기본값 정의
- images 배열: 문자열 프롬프트 또는 상세 객체 포함

상세 이미지 객체 필드:
- prompt: 이미지 설명 (필수)
- filename: 출력 파일명
- resolution: 1K, 2K, 4K
- aspect_ratio: 종횡비
- style: 스타일 프리픽스

### 모델 사양

모델명: gemini-3-pro-image-preview (Nano Banana Pro)

지원 해상도:
- 1K: 표준 품질, 빠른 생성
- 2K: 품질과 속도의 균형 (기본값)
- 4K: 최고 품질, 상세한 출력

지원 종횡비:
- 정사각형: 1:1
- 세로: 2:3, 3:4, 4:5, 9:16
- 가로: 3:2, 4:3, 5:4, 16:9, 21:9

특수 기능:
- Thinking 모드 기본 활성화
- 최대 14개 참조 이미지 혼합 지원
- 로고용 고급 텍스트 렌더링
- Google 검색 그라운딩 옵션

### 스타일 프리픽스 예제

일반적인 스타일 프리픽스:
- "photorealistic, 8K, detailed" - 사진 같은 품질
- "minimalist, vector art, clean" - 간단한 그래픽
- "watercolor painting, soft colors" - 예술적 스타일
- "3D render, cinematic lighting" - 3D 외관
- "anime style, vibrant colors" - 애니메이션 풍
- "professional photography, studio lighting" - 제품 촬영
- "digital art, fantasy, ethereal" - 창의적 일러스트레이션
- "infographic, data visualization" - 차트 및 다이어그램

### 오류 처리

스크립트에 지수 백오프를 통한 자동 재시도 포함:
- 최대 5회 재시도 (설정 가능)
- 기본 지연: 2초, 매 시도마다 2배 증가
- 최대 지연: 60-120초
- 할당량 오류 시 지터 추가로 일시적 폭주 방지

재시도 불가능한 오류 (즉시 실패):
- 잘못된 API 키
- 잘못된 인수 (잘못된 종횡비/해상도)

---

## 고급 패턴

### 프로그래밍 방식 사용

generate_image 모듈 활용:
- prompt: 이미지 설명
- output_path: 출력 파일 경로
- aspect_ratio: 종횡비
- resolution: 해상도 (1K, 2K, 4K)
- style: 스타일 프리픽스
- enable_grounding: Google 검색 그라운딩 활성화
- max_retries: 최대 재시도 횟수
- verbose: 상세 로그 출력

결과 객체:
- success: 성공 여부 (bool)
- output_path: 생성된 파일 경로
- generation_time_seconds: 생성 소요 시간

batch_generate 모듈 활용:
- ImageTask 객체: prompt, output_path, resolution, aspect_ratio, style, task_id
- run_batch_generation: tasks, concurrency, max_retries, verbose
- save_batch_report: 결과를 JSON으로 저장

---

## 설정

환경 변수:
- GOOGLE_API_KEY: Google AI API 키 (필수)
- python-dotenv 설치 시 .env 파일 지원

API 키 발급: https://aistudio.google.com/apikey

의존성:
- google-genai
- python-dotenv
- pyyaml

---

## AI 콘텐츠 생성 워크플로우

### 문서 생성 파이프라인

사양 데이터에서 완전한 문서 생성:

Phase 1 - API 참조 생성:
- 입력: API 엔드포인트, 요구사항
- 모델: claude-3-5-sonnet
- 온도: 0.3 (일관성 중시)
- 출력: 포괄적 API 문서

Phase 2 - 사용 예제 생성:
- 입력: Phase 1 API 문서
- 포함: 다중 언어 예제 (Python, TypeScript, curl)
- 온도: 0.5 (적절한 창의성)

Phase 3 - 튜토리얼 생성:
- 입력: API 문서 및 예제
- 대상: 초급~중급 개발자
- 온도: 0.7 (설명적 스타일)

Phase 4 - 다이어그램 설명 생성:
- 포함: 아키텍처, 시퀀스, 흐름 다이어그램
- 형식: Mermaid 다이어그램 구문

### 다국어 콘텐츠 생성

기본 콘텐츠를 여러 언어로 생성:

번역 요구사항:
- 기술적 정확성 유지
- 문화적 참조 적응
- 서식 유지
- 코드 예제는 영어로 유지

### 디자인 자산 생성

완전한 디자인 자산 세트 생성:

프로세스:
- 각 사양에 대해 기본 이미지 생성
- 스타일 가이드 (시각적 스타일, 색상 팔레트) 적용
- 요청 시 변형 생성 (variation_strength 조정)

결과 메타데이터:
- 생성 타임스탬프
- 적용된 스타일 가이드
- 원본 사양

### 콘텐츠 품질 평가

포괄적 콘텐츠 품질 분석:

Phase 1 - 가독성 분석:
- 분석 유형: readability
- 메트릭 포함

Phase 2 - 기술적 정확성:
- 분석 유형: technical_accuracy
- 기준: 기술 표준

Phase 3 - 완전성 점검:
- 분석 유형: completeness
- 필수 섹션 확인

Phase 4 - 스타일 일관성:
- 분석 유형: style_consistency
- 스타일 가이드 적용

최종 결과:
- 전체 품질 점수 산출
- 각 분야별 상세 분석
- 개선 권장사항

---

## 배치 처리

### 자동화된 콘텐츠 파이프라인

여러 콘텐츠 항목 병렬 처리:

설정 옵션:
- batch_size: 배치당 항목 수 (기본: 5)
- max_retries: 최대 재시도 횟수 (기본: 3)
- batch_delay: 배치 간 지연 (기본: 1.0초)

처리 흐름:
- 입력 항목을 배치로 분할
- 각 배치 내 항목 동시 처리
- 성공/실패 결과 분리
- 배치 간 속도 제한 적용

결과 객체:
- total_items: 총 항목 수
- successful: 성공 수
- failed: 실패 수
- results: 성공 결과 목록
- errors: 오류 목록 (항목, 오류 메시지, 타임스탬프)
- processing_time: 총 처리 시간

단일 항목 처리:
- 지수 백오프로 재시도 로직 적용
- 각 시도 간 대기 시간: 2의 시도 횟수 제곱 초

---

## 사용 예제

### 콘텐츠 생성

API 문서 생성:
- 프롬프트: REST API 문서 생성 요청
- 엔드포인트 목록: 인증, 사용자 프로필 조회/수정
- 포함 항목: 요청/응답 예제, 오류 코드, 인증

코드 예제 생성:
- 프롬프트: JWT 인증 코드 예제
- 언어: Python, TypeScript
- 포함: 토큰 생성, 검증, 갱신 흐름

### 이미지 생성

히어로 이미지:
- 프롬프트: Modern SaaS dashboard hero image, clean UI, gradient background
- 크기: 1920x1080
- 스타일: digital_art
- 품질: high

아이콘 세트:
- 각 아이콘 설명에 대해 반복 생성
- 크기: 256x256
- 스타일: vector
- 배경: transparent

아이콘 예제:
- User profile icon, minimal line art
- Settings gear icon, modern design
- Notification bell icon, clean style

### 텍스트 분석

감정 분석:
- 입력: 사용자 피드백
- 분석 유형: sentiment
- 점수 포함

요약:
- 입력: 긴 문서
- 분석 유형: summary
- 최대 길이: 200
- 핵심 포인트 포함

엔티티 추출:
- 입력: 기술 문서
- 분석 유형: entity_extraction
- 엔티티 유형: technology, framework, tool, concept

---

## 보완 Skill

보완 Skill:
- do-docs-generation: 이미지를 포함한 자동화된 문서 생성
- do-domain-uiux: UI/UX 디자인 자산 생성
- do-domain-frontend: 프론트엔드 시각적 컴포넌트
- do-workflow-templates: 템플릿 기반 이미지 워크플로우

통합 포인트:
- 문서 일러스트레이션 파이프라인
- 디자인 시스템 자산 생성
- 마케팅 자료 생성
- 자동화된 시각적 콘텐츠 워크플로우

---

## 리소스

Skill 파일:
- SKILL.md: 메인 문서 (이 파일)
- scripts/generate_image.py: 단일 이미지 생성 CLI
- scripts/batch_generate.py: 배치 생성 CLI
- examples.md: 추가 워크플로우 예제

외부 문서:
- Google AI Studio: https://aistudio.google.com/
- Gemini API Reference: https://ai.google.dev/gemini-api/docs

---

Version: 1.1.0
Last Updated: 2025-12-23
Status: Active
