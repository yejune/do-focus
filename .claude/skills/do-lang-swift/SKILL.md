---
name: do-lang-swift
description: Swift 6+ 개발 전문가 - SwiftUI, Combine, Swift Concurrency 및 iOS 패턴. iOS 앱, macOS 앱, Apple 플랫폼 애플리케이션 구축 시 사용.
version: 1.0.0
category: language
tags: [swift, swiftui, ios, macos, combine, concurrency]
context7-libraries: [/apple/swift, /apple/swift-evolution]
related-skills: [do-lang-kotlin, do-lang-flutter]
updated: 2025-12-07
status: active
allowed-tools: [read, grep, glob, bash, edit, write]
user-invocable: false
---

## 빠른 참조

Swift 6+ 개발 전문가 - iOS/macOS용 SwiftUI, Combine, Swift Concurrency.

자동 트리거: Swift 파일 (`.swift`), iOS/macOS 프로젝트, Xcode 워크스페이스

핵심 기능:
- Swift 6.0: Typed throws, actors, 기본 data-race 안전성
- SwiftUI 6: @Observable, NavigationStack, 현대적 선언형 UI
- Combine: Publishers와 Subscribers를 활용한 반응형 프로그래밍
- Swift Concurrency: async/await, actors, TaskGroup
- XCTest: 단위 테스트, UI 테스트, 비동기 테스트 지원
- Swift Package Manager: 의존성 관리

버전 요구사항:
- Swift: 6.0+
- Xcode: 16.0+
- iOS: 17.0+ (권장), 최소 15.0
- macOS: 14.0+ (권장)

---

## Swift 6.0 핵심 기능

### Typed Throws (에러 타입 명시)

```swift
enum NetworkError: Error {
    case invalidURL
    case requestFailed(statusCode: Int)
    case decodingFailed
}

func fetchData() throws(NetworkError) -> Data {
    guard let url = URL(string: "https://api.example.com") else {
        throw .invalidURL
    }
    // 구현
}

// 호출자는 정확한 에러 타입을 알 수 있음
do {
    let data = try fetchData()
} catch .invalidURL {
    print("Invalid URL")
} catch .requestFailed(let code) {
    print("Request failed: \(code)")
} catch .decodingFailed {
    print("Decoding failed")
}
```

### Complete Concurrency (Data-Race 안전성)

```swift
// Swift 6은 기본적으로 data-race 안전성 강제
actor UserCache {
    private var cache: [String: User] = [:]

    func get(_ id: String) -> User? { cache[id] }
    func set(_ id: String, user: User) { cache[id] = user }
    func clear() { cache.removeAll() }
}

// Actor 간 데이터 전송 시 Sendable 준수 필수
struct User: Codable, Identifiable, Sendable {
    let id: String
    let name: String
    let email: String
}

// UI 관련 코드는 MainActor 사용
@MainActor
final class UserViewModel: ObservableObject {
    @Published private(set) var user: User?
    private let cache = UserCache()

    func loadUser(_ id: String) async throws {
        if let cached = await cache.get(id) {
            self.user = cached
            return
        }
        let user = try await api.fetchUser(id)
        await cache.set(id, user: user)
        self.user = user
    }
}
```

---

## SwiftUI 6 패턴

### @Observable 매크로 (iOS 17+)

```swift
import Observation

@Observable
class ProfileViewModel {
    var user: User?
    var isLoading = false
    var errorMessage: String?

    func loadProfile() async {
        isLoading = true
        defer { isLoading = false }

        do {
            user = try await api.fetchCurrentUser()
        } catch {
            errorMessage = error.localizedDescription
        }
    }
}

struct ProfileView: View {
    @State private var viewModel = ProfileViewModel()

    var body: some View {
        NavigationStack {
            Group {
                if viewModel.isLoading {
                    ProgressView()
                } else if let user = viewModel.user {
                    UserDetailView(user: user)
                } else if let error = viewModel.errorMessage {
                    ErrorView(message: error)
                }
            }
            .task { await viewModel.loadProfile() }
            .navigationTitle("Profile")
        }
    }
}
```

### 모던 네비게이션 (NavigationStack)

```swift
@Observable
class NavigationRouter {
    var path = NavigationPath()

    func push<D: Hashable>(_ destination: D) {
        path.append(destination)
    }

    func pop() {
        guard !path.isEmpty else { return }
        path.removeLast()
    }

    func popToRoot() {
        path.removeLast(path.count)
    }
}

struct ContentView: View {
    @State private var router = NavigationRouter()

    var body: some View {
        NavigationStack(path: $router.path) {
            HomeView()
                .navigationDestination(for: User.self) { user in
                    UserDetailView(user: user)
                }
                .navigationDestination(for: Post.self) { post in
                    PostDetailView(post: post)
                }
        }
        .environment(router)
    }
}
```

---

## Swift Concurrency 패턴

### Async/Await 에러 처리

```swift
protocol UserServiceProtocol: Sendable {
    func fetchUser(_ id: String) async throws(NetworkError) -> User
    func updateUser(_ user: User) async throws(NetworkError) -> User
}

actor UserService: UserServiceProtocol {
    private let session: URLSession
    private let decoder: JSONDecoder

    init(session: URLSession = .shared) {
        self.session = session
        self.decoder = JSONDecoder()
        decoder.keyDecodingStrategy = .convertFromSnakeCase
    }

    func fetchUser(_ id: String) async throws(NetworkError) -> User {
        let url = URL(string: "https://api.example.com/users/\(id)")!

        do {
            let (data, response) = try await session.data(from: url)
            guard let httpResponse = response as? HTTPURLResponse,
                  200..<300 ~= httpResponse.statusCode else {
                throw NetworkError.requestFailed(
                    statusCode: (response as? HTTPURLResponse)?.statusCode ?? 0
                )
            }
            return try decoder.decode(User.self, from: data)
        } catch is DecodingError {
            throw NetworkError.decodingFailed
        } catch let error as NetworkError {
            throw error
        } catch {
            throw NetworkError.requestFailed(statusCode: 0)
        }
    }
}
```

### TaskGroup 병렬 실행

```swift
func loadDashboard() async throws -> Dashboard {
    // 고정 개수 작업에는 async let 사용
    async let user = api.fetchUser()
    async let posts = api.fetchPosts()
    async let notifications = api.fetchNotifications()

    return try await Dashboard(
        user: user,
        posts: posts,
        notifications: notifications
    )
}

// 동적 병렬 처리에는 TaskGroup 사용
func loadAllUsers(_ ids: [String]) async throws -> [User] {
    try await withThrowingTaskGroup(of: User.self) { group in
        for id in ids {
            group.addTask { try await api.fetchUser(id) }
        }
        return try await group.reduce(into: []) { $0.append($1) }
    }
}
```

### Actor 기반 스레드 안전 캐싱

```swift
actor ImageCache {
    private var cache: [URL: UIImage] = [:]
    private var inProgress: [URL: Task<UIImage, Error>] = [:]

    func image(for url: URL) async throws -> UIImage {
        // 캐시된 이미지 반환
        if let cached = cache[url] { return cached }

        // 이미 다운로드 중인 작업 반환
        if let task = inProgress[url] {
            return try await task.value
        }

        // 새 다운로드 작업 시작
        let task = Task { try await downloadImage(url) }
        inProgress[url] = task

        do {
            let image = try await task.value
            cache[url] = image
            inProgress[url] = nil
            return image
        } catch {
            inProgress[url] = nil
            throw error
        }
    }

    private func downloadImage(_ url: URL) async throws -> UIImage {
        let (data, _) = try await URLSession.shared.data(from: url)
        guard let image = UIImage(data: data) else {
            throw ImageError.invalidData
        }
        return image
    }
}
```

---

## Combine 프레임워크

### Publisher와 Subscriber 패턴

```swift
import Combine

class SearchViewModel: ObservableObject {
    @Published var searchText = ""
    @Published private(set) var results: [SearchResult] = []
    @Published private(set) var isSearching = false

    private var cancellables = Set<AnyCancellable>()
    private let searchService: SearchServiceProtocol

    init(searchService: SearchServiceProtocol) {
        self.searchService = searchService
        setupSearchPipeline()
    }

    private func setupSearchPipeline() {
        $searchText
            .debounce(for: .milliseconds(300), scheduler: DispatchQueue.main)
            .removeDuplicates()
            .filter { $0.count >= 2 }
            .handleEvents(receiveOutput: { [weak self] _ in
                self?.isSearching = true
            })
            .flatMap { [searchService] query in
                searchService.search(query)
                    .catch { _ in Just([]) }
            }
            .receive(on: DispatchQueue.main)
            .sink { [weak self] results in
                self?.results = results
                self?.isSearching = false
            }
            .store(in: &cancellables)
    }
}
```

---

## XCTest 단위 테스트

### MainActor 비동기 테스트

```swift
@MainActor
final class UserViewModelTests: XCTestCase {
    var sut: UserViewModel!
    var mockAPI: MockUserAPI!

    override func setUp() {
        mockAPI = MockUserAPI()
        sut = UserViewModel(api: mockAPI)
    }

    override func tearDown() {
        sut = nil
        mockAPI = nil
    }

    func testLoadUserSuccess() async throws {
        // Given
        let expectedUser = User(id: "1", name: "Test User", email: "test@example.com")
        mockAPI.mockUser = expectedUser

        // When
        try await sut.loadUser("1")

        // Then
        XCTAssertEqual(sut.user?.id, expectedUser.id)
        XCTAssertEqual(sut.user?.name, expectedUser.name)
        XCTAssertFalse(sut.isLoading)
        XCTAssertNil(sut.errorMessage)
    }

    func testLoadUserNetworkError() async {
        // Given
        mockAPI.error = NetworkError.requestFailed(statusCode: 500)

        // When
        do {
            try await sut.loadUser("1")
            XCTFail("Expected error to be thrown")
        } catch {
            // Then
            XCTAssertNil(sut.user)
            XCTAssertNotNil(sut.errorMessage)
        }
    }
}

// Mock 구현
class MockUserAPI: UserAPIProtocol {
    var mockUser: User?
    var error: Error?

    func fetchUser(_ id: String) async throws -> User {
        if let error = error { throw error }
        guard let user = mockUser else {
            throw NetworkError.requestFailed(statusCode: 404)
        }
        return user
    }
}
```

---

## Swift Package Manager

### Package.swift 설정

```swift
// swift-tools-version: 6.0
import PackageDescription

let package = Package(
    name: "MyApp",
    platforms: [
        .iOS(.v17),
        .macOS(.v14)
    ],
    products: [
        .library(name: "MyAppCore", targets: ["MyAppCore"]),
        .executable(name: "MyAppCLI", targets: ["MyAppCLI"])
    ],
    dependencies: [
        .package(url: "https://github.com/Alamofire/Alamofire.git", from: "5.9.0"),
        .package(url: "https://github.com/onevcat/Kingfisher.git", from: "7.12.0"),
        .package(url: "https://github.com/pointfreeco/swift-composable-architecture", from: "1.15.0")
    ],
    targets: [
        .target(
            name: "MyAppCore",
            dependencies: [
                "Alamofire",
                "Kingfisher",
                .product(name: "ComposableArchitecture", package: "swift-composable-architecture")
            ],
            swiftSettings: [
                .enableExperimentalFeature("StrictConcurrency")
            ]
        ),
        .testTarget(
            name: "MyAppCoreTests",
            dependencies: ["MyAppCore"]
        )
    ]
)
```

---

## Context7 라이브러리 매핑

핵심 Swift:
- `/apple/swift` - Swift 언어 및 표준 라이브러리
- `/apple/swift-evolution` - Swift 진화 제안
- `/apple/swift-package-manager` - SwiftPM 문서

인기 라이브러리:
- `/Alamofire/Alamofire` - HTTP 네트워킹
- `/onevcat/Kingfisher` - 이미지 다운로드 및 캐싱
- `/realm/realm-swift` - 모바일 데이터베이스
- `/pointfreeco/swift-composable-architecture` - TCA 아키텍처
- `/Quick/Quick` - BDD 테스팅 프레임워크
- `/Quick/Nimble` - Matcher 프레임워크

---

## 연관 스킬

- `do-lang-kotlin` - 크로스 플랫폼 프로젝트의 Android 대응
- `do-lang-flutter` - 크로스 플랫폼 모바일용 Flutter/Dart
- `do-domain-backend` - API 통합 및 백엔드 통신

---

Version: 1.0.0
Last Updated: 2025-12-07
Status: Production Ready
