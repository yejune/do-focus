---
description: 사용자 이름과 언어 설정
allowed-tools: AskUserQuestion, Write, Read
---

# /do:setup - 사용자 설정

## 실행 단계

### Step 1: 현재 설정 로드

`.do/config/config.yaml` 파일을 읽어 현재 설정 확인 (user 섹션).

### Step 2: AskUserQuestion으로 수집

**Question 1 - 이름:**
- question: "어떤 이름으로 불러드릴까요?"
- header: "Name"
- options:
  - label: "Keep current" (현재 값이 있으면)
  - 사용자가 직접 입력

**Question 2 - 언어:**
- question: "어떤 언어로 대화할까요?"
- header: "Language"
- options:
  - label: "한국어" → ko
  - label: "English" → en
  - label: "日本語" → ja
  - label: "中文" → zh

### Step 3: 설정 저장

`.do/config/config.yaml`의 user 섹션과 language.conversation 업데이트:

```yaml
user:
  name: "{입력받은 이름}"

language:
  conversation: "{선택한 언어 코드}"
  code_comments: en
  documentation: "{선택한 언어 코드}"
```

기존 config.yaml의 다른 섹션(do, style 등)은 유지.

### Step 4: 완료 메시지

```
설정 완료!
- 이름: {name}
- 언어: {language}
```

## 중요

- AskUserQuestion의 options는 최대 4개
- 이모지 금지
- 간결하게
