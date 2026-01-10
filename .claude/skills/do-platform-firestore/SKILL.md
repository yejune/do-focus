---
name: do-platform-firestore
description: Firebase Firestore 전문가 - NoSQL 패턴, 실시간 동기화, 오프라인 캐싱, Security Rules. 모바일 우선 앱, 오프라인 지원, 실시간 listener, Firestore 보안 설정에 활용.
version: 1.0.0
category: platform
updated: 2025-01-06
status: active
tags:
  - firestore
  - firebase
  - nosql
  - realtime
  - offline
  - mobile
context7-libraries:
  - /firebase/firebase-docs
related-skills:
  - do-platform-firebase-auth
  - do-lang-flutter
  - do-lang-typescript
allowed-tools: Read, Write, Bash, Grep, Glob
user-invocable: false
---

# do-platform-firestore: Firebase Firestore 전문가

## 역할

Firebase Firestore NoSQL 문서 데이터베이스 전문가. 실시간 동기화, 오프라인 우선 아키텍처, Security Rules, Cloud Functions 트리거, 모바일 최적화 SDK 구현을 담당.

---

## 핵심 기능

### 실시간 동기화
연결된 모든 클라이언트에 변경사항 자동 동기화

### 오프라인 캐싱
IndexedDB 기반 로컬 지속성 및 온라인 복귀 시 자동 동기화

### Security Rules
선언형 필드 수준 접근 제어

### Cloud Functions
문서 변경에 반응하는 서버 측 트리거

### Composite Indexes
복잡한 쿼리 최적화를 위한 복합 인덱스

---

## 사용 시점

- 오프라인 지원이 필요한 모바일 우선 애플리케이션
- 실시간 협업 기능 구현
- 크로스 플랫폼 앱 개발 (iOS, Android, Web, Flutter)
- Google Cloud 통합이 필요한 프로젝트
- 유연하고 변화하는 데이터 구조가 필요한 앱

---

## Firestore 초기화

### 오프라인 지속성 설정

```typescript
import { initializeApp } from 'firebase/app'
import {
  initializeFirestore,
  persistentLocalCache,
  persistentMultipleTabManager,
  CACHE_SIZE_UNLIMITED
} from 'firebase/firestore'

const app = initializeApp({
  apiKey: process.env.NEXT_PUBLIC_FIREBASE_API_KEY,
  authDomain: process.env.NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN,
  projectId: process.env.NEXT_PUBLIC_FIREBASE_PROJECT_ID
})

export const db = initializeFirestore(app, {
  localCache: persistentLocalCache({
    tabManager: persistentMultipleTabManager(),
    cacheSizeBytes: CACHE_SIZE_UNLIMITED
  })
})
```

---

## 실시간 Listener

### 메타데이터 포함 구독

```typescript
import { collection, query, where, orderBy, onSnapshot } from 'firebase/firestore'

export function subscribeToDocuments(userId: string, callback: (docs: any[]) => void) {
  const q = query(
    collection(db, 'documents'),
    where('collaborators', 'array-contains', userId),
    orderBy('createdAt', 'desc')
  )

  return onSnapshot(q, { includeMetadataChanges: true }, (snapshot) => {
    callback(snapshot.docs.map((doc) => ({
      id: doc.id,
      ...doc.data(),
      _pending: doc.metadata.hasPendingWrites,
      _fromCache: doc.metadata.fromCache
    })))
  })
}
```

메타데이터 플래그 설명:
- `hasPendingWrites`: 서버에 아직 커밋되지 않은 로컬 변경사항 존재
- `fromCache`: 데이터가 로컬 캐시에서 제공됨

---

## Security Rules

### 기본 구조

```javascript
rules_version = '2';
service cloud.firestore {
  match /databases/{database}/documents {
    match /users/{userId} {
      allow read, write: if request.auth.uid == userId;
    }

    match /documents/{docId} {
      allow read: if resource.data.isPublic == true
        || request.auth.uid == resource.data.ownerId
        || request.auth.uid in resource.data.collaborators;
      allow create: if request.auth != null
        && request.resource.data.ownerId == request.auth.uid;
      allow update, delete: if request.auth.uid == resource.data.ownerId;
    }
  }
}
```

### 역할 기반 접근 제어 (Custom Claims 활용)

```javascript
rules_version = '2';
service cloud.firestore {
  match /databases/{database}/documents {
    function isSignedIn() { return request.auth != null; }
    function isAdmin() { return request.auth.token.admin == true; }

    match /organizations/{orgId} {
      function isMember() {
        return exists(/databases/$(database)/documents/organizations/$(orgId)/members/$(request.auth.uid));
      }
      function getMemberRole() {
        return get(/databases/$(database)/documents/organizations/$(orgId)/members/$(request.auth.uid)).data.role;
      }
      function isOrgAdmin() { return isMember() && getMemberRole() in ['admin', 'owner']; }

      allow read: if isSignedIn() && isMember();
      allow update: if isOrgAdmin();
      allow delete: if getMemberRole() == 'owner';

      match /members/{memberId} {
        allow read: if isMember();
        allow write: if isOrgAdmin();
      }

      match /projects/{projectId} {
        allow read: if isMember();
        allow create: if isMember() && getMemberRole() in ['admin', 'owner', 'editor'];
        allow update, delete: if isOrgAdmin() || resource.data.createdBy == request.auth.uid;
      }
    }
  }
}
```

---

## Composite Indexes 설정

firestore.indexes.json 파일 구성:

```json
{
  "indexes": [
    {
      "collectionGroup": "documents",
      "queryScope": "COLLECTION",
      "fields": [
        { "fieldPath": "organizationId", "order": "ASCENDING" },
        { "fieldPath": "createdAt", "order": "DESCENDING" }
      ]
    },
    {
      "collectionGroup": "documents",
      "queryScope": "COLLECTION",
      "fields": [
        { "fieldPath": "tags", "arrayConfig": "CONTAINS" },
        { "fieldPath": "createdAt", "order": "DESCENDING" }
      ]
    }
  ]
}
```

---

## Cloud Functions V2 트리거

### 문서 업데이트 트리거

```typescript
import { onDocumentUpdated } from 'firebase-functions/v2/firestore'
import { getFirestore, FieldValue } from 'firebase-admin/firestore'

const db = getFirestore()

export const onDocumentUpdate = onDocumentUpdated(
  { document: 'documents/{docId}', region: 'us-central1' },
  async (event) => {
    const before = event.data?.before.data()
    const after = event.data?.after.data()
    if (!before || !after) return

    const batch = db.batch()
    batch.set(db.collection('changes').doc(), {
      documentId: event.params.docId,
      before, after,
      changedAt: FieldValue.serverTimestamp()
    })
    batch.update(db.doc('stats/documents'), {
      totalModifications: FieldValue.increment(1)
    })
    await batch.commit()
  }
)
```

### Callable Function

```typescript
import { onCall, HttpsError } from 'firebase-functions/v2/https'
import { getFirestore, FieldValue } from 'firebase-admin/firestore'

const db = getFirestore()

export const inviteToOrganization = onCall({ region: 'us-central1' }, async (request) => {
  if (!request.auth) throw new HttpsError('unauthenticated', 'Must be signed in')

  const { organizationId, email, role } = request.data
  const memberDoc = await db.doc(`organizations/${organizationId}/members/${request.auth.uid}`).get()

  if (!memberDoc.exists || !['admin', 'owner'].includes(memberDoc.data()?.role)) {
    throw new HttpsError('permission-denied', 'Must be organization admin')
  }

  const invitation = await db.collection('invitations').add({
    organizationId, email, role,
    invitedBy: request.auth.uid,
    createdAt: FieldValue.serverTimestamp(),
    status: 'pending'
  })

  return { invitationId: invitation.id }
})
```

### 스케줄 트리거

```typescript
import { onSchedule } from 'firebase-functions/v2/scheduler'
import { getFirestore } from 'firebase-admin/firestore'

const db = getFirestore()

export const dailyCleanup = onSchedule(
  { schedule: '0 0 * * *', timeZone: 'UTC', region: 'us-central1' },
  async () => {
    const thirtyDaysAgo = new Date()
    thirtyDaysAgo.setDate(thirtyDaysAgo.getDate() - 30)

    const oldDocs = await db.collection('tempFiles')
      .where('createdAt', '<', thirtyDaysAgo).limit(500).get()

    const batch = db.batch()
    oldDocs.docs.forEach((doc) => batch.delete(doc.ref))
    await batch.commit()
  }
)
```

---

## 오프라인 우선 React Hook

```typescript
import { useEffect, useState } from 'react'
import {
  collection, query, where, orderBy, onSnapshot,
  addDoc, updateDoc, doc, serverTimestamp
} from 'firebase/firestore'

export function useTasks(userId: string) {
  const [tasks, setTasks] = useState<any[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!userId) return
    const q = query(
      collection(db, 'tasks'),
      where('userId', '==', userId),
      orderBy('createdAt', 'desc')
    )

    return onSnapshot(q, { includeMetadataChanges: true }, (snapshot) => {
      setTasks(snapshot.docs.map((doc) => ({
        id: doc.id,
        ...doc.data(),
        _pending: doc.metadata.hasPendingWrites,
        _fromCache: doc.metadata.fromCache
      })))
      setLoading(false)
    })
  }, [userId])

  const addTask = (title: string) => addDoc(collection(db, 'tasks'), {
    title, completed: false, userId, createdAt: serverTimestamp()
  })

  const toggleTask = (taskId: string, completed: boolean) =>
    updateDoc(doc(db, 'tasks', taskId), { completed, updatedAt: serverTimestamp() })

  return { tasks, loading, addTask, toggleTask }
}
```

---

## Batch 및 Transaction

### Batch 작업

```typescript
import { writeBatch, doc, serverTimestamp } from 'firebase/firestore'

async function batchUpdate(updates: Array<{ id: string; data: any }>) {
  const batch = writeBatch(db)
  updates.forEach(({ id, data }) => {
    batch.update(doc(db, 'documents', id), { ...data, updatedAt: serverTimestamp() })
  })
  await batch.commit()
}
```

### Transaction

```typescript
import { runTransaction, doc, increment } from 'firebase/firestore'

async function transferCredits(fromUserId: string, toUserId: string, amount: number) {
  await runTransaction(db, async (transaction) => {
    const fromRef = doc(db, 'users', fromUserId)
    const toRef = doc(db, 'users', toUserId)
    const fromDoc = await transaction.get(fromRef)

    if (!fromDoc.exists()) throw new Error('Sender not found')
    if (fromDoc.data().credits < amount) throw new Error('Insufficient credits')

    transaction.update(fromRef, { credits: increment(-amount) })
    transaction.update(toRef, { credits: increment(amount) })
  })
}
```

Batch와 Transaction 차이점:
- Batch: 최대 500개 작업, 원자적 커밋, 읽기 불가
- Transaction: 읽기/쓰기 모두 가능, 충돌 시 자동 재시도, 최대 20초 제한

---

## 성능 특성

### 지연시간
- 읽기: 50-200ms (리전에 따라 상이)
- 쓰기: 100-300ms
- 실시간 전파: 100-500ms
- 오프라인 동기화: 연결 복구 시 자동

### 무료 티어 (2024 기준)
- 저장소: 1GB
- 일일 읽기: 50,000
- 일일 쓰기: 20,000
- 일일 삭제: 20,000

---

## 모범 사례

### 데이터 모델링
- 자주 함께 읽는 데이터는 같은 문서에 저장
- 컬렉션 그룹 쿼리를 활용한 중첩 데이터 접근
- 역정규화를 통한 읽기 최적화

### 쿼리 최적화
- 필요한 필드만 선택하여 대역폭 절약
- 적절한 인덱스 설정으로 쿼리 성능 향상
- limit() 사용으로 불필요한 데이터 전송 방지

### 보안
- 모든 경로에 Security Rules 적용
- 클라이언트 입력 검증
- Custom Claims를 통한 역할 기반 접근 제어

---

## 관련 Skills

- do-platform-firebase-auth - Firebase Authentication 통합
- do-lang-flutter - Flutter SDK 패턴
- do-lang-typescript - TypeScript 클라이언트 패턴

---

Status: Production Ready
Platform: Firebase Firestore
Last Updated: 2025-01-06
