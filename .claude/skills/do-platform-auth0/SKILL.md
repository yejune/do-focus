---
name: do-platform-auth0
description: Auth0 보안 전문가로서 공격 방어, 다중 인증, 토큰 보안, 발신자 제약(DPoP/mTLS), 규정 준수(FAPI, GDPR, HIPAA)를 다룸. Auth0 보안 기능 구현, MFA 설정, 토큰 보안, 규정 요구사항 충족 시 사용.
version: 1.0.0
category: security
tags: [auth0, security, mfa, attack-protection, tokens, dpop, mtls, compliance, fapi, gdpr]
updated: 2026-01-06
status: active
allowed-tools: Read, Write, Edit, Grep, Glob, WebFetch, WebSearch, Bash
context7-libraries: /auth0/docs
---

# Auth0 보안 전문가

Auth0 구현을 위한 종합 보안 기술. 공격 방어, 다중 인증(MFA), 토큰 보안, 발신자 제약(DPoP/mTLS), 규정 준수(FAPI, GDPR, HIPAA)를 포함.

---

## 빠른 참조

### 보안 기능 카테고리

공격 방어:
- Bot Detection: 의심스러운 트래픽에 CAPTCHA 챌린지
- Breached Password Detection: 손상된 자격증명 차단
- Brute Force Protection: 계정별 로그인 실패 시도 제한
- Suspicious IP Throttling: 고속 공격 속도 제한

다중 인증:
- Auth0 Guardian 푸시 알림
- 일회용 비밀번호(OTP/TOTP)
- 보안 키 및 생체인식 WebAuthn
- SMS/음성 인증, Adaptive MFA

토큰 보안:
- JWT 구조 및 검증
- 스코프 기반 Access Token 관리
- Refresh Token 순환 및 만료
- 토큰 취소 전략

발신자 제약:
- DPoP: 애플리케이션 계층 토큰 바인딩
- mTLS: 전송 계층 인증서 바인딩

규정 준수: FAPI, GDPR, HIPAA/HITECH, PCI DSS, ISO 27001, SOC 2

### Dashboard 탐색

공격 방어: Dashboard > Security > Attack Protection
MFA 설정: Dashboard > Security > Multi-factor Auth
보안 센터: Dashboard > Security > Security Center

### 필수 설정 체크리스트

1. 적절한 민감도로 Bot Detection 활성화
2. Breached Password Detection 활성화
3. Brute Force Protection 임계값 설정
4. Suspicious IP Throttling 활성화
5. 최소 하나의 MFA 요소 설정
6. 토큰 만료 정책 구성

---

## 공격 방어 구현

### Bot Detection

위험 신호: 트래픽 품질 패턴 기반 IP 평판 분석.

메커니즘: 봇 활동이 의심되는 IP에서 로그인 시도 시 인증 챌린지 트리거. 로그인, 가입, 비밀번호 재설정 트래픽 패턴을 분석하는 통계 모델 사용.

설정 옵션:
- 민감도 수준: Low, Medium(기본), High
- 응답 유형: Auth Challenge(CAPTCHA 없는 JS 검증), Simple CAPTCHA, 서드파티 연동(reCAPTCHA)
- IP 허용목록: 최대 100개 주소 또는 CIDR 범위

지원 플로우: Universal Login, Classic Login, Lock.js v12.4.0+, Auth0.swift 1.28.0+, Auth0.Android 1.25.0+ 네이티브 앱

미지원 플로우: Enterprise 연결, 소셜 로그인, 교차 출처 인증

### Breached Password Detection

위험 신호: 다크웹 데이터베이스 및 서드파티 유출 데이터에서 발견된 손상된 비밀번호.

탐지 방법:

Standard Detection:
- 공개 유출 데이터 추적
- 탐지 시간: 유출 공개 후 7-13개월
- B2B/B2C Professional 또는 Enterprise 플랜

Credential Guard(Enterprise 추가 기능):
- 전용 보안팀을 통한 비공개 유출 데이터 접근
- 탐지 시간: 12-36시간
- 200개 이상 국가 커버

대응 시나리오:
- 신규 가입 시 손상된 자격증명 차단
- 로그인 시 손상된 사용자 계정 차단
- 비밀번호 재설정 시 손상된 자격증명 차단

알림 옵션: 사용자 알림, 관리자 경고(즉시, 일간, 주간, 월간)

테스트: AUTH0-TEST-로 시작하는 비밀번호로 실제 경고 없이 탐지 확인 가능.

### Brute Force Protection

위험 신호: 특정 계정을 대상으로 하는 로그인 시도 속도.

메커니즘: 단일 IP 주소에서 정의된 기간 내 반복된 로그인 실패 시도 식별. 트리거 시 해당 사용자로의 의심스러운 IP 로그인 차단.

설정:
- Brute Force 임계값: 기본 10회 실패(1-100 설정 가능)
- IP 허용목록: 신뢰할 수 있는 IP 또는 CIDR 범위 제외
- 대응 옵션: Brute Force 로그인 차단(IP 기반), 계정 잠금(모든 IP), 사용자 알림

차단 해제 조건:
- 마지막 실패 후 30일 경과
- 연결된 모든 계정에서 비밀번호 변경
- 관리자가 차단 제거 또는 임계값 상향
- 알림 이메일의 차단 해제 링크 선택

### Suspicious IP Throttling

위험 신호: 여러 계정에 걸친 단일 IP의 로그인 시도 속도.

메커니즘: 고속 로그인 또는 가입 시도를 보이는 IP 주소의 트래픽 자동 차단. HTTP 429(Too Many Requests) 상태 코드로 응답.

속도 탐지 방식:

로그인 시도:
- IP 주소별 일일 실패 로그인 시도 추적
- 임계값 초과 시 후속 시도 조절
- 24시간에 걸쳐 균등 분배
- 예: 100 비율은 약 15분마다 한 번 시도 허용

가입 시도:
- 1분 내 모든 시도(성공/실패) 집계
- IP가 한도 초과 시 추가 가입 차단
- 예: 72,000 비율은 초당 약 한 번 시도 허용

중요 사항: 잘못된 요청과 스키마 검증 오류는 임계값에 포함되지 않음. 새 테넌트에서 기본 활성화.

### 모니터링 모드

대응 설정 없이 기능을 활성화하여 모니터링 모드 활성화. 활성 차단 메커니즘 배포 전 분석 및 의사결정을 위해 테넌트 로그에 이벤트 기록.

### 설정 모범 사례

초기 배포:
1. 먼저 모니터링 모드로 기능 활성화
2. 오탐 패턴에 대한 테넌트 로그 분석
3. 신뢰 가능한 소스에 IP 허용목록 설정
4. 점진적으로 대응 조치 활성화
5. 적절한 알림 빈도 설정

권장 시작 설정:
- Bot Detection: Medium 민감도, Auth Challenge
- Breached Password Detection: 가입/로그인 시 차단, 사용자 알림
- Brute Force Protection: 10회 시도, IP 차단, 사용자 알림
- Suspicious IP Throttling: 기본 임계값, 관리자 알림

---

## 다중 인증 구현

### 지원 MFA 요소

독립 요소(최소 하나 필수):
- WebAuthn 보안 키: 물리적 보안 키(YubiKey 등), FIDO2/U2F 표준, 피싱 방지
- 일회용 비밀번호(OTP/TOTP): 시간 기반 일회용 비밀번호, Google Authenticator/Authy 호환
- 푸시 알림(Auth0 Guardian): 원탭 승인/거부, 컨텍스트가 포함된 풍부한 알림
- 전화 메시지: SMS 인증 코드, 음성 통화 인증, 스마트폰 미사용자 대체
- Cisco Duo Security: 기업 Duo 통합, 통합 보안 플랫폼

종속 요소(독립 요소 먼저 필요):
- 장치 생체인식 WebAuthn: Face ID, Touch ID, Windows Hello
- 이메일 인증: 이메일을 통한 일회용 코드
- 복구 코드: 미리 생성된 백업 코드, 일회용

### MFA 정책

Never: MFA 불필요, 선택적 등록, 최저 보안/최고 편의

Always: 모든 로그인에 MFA 필수, 모든 사용자 완료 필요, 최고 보안

Adaptive MFA(Enterprise): 위험 기반 MFA 챌린지, 위험 감지 시만 챌린지, 보안과 편의의 균형

### Adaptive MFA 위험 신호

트랜잭션별 위험 신호 평가:
- NewDevice: 지난 30일간 미사용 기기
- ImpossibleTravel: 지리적 이상
- UntrustedIP: 의심스러운 활동 이력

고위험 트랜잭션은 기존 MFA 세션과 관계없이 검증 필요.

### WebAuthn 구현

인증 유형:
- 보안 키: USB, NFC, Bluetooth 연결 물리적 기기
- 장치 생체인식: Face ID, Touch ID, Windows Hello, Android 지문/얼굴

사용자 검증 설정:
- Discouraged: 사용자 검증 없음
- Preferred: 가능 시 사용자 검증(기본)
- Required: 사용자 검증 필수(고보안 권장)

장점:
- 피싱 방지: 자격증명이 특정 도메인에 바인딩
- 공유 비밀 없음: 개인 키가 기기를 떠나지 않음
- 단일 단계 다중 인증: 소유(기기)와 사용자 검증 결합

### Step-Up 인증

민감한 작업에 대한 향상된 검증. API는 스코프 사용, 웹 앱은 ID 토큰 클레임 확인.

---

## 토큰 보안

### JWT 기초

구조(RFC 7519):
- Header: 토큰 유형(JWT), 서명 알고리즘(alg)
- Payload: 클레임(iss, sub, aud, exp, iat, jti)
- Signature: 무결성 및 진위 보장

주요 알고리즘:
- HS256(대칭): 단일 공유 비밀, 간단하지만 보호 필요
- RS256(비대칭): 공개/개인 키 쌍, 권장, 분산 시스템에 적합
- ES256(비대칭): 타원 곡선, 현대 애플리케이션용

검증 필수 항목:
- 서명 검증
- exp: 토큰 미만료
- iss: 발급자 일치
- aud: 대상 애플리케이션 포함

보안 모범 사례:
- 서명되지 않은 토큰 신뢰 금지
- 민감한 데이터를 페이로드에 저장 금지(base64 인코딩만, 암호화 아님)
- 항상 HTTPS 사용
- "none" 알고리즘 수락 금지
- 검증 시 예상 알고리즘 명시

### Access Token

API 접근 권한을 스코프와 함께 부여.

유형:
- Opaque: 내성(introspection) 필요
- JWT: 자체 포함

주요 클레임: iss, sub, aud, scope, exp

기본 수명: 86400초(24시간)

### Refresh Token

세션 연속성 활성화. 사용자당 애플리케이션당 최대 200개 활성.

보안 기능:
- 순환: 이전 토큰 무효화
- 만료: 유휴/절대 만료
- 취소: Management API를 통한 취소

### 토큰 모범 사례

- 서명 키를 중요 자격증명으로 취급
- 공개 키 검증을 위해 HS256보다 RS256 선호
- 가능하면 서버 측에 토큰 저장
- 만료까지 캐시 및 재사용

JWKS 엔드포인트: https://{your-tenant}.auth0.com/.well-known/jwks.json

---

## 발신자 제약

### DPoP(애플리케이션 계층)

목적: 토큰을 클라이언트 생성 비대칭 키 쌍에 바인딩하여 토큰 도난 및 오용 방지.

현재 상태: Auth0에서 Early Access.

작동 방식:

1. 키 쌍 생성:
   - 개인 키: 비밀 유지, 전송 금지
   - 공개 키: DPoP Proof JWT에 포함
   - 권장 알고리즘: ES256

2. DPoP Proof JWT:
   - Header: typ(dpop+jwt), alg, jwk(공개 키)
   - Payload: jti(고유 식별자), htm(HTTP 메서드), htu(HTTP URI), iat, ath(API 호출용 Access Token 해시)

3. 토큰 요청:
   - 토큰 엔드포인트용 새 DPoP Proof 생성
   - DPoP 헤더로 전송
   - 공개 클라이언트는 use_dpop_nonce 오류 처리 필요

4. 토큰 바인딩:
   - Access Token에 cnf(확인) 클레임 포함
   - jkt(JWK Thumbprint) 값: 공개 키의 SHA-256 해시
   - 키 보유자만 토큰 사용 가능

5. API 요청:
   - Authorization: DPoP {access_token}
   - DPoP: {dpop_proof_jwt}
   - 각 요청마다 새로운 Proof 생성

DPoP 장점:
- 애플리케이션 계층(TLS 변경 불필요)
- 공개 클라이언트와 작동
- PKI 인프라 불필요
- 더 쉬운 배포

### mTLS(전송 계층)

목적: X.509 인증서에 토큰 바인딩.

프로세스:
1. 클라이언트가 mTLS 연결 수립
2. Auth0가 인증서 SHA-256 지문 계산
3. 토큰 cnf 클레임에 x5t#S256로 삽입
4. 리소스 서버가 지문 검증

요구사항:
- 기밀 클라이언트만
- Enterprise Plan + HRI 추가 기능
- PKI 인프라

mTLS 장점:
- 전송 계층 바인딩
- 기존 인프라 활용
- 기밀 클라이언트에 더 간단

선택 기준:
- 공개 클라이언트(SPA, 모바일): DPoP
- PKI 없음, 유연성 필요: DPoP
- 기밀 클라이언트, PKI 존재: mTLS
- 전송 계층 바인딩 선호: mTLS

---

## 규정 준수

### Highly Regulated Identity(HRI)

요구사항: Enterprise Plan + HRI 추가 기능

핵심 기능:

Strong Customer Authentication(SCA):
- 다른 카테고리의 최소 두 독립 요소 필요
- 지식(비밀번호, PIN), 소유(기기, 토큰), 고유성(생체인식)

Dynamic Linking:
- 인가에 트랜잭션 세부정보 바인딩
- 사용자가 정확히 무엇을 인가하는지 인지
- 다른 트랜잭션에 재사용 불가

PAR(Pushed Authorization Requests):
- 인가 매개변수를 Auth0에 직접 전송
- 참조 URI 수신
- 브라우저에서 매개변수 노출 방지

JAR(JWT-Secured Authorization Requests):
- 인가 요청을 JWT로 서명
- 요청 무결성 보호
- 선택적 암호화로 기밀성 확보

JWE(JSON Web Encryption):
- Access Token 페이로드 암호화
- 민감한 데이터 기밀성 보호

클라이언트 인증 방식:
- Private Key JWT: 비대칭 키, 개인 키 미전송, 최대 2개 공개 키 등록
- mTLS for OAuth: X.509 인증서, 가장 강력한 보호

### GDPR 준수

적용 대상:
- EU 거주자 개인 데이터 처리 기업
- EU에 설립된 기업
- EU 개인을 모니터링하는 비EU 기업
- EU에 서비스를 제공하는 비EU 기업

역할 분리:
- Auth0: 데이터 처리자(고객 지시에 따라 처리)
- 고객: 데이터 컨트롤러(처리 목적 결정, 동의 관리)

사용자 권리 구현:

접근권: Management API로 사용자 데이터 내보내기, 기계 판독 가능 형식

수정권: Dashboard 또는 Management API로 업데이트

삭제권(잊혀질 권리): Dashboard 또는 API 삭제, 프로필 및 메타데이터 제거

이동권: JSON 형식 내보내기, 표준 구조

동의 관리:
- 동의 체크박스 구현
- 사용자 메타데이터에 동의 저장
- 동의 타임스탬프 및 버전 기록
- 철회 메커니즘 제공

### 인증 및 표준

- ISO 27001/27017/27018
- SOC 2 Type 2
- CSA STAR
- FAPI 1 Advanced OP
- HIPAA BAA 가능
- PCI DSS 준수 모델

---

## 고급 패턴

### 보안 센터 모니터링

위치: Dashboard > Security > Security Center

위협 카테고리:
- 자격증명 스터핑: 기계 기반 침해 시도
- 가입 공격: 자동화된 계정 생성
- MFA 우회: 우회 시도

필터링: 기간(최대 14일), 애플리케이션, 연결. 분/시간/일별 자동 집계.

메트릭: Bot Detection 횟수, IP 조절 이벤트, Brute Force 트리거, 유출 비밀번호 알림, MFA 성공/실패율.

### 애플리케이션 자격증명

Client Secret(기본):
- 대칭, 간단하지만 가로채기 취약

Private Key JWT(Enterprise):
- 비대칭 키 쌍
- 개인 키 미전송
- 단기 어설션
- 향상된 보안에 권장

mTLS for OAuth(HRI):
- X.509 인증서
- 가장 강력한 보호

키 관리: 무중단 순환을 위해 최대 2개 공개 키 등록. 알고리즘: RS256, RS384, PS256.

### 지속적 세션 보호

Auth0 Actions를 사용하여 토큰 갱신 이벤트 중 세션 컨텍스트 활용.

기능: IP/ASN 모니터링, 기기 추적, 만료 관리, 이상 탐지.

동적 관리: 사용자 속성, 조직, 역할별 수명 커스터마이징.

---

## 사용 가이드

이 스킬은 다음 상황에서 사용:
- 공격 방어 설정 및 구성
- 다중 인증 방식 선택 및 구현
- 토큰 보안 전략 수립
- 발신자 제약(DPoP/mTLS) 구현
- FAPI, GDPR, HIPAA 규정 준수 검증

종합 보안 검토가 필요한 경우 expert-security 에이전트 사용 권장.

---

## 리소스

공식 문서:
- https://auth0.com/docs/secure
- https://auth0.com/docs/secure/attack-protection
- https://auth0.com/docs/secure/multi-factor-authentication
- https://auth0.com/docs/secure/tokens
- https://auth0.com/docs/secure/sender-constraining
- https://auth0.com/docs/secure/data-privacy-and-compliance
