---
name: do-platform-supabase
description: PostgreSQL 16, pgvector, RLS, 실시간 구독, Edge Functions를 다루는 Supabase 전문가. Supabase 백엔드를 활용한 풀스택 앱 개발 시 사용
version: 1.0.0
category: platform
tags: [supabase, postgresql, pgvector, realtime, rls, edge-functions]
context7-libraries: [/supabase/supabase]
related-skills: [do-platform-neon, do-lang-typescript]
updated: 2025-01-06
status: active
allowed-tools: [Read, Write, Bash, Grep, Glob]
user-invocable: false
---

# Supabase 플랫폼 전문가

## 빠른 참조

PostgreSQL 16 기반 풀스택 플랫폼. AI/벡터 검색용 pgvector, 멀티테넌트용 Row-Level Security, 실시간 구독, Deno Edge Functions, 이미지 변환 Storage 통합 제공.

### 핵심 기능

- PostgreSQL 16: SQL, JSONB, 고급 데이터 타입 지원
- pgvector: HNSW/IVFFlat 인덱스를 통한 AI 임베딩 유사도 검색
- Row-Level Security: 데이터베이스 레벨 멀티테넌트 데이터 격리
- Real-time: Postgres Changes와 Presence를 통한 라이브 데이터 동기화
- Edge Functions: Deno 런타임 서버리스 함수
- Storage: 자동 이미지 변환 지원 파일 저장소
- Auth: JWT 기반 인증 통합

### 사용 시점

- 데이터 격리가 필요한 멀티테넌트 SaaS 애플리케이션
- 벡터 임베딩과 유사도 검색이 필요한 AI/ML 애플리케이션
- Presence, 라이브 업데이트 등 실시간 협업 기능
- 인증, 데이터베이스, 저장소가 필요한 풀스택 애플리케이션
- PostgreSQL 고유 기능이 필요한 프로젝트

---

## 구현 가이드

### PostgreSQL 16 + pgvector 설정

익스텐션 활성화 및 임베딩 테이블 생성:
```sql
-- 필수 익스텐션 활성화
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS vector;

-- 시맨틱 검색용 임베딩 테이블
CREATE TABLE documents (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  content TEXT NOT NULL,
  embedding vector(1536),  -- OpenAI ada-002 차원
  metadata JSONB DEFAULT '{}',
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- HNSW 인덱스: 빠른 유사도 검색 (권장)
CREATE INDEX idx_documents_embedding ON documents
USING hnsw (embedding vector_cosine_ops)
WITH (m = 16, ef_construction = 64);

-- IVFFlat 인덱스: 대용량 데이터셋 (수백만 건)
-- CREATE INDEX idx_documents_ivf ON documents
-- USING ivfflat (embedding vector_cosine_ops)
-- WITH (lists = 100);
```

시맨틱 검색 함수:
```sql
CREATE OR REPLACE FUNCTION search_documents(
  query_embedding vector(1536),
  match_threshold FLOAT DEFAULT 0.8,
  match_count INT DEFAULT 10
) RETURNS TABLE (id UUID, content TEXT, similarity FLOAT)
LANGUAGE plpgsql AS $$
BEGIN
  RETURN QUERY SELECT d.id, d.content,
    1 - (d.embedding <=> query_embedding) AS similarity
  FROM documents d
  WHERE 1 - (d.embedding <=> query_embedding) > match_threshold
  ORDER BY d.embedding <=> query_embedding
  LIMIT match_count;
END; $$;
```

하이브리드 검색 (벡터 + 전문 검색):
```sql
CREATE OR REPLACE FUNCTION hybrid_search(
  query_text TEXT,
  query_embedding vector(1536),
  match_count INT DEFAULT 10,
  full_text_weight FLOAT DEFAULT 0.3,
  semantic_weight FLOAT DEFAULT 0.7
) RETURNS TABLE (id UUID, content TEXT, score FLOAT) AS $$
BEGIN
  RETURN QUERY
  WITH semantic AS (
    SELECT e.id, e.content, 1 - (e.embedding <=> query_embedding) AS similarity
    FROM documents e ORDER BY e.embedding <=> query_embedding LIMIT match_count * 2
  ),
  full_text AS (
    SELECT e.id, e.content,
      ts_rank(to_tsvector('english', e.content), plainto_tsquery('english', query_text)) AS rank
    FROM documents e
    WHERE to_tsvector('english', e.content) @@ plainto_tsquery('english', query_text)
    LIMIT match_count * 2
  )
  SELECT COALESCE(s.id, f.id), COALESCE(s.content, f.content),
    (COALESCE(s.similarity, 0) * semantic_weight + COALESCE(f.rank, 0) * full_text_weight)
  FROM semantic s FULL OUTER JOIN full_text f ON s.id = f.id
  ORDER BY 3 DESC LIMIT match_count;
END; $$ LANGUAGE plpgsql;
```

### Row-Level Security (RLS) 패턴

기본 테넌트 격리:
```sql
ALTER TABLE projects ENABLE ROW LEVEL SECURITY;

-- JWT 클레임 기반 정책
CREATE POLICY "tenant_isolation" ON projects FOR ALL
  USING (tenant_id = (auth.jwt() ->> 'tenant_id')::UUID);

-- 소유자 기반 접근
CREATE POLICY "owner_access" ON projects FOR ALL
  USING (owner_id = auth.uid());
```

계층적 접근 권한을 가진 멀티테넌트:
```sql
-- 조직 기반 접근
CREATE POLICY "org_member_select" ON organizations FOR SELECT
  USING (id IN (SELECT org_id FROM org_members WHERE user_id = auth.uid()));

-- 역할 기반 수정 권한
CREATE POLICY "org_admin_modify" ON organizations FOR UPDATE
  USING (id IN (
    SELECT org_id FROM org_members
    WHERE user_id = auth.uid() AND role IN ('owner', 'admin')
  ));

-- 조직 멤버십을 통한 프로젝트 접근 권한 전파
CREATE POLICY "project_access" ON projects FOR ALL
  USING (org_id IN (SELECT org_id FROM org_members WHERE user_id = auth.uid()));

-- 서버사이드 작업용 서비스 역할 우회
CREATE POLICY "service_bypass" ON organizations FOR ALL TO service_role USING (true);
```

### Real-time 구독

테이블 변경 구독:
```typescript
import { createClient } from '@supabase/supabase-js'

const supabase = createClient(SUPABASE_URL, SUPABASE_ANON_KEY)

// 테이블의 모든 변경사항 구독
const channel = supabase.channel('db-changes')
  .on('postgres_changes',
    { event: '*', schema: 'public', table: 'messages' },
    (payload) => console.log('Change:', payload)
  )
  .subscribe()

// 특정 조건 필터링
supabase.channel('project-updates')
  .on('postgres_changes',
    { event: 'UPDATE', schema: 'public', table: 'projects', filter: `id=eq.${projectId}` },
    (payload) => handleProjectUpdate(payload.new)
  )
  .subscribe()
```

Presence 추적:
```typescript
interface PresenceState {
  user_id: string
  online_at: string
  typing?: boolean
  cursor?: { x: number; y: number }
}

const channel = supabase.channel('room:collaborative-doc', {
  config: { presence: { key: userId } }
})

channel
  .on('presence', { event: 'sync' }, () => {
    const state = channel.presenceState<PresenceState>()
    console.log('Online users:', Object.keys(state))
  })
  .on('presence', { event: 'join' }, ({ key, newPresences }) => {
    console.log('User joined:', key, newPresences)
  })
  .on('presence', { event: 'leave' }, ({ key, leftPresences }) => {
    console.log('User left:', key, leftPresences)
  })
  .subscribe(async (status) => {
    if (status === 'SUBSCRIBED') {
      await channel.track({ user_id: userId, online_at: new Date().toISOString() })
    }
  })

// Presence 상태 업데이트
await channel.track({ typing: true })
await channel.track({ cursor: { x: 100, y: 200 } })
```

### Edge Functions

인증 포함 기본 Edge Function:
```typescript
// supabase/functions/api/index.ts
import { serve } from 'https://deno.land/std@0.168.0/http/server.ts'
import { createClient } from 'https://esm.sh/@supabase/supabase-js@2'

const corsHeaders = {
  'Access-Control-Allow-Origin': '*',
  'Access-Control-Allow-Headers': 'authorization, x-client-info, apikey, content-type'
}

serve(async (req) => {
  if (req.method === 'OPTIONS') {
    return new Response('ok', { headers: corsHeaders })
  }

  const supabase = createClient(
    Deno.env.get('SUPABASE_URL')!,
    Deno.env.get('SUPABASE_SERVICE_ROLE_KEY')!
  )

  // JWT 토큰 검증
  const authHeader = req.headers.get('authorization')
  if (!authHeader) {
    return new Response(JSON.stringify({ error: 'Unauthorized' }),
      { status: 401, headers: { ...corsHeaders, 'Content-Type': 'application/json' } })
  }

  const { data: { user }, error } = await supabase.auth.getUser(
    authHeader.replace('Bearer ', '')
  )

  if (error || !user) {
    return new Response(JSON.stringify({ error: 'Invalid token' }),
      { status: 401, headers: { ...corsHeaders, 'Content-Type': 'application/json' } })
  }

  // 요청 처리
  const body = await req.json()
  return new Response(JSON.stringify({ success: true, user_id: user.id }),
    { headers: { ...corsHeaders, 'Content-Type': 'application/json' } })
})
```

Edge Functions 속도 제한:
```typescript
async function checkRateLimit(
  supabase: SupabaseClient, identifier: string, limit: number, windowSeconds: number
): Promise<boolean> {
  const windowStart = new Date(Date.now() - windowSeconds * 1000).toISOString()
  const { count } = await supabase
    .from('rate_limits')
    .select('*', { count: 'exact', head: true })
    .eq('identifier', identifier)
    .gte('created_at', windowStart)

  if (count && count >= limit) return false
  await supabase.from('rate_limits').insert({ identifier })
  return true
}
```

### Storage 이미지 변환

```typescript
const supabase = createClient(SUPABASE_URL, SUPABASE_ANON_KEY)

async function uploadImage(file: File, userId: string) {
  const fileName = `${userId}/${Date.now()}-${file.name}`

  const { data, error } = await supabase.storage
    .from('images')
    .upload(fileName, file, { cacheControl: '3600', upsert: false })

  if (error) throw error

  // 변환된 URL 가져오기
  const { data: { publicUrl } } = supabase.storage
    .from('images')
    .getPublicUrl(fileName, {
      transform: { width: 800, height: 600, resize: 'contain' }
    })

  const { data: { publicUrl: thumbnailUrl } } = supabase.storage
    .from('images')
    .getPublicUrl(fileName, {
      transform: { width: 200, height: 200, resize: 'cover' }
    })

  return { originalPath: data.path, publicUrl, thumbnailUrl }
}
```

---

## 고급 패턴

### 멀티테넌트 SaaS 아키텍처

전체 스키마 설정:
```sql
-- 조직 (테넌트)
CREATE TABLE organizations (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name TEXT NOT NULL,
  slug TEXT UNIQUE NOT NULL,
  plan TEXT DEFAULT 'free' CHECK (plan IN ('free', 'pro', 'enterprise')),
  settings JSONB DEFAULT '{}',
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 역할이 있는 조직 멤버
CREATE TABLE organization_members (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
  user_id UUID NOT NULL,
  role TEXT NOT NULL CHECK (role IN ('owner', 'admin', 'member', 'viewer')),
  joined_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(organization_id, user_id)
);

-- 조직 내 프로젝트
CREATE TABLE projects (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  owner_id UUID NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 모든 테이블에 RLS 활성화
ALTER TABLE organizations ENABLE ROW LEVEL SECURITY;
ALTER TABLE organization_members ENABLE ROW LEVEL SECURITY;
ALTER TABLE projects ENABLE ROW LEVEL SECURITY;

-- 포괄적인 RLS 정책 생성
CREATE POLICY "org_member_select" ON organizations FOR SELECT
  USING (id IN (SELECT organization_id FROM organization_members WHERE user_id = auth.uid()));

CREATE POLICY "org_admin_update" ON organizations FOR UPDATE
  USING (id IN (SELECT organization_id FROM organization_members
    WHERE user_id = auth.uid() AND role IN ('owner', 'admin')));

CREATE POLICY "project_member_access" ON projects FOR ALL
  USING (organization_id IN (SELECT organization_id FROM organization_members WHERE user_id = auth.uid()));
```

### TypeScript 클라이언트 패턴

서버사이드 클라이언트 (Next.js App Router):
```typescript
import { createServerClient } from '@supabase/ssr'
import { cookies } from 'next/headers'
import { Database } from './database.types'

export function createServerSupabase() {
  const cookieStore = cookies()
  return createServerClient<Database>(
    process.env.NEXT_PUBLIC_SUPABASE_URL!,
    process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY!,
    {
      cookies: {
        get(name: string) { return cookieStore.get(name)?.value },
        set(name, value, options) { cookieStore.set({ name, value, ...options }) },
        remove(name, options) { cookieStore.set({ name, value: '', ...options }) }
      }
    }
  )
}
```

서비스 레이어 패턴:
```typescript
import { supabase } from './supabase/client'

export class DocumentService {
  async create(projectId: string, title: string, content: string) {
    const { data: { user } } = await supabase.auth.getUser()
    const { data, error } = await supabase
      .from('documents')
      .insert({ project_id: projectId, title, content, created_by: user!.id })
      .select().single()

    if (error) throw error

    // 비동기로 임베딩 생성
    await supabase.functions.invoke('generate-embedding',
      { body: { documentId: data.id, content } })

    return data
  }

  async semanticSearch(projectId: string, query: string) {
    const { data: embeddingData } = await supabase.functions.invoke(
      'get-embedding', { body: { text: query } })

    const { data, error } = await supabase.rpc('search_documents', {
      p_project_id: projectId,
      p_query_embedding: embeddingData.embedding,
      p_match_threshold: 0.7,
      p_match_count: 10
    })

    if (error) throw error
    return data
  }

  subscribeToChanges(projectId: string, callback: (payload: any) => void) {
    return supabase.channel(`documents:${projectId}`)
      .on('postgres_changes',
        { event: '*', schema: 'public', table: 'documents', filter: `project_id=eq.${projectId}` },
        callback)
      .subscribe()
  }
}
```

---

## 모범 사례

### 성능 최적화

- 벡터 인덱스: HNSW 사용 (빠른 응답), IVFFlat (수백만 건 데이터)
- 연결 풀링: 서버리스 환경에서 Supavisor 사용
- 쿼리 최적화: 필요한 컬럼만 select, limit 적용

### 보안

- RLS 필수: 모든 테이블에 Row-Level Security 활성화
- JWT 검증: 모든 API 요청에서 토큰 검증
- service_role 제한: Edge Functions 내부에서만 사용

### 마이그레이션

- Supabase CLI 활용: supabase migration new, supabase db push
- 버전 관리: 모든 스키마 변경 마이그레이션 파일로 관리
- 테스트 우선: 로컬에서 마이그레이션 테스트 후 배포

---

## 관련 스킬

- do-platform-neon: PostgreSQL 대체 옵션
- do-lang-typescript: Supabase TypeScript 클라이언트 패턴
- do-domain-backend: 백엔드 아키텍처 통합
- do-workflow-testing: Supabase 테스트 주도 개발

---

Status: Production Ready
Last Updated: 2025-01-06
Coverage: PostgreSQL 16, pgvector, RLS, Real-time, Edge Functions, Storage
