---
name: do-formats-data
description: TOON encoding, JSON/YAML 최적화, 직렬화 패턴, 데이터 검증을 다루는 데이터 포맷 전문가
version: 1.0.0
category: library
allowed-tools: Read, Write, Edit, Grep, Glob
tags:
  - formats
  - data
  - toon
  - serialization
  - validation
updated: 2026-01-06
status: stable
author: Do Team
user-invocable: false
---

# 데이터 포맷 전문가

## 빠른 참조

고급 데이터 포맷 관리 - TOON encoding, JSON/YAML 최적화, 직렬화 패턴, 데이터 검증

핵심 기능:
- TOON Encoding: JSON 대비 40-60% 토큰 감소
- JSON/YAML 최적화: 효율적인 직렬화 및 파싱 패턴
- 데이터 검증: 스키마 검증, 타입 체크, 오류 처리
- 포맷 변환: 데이터 포맷 간 원활한 변환
- 성능: 최적화된 데이터 구조 및 캐싱 전략

사용 시점:
- LLM 토큰 예산 내에서 데이터 전송 최적화
- 고성능 직렬화/역직렬화
- 스키마 검증 및 데이터 무결성
- 대규모 데이터셋 처리

---

## 핵심 개념

### TOON (Token-Optimized Object Notation)

LLM 토큰 사용에 최적화된 커스텀 바이너리 호환 포맷:
- 타입 마커: # (숫자), ! (불린), @ (타임스탬프), ~ (null)
- JSON 대비 40-60% 크기 감소
- 무손실 왕복 인코딩/디코딩

### 성능 최적화

- orjson 기반 초고속 JSON 처리 (표준 대비 2-5배)
- ijson을 사용한 대규모 데이터셋 스트리밍 처리
- LRU 제거 및 메모리 관리를 갖춘 지능형 캐싱

### 데이터 검증

- 커스텀 규칙과 패턴을 갖춘 타입 안전 검증
- 스키마 진화 및 마이그레이션 지원
- 교차 필드 검증 및 의존성 체크

---

## 구현 가이드

### TOON Encoding

```
class TOONEncoder {
    markers = { number: '#', boolean: '!', timestamp: '@', null: '~' }

    encode(data) {
        return encodeValue(data)
    }

    encodeValue(value) {
        if value is null: return '~'
        if value is boolean: return '!' + (value ? '1' : '0')
        if value is number: return '#' + value
        if value is datetime: return '@' + value.toISOString()
        if value is array: return '[' + value.map(encodeValue).join('|') + ']'
        if value is object: return '{' + objectToPairs(value) + '}'
        return escapeString(value)
    }

    decode(toonString) {
        return parseValue(toonString.trim())
    }

    parseValue(s) {
        if s == '~': return null
        if s.startsWith('!'): return s[1] == '1'
        if s.startsWith('#'): return parseNumber(s.slice(1))
        if s.startsWith('@'): return parseISO(s.slice(1))
        if s.startsWith('['): return parseArray(s)
        if s.startsWith('{'): return parseObject(s)
        return unescapeString(s)
    }
}

// 사용 예시
encoder = new TOONEncoder()
data = {user: {id: 123, name: "John", active: true}, tags: ["read", "write"]}
toonData = encoder.encode(data)
// 결과: {user:{id:#123,name:John,active:!1},tags:[read|write]}
```

### JSON 최적화

```
class JSONOptimizer {
    serializeFast(obj) {
        // orjson 옵션: SERIALIZE_NUMPY, SERIALIZE_DATACLASS, UTC_Z
        return fastJSON.stringify(obj)
    }

    deserializeFast(jsonBytes) {
        return fastJSON.parse(jsonBytes)
    }

    streamParse(filePath, callback) {
        // ijson으로 메모리 효율적 대용량 파일 스트리밍
        for item in streamItems(filePath) {
            callback(item)
        }
    }
}

// 사용 예시
optimizer = new JSONOptimizer()
data = {users: generateUsers(10000)}
jsonBytes = optimizer.serializeFast(data)  // 표준 대비 5배
```

### 데이터 검증

```
class DataValidator {
    typeValidators = {
        string: (v) => typeof v == 'string',
        integer: (v) => Number.isInteger(v),
        email: (v) => emailPattern.test(v),
        array: (v) => Array.isArray(v),
        object: (v) => typeof v == 'object'
    }

    createSchema(rules) {
        schema = {}
        for field, config in rules {
            schema[field] = {
                type: config.type ?? 'string',
                required: config.required ?? true,
                minLength: config.min_length,
                maxLength: config.max_length,
                minValue: config.min_value,
                maxValue: config.max_value,
                pattern: config.pattern ? compileRegex(config.pattern) : null
            }
        }
        return schema
    }

    validate(data, schema) {
        errors = {}
        sanitized = {}

        for field, rule in schema {
            value = data[field]

            if rule.required and value is undefined {
                errors[field] = field + ' is required'
                continue
            }
            if value is undefined: continue

            if not typeValidators[rule.type](value) {
                errors[field] = field + ' must be ' + rule.type
                continue
            }

            // 길이/값 범위 검증
            if rule.minLength and value.length < rule.minLength {
                errors[field] = 'minimum length is ' + rule.minLength
            }
            if rule.pattern and not rule.pattern.test(value) {
                errors[field] = 'pattern mismatch'
            }

            sanitized[field] = sanitizeValue(value, rule.type)
        }

        return { valid: isEmpty(errors), errors, sanitizedData: sanitized }
    }
}

// 사용 예시
validator = new DataValidator()
userSchema = validator.createSchema({
    username: {type: 'string', required: true, min_length: 3},
    email: {type: 'email', required: true},
    age: {type: 'integer', required: false, min_value: 13}
})
result = validator.validate({username: 'john_doe', email: 'john@example.com'}, userSchema)
```

---

## 고급 기능

### 스트림 처리

```
class StreamProcessor {
    processJsonStream(filePath, processor) {
        stream = openFile(filePath, 'rb')
        parser = createJsonParser(stream)
        for item in parser.items('item') {
            processor(item)
        }
    }

    batchProcess(filePath, processor, batchSize = 1000) {
        batch = []
        for item in streamItems(filePath) {
            batch.push(processor(item))
            if batch.length >= batchSize {
                yield batch
                batch = []
            }
        }
        if batch.length > 0: yield batch
    }
}

// 사용: 대용량 파일 처리
processor = new StreamProcessor()
for batch in processor.batchProcess("data.json", processItem, 5000) {
    saveBatch(batch)
}
```

### 스키마 진화

```
class SchemaEvolution {
    schemas = {}
    migrations = {}

    registerSchema(version, schema) {
        this.schemas[version] = { schema, timestamp: now(), version }
    }

    addMigration(fromVersion, toVersion, migrationFn) {
        this.migrations[fromVersion + ':' + toVersion] = migrationFn
    }

    migrateData(data, fromVersion, toVersion) {
        key = fromVersion + ':' + toVersion
        if key not in this.migrations {
            throw 'No migration path'
        }
        return this.migrations[key](data)
    }
}

// 사용: 버전 간 데이터 마이그레이션
evolution = new SchemaEvolution()
evolution.registerSchema('v1', {name: {type: 'string'}})
evolution.registerSchema('v2', {full_name: {type: 'string'}, email: {type: 'email'}})
evolution.addMigration('v1', 'v2', (data) => ({full_name: data.name, email: null}))

newData = evolution.migrateData({name: 'John Doe'}, 'v1', 'v2')
```

### 지능형 캐싱

```
class SmartCache {
    maxMemoryMB = 50
    cache = {}
    accessTimes = {}
    sizes = {}

    generateKey(data) {
        return hash(fastJSON.stringify(data, {sortKeys: true})).slice(0, 16)
    }

    get(key) {
        if key in this.cache {
            this.accessTimes[key] = now()
            return this.cache[key]
        }
        return null
    }

    set(key, value) {
        size = fastJSON.stringify(value).length
        while totalSize() + size > this.maxMemoryMB * 1024 * 1024 {
            evictLRU()
        }
        this.cache[key] = value
        this.sizes[key] = size
        this.accessTimes[key] = now()
    }

    evictLRU() {
        oldestKey = min(this.accessTimes.keys(), by: k => this.accessTimes[k])
        delete this.cache[oldestKey]
        delete this.sizes[oldestKey]
        delete this.accessTimes[oldestKey]
    }
}
```

---

## 통합 패턴

### LLM 데이터 최적화

```
class LLMDataPreparer {
    maxTokens = 4000
    encoder = new TOONEncoder()

    prepareForLLM(data) {
        encoded = this.encoder.encode(data)
        if encoded.split(' ').length > this.maxTokens {
            reduced = filterObject(data, ['id', 'name', 'type', 'status'])
            encoded = this.encoder.encode(reduced)
        }
        return encoded
    }

    parseLLMResponse(response) {
        return this.encoder.decode(response)
    }
}
```

### API 응답 최적화

```
encoder = new TOONEncoder()
optimizer = new JSONOptimizer()

async function getUser(userId, format = 'json') {
    user = await fetchUser(userId)
    if format == 'toon' {
        return Response(encoder.encode(user), {contentType: 'application/x-toon'})
    }
    return Response(optimizer.serializeFast(user), {contentType: 'application/json'})
}
```

---

## CLI 사용

```
# TOON 포맷으로 인코딩
do-formats encode-toon --input data.json --output data.toon

# 스키마 기반 데이터 검증
do-formats validate --schema schema.json --data data.json

# 포맷 간 변환
do-formats convert --input data.json --output data.yaml --format yaml

# JSON 구조 최적화
do-formats optimize-json --input large-data.json --output optimized.json
```

---

## 문제 해결

### TOON 인코딩 이슈

"Invalid type marker" 디코드 실패:
- 입력 데이터에서 지원되지 않는 타입 확인
- 동일한 TOONEncoder 설정으로 인코딩/디코딩 확인
- try/catch로 JSON 폴백 구현

토큰 절감 효과 미달:
- 데이터 구조 점검 (중첩 객체에서 효과 극대화)
- 반복 키에 대해 키 압축 활성화

### JSON 성능 이슈

직렬화 속도 저하:
- orjson 설치 확인
- 가능하면 bytes 출력 사용

대용량 파일 메모리 부족:
- ijson으로 스트리밍 파서 사용
- StreamProcessor로 청크 처리

### 검증 이슈

커스텀 검증기 미작동:
- 스키마 생성 전 검증기 등록 확인
- 검증기 함수 시그니처 확인

스키마 진화로 기존 데이터 깨짐:
- 버전별 마이그레이션에 SchemaEvolution 사용
- 역호환 기본값 추가

---

## 모범 사례

대규모 데이터셋 처리:
- 라인 단위 처리에 NDJSON 포맷 사용
- 독립 항목에 병렬 처리 활성화
- 컴파일된 정규식 패턴 및 스키마 캐싱

메모리 관리:
- 스트리밍에 리스트 대신 제너레이터 사용
- 주기적으로 캐시 정리
- 설정에 명시적 메모리 제한 설정

---

## 연관 스킬

- do-domain-backend - 백엔드 데이터 직렬화 및 API 응답
- do-domain-database - 데이터베이스 데이터 포맷 최적화
- do-foundation-core - MCP 데이터 직렬화 및 전송 패턴
- do-docs-generation - 문서 데이터 포맷팅

---

## 기술 스택

핵심 라이브러리:
- orjson: 초고속 JSON 파싱 및 직렬화
- PyYAML: YAML 처리 (C 기반 로더)
- ijson: 대용량 파일용 스트리밍 JSON 파서
- regex: 고급 정규식 지원

검증 라이브러리:
- jsonschema, pydantic, marshmallow, cerberus

---

Status: Production Ready
Last Updated: 2026-01-06
Maintained by: Do Team
