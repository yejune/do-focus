---
name: do-lang-javascript
description: JavaScript ES2024+ 개발 전문가 - Node.js 22 LTS, 현대 런타임 (Deno, Bun), 테스팅 (Vitest, Jest), 린팅 (ESLint 9, Biome), 백엔드 프레임워크 (Express, Fastify, Hono)를 다룬다. JavaScript API, 웹 애플리케이션, Node.js 프로젝트 개발 시 사용.
version: 1.0.0
updated: 2026-01-06
status: active
allowed-tools: Read, Grep, Glob, Bash, mcp__context7__resolve-library-id, mcp__context7__query-docs
user-invocable: false
---

## 빠른 참고 (30초)

JavaScript ES2024+ 개발 전문가 - Node.js 22 LTS, 다양한 런타임, 현대적 도구 체인 통합.

자동 트리거: `.js`, `.mjs`, `.cjs` 파일, `package.json`, Node.js 프로젝트, JavaScript 관련 논의

핵심 스택:
- ES2024+: Set 메서드, Promise.withResolvers, 불변 배열, import attributes
- Node.js 22 LTS: 네이티브 TypeScript, 내장 WebSocket, 안정화된 watch 모드
- 런타임: Node.js 20/22 LTS, Deno 2.x, Bun 1.x
- 테스팅: Vitest, Jest, Node.js test runner
- 린팅: ESLint 9 flat config, Biome
- 번들러: Vite, esbuild, Rollup
- 프레임워크: Express, Fastify, Hono, Koa

빠른 명령어:
```bash
# Vite 프로젝트 생성
npm create vite@latest my-app -- --template vanilla

# 현대적 도구로 초기화
npm init -y && npm install -D vitest eslint @eslint/js

# Node.js watch 모드 실행
node --watch server.js

# Node.js 22+에서 TypeScript 직접 실행
node --experimental-strip-types app.ts
```

---

## 구현 가이드 (5분)

### ES2024 핵심 기능

Set 연산:
```javascript
const setA = new Set([1, 2, 3, 4]);
const setB = new Set([3, 4, 5, 6]);

setA.intersection(setB);      // Set {3, 4}
setA.union(setB);             // Set {1, 2, 3, 4, 5, 6}
setA.difference(setB);        // Set {1, 2}
setA.symmetricDifference(setB); // Set {1, 2, 5, 6}
setA.isSubsetOf(setB);        // false
setA.isSupersetOf(setB);      // false
setA.isDisjointFrom(setB);    // false
```

Promise.withResolvers():
```javascript
function createDeferred() {
  const { promise, resolve, reject } = Promise.withResolvers();
  return { promise, resolve, reject };
}

const deferred = createDeferred();
setTimeout(() => deferred.resolve('done'), 1000);
const result = await deferred.promise;
```

불변 배열 메서드:
```javascript
const original = [3, 1, 4, 1, 5];

// 새 메서드는 새 배열 반환 (원본 변경 없음)
const sorted = original.toSorted();           // [1, 1, 3, 4, 5]
const reversed = original.toReversed();       // [5, 1, 4, 1, 3]
const spliced = original.toSpliced(1, 2, 9);  // [3, 9, 1, 5]
const changed = original.with(2, 99);         // [3, 1, 99, 1, 5]

console.log(original); // [3, 1, 4, 1, 5] - 변경 없음
```

Object.groupBy와 Map.groupBy:
```javascript
const items = [
  { type: 'fruit', name: 'apple' },
  { type: 'vegetable', name: 'carrot' },
  { type: 'fruit', name: 'banana' },
];

const grouped = Object.groupBy(items, item => item.type);
// { fruit: [{...}, {...}], vegetable: [{...}] }

const mapGrouped = Map.groupBy(items, item => item.type);
// Map { 'fruit' => [...], 'vegetable' => [...] }
```

### ES2025 기능

Import Attributes (JSON 모듈):
```javascript
import config from './config.json' with { type: 'json' };
import styles from './styles.css' with { type: 'css' };

console.log(config.apiUrl);
```

RegExp.escape:
```javascript
const userInput = 'hello (world)';
const safePattern = RegExp.escape(userInput);
// "hello\\ \\(world\\)"
const regex = new RegExp(safePattern);
```

### Node.js 22 LTS 기능

내장 WebSocket 클라이언트:
```javascript
const ws = new WebSocket('wss://example.com/socket');

ws.addEventListener('open', () => {
  ws.send(JSON.stringify({ type: 'hello' }));
});

ws.addEventListener('message', (event) => {
  const data = JSON.parse(event.data);
  console.log('Received:', data);
});
```

네이티브 TypeScript 지원 (실험적):
```bash
# Node.js 22.6+에서 .ts 파일 직접 실행
node --experimental-strip-types app.ts

# Node.js 22.18+에서는 타입 제거가 기본 활성화
node app.ts
```

Watch 모드 (안정화):
```bash
# 파일 변경 시 자동 재시작
node --watch server.js

# 특정 파일 감시
node --watch-path=./src --watch-path=./config server.js
```

권한 모델:
```bash
# 파일 시스템 접근 제한
node --permission --allow-fs-read=/app/data server.js

# 네트워크 접근 제한
node --permission --allow-net=api.example.com server.js
```

### 백엔드 프레임워크

Express (전통적):
```javascript
import express from 'express';

const app = express();
app.use(express.json());

app.get('/api/users', async (req, res) => {
  const users = await db.users.findAll();
  res.json(users);
});

app.post('/api/users', async (req, res) => {
  const user = await db.users.create(req.body);
  res.status(201).json(user);
});

app.listen(3000, () => console.log('Server running on port 3000'));
```

Fastify (고성능):
```javascript
import Fastify from 'fastify';

const fastify = Fastify({ logger: true });

const userSchema = {
  body: {
    type: 'object',
    required: ['name', 'email'],
    properties: {
      name: { type: 'string', minLength: 2 },
      email: { type: 'string', format: 'email' },
    },
  },
};

fastify.post('/api/users', { schema: userSchema }, async (request, reply) => {
  const user = await db.users.create(request.body);
  return reply.code(201).send(user);
});

await fastify.listen({ port: 3000 });
```

Hono (Edge 우선):
```javascript
import { Hono } from 'hono';
import { cors } from 'hono/cors';
import { logger } from 'hono/logger';
import { validator } from 'hono/validator';

const app = new Hono();

app.use('*', logger());
app.use('/api/*', cors());

app.get('/api/users', async (c) => {
  const users = await db.users.findAll();
  return c.json(users);
});

app.post('/api/users',
  validator('json', (value, c) => {
    if (!value.name || !value.email) {
      return c.json({ error: 'Invalid input' }, 400);
    }
    return value;
  }),
  async (c) => {
    const user = await db.users.create(c.req.valid('json'));
    return c.json(user, 201);
  }
);

export default app;
```

### Vitest 테스팅

설정:
```javascript
// vitest.config.js
import { defineConfig } from 'vitest/config';

export default defineConfig({
  test: {
    globals: true,
    environment: 'node',
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html'],
    },
  },
});
```

테스트 예제:
```javascript
// user.test.js
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { createUser, getUser } from './user.js';

describe('User Service', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should create a user', async () => {
    const user = await createUser({ name: 'John', email: 'john@example.com' });
    expect(user).toMatchObject({ name: 'John', email: 'john@example.com' });
    expect(user.id).toBeDefined();
  });

  it('should throw on invalid email', async () => {
    await expect(createUser({ name: 'John', email: 'invalid' }))
      .rejects.toThrow('Invalid email');
  });
});
```

### ESLint 9 Flat Config

```javascript
// eslint.config.js
import js from '@eslint/js';
import globals from 'globals';

export default [
  js.configs.recommended,
  {
    languageOptions: {
      ecmaVersion: 2025,
      sourceType: 'module',
      globals: {
        ...globals.node,
        ...globals.es2025,
      },
    },
    rules: {
      'no-unused-vars': ['error', { argsIgnorePattern: '^_' }],
      'no-console': ['warn', { allow: ['warn', 'error'] }],
      'prefer-const': 'error',
      'no-var': 'error',
    },
  },
];
```

### Biome (올인원)

```json
{
  "$schema": "https://biomejs.dev/schemas/1.9.0/schema.json",
  "organizeImports": { "enabled": true },
  "linter": {
    "enabled": true,
    "rules": { "recommended": true }
  },
  "formatter": {
    "enabled": true,
    "indentStyle": "space",
    "indentWidth": 2
  },
  "javascript": {
    "formatter": {
      "quoteStyle": "single",
      "semicolons": "always"
    }
  }
}
```

---

## 고급 패턴

포괄적인 문서 (고급 비동기 패턴, 모듈 시스템 상세, 성능 최적화, 프로덕션 배포 설정 포함):

- [reference.md](reference.md) - 완전한 API 레퍼런스, Context7 라이브러리 매핑, 패키지 관리자 비교
- [examples.md](examples.md) - 프로덕션 수준 코드 예제, 풀스택 패턴, 테스팅 템플릿

### Context7 통합

```javascript
// Node.js - mcp__context7__query-docs("/nodejs/node", "esm modules async")
// Express - mcp__context7__query-docs("/expressjs/express", "middleware routing")
// Fastify - mcp__context7__query-docs("/fastify/fastify", "plugins hooks")
// Hono - mcp__context7__query-docs("/honojs/hono", "middleware validators")
// Vitest - mcp__context7__query-docs("/vitest-dev/vitest", "mocking coverage")
```

---

## 관련 기술

- `do-lang-typescript` - TypeScript 통합, JSDoc 타입 검사
- `do-domain-backend` - API 설계, 마이크로서비스 아키텍처
- `do-domain-database` - 데이터베이스 통합, ORM 패턴
- `do-workflow-testing` - TDD 워크플로우, 테스팅 전략
- `do-foundation-quality` - 코드 품질 표준

---

## 빠른 문제 해결

모듈 시스템 문제:
```bash
# package.json type 확인
cat package.json | grep '"type"'

# ESM: "type": "module" - import/export 사용
# CommonJS: "type": "commonjs" 또는 생략 - require/module.exports 사용
```

Node.js 버전 확인:
```bash
node --version  # 20.x 또는 22.x LTS 권장
npm --version   # 10.x+ 권장
```

일반적인 수정:
```bash
# npm 캐시 정리
npm cache clean --force

# node_modules 삭제 후 재설치
rm -rf node_modules package-lock.json && npm install

# 권한 문제 해결
npm config set prefix ~/.npm-global
```

ESM/CommonJS 상호 운용:
```javascript
// ESM에서 CommonJS 가져오기
import pkg from 'commonjs-package';
const { namedExport } = pkg;

// CommonJS에서 동적 import
const { default: esmModule } = await import('esm-package');
```

---

Last Updated: 2026-01-06
Status: Active (v1.0.0)
