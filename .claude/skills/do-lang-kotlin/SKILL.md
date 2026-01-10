---
name: do-lang-kotlin
description: Kotlin 2.0+ 개발 전문가 - Ktor, coroutines, Compose Multiplatform, Kotlin 관용 패턴
version: 1.0.0
category: language
tags: kotlin, ktor, coroutines, compose, android, multiplatform
context7-libraries: /ktorio/ktor, /jetbrains/compose-multiplatform, /jetbrains/exposed
related-skills: do-lang-java, do-lang-swift
updated: 2025-12-07
status: active
allowed-tools:
  - Read
  - Grep
  - Glob
  - Bash
  - Edit
  - Write
user-invocable: false
---

## 빠른 참조

Kotlin 2.0+ 전문가 - K2 compiler, coroutines, Ktor, Compose Multiplatform, Context7 통합

자동 트리거: Kotlin 파일 (`.kt`, `.kts`), Gradle Kotlin DSL (`build.gradle.kts`, `settings.gradle.kts`)

핵심 기능:
- Kotlin 2.0: K2 compiler, coroutines, Flow, sealed classes, value classes
- Ktor 3.0: 비동기 HTTP 서버/클라이언트, WebSocket, JWT 인증
- Exposed 0.55: Coroutines 지원 Kotlin SQL 프레임워크
- Spring Boot (Kotlin): WebFlux 기반 Kotlin 관용적 Spring
- Compose Multiplatform: Desktop, iOS, Web, Android UI 공유
- Testing: JUnit 5, MockK, Kotest, Turbine (Flow 테스트)

---

## 구현 가이드

### Kotlin 2.0 기능

Coroutines와 Flow:
```kotlin
// 구조화된 동시성 - async/await
suspend fun fetchUserWithOrders(userId: Long): UserWithOrders = coroutineScope {
    val userDeferred = async { userRepository.findById(userId) }
    val ordersDeferred = async { orderRepository.findByUserId(userId) }
    UserWithOrders(userDeferred.await(), ordersDeferred.await())
}

// Flow 기반 reactive 스트림
fun observeUsers(): Flow<User> = flow {
    while (true) {
        emit(userRepository.findLatest())
        delay(1000)
    }
}.flowOn(Dispatchers.IO)
```

Sealed Classes와 Value Classes:
```kotlin
sealed interface Result<out T> {
    data class Success<T>(val data: T) : Result<T>
    data class Error(val exception: Throwable) : Result<Nothing>
    data object Loading : Result<Nothing>
}

@JvmInline
value class UserId(val value: Long) {
    init { require(value > 0) { "UserId must be positive" } }
}

@JvmInline
value class Email(val value: String) {
    init { require(value.contains("@")) { "Invalid email format" } }
}
```

### Ktor 3.0 서버

애플리케이션 설정:
```kotlin
fun main() {
    embeddedServer(Netty, port = 8080, host = "0.0.0.0") {
        configureKoin()
        configureSecurity()
        configureRouting()
        configureContentNegotiation()
    }.start(wait = true)
}

fun Application.configureKoin() {
    install(Koin) { modules(appModule) }
}

val appModule = module {
    single<Database> { DatabaseFactory.create() }
    single<UserRepository> { UserRepositoryImpl(get()) }
    single<UserService> { UserServiceImpl(get()) }
}

fun Application.configureSecurity() {
    install(Authentication) {
        jwt("auth-jwt") {
            realm = "User API"
            verifier(JwtConfig.verifier)
            validate { credential ->
                if (credential.payload.audience.contains("api"))
                    JWTPrincipal(credential.payload) else null
            }
        }
    }
}

fun Application.configureContentNegotiation() {
    install(ContentNegotiation) {
        json(Json { prettyPrint = true; ignoreUnknownKeys = true })
    }
}
```

라우팅과 인증:
```kotlin
fun Application.configureRouting() {
    val userService by inject<UserService>()

    routing {
        route("/api/v1") {
            post("/auth/register") {
                val request = call.receive<CreateUserRequest>()
                val user = userService.create(request)
                call.respond(HttpStatusCode.Created, user.toDto())
            }

            authenticate("auth-jwt") {
                route("/users") {
                    get {
                        val page = call.parameters["page"]?.toIntOrNull() ?: 0
                        val size = call.parameters["size"]?.toIntOrNull() ?: 20
                        call.respond(userService.findAll(page, size).map { it.toDto() })
                    }

                    get("/{id}") {
                        val id = call.parameters["id"]?.toLongOrNull()
                            ?: return@get call.respond(HttpStatusCode.BadRequest)
                        userService.findById(id)?.let { call.respond(it.toDto()) }
                            ?: call.respond(HttpStatusCode.NotFound)
                    }
                }
            }
        }
    }
}
```

### Exposed SQL 프레임워크

테이블과 엔티티:
```kotlin
object Users : LongIdTable("users") {
    val name = varchar("name", 100)
    val email = varchar("email", 255).uniqueIndex()
    val passwordHash = varchar("password_hash", 255)
    val status = enumerationByName<UserStatus>("status", 20)
    val createdAt = timestamp("created_at").defaultExpression(CurrentTimestamp())
}

class UserEntity(id: EntityID<Long>) : LongEntity(id) {
    companion object : LongEntityClass<UserEntity>(Users)
    var name by Users.name
    var email by Users.email
    var passwordHash by Users.passwordHash
    var status by Users.status
    var createdAt by Users.createdAt

    fun toModel() = User(id.value, name, email, passwordHash, status, createdAt)
}
```

Coroutines 지원 Repository:
```kotlin
class UserRepositoryImpl(private val database: Database) : UserRepository {
    override suspend fun findById(id: Long): User? = dbQuery {
        UserEntity.findById(id)?.toModel()
    }

    override suspend fun save(user: User): User = dbQuery {
        UserEntity.new {
            name = user.name
            email = user.email
            passwordHash = user.passwordHash
            status = user.status
        }.toModel()
    }

    private suspend fun <T> dbQuery(block: suspend () -> T): T =
        newSuspendedTransaction(Dispatchers.IO, database) { block() }
}
```

### Spring Boot와 Kotlin

WebFlux Controller:
```kotlin
@RestController
@RequestMapping("/api/users")
class UserController(private val userService: UserService) {

    @GetMapping
    suspend fun listUsers(
        @RequestParam(defaultValue = "0") page: Int,
        @RequestParam(defaultValue = "20") size: Int
    ): Flow<UserDto> = userService.findAll(page, size).map { it.toDto() }

    @GetMapping("/{id}")
    suspend fun getUser(@PathVariable id: Long): ResponseEntity<UserDto> =
        userService.findById(id)?.let { ResponseEntity.ok(it.toDto()) }
            ?: ResponseEntity.notFound().build()

    @PostMapping
    suspend fun createUser(@Valid @RequestBody request: CreateUserRequest): ResponseEntity<UserDto> {
        val user = userService.create(request)
        return ResponseEntity.created(URI.create("/api/users/${user.id}")).body(user.toDto())
    }
}
```

---

## 고급 패턴

### Compose Multiplatform

공유 UI 컴포넌트:
```kotlin
@Composable
fun UserListScreen(viewModel: UserListViewModel, onUserClick: (Long) -> Unit) {
    val uiState by viewModel.uiState.collectAsState()

    when (val state = uiState) {
        is UiState.Loading -> LoadingIndicator()
        is UiState.Success -> UserList(state.users, onUserClick)
        is UiState.Error -> ErrorMessage(state.message, onRetry = viewModel::retry)
    }
}

@Composable
fun UserCard(user: User, onClick: () -> Unit) {
    Card(
        modifier = Modifier.fillMaxWidth().clickable(onClick = onClick),
        elevation = CardDefaults.cardElevation(4.dp)
    ) {
        Row(modifier = Modifier.padding(16.dp)) {
            AsyncImage(model = user.avatarUrl, contentDescription = user.name,
                modifier = Modifier.size(48.dp).clip(CircleShape))
            Spacer(Modifier.width(16.dp))
            Column {
                Text(user.name, style = MaterialTheme.typography.titleMedium)
                Text(user.email, style = MaterialTheme.typography.bodySmall)
            }
        }
    }
}
```

### MockK 테스팅

```kotlin
class UserServiceTest {
    private val userRepository = mockk<UserRepository>()
    private val userService = UserService(userRepository)

    @Test
    fun `should fetch user concurrently`() = runTest {
        val testUser = User(1L, "John", "john@example.com")
        coEvery { userRepository.findById(1L) } coAnswers { delay(100); testUser }

        val result = userService.findById(1L)
        assertThat(result).isEqualTo(testUser)
    }

    @Test
    fun `should handle Flow emissions`() = runTest {
        val users = listOf(User(1L, "John", "john@example.com"))
        coEvery { userRepository.findAllAsFlow() } returns users.asFlow()

        userService.streamUsers().toList().also { result ->
            assertThat(result).hasSize(1)
        }
    }
}
```

### Gradle 빌드 설정

```kotlin
plugins {
    kotlin("jvm") version "2.0.20"
    kotlin("plugin.serialization") version "2.0.20"
    id("io.ktor.plugin") version "3.0.0"
}

kotlin { jvmToolchain(21) }

dependencies {
    implementation("io.ktor:ktor-server-core-jvm")
    implementation("io.ktor:ktor-server-netty-jvm")
    implementation("io.ktor:ktor-server-content-negotiation-jvm")
    implementation("io.ktor:ktor-server-auth-jwt-jvm")
    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-core:1.9.0")
    implementation("org.jetbrains.exposed:exposed-core:0.55.0")
    implementation("org.jetbrains.exposed:exposed-dao:0.55.0")
    implementation("org.postgresql:postgresql:42.7.3")

    testImplementation("io.mockk:mockk:1.13.12")
    testImplementation("org.jetbrains.kotlinx:kotlinx-coroutines-test:1.9.0")
    testImplementation("app.cash.turbine:turbine:1.1.0")
}
```

---

## Context7 통합

최신 문서 라이브러리 매핑:
- `/ktorio/ktor` - Ktor 3.0 서버/클라이언트 문서
- `/jetbrains/exposed` - Exposed SQL 프레임워크
- `/JetBrains/kotlin` - Kotlin 2.0 언어 레퍼런스
- `/Kotlin/kotlinx.coroutines` - Coroutines 라이브러리
- `/jetbrains/compose-multiplatform` - Compose Multiplatform
- `/arrow-kt/arrow` - Arrow 함수형 프로그래밍

---

## 사용 시나리오

**Kotlin 사용**:
- Android 애플리케이션 개발 (공식 언어)
- Ktor 기반 현대적 서버 애플리케이션
- 간결하고 표현력 있는 문법 선호 시
- Coroutines와 Flow 기반 반응형 서비스
- 멀티플랫폼 앱 (iOS, Desktop, Web)
- Java 상호운용성 필요 시

**대안 검토**:
- 레거시 Java 코드베이스 - 최소 변경 필요 시
- 빅데이터 파이프라인 - Scala with Spark 권장

---

## 연관 스킬

- `do-lang-java` - Java 상호운용성, Spring Boot 패턴
- `do-domain-backend` - REST API, GraphQL, 마이크로서비스 아키텍처
- `do-domain-database` - JPA, Exposed, R2DBC 패턴
- `do-workflow-testing` - JUnit 5, MockK, TestContainers 통합
- `do-infra-docker` - JVM 컨테이너 최적화

---

## 문제 해결

K2 Compiler: gradle.properties에 `kotlin.experimental.tryK2=true` 추가, `.gradle` 삭제 후 리빌드

Coroutines: suspend 컨텍스트에서 `runBlocking` 사용 금지, 블로킹 작업은 `Dispatchers.IO` 사용

Ktor: `ContentNegotiation` 설치 확인, JWT verifier 설정 검증, 라우팅 계층 구조 확인

Exposed: 모든 DB 작업은 트랜잭션 내에서 실행, 트랜잭션 외부의 엔티티 로딩 주의

---

Last Updated: 2025-12-07
Status: Production Ready (v1.0.0)
