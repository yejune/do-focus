---
name: do-lang-java
description: Java 21 LTS 개발 전문가로서 Spring Boot 3.3, Virtual Threads, Pattern Matching 및 엔터프라이즈 패턴을 다룬다. 엔터프라이즈 애플리케이션, 마이크로서비스, Spring 프로젝트 구축 시 사용
version: 1.0.0
category: language
tags: [java, spring-boot, jpa, hibernate, virtual-threads, enterprise]
context7-libraries: [/spring-projects/spring-boot, /spring-projects/spring-framework, /spring-projects/spring-security]
related-skills: [do-lang-kotlin, do-domain-backend]
updated: 2025-12-07
status: active
user-invocable: false
---

## 빠른 참조

Java 21 LTS 전문가 - Spring Boot 3.3, Virtual Threads, 모던 Java 기능을 활용한 엔터프라이즈 개발

자동 활성화: Java 파일(`.java`), 빌드 파일(`pom.xml`, `build.gradle`, `build.gradle.kts`)

핵심 기능:
- Java 21 LTS: Virtual Threads, Pattern Matching, Record Patterns, Sealed Classes
- Spring Boot 3.3: REST Controllers, Services, Repositories, WebFlux Reactive
- Spring Security 6: JWT 인증, OAuth2, 역할 기반 접근 제어
- JPA/Hibernate 7: Entity 매핑, 관계, 쿼리, 트랜잭션
- JUnit 5: 단위 테스트, Mocking, TestContainers 통합
- 빌드 도구: Maven 3.9, Gradle 8.5 Kotlin DSL

---

## 구현 가이드

### Java 21 LTS 기능

Virtual Threads (Project Loom):
```java
// 경량 동시 프로그래밍
try (var executor = Executors.newVirtualThreadPerTaskExecutor()) {
    IntStream.range(0, 10_000).forEach(i ->
        executor.submit(() -> {
            Thread.sleep(Duration.ofSeconds(1));
            return i;
        })
    );
}

// Structured Concurrency (preview)
try (var scope = new StructuredTaskScope.ShutdownOnFailure()) {
    Supplier<User> user = scope.fork(() -> fetchUser(userId));
    Supplier<List<Order>> orders = scope.fork(() -> fetchOrders(userId));
    scope.join().throwIfFailed();
    return new UserWithOrders(user.get(), orders.get());
}
```

Switch용 Pattern Matching:
```java
String describe(Object obj) {
    return switch (obj) {
        case Integer i when i > 0 -> "양수: " + i;
        case Integer i -> "비양수: " + i;
        case String s -> "문자열 길이 " + s.length();
        case List<?> list -> "목록 " + list.size() + "개 요소";
        case null -> "널 값";
        default -> "알 수 없는 타입";
    };
}
```

Record Patterns와 Sealed Classes:
```java
record Point(int x, int y) {}
record Rectangle(Point topLeft, Point bottomRight) {}

int area(Rectangle rect) {
    return switch (rect) {
        case Rectangle(Point(var x1, var y1), Point(var x2, var y2)) ->
            Math.abs((x2 - x1) * (y2 - y1));
    };
}

public sealed interface Shape permits Circle, Rectangle {
    double area();
}
public record Circle(double radius) implements Shape {
    public double area() { return Math.PI * radius * radius; }
}
```

### Spring Boot 3.3

REST Controller:
```java
@RestController
@RequestMapping("/api/users")
@RequiredArgsConstructor
public class UserController {
    private final UserService userService;

    @GetMapping("/{id}")
    public ResponseEntity<UserDto> getUser(@PathVariable Long id) {
        return userService.findById(id)
            .map(ResponseEntity::ok)
            .orElse(ResponseEntity.notFound().build());
    }

    @PostMapping
    public ResponseEntity<UserDto> createUser(@Valid @RequestBody CreateUserRequest request) {
        UserDto user = userService.create(request);
        URI location = URI.create("/api/users/" + user.id());
        return ResponseEntity.created(location).body(user);
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Void> deleteUser(@PathVariable Long id) {
        return userService.delete(id)
            ? ResponseEntity.noContent().build()
            : ResponseEntity.notFound().build();
    }
}
```

Service 계층:
```java
@Service
@RequiredArgsConstructor
@Transactional(readOnly = true)
public class UserService {
    private final UserRepository userRepository;
    private final PasswordEncoder passwordEncoder;

    public Optional<User> findById(Long id) {
        return userRepository.findById(id);
    }

    @Transactional
    public User create(CreateUserRequest request) {
        if (userRepository.existsByEmail(request.email())) {
            throw new DuplicateEmailException(request.email());
        }
        var user = User.builder()
            .name(request.name())
            .email(request.email())
            .passwordHash(passwordEncoder.encode(request.password()))
            .status(UserStatus.PENDING)
            .build();
        return userRepository.save(user);
    }
}
```

### Spring Security 6

보안 설정:
```java
@Configuration
@EnableWebSecurity
public class SecurityConfig {
    @Bean
    public SecurityFilterChain filterChain(HttpSecurity http) throws Exception {
        return http
            .authorizeHttpRequests(auth -> auth
                .requestMatchers("/api/public/**").permitAll()
                .requestMatchers("/api/admin/**").hasRole("ADMIN")
                .anyRequest().authenticated()
            )
            .oauth2ResourceServer(oauth2 -> oauth2.jwt(Customizer.withDefaults()))
            .sessionManagement(session ->
                session.sessionCreationPolicy(SessionCreationPolicy.STATELESS))
            .csrf(csrf -> csrf.disable())
            .build();
    }

    @Bean
    public PasswordEncoder passwordEncoder() {
        return new BCryptPasswordEncoder();
    }
}
```

### JPA/Hibernate 패턴

Entity 정의:
```java
@Entity
@Table(name = "users")
@Getter @Setter
@NoArgsConstructor @Builder
public class User {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(nullable = false)
    private String name;

    @Column(nullable = false, unique = true)
    private String email;

    @Enumerated(EnumType.STRING)
    private UserStatus status;

    @OneToMany(mappedBy = "user", cascade = CascadeType.ALL, orphanRemoval = true)
    private List<Order> orders = new ArrayList<>();
}
```

커스텀 쿼리 Repository:
```java
public interface UserRepository extends JpaRepository<User, Long> {
    Optional<User> findByEmail(String email);
    boolean existsByEmail(String email);

    @Query("SELECT u FROM User u LEFT JOIN FETCH u.orders WHERE u.id = :id")
    Optional<User> findByIdWithOrders(@Param("id") Long id);

    Page<User> findByNameContainingIgnoreCase(String name, Pageable pageable);
}
```

Record DTO:
```java
public record UserDto(Long id, String name, String email, UserStatus status) {
    public static UserDto from(User user) {
        return new UserDto(user.getId(), user.getName(), user.getEmail(), user.getStatus());
    }
}

public record CreateUserRequest(
    @NotBlank @Size(min = 2, max = 100) String name,
    @NotBlank @Email String email,
    @NotBlank @Size(min = 8) String password
) {}
```

---

## 고급 패턴

### Virtual Threads 통합

```java
@Service
@RequiredArgsConstructor
public class AsyncUserService {
    public UserWithDetails fetchUserDetails(Long userId) throws Exception {
        try (var scope = new StructuredTaskScope.ShutdownOnFailure()) {
            Supplier<User> userTask = scope.fork(() -> userRepo.findById(userId).orElseThrow());
            Supplier<List<Order>> ordersTask = scope.fork(() -> orderRepo.findByUserId(userId));
            scope.join().throwIfFailed();
            return new UserWithDetails(userTask.get(), ordersTask.get());
        }
    }

    public void processUsersInParallel(List<Long> userIds) {
        try (var executor = Executors.newVirtualThreadPerTaskExecutor()) {
            userIds.stream().map(id -> executor.submit(() -> processUser(id))).toList();
        }
    }
}
```

### 빌드 설정

Maven 3.9:
```xml
<project>
    <parent>
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-starter-parent</artifactId>
        <version>3.3.0</version>
    </parent>
    <properties><java.version>21</java.version></properties>
    <dependencies>
        <dependency><groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-web</artifactId></dependency>
        <dependency><groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-data-jpa</artifactId></dependency>
    </dependencies>
</project>
```

Gradle 8.5 (Kotlin DSL):
```kotlin
plugins {
    id("org.springframework.boot") version "3.3.0"
    id("io.spring.dependency-management") version "1.1.4"
    java
}
java { toolchain { languageVersion = JavaLanguageVersion.of(21) } }
dependencies {
    implementation("org.springframework.boot:spring-boot-starter-web")
    implementation("org.springframework.boot:spring-boot-starter-data-jpa")
    testImplementation("org.springframework.boot:spring-boot-starter-test")
}
```

### JUnit 5 테스트

단위 테스트:
```java
@ExtendWith(MockitoExtension.class)
class UserServiceTest {
    @Mock private UserRepository userRepository;
    @InjectMocks private UserService userService;

    @Test
    void shouldCreateUser() {
        when(userRepository.existsByEmail(anyString())).thenReturn(false);
        when(userRepository.save(any())).thenReturn(User.builder().id(1L).build());
        var result = userService.create(new CreateUserRequest("John", "john@test.com", "pass"));
        assertThat(result.getId()).isEqualTo(1L);
    }
}
```

TestContainers 통합 테스트:
```java
@Testcontainers
@SpringBootTest
class UserRepositoryTest {
    @Container
    static PostgreSQLContainer<?> postgres = new PostgreSQLContainer<>("postgres:16-alpine");

    @DynamicPropertySource
    static void props(DynamicPropertyRegistry registry) {
        registry.add("spring.datasource.url", postgres::getJdbcUrl);
    }

    @Autowired UserRepository repo;

    @Test
    void shouldSaveUser() {
        var user = repo.save(User.builder().name("John").email("john@test.com").build());
        assertThat(user.getId()).isNotNull();
    }
}
```

---

## Context7 통합

최신 문서용 라이브러리 매핑:
- `/spring-projects/spring-boot` - Spring Boot 3.3 문서
- `/spring-projects/spring-framework` - Spring Framework 코어
- `/spring-projects/spring-security` - Spring Security 6
- `/hibernate/hibernate-orm` - Hibernate 7 ORM 패턴
- `/junit-team/junit5` - JUnit 5 테스트 프레임워크

---

## 연관 스킬

- `do-lang-kotlin` - Kotlin 상호운용성, Spring Kotlin 확장
- `do-domain-backend` - REST API, GraphQL, 마이크로서비스 아키텍처
- `do-domain-database` - JPA, Hibernate, R2DBC 패턴
- `do-foundation-quality` - JUnit 5, Mockito, TestContainers 통합
- `do-infra-docker` - JVM 컨테이너 최적화

---

## 문제 해결

일반적인 이슈:
- 버전 불일치: `java -version` 확인, `JAVA_HOME`이 Java 21 가리키는지 확인
- 컴파일 오류: `mvn clean compile -X` 또는 `gradle build --info` 실행
- Virtual Thread 이슈: Java 21+ 확인, 필요시 `--enable-preview` 사용
- JPA Lazy Loading: `@Transactional` 또는 `JOIN FETCH` 쿼리 사용

성능 팁:
- Virtual Threads 활성화: `spring.threads.virtual.enabled=true`
- 빠른 시작을 위해 GraalVM Native Image 사용
- HikariCP로 커넥션 풀링 설정

---

## 고급 문서

상세 참조 자료:
- [reference.md](reference.md) - Java 21 기능, Context7 매핑, 성능
- [examples.md](examples.md) - 프로덕션 준비된 Spring Boot 예제

---

최종 업데이트: 2025-12-07
상태: Production Ready (v1.0.0)
