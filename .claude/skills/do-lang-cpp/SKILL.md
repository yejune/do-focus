---
name: do-lang-cpp
description: Modern C++ (C++23/C++20) development specialist covering RAII, smart pointers, concepts, ranges, modules, and CMake. Use when developing high-performance applications, games, system software, or embedded systems.
version: 1.0.0
updated: 2026-01-06
status: active
allowed-tools: Read, Grep, Glob, Bash, mcp__context7__resolve-library-id, mcp__context7__query-docs
user-invocable: false
---

## 빠른 참조

Modern C++ (C++23/C++20) 개발 전문가 - RAII, 스마트 포인터, 개념, 범위, 모듈, CMake

자동 트리거: `.cpp`, `.hpp`, `.h`, `CMakeLists.txt`, `vcpkg.json`, `conanfile.txt`

핵심 기능:
- C++23: std::expected, std::print, deducing this
- C++20: Concepts, Ranges, Modules, std::format
- 메모리 안전성: RAII, Rule of 5, 스마트 포인터
- 빌드 시스템: CMake 3.28+, FetchContent
- 동시성: std::jthread, std::async, std::latch/barrier
- 테스트: Google Test, Catch2
- 패키지: vcpkg, Conan 2.0

### 핵심 패턴

스마트 포인터 팩토리:
```cpp
class Widget {
public:
    static auto create(int value) -> std::unique_ptr<Widget> {
        return std::make_unique<Widget>(value);
    }
    explicit Widget(int v) : value_(v) {}
private:
    int value_;
};
```

Concepts 제약:
```cpp
template<typename T>
concept Numeric = std::integral<T> || std::floating_point<T>;

template<Numeric T>
auto square(T value) -> T { return value * value; }
```

Ranges 파이프라인:
```cpp
auto result = std::views::iota(1, 100)
    | std::views::filter([](int n) { return n % 2 == 0; })
    | std::views::transform([](int n) { return n * n; })
    | std::views::take(10);
```

---

## 구현 가이드

### C++23 신기능

std::expected 에러 처리:
```cpp
enum class ParseError { InvalidFormat, OutOfRange };

auto parse_int(std::string_view str) -> std::expected<int, ParseError> {
    try {
        return std::stoi(std::string(str));
    } catch (const std::invalid_argument&) {
        return std::unexpected(ParseError::InvalidFormat);
    } catch (const std::out_of_range&) {
        return std::unexpected(ParseError::OutOfRange);
    }
}

auto result = parse_int("42");
if (result) std::println("Value: {}", *result);
```

Deducing This (명시적 객체 매개변수):
```cpp
class Builder {
    std::string data_;
public:
    template<typename Self>
    auto append(this Self&& self, std::string_view s) -> Self&& {
        self.data_ += s;
        return std::forward<Self>(self);
    }
    auto build() const -> std::string { return data_; }
};
```

### C++20 기능

Concepts와 제약:
```cpp
template<typename T>
concept Hashable = requires(T a) {
    { std::hash<T>{}(a) } -> std::convertible_to<std::size_t>;
};

template<typename Container>
    requires std::ranges::range<Container>
auto make_default_filled(std::size_t count) -> Container {
    Container c;
    c.resize(count);
    return c;
}
```

모듈 (C++20):
```cpp
// math.cppm
export module math;

export namespace math {
    template<typename T>
    constexpr auto add(T a, T b) -> T { return a + b; }
}

// main.cpp
import math;
int main() {
    auto result = math::add(10, 20);
}
```

Ranges 라이브러리:
```cpp
struct Person { std::string name; int age; };

void process_people(std::vector<Person>& people) {
    auto adults = people
        | std::views::filter([](const Person& p) { return p.age >= 18; })
        | std::views::transform([](const Person& p) { return p.name; });

    for (const auto& name : adults) std::println("{}", name);
    std::ranges::sort(people, {}, &Person::age);
}
```

### RAII와 리소스 관리

Rule of Five:
```cpp
class Resource {
    int* data_;
    std::size_t size_;

public:
    explicit Resource(std::size_t size) : data_(new int[size]), size_(size) {}
    ~Resource() { delete[] data_; }

    Resource(const Resource& other)
        : data_(new int[other.size_]), size_(other.size_) {
        std::copy(other.data_, other.data_ + size_, data_);
    }

    auto operator=(const Resource& other) -> Resource& {
        if (this != &other) { Resource temp(other); swap(temp); }
        return *this;
    }

    Resource(Resource&& other) noexcept
        : data_(std::exchange(other.data_, nullptr))
        , size_(std::exchange(other.size_, 0)) {}

    auto operator=(Resource&& other) noexcept -> Resource& {
        if (this != &other) {
            delete[] data_;
            data_ = std::exchange(other.data_, nullptr);
            size_ = std::exchange(other.size_, 0);
        }
        return *this;
    }

    void swap(Resource& other) noexcept {
        std::swap(data_, other.data_);
        std::swap(size_, other.size_);
    }
};
```

스마트 포인터 패턴:
```cpp
class Node : public std::enable_shared_from_this<Node> {
public:
    std::vector<std::shared_ptr<Node>> children;
    std::weak_ptr<Node> parent;  // 순환 참조 방지

    void add_child(std::shared_ptr<Node> child) {
        child->parent = weak_from_this();
        children.push_back(std::move(child));
    }
};
```

### CMake 모던 패턴

CMakeLists.txt (C++23):
```cmake
cmake_minimum_required(VERSION 3.28)
project(MyProject VERSION 1.0.0 LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 23)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

add_compile_options(
    $<$<CXX_COMPILER_ID:GNU,Clang>:-Wall -Wextra -Wpedantic>
    $<$<CXX_COMPILER_ID:MSVC>:/W4>
)

include(FetchContent)
FetchContent_Declare(fmt GIT_REPOSITORY https://github.com/fmtlib/fmt GIT_TAG 10.2.1)
FetchContent_Declare(googletest GIT_REPOSITORY https://github.com/google/googletest GIT_TAG v1.14.0)
FetchContent_MakeAvailable(fmt googletest)

add_library(mylib STATIC src/core.cpp)
target_include_directories(mylib PUBLIC include)
target_link_libraries(mylib PUBLIC fmt::fmt)

add_executable(myapp src/main.cpp)
target_link_libraries(myapp PRIVATE mylib)

enable_testing()
add_executable(mylib_tests tests/core_test.cpp)
target_link_libraries(mylib_tests PRIVATE mylib GTest::gtest_main)
include(GoogleTest)
gtest_discover_tests(mylib_tests)
```

### 동시성

std::jthread와 Stop Token:
```cpp
void worker(std::stop_token stoken) {
    while (!stoken.stop_requested()) {
        std::this_thread::sleep_for(std::chrono::milliseconds(100));
    }
}

int main() {
    std::jthread worker_thread(worker);
    std::this_thread::sleep_for(std::chrono::seconds(1));
    worker_thread.request_stop();
}
```

동기화 프리미티브:
```cpp
void parallel_init(std::latch& ready, int id) { ready.count_down(); }

void parallel_compute(std::barrier<>& sync, int id) {
    for (int i = 0; i < 10; ++i) sync.arrive_and_wait();
}

std::counting_semaphore<4> pool(4);
void limited_resource() { pool.acquire(); pool.release(); }
```

---

## 고급 패턴

### 템플릿 메타프로그래밍

가변 템플릿:
```cpp
template<typename... Args>
auto sum(Args... args) { return (args + ...); }

template<typename... Args>
void print_all(Args&&... args) {
    ((std::cout << std::forward<Args>(args) << " "), ...);
}
```

SFINAE와 if constexpr:
```cpp
template<typename T>
auto to_string(const T& value) -> std::string {
    if constexpr (std::is_arithmetic_v<T>) return std::to_string(value);
    else if constexpr (requires { value.to_string(); }) return value.to_string();
    else return "unknown";
}
```

### 테스트: Google Test

```cpp
class CalculatorTest : public ::testing::Test {
protected:
    Calculator calc;
    void SetUp() override { calc = Calculator{}; }
};

TEST_F(CalculatorTest, Addition) { EXPECT_EQ(calc.add(2, 3), 5); }
TEST_F(CalculatorTest, DivisionByZero) {
    EXPECT_THROW(calc.divide(1, 0), std::invalid_argument);
}

// 매개변수화 테스트
class AdditionTest : public ::testing::TestWithParam<std::tuple<int, int, int>> {};

TEST_P(AdditionTest, Works) {
    auto [a, b, expected] = GetParam();
    EXPECT_EQ(Calculator{}.add(a, b), expected);
}

INSTANTIATE_TEST_SUITE_P(Basics, AdditionTest,
    ::testing::Values(std::make_tuple(1, 1, 2), std::make_tuple(0, 0, 0)));
```

### 스레드 풀 구현

```cpp
class ThreadPool {
public:
    explicit ThreadPool(size_t num_threads) : stop_(false) {
        for (size_t i = 0; i < num_threads; ++i) {
            workers_.emplace_back([this] {
                while (true) {
                    std::function<void()> task;
                    {
                        std::unique_lock lock(mutex_);
                        cv_.wait(lock, [this] { return stop_ || !tasks_.empty(); });
                        if (stop_ && tasks_.empty()) return;
                        task = std::move(tasks_.front());
                        tasks_.pop();
                    }
                    task();
                }
            });
        }
    }

    template<typename F, typename... Args>
    auto enqueue(F&& f, Args&&... args) -> std::future<std::invoke_result_t<F, Args...>> {
        using return_type = std::invoke_result_t<F, Args...>;
        auto task = std::make_shared<std::packaged_task<return_type()>>(
            std::bind(std::forward<F>(f), std::forward<Args>(args)...));
        auto result = task->get_future();
        { std::unique_lock lock(mutex_); tasks_.emplace([task] { (*task)(); }); }
        cv_.notify_one();
        return result;
    }

    ~ThreadPool() {
        { std::unique_lock lock(mutex_); stop_ = true; }
        cv_.notify_all();
        for (auto& worker : workers_) worker.join();
    }

private:
    std::vector<std::thread> workers_;
    std::queue<std::function<void()>> tasks_;
    std::mutex mutex_;
    std::condition_variable cv_;
    bool stop_;
};
```

### 의존성 주입 패턴

```cpp
class ILogger {
public:
    virtual ~ILogger() = default;
    virtual void log(std::string_view message) = 0;
};

class ConsoleLogger : public ILogger {
public:
    void log(std::string_view message) override { std::println("{}", message); }
};

class UserService {
    std::shared_ptr<ILogger> logger_;
public:
    explicit UserService(std::shared_ptr<ILogger> logger) : logger_(std::move(logger)) {}
    void create_user(std::string_view name) {
        logger_->log(std::format("Creating user: {}", name));
    }
};
```

### 빌드 시스템 패턴

vcpkg 매니페스트:
```json
{
    "name": "myproject",
    "version": "1.0.0",
    "dependencies": ["fmt", "nlohmann-json", "spdlog", { "name": "gtest", "features": ["gmock"] }]
}
```

---

## Context7 라이브러리 매핑

```
/microsoft/vcpkg - 패키지 관리자
/google/googletest - Google Test 프레임워크
/fmtlib/fmt - 포맷팅 라이브러리
/nlohmann/json - JSON for Modern C++
/gabime/spdlog - 고속 로깅 라이브러리
```

---

## 관련 스킬

- `do-lang-rust` - 시스템 프로그래밍 비교
- `do-domain-backend` - 백엔드 아키텍처
- `do-workflow-testing` - TDD 및 테스트 전략

---

## 문제 해결

버전 확인:
```bash
g++ --version       # GCC 13+ (C++23)
clang++ --version   # Clang 17+ (C++23)
cmake --version     # CMake 3.28+
```

컴파일 플래그:
```bash
g++ -std=c++23 -Wall -Wextra -Wpedantic -O2 main.cpp -o main
g++ -std=c++23 -fsanitize=address,undefined -g main.cpp -o main
```

vcpkg 통합:
```bash
git clone https://github.com/microsoft/vcpkg
./vcpkg/bootstrap-vcpkg.sh
./vcpkg/vcpkg install fmt nlohmann-json gtest
cmake -B build -DCMAKE_TOOLCHAIN_FILE=./vcpkg/scripts/buildsystems/vcpkg.cmake
```

---

Last Updated: 2026-01-06
Status: Active (v1.0.0)
