---
name: do-lang-flutter
description: Flutter 3.24+ / Dart 3.5+ 개발 전문가로 Riverpod, go_router, 크로스플랫폼 패턴을 다룬다. 크로스플랫폼 모바일 앱, 데스크톱 앱, 웹 애플리케이션을 Flutter로 구축할 때 사용.
version: 1.0.0
category: language
tags: [flutter, dart, riverpod, cross-platform, mobile, desktop]
context7-libraries: [/flutter/flutter, /rrousselgit/riverpod, /flutter/packages]
related-skills: [do-lang-swift, do-lang-kotlin]
updated: 2025-12-07
status: active
user-invocable: false
---

## 빠른 참조 (30초)

Flutter/Dart 개발 전문가 - Dart 3.5+, Flutter 3.24+ 최신 패턴.

자동 트리거: Flutter 프로젝트 (`.dart` 파일, `pubspec.yaml`), 크로스플랫폼 앱, 위젯 개발

핵심 역량:
- Dart 3.5: 패턴 매칭, 레코드, sealed 클래스, extension types
- Flutter 3.24: 위젯 트리, Material 3, 적응형 레이아웃
- Riverpod: 코드 생성 기반 상태 관리
- go_router: 선언적 네비게이션 및 딥 링킹
- Platform Channels: 네이티브 iOS/Android 통합
- 테스팅: flutter_test, widget_test, integration_test

## 구현 가이드

### Dart 3.5 언어 기능

Sealed 클래스와 패턴 매칭:
```dart
sealed class Result<T> {
  const Result();
}
class Success<T> extends Result<T> {
  final T data;
  const Success(this.data);
}
class Failure<T> extends Result<T> {
  final String error;
  const Failure(this.error);
}

// 완전한 switch 표현식
String handleResult(Result<User> result) => switch (result) {
  Success(:final data) => 'User: ${data.name}',
  Failure(:final error) => 'Error: $error',
};

// 패턴 가드절
String describeUser(User user) => switch (user) {
  User(age: var a) when a < 18 => 'Minor',
  User(age: var a) when a >= 65 => 'Senior',
  User(name: var n, age: var a) => '$n, age $a',
};
```

레코드와 구조 분해:
```dart
typedef UserRecord = ({String name, int age, String email});

// 다중 반환값
(String name, int age) parseUser(Map<String, dynamic> json) {
  return (json['name'] as String, json['age'] as int);
}

// 구조 분해
void processUser(Map<String, dynamic> json) {
  final (name, age) = parseUser(json);
  print('$name is $age years old');
}

// 컬렉션 레코드 패턴
void processUsers(List<UserRecord> users) {
  for (final (:name, :age, :email) in users) {
    print('$name ($age): $email');
  }
}
```

Extension Types:
```dart
extension type UserId(String value) {
  factory UserId.generate() => UserId(Uuid().v4());
  bool get isValid => value.isNotEmpty;
}

extension type Email(String value) {
  bool get isValid => value.contains('@') && value.contains('.');
  String get domain => value.split('@').last;
}
```

### Riverpod 상태 관리

Provider 정의:
```dart
import 'package:riverpod_annotation/riverpod_annotation.dart';
part 'providers.g.dart';

@riverpod
UserRepository userRepository(Ref ref) {
  return UserRepository(ref.read(dioProvider));
}

@riverpod
Future<User> user(Ref ref, String userId) async {
  return ref.watch(userRepositoryProvider).getUser(userId);
}

@riverpod
class UserNotifier extends _$UserNotifier {
  @override
  FutureOr<User?> build() => null;

  Future<void> loadUser(String id) async {
    state = const AsyncLoading();
    state = await AsyncValue.guard(
      () => ref.read(userRepositoryProvider).getUser(id),
    );
  }

  Future<void> updateUser(User user) async {
    state = const AsyncLoading();
    state = await AsyncValue.guard(
      () => ref.read(userRepositoryProvider).updateUser(user),
    );
  }
}

@riverpod
Stream<List<Message>> messages(Ref ref, String chatId) {
  return ref.watch(chatRepositoryProvider).watchMessages(chatId);
}
```

Widget 통합:
```dart
class UserScreen extends ConsumerWidget {
  final String userId;
  const UserScreen({required this.userId, super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final userAsync = ref.watch(userProvider(userId));

    return Scaffold(
      appBar: AppBar(title: const Text('User Profile')),
      body: userAsync.when(
        data: (user) => UserProfile(user: user),
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (error, stack) => Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Text('Error: $error'),
              ElevatedButton(
                onPressed: () => ref.invalidate(userProvider(userId)),
                child: const Text('Retry'),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
```

ConsumerStatefulWidget:
```dart
class EditUserScreen extends ConsumerStatefulWidget {
  const EditUserScreen({super.key});
  @override
  ConsumerState<EditUserScreen> createState() => _EditUserScreenState();
}

class _EditUserScreenState extends ConsumerState<EditUserScreen> {
  final _formKey = GlobalKey<FormState>();
  late TextEditingController _nameController;

  @override
  void initState() {
    super.initState();
    _nameController = TextEditingController();
  }

  @override
  void dispose() {
    _nameController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    ref.listen(userNotifierProvider, (prev, next) {
      next.whenOrNull(
        error: (e, _) => ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Error: $e')),
        ),
      );
    });

    final isLoading = ref.watch(userNotifierProvider).isLoading;

    return Form(
      key: _formKey,
      child: Column(
        children: [
          TextFormField(
            controller: _nameController,
            validator: (v) => v?.isEmpty ?? true ? 'Required' : null,
          ),
          ElevatedButton(
            onPressed: isLoading ? null : () async {
              if (_formKey.currentState!.validate()) {
                await ref.read(userNotifierProvider.notifier)
                    .updateUser(User(name: _nameController.text));
              }
            },
            child: isLoading
                ? const CircularProgressIndicator()
                : const Text('Save'),
          ),
        ],
      ),
    );
  }
}
```

### go_router 네비게이션

Router 구성:
```dart
final router = GoRouter(
  initialLocation: '/',
  routes: [
    GoRoute(
      path: '/',
      name: 'home',
      builder: (context, state) => const HomeScreen(),
      routes: [
        GoRoute(
          path: 'user/:id',
          name: 'user-detail',
          builder: (context, state) => UserDetailScreen(
            userId: state.pathParameters['id']!,
          ),
        ),
      ],
    ),
    ShellRoute(
      builder: (context, state, child) => MainShell(child: child),
      routes: [
        GoRoute(
          path: '/feed',
          pageBuilder: (_, __) => const NoTransitionPage(child: FeedScreen()),
        ),
        GoRoute(
          path: '/search',
          pageBuilder: (_, __) => const NoTransitionPage(child: SearchScreen()),
        ),
        GoRoute(
          path: '/profile',
          pageBuilder: (_, __) => const NoTransitionPage(child: ProfileScreen()),
        ),
      ],
    ),
  ],
  redirect: (context, state) {
    final isLoggedIn = authNotifier.isLoggedIn;
    final isLoggingIn = state.matchedLocation == '/login';
    if (!isLoggedIn && !isLoggingIn) return '/login';
    if (isLoggedIn && isLoggingIn) return '/';
    return null;
  },
  errorBuilder: (context, state) => ErrorScreen(error: state.error),
);

// 네비게이션 메서드
void navigateToUser(BuildContext context, String userId) {
  context.go('/user/$userId');
}

void goBack(BuildContext context) {
  if (context.canPop()) context.pop();
  else context.go('/');
}
```

### Platform Channels

Dart 구현:
```dart
class NativeBridge {
  static const _channel = MethodChannel('com.example.app/native');
  static const _eventChannel = EventChannel('com.example.app/events');

  Future<String> getPlatformVersion() async {
    try {
      final version = await _channel.invokeMethod<String>('getPlatformVersion');
      return version ?? 'Unknown';
    } on PlatformException catch (e) {
      throw NativeBridgeException('Failed: ${e.message}');
    }
  }

  Future<void> shareContent({required String text, String? title}) async {
    await _channel.invokeMethod('share', {
      'text': text,
      if (title != null) 'title': title,
    });
  }

  Stream<BatteryState> watchBatteryState() {
    return _eventChannel.receiveBroadcastStream().map((event) {
      final data = event as Map<dynamic, dynamic>;
      return BatteryState(
        level: data['level'] as int,
        isCharging: data['isCharging'] as bool,
      );
    });
  }

  void setupMethodCallHandler() {
    _channel.setMethodCallHandler((call) async {
      switch (call.method) {
        case 'onNativeEvent':
          return true;
        default:
          throw MissingPluginException('Not implemented: ${call.method}');
      }
    });
  }
}

class BatteryState {
  final int level;
  final bool isCharging;
  const BatteryState({required this.level, required this.isCharging});
}
```

### Widget 패턴

적응형 레이아웃:
```dart
class AdaptiveScaffold extends StatelessWidget {
  final Widget child;
  final List<NavigationDestination> destinations;
  final int selectedIndex;
  final ValueChanged<int> onDestinationSelected;

  const AdaptiveScaffold({
    required this.child,
    required this.destinations,
    required this.selectedIndex,
    required this.onDestinationSelected,
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    final width = MediaQuery.sizeOf(context).width;

    if (width < 600) {
      return Scaffold(
        body: child,
        bottomNavigationBar: NavigationBar(
          selectedIndex: selectedIndex,
          onDestinationSelected: onDestinationSelected,
          destinations: destinations,
        ),
      );
    }

    if (width < 840) {
      return Scaffold(
        body: Row(children: [
          NavigationRail(
            selectedIndex: selectedIndex,
            onDestinationSelected: onDestinationSelected,
            destinations: destinations.map((d) => NavigationRailDestination(
              icon: d.icon, selectedIcon: d.selectedIcon, label: Text(d.label),
            )).toList(),
          ),
          const VerticalDivider(thickness: 1, width: 1),
          Expanded(child: child),
        ]),
      );
    }

    return Scaffold(
      body: Row(children: [
        NavigationDrawer(
          selectedIndex: selectedIndex,
          onDestinationSelected: onDestinationSelected,
          children: destinations.map((d) => NavigationDrawerDestination(
            icon: d.icon, selectedIcon: d.selectedIcon ?? d.icon, label: Text(d.label),
          )).toList(),
        ),
        const VerticalDivider(thickness: 1, width: 1),
        Expanded(child: child),
      ]),
    );
  }
}
```

### 테스팅

Widget 테스트 예제:
```dart
void main() {
  testWidgets('UserScreen displays data', (tester) async {
    final container = ProviderContainer(overrides: [
      userRepositoryProvider.overrideWithValue(MockUserRepository()),
    ]);

    await tester.pumpWidget(
      UncontrolledProviderScope(
        container: container,
        child: const MaterialApp(home: UserScreen(userId: '1')),
      ),
    );

    expect(find.byType(CircularProgressIndicator), findsOneWidget);
    await tester.pumpAndSettle();
    expect(find.text('Test User'), findsOneWidget);
  });
}
```

포괄적인 테스팅 패턴은 [examples.md](examples.md) 참조.

## 고급 패턴

다음 주제에 대한 포괄적인 가이드:
- Riverpod과 Clean Architecture
- 계산량 많은 작업의 Isolates 처리
- 커스텀 렌더 객체와 페인팅
- FFI와 플랫폼별 플러그인
- 성능 최적화 및 프로파일링

자세한 내용: [reference.md](reference.md) 및 [examples.md](examples.md)

## Context7 라이브러리 매핑

Flutter/Dart 코어:
- `/flutter/flutter` - Flutter 프레임워크
- `/dart-lang/sdk` - Dart SDK

상태 관리:
- `/rrousselGit/riverpod` - Riverpod 상태 관리
- `/felangel/bloc` - BLoC 패턴

네비게이션 및 스토리지:
- `/flutter/packages` - go_router 및 공식 패키지
- `/cfug/dio` - HTTP 클라이언트
- `/isar/isar` - NoSQL 데이터베이스

## 연관 스킬

- `do-lang-swift` - Platform Channels용 iOS 네이티브 통합
- `do-lang-kotlin` - Platform Channels용 Android 네이티브 통합
- `do-domain-backend` - API 통합 및 백엔드 통신
- `do-foundation-quality` - 모바일 보안 베스트 프랙티스

---

Version: 1.0.0
Last Updated: 2025-12-07
Status: Production Ready
