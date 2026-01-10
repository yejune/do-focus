---
name: do-platform-firebase-auth
description: Firebase Authentication 전문가 - Google 에코시스템, 소셜 인증, 휴대폰 인증, 모바일 중심 패턴을 다룸. Firebase 기반 또는 Google 에코시스템 앱 구축 시 사용.
version: 1.0.0
category: platform
updated: 2025-12-07
status: active
tags: firebase, google, social-auth, mobile, authentication
context7-libraries: /firebase/firebase-docs
related-skills: do-platform-firestore, do-lang-flutter
allowed-tools: Read, Write, Bash, Grep, Glob
user-invocable: false
---

# Firebase Authentication 전문가

Google 에코시스템 통합, 소셜 인증 제공자, 휴대폰 인증, 익명 인증, 커스텀 클레임, Security Rules 통합을 포괄하는 Firebase Authentication 구현 가이드.

## 빠른 참조

Firebase Auth 핵심 기능:

- Google Sign-In: Cloud Identity를 통한 네이티브 Google 에코시스템 통합
- Social Auth: Facebook, Twitter/X, GitHub, Apple, Microsoft, Yahoo
- Phone Auth: 국제 지원을 포함한 SMS 기반 인증
- Anonymous Auth: 계정 연결 기능이 있는 게스트 접근
- Custom Claims: 역할 기반 접근 제어 및 관리자 권한
- Security Rules: Firestore, Storage, Realtime Database 통합

Context7 라이브러리 접근:

- Firebase Documentation: /firebase/firebase-docs
- resolve-library-id에 "firebase"를 사용한 후 get-library-docs로 최신 API 확인

플랫폼 SDK 지원:

- Web: firebase/auth modular SDK (v9+)
- iOS: FirebaseAuth with Swift/SwiftUI
- Android: firebase-auth with Kotlin
- Flutter: firebase_auth 패키지
- React Native: @react-native-firebase/auth

결정 가이드:

- Google 에코시스템 통합 필요? Firebase Auth 사용
- 모바일 중심 애플리케이션? Firebase Auth 사용
- 서버리스 Cloud Functions 필요? Firebase Auth 사용
- 익명 게스트 접근 필요? Firebase Auth 사용
- 기존 Firebase 인프라 보유? Firebase Auth 사용

---

## 구현 가이드

### Google Sign-In 통합

Google Sign-In은 Google API와 서비스에 대한 접근 권한과 함께 Google 에코시스템 내에서 원활한 인증 제공.

Web 구현:

Step 1: Firebase Console의 Authentication에서 Google Sign-In 제공자 활성화
Step 2: Google Cloud Console에서 OAuth consent screen 구성
Step 3: Firebase Auth SDK 가져오기 및 구성

```typescript
import { getAuth, signInWithPopup, GoogleAuthProvider } from 'firebase/auth';
const auth = getAuth();
const provider = new GoogleAuthProvider();
provider.addScope('https://www.googleapis.com/auth/calendar.readonly');
const result = await signInWithPopup(auth, provider);
const credential = GoogleAuthProvider.credentialFromResult(result);
```

Flutter 구현:

```dart
Future<UserCredential> signInWithGoogle() async {
  final googleUser = await GoogleSignIn().signIn();
  final googleAuth = await googleUser?.authentication;
  final credential = GoogleAuthProvider.credential(
    accessToken: googleAuth?.accessToken, idToken: googleAuth?.idToken);
  return await FirebaseAuth.instance.signInWithCredential(credential);
}
```

모바일 구성 요구사항:

- iOS: reverse client ID로 Info.plist에 URL schemes 구성
- Android: Firebase 프로젝트 설정에 SHA-1 지문 추가
- Web: Firebase Console에서 승인된 도메인 구성

### 소셜 인증 제공자

Firebase Auth는 통합 API로 주요 소셜 ID 제공자 지원.

Facebook 로그인:

```typescript
import { FacebookAuthProvider, signInWithPopup } from 'firebase/auth';
const provider = new FacebookAuthProvider();
provider.addScope('email');
provider.addScope('public_profile');
const result = await signInWithPopup(auth, provider);
```

Apple Sign-In (타사 로그인이 있는 iOS 앱에 필수):

```swift
let provider = OAuthProvider(providerID: "apple.com")
provider.scopes = ["email", "fullName"]
provider.getCredentialWith(nil) { credential, error in
    if let credential = credential {
        Auth.auth().signIn(with: credential) { result, error in }
    }
}
```

Twitter/X 및 GitHub 인증:

```typescript
// Twitter
const twitterProvider = new TwitterAuthProvider();
await signInWithPopup(auth, twitterProvider);

// GitHub
const githubProvider = new GithubAuthProvider();
githubProvider.addScope('repo');
await signInWithPopup(auth, githubProvider);
```

### 휴대폰 번호 인증

국제 지원 및 reCAPTCHA 인증을 포함한 SMS 기반 휴대폰 인증.

Web 구현:

```typescript
import { RecaptchaVerifier, signInWithPhoneNumber } from 'firebase/auth';
const recaptchaVerifier = new RecaptchaVerifier(auth, 'recaptcha-container', {
  size: 'normal', callback: () => { /* reCAPTCHA 해결됨 */ }
});
const confirmationResult = await signInWithPhoneNumber(auth, '+1234567890', recaptchaVerifier);
const credential = await confirmationResult.confirm(verificationCode);
```

Flutter 구현:

```dart
await FirebaseAuth.instance.verifyPhoneNumber(
  phoneNumber: '+1234567890',
  verificationCompleted: (credential) async {
    await FirebaseAuth.instance.signInWithCredential(credential);
  },
  verificationFailed: (e) => print('Failed: ${e.message}'),
  codeSent: (verificationId, resendToken) { /* verificationId 저장 */ },
  codeAutoRetrievalTimeout: (verificationId) {},
);
```

휴대폰 인증 모범 사례:

- 전화번호에 E.164 형식 사용
- Android에서 자동 인증 처리
- 수동 코드 입력 폴백 제공
- SMS 비용 및 속도 제한 고려

### 익명 인증

익명 인증으로 계정 업그레이드 경로와 함께 게스트 접근 활성화.

```typescript
// 익명 로그인
const result = await signInAnonymously(auth);
console.log('Anonymous UID:', result.user.uid);

// 계정 연결 (영구 계정으로 업그레이드)
import { linkWithCredential, EmailAuthProvider } from 'firebase/auth';
const credential = EmailAuthProvider.credential(email, password);
const linked = await linkWithCredential(auth.currentUser, credential);

// 소셜 제공자와 연결
const googleProvider = new GoogleAuthProvider();
await linkWithPopup(auth.currentUser, googleProvider);
```

익명 인증 사용 사례:

- 결제 전 전자상거래 장바구니
- 가입 전 콘텐츠 미리보기
- 계정 생성 전 게임 진행 저장

### 커스텀 클레임과 토큰

커스텀 클레임으로 역할 기반 접근 제어 및 관리자 권한 활성화.

서버 측 (Admin SDK):

```typescript
import { getAuth } from 'firebase-admin/auth';
await getAuth().setCustomUserClaims(uid, { admin: true });
await getAuth().setCustomUserClaims(uid, {
  role: 'editor', organizationId: 'org_123', permissions: ['read', 'write']
});
```

클라이언트 측:

```typescript
const idTokenResult = await auth.currentUser.getIdTokenResult();
if (idTokenResult.claims.admin === true) { /* 관리자 UI 표시 */ }

// 클레임 업데이트 후 토큰 강제 새로고침
await auth.currentUser.getIdToken(true);
```

커스텀 클레임 모범 사례:

- 클레임은 작게 유지 (총 1000바이트 미만)
- 사용자 데이터가 아닌 접근 제어용으로 클레임 사용
- 클레임 변경 후 토큰 새로고침 필요

### Security Rules 통합

Firebase Security Rules는 접근 제어를 위해 인증 상태 사용.

Firestore Security Rules:

```javascript
rules_version = '2';
service cloud.firestore {
  match /databases/{database}/documents {
    match /users/{userId} {
      allow read, write: if request.auth != null && request.auth.uid == userId;
    }
    match /admin/{document=**} {
      allow read, write: if request.auth.token.admin == true;
    }
    match /organizations/{orgId}/documents/{docId} {
      allow read, write: if request.auth.token.organizationId == orgId;
    }
  }
}
```

Cloud Storage Security Rules:

```javascript
rules_version = '2';
service firebase.storage {
  match /b/{bucket}/o {
    match /users/{userId}/{allPaths=**} {
      allow read, write: if request.auth != null && request.auth.uid == userId;
    }
    match /public/{allPaths=**} {
      allow read: if true;
      allow write: if request.auth != null;
    }
  }
}
```

---

## 고급 패턴

### Cloud Functions Auth 트리거

Firebase Cloud Functions는 인증 생명주기 이벤트에 응답.

```typescript
import { auth } from 'firebase-functions';
import { getFirestore } from 'firebase-admin/firestore';

export const onUserCreate = auth.user().onCreate(async (user) => {
  await getFirestore().collection('users').doc(user.uid).set({
    email: user.email, displayName: user.displayName,
    createdAt: FieldValue.serverTimestamp(), role: 'member'
  });
});
```

블로킹 함수:

```typescript
import { beforeUserCreated } from 'firebase-functions/v2/identity';
export const validateUserCreate = beforeUserCreated((event) => {
  if (!event.data.email?.endsWith('@company.com')) {
    throw new HttpsError('invalid-argument', 'Unauthorized email domain');
  }
  return { customClaims: { role: 'employee' } };
});
```

### 다중 요소 인증

Firebase Auth는 SMS 기반 2차 인증 지원.

```typescript
import { multiFactor, PhoneMultiFactorGenerator, PhoneAuthProvider } from 'firebase/auth';

// MFA 등록
const mfUser = multiFactor(auth.currentUser);
const session = await mfUser.getSession();
const verificationId = await new PhoneAuthProvider(auth)
  .verifyPhoneNumber({ phoneNumber: '+1234567890', session }, recaptchaVerifier);
const credential = PhoneAuthProvider.credential(verificationId, code);
await mfUser.enroll(PhoneMultiFactorGenerator.assertion(credential), 'Phone');

// MFA 챌린지와 함께 로그인
try { await signInWithEmailAndPassword(auth, email, password); }
catch (error) {
  if (error.code === 'auth/multi-factor-auth-required') {
    const resolver = getMultiFactorResolver(auth, error);
    // resolver.hints로 MFA 인증 완료
  }
}
```

### 세션 관리

```typescript
import { setPersistence, browserLocalPersistence, browserSessionPersistence,
         inMemoryPersistence, onAuthStateChanged } from 'firebase/auth';

// 지속성 옵션
await setPersistence(auth, browserLocalPersistence);  // 기본값 - 세션 간 유지
await setPersistence(auth, browserSessionPersistence); // 탭 닫으면 초기화
await setPersistence(auth, inMemoryPersistence);       // 메모리에만 유지

// 인증 상태 리스너
const unsubscribe = onAuthStateChanged(auth, (user) => {
  if (user) { console.log('Signed in:', user.uid); }
  else { console.log('Signed out'); }
});
```

### Firebase Auth 에뮬레이터

```typescript
import { connectAuthEmulator } from 'firebase/auth';
if (process.env.NODE_ENV === 'development') {
  connectAuthEmulator(auth, 'http://localhost:9099');
}
```

에뮬레이터 기능:

- 실제 SMS 없이 휴대폰 인증 테스트
- 프로그래밍 방식으로 테스트 사용자 생성
- 테스트 간 인증 상태 초기화
- 토큰 생성 디버그

---

## 리소스

Context7 문서 접근:

resolve-library-id에 "firebase"를 사용한 후 get-library-docs에 "authentication" 주제로 포괄적인 API 문서 확인.

Firebase 공식 리소스:

- Firebase Console: console.firebase.google.com
- Authentication Documentation: firebase.google.com/docs/auth
- Security Rules Reference: firebase.google.com/docs/rules

관련 스킬:

- do-platform-firestore: 인증 기반 보안이 있는 Firestore 데이터베이스 통합
- do-lang-flutter: 모바일 Firebase Auth용 Flutter SDK
- do-lang-typescript: Firebase SDK용 TypeScript 패턴
- do-domain-backend: Firebase Admin SDK를 사용한 백엔드 아키텍처

---

Status: Production Ready
Generated with: Do Skill Factory v1.0
Last Updated: 2025-12-07
Provider Coverage: Firebase Authentication Only
