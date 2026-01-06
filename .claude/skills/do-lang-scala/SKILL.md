---
name: do-lang-scala
description: Scala 3.4+ 개발 전문가 - Akka, Cats Effect, ZIO, Spark 패턴을 다룸. 분산 시스템, 빅 데이터 파이프라인, 함수형 프로그래밍 애플리케이션 구축 시 사용.
version: 1.0.0
category: language
tags:
  - scala
  - akka
  - cats-effect
  - zio
  - spark
  - functional-programming
context7-libraries:
  - /akka/akka
  - /typelevel/cats-effect
  - /zio/zio
  - /apache/spark
related-skills:
  - do-lang-java
  - do-domain-database
updated: 2025-12-07
status: active
allowed-tools:
  - mcp__context7__resolve-library-id
  - mcp__context7__query-docs
---

## 빠른 참조

Scala 3.4+ 개발 전문가 - 함수형 프로그래밍, 이펙트 시스템, 빅 데이터 처리.

자동 트리거: Scala 파일 (`.scala`, `.sc`), 빌드 파일 (`build.sbt`)

핵심 기능:
- Scala 3.4: Given/using, extension methods, enums, opaque types
- Akka 2.9: Typed actors, streams, clustering
- Cats Effect 3.5: Pure FP runtime, fibers
- ZIO 2.1: Effect system, layers, streaming
- Apache Spark 3.5: DataFrame API, structured streaming

주요 라이브러리:
- HTTP: Http4s 0.24, Tapir 1.10
- JSON: Circe 0.15, ZIO JSON 0.6
- DB: Doobie 1.0, Slick 3.5
- Stream: FS2 3.10, ZIO Streams 2.1

---

## Scala 3.4 핵심 기능

### Extension Methods

```scala
extension (s: String)
  def words: List[String] = s.split("\\s+").toList
  def truncate(maxLen: Int): String =
    if s.length <= maxLen then s else s.take(maxLen - 3) + "..."
```

### Given과 Using

```scala
trait JsonEncoder[A]:
  def encode(value: A): String

given JsonEncoder[String] with
  def encode(value: String): String = s"\"$value\""

given [A](using encoder: JsonEncoder[A]): JsonEncoder[List[A]] with
  def encode(value: List[A]): String =
    value.map(encoder.encode).mkString("[", ",", "]")

def toJson[A](value: A)(using encoder: JsonEncoder[A]): String =
  encoder.encode(value)
```

### Enum 타입과 ADT

```scala
enum Color(val hex: String):
  case Red extends Color("#FF0000")
  case Green extends Color("#00FF00")
  case Custom(override val hex: String) extends Color(hex)

enum Result[+E, +A]:
  case Success(value: A)
  case Failure(error: E)

  def map[B](f: A => B): Result[E, B] = this match
    case Success(a) => Success(f(a))
    case Failure(e) => Failure(e)
```

### Opaque Types

```scala
object UserId:
  opaque type UserId = Long
  def apply(id: Long): UserId = id
  extension (id: UserId)
    def value: Long = id
    def asString: String = id.toString
```

### Union과 Intersection Types

```scala
type StringOrInt = String | Int

def describe(value: StringOrInt): String = value match
  case s: String => s"String: $s"
  case i: Int => s"Int: $i"

type Person = HasName & HasAge
```

---

## Cats Effect 3.5

### 기본 IO 연산

```scala
import cats.effect.*
import cats.syntax.all.*

def program: IO[Unit] =
  for
    _ <- IO.println("Enter your name:")
    name <- IO.readLine
    _ <- IO.println(s"Hello, $name!")
  yield ()

def withFile[A](path: String)(use: BufferedReader => IO[A]): IO[A] =
  Resource.make(IO(new BufferedReader(new FileReader(path))))(r => IO(r.close())).use(use)
```

### 동시성 프로그래밍

```scala
import cats.effect.std.*

// 병렬 실행
def fetchUserData(userId: Long): IO[UserData] =
  (fetchUser(userId), fetchOrders(userId), fetchPreferences(userId)).parMapN(UserData.apply)

// Semaphore로 속도 제한
def rateLimitedRequests[A](tasks: List[IO[A]], max: Int): IO[List[A]] =
  Semaphore[IO](max).flatMap { sem =>
    tasks.parTraverse(task => sem.permit.use(_ => task))
  }
```

### FS2 스트리밍

```scala
import fs2.*
import fs2.io.file.*

def processLargeFile(path: Path): Stream[IO, String] =
  Files[IO].readUtf8Lines(path)
    .filter(_.nonEmpty)
    .map(_.toLowerCase)
    .evalTap(line => IO.println(s"Processing: $line"))
```

---

## ZIO 2.1

### 기본 ZIO 연산

```scala
import zio.*

val program: ZIO[Any, Nothing, Unit] =
  for
    _ <- Console.printLine("Enter your name:")
    name <- Console.readLine
    _ <- Console.printLine(s"Hello, $name!")
  yield ()

def fetchUser(id: Long): ZIO[UserRepository, UserError, User] =
  for
    repo <- ZIO.service[UserRepository]
    user <- ZIO.fromOption(repo.findById(id)).orElseFail(UserNotFound(id))
  yield user
```

### ZIO Layers

```scala
trait UserRepository:
  def findById(id: Long): Task[Option[User]]
  def save(user: User): Task[User]

case class UserRepositoryLive(db: Database) extends UserRepository:
  def findById(id: Long): Task[Option[User]] =
    ZIO.attempt(db.query(s"SELECT * FROM users WHERE id = $id")).map(_.headOption)

object UserRepositoryLive:
  val layer: ZLayer[Database, Nothing, UserRepository] =
    ZLayer.fromFunction(UserRepositoryLive.apply)

val appLayer = Database.layer >>> UserRepositoryLive.layer ++ EmailServiceLive.layer
```

### ZIO Streaming

```scala
import zio.stream.*

def processEvents: ZStream[Any, Throwable, ProcessedEvent] =
  ZStream.fromQueue(eventQueue)
    .filter(_.isValid)
    .mapZIO(enrichEvent)
    .grouped(100)
    .mapZIO(batchProcess)
    .flattenIterables
```

---

## Akka Typed Actors

### 액터 정의

```scala
import akka.actor.typed.*
import akka.actor.typed.scaladsl.*

object UserActor:
  sealed trait Command
  case class GetUser(id: Long, replyTo: ActorRef[Option[User]]) extends Command
  case class CreateUser(request: CreateUserRequest, replyTo: ActorRef[User]) extends Command

  def apply(repository: UserRepository): Behavior[Command] =
    Behaviors.receiveMessage {
      case GetUser(id, replyTo) =>
        replyTo ! repository.findById(id)
        Behaviors.same
      case CreateUser(request, replyTo) =>
        replyTo ! repository.save(User.from(request))
        Behaviors.same
    }
```

### Akka Streams

```scala
import akka.stream.*
import akka.stream.scaladsl.*

val source: Source[Int, NotUsed] = Source(1 to 1000)
val flow: Flow[Int, String, NotUsed] =
  Flow[Int].filter(_ % 2 == 0).map(_ * 2).map(_.toString)
val sink: Sink[String, Future[Done]] = Sink.foreach(println)

val graph = source.via(flow).toMat(sink)(Keep.right)

val throttledSource = source
  .throttle(100, 1.second)
  .buffer(1000, OverflowStrategy.backpressure)
```

---

## Apache Spark 3.5

### DataFrame 연산

```scala
import org.apache.spark.sql.{DataFrame, SparkSession}
import org.apache.spark.sql.functions.*

val spark = SparkSession.builder()
  .appName("Data Analysis")
  .config("spark.sql.adaptive.enabled", "true")
  .getOrCreate()

import spark.implicits.*

val userMetrics = orders
  .groupBy("user_id")
  .agg(
    sum("amount").as("total_spent"),
    count("*").as("order_count"),
    avg("amount").as("avg_order_value")
  )
  .join(users, Seq("user_id"), "left")
  .withColumn("customer_tier",
    when(col("total_spent") > 10000, "platinum")
      .when(col("total_spent") > 1000, "gold")
      .otherwise("standard")
  )
```

### Structured Streaming

```scala
val streamingOrders = spark.readStream
  .format("kafka")
  .option("kafka.bootstrap.servers", "localhost:9092")
  .option("subscribe", "orders")
  .load()
  .selectExpr("CAST(value AS STRING)")

val aggregated = streamingOrders
  .withWatermark("timestamp", "10 minutes")
  .groupBy(window($"timestamp", "1 hour"), $"product_category")
  .agg(sum("amount").as("hourly_sales"))
```

---

## 빌드 설정 (SBT 1.10)

```scala
ThisBuild / scalaVersion := "3.4.2"
ThisBuild / organization := "com.example"

lazy val root = (project in file("."))
  .settings(
    name := "scala-service",
    libraryDependencies ++= Seq(
      "org.typelevel" %% "cats-effect" % "3.5.4",
      "dev.zio" %% "zio" % "2.1.0",
      "com.typesafe.akka" %% "akka-actor-typed" % "2.9.0",
      "org.http4s" %% "http4s-ember-server" % "0.24.0",
      "io.circe" %% "circe-generic" % "0.15.0",
      "org.tpolecat" %% "doobie-core" % "1.0.0-RC4",
      "org.scalatest" %% "scalatest" % "3.2.18" % Test
    ),
    scalacOptions ++= Seq("-deprecation", "-feature", "-Xfatal-warnings")
  )
```

---

## 테스팅

### ScalaTest

```scala
class UserServiceSpec extends AnyFlatSpec with Matchers:
  "UserService" should "create user successfully" in {
    val result = service.createUser(CreateUserRequest("John", "john@example.com"))
    result.name shouldBe "John"
  }
```

### MUnit with Cats Effect

```scala
class UserServiceSuite extends CatsEffectSuite:
  test("should fetch user") {
    UserService.findById(1L).map { result =>
      assertEquals(result.name, "John")
    }
  }
```

### ZIO Test

```scala
object UserServiceSpec extends ZIOSpecDefault:
  def spec = suite("UserService")(
    test("should find user") {
      for result <- UserService.findById(1L)
      yield assertTrue(result.name == "John")
    }
  )
```

---

## 트러블슈팅

일반적인 문제:
- Implicit 해석: `scalac -explain`으로 상세 에러 확인
- 타입 추론 실패: 명시적 타입 주석 추가
- SBT 느린 컴파일: `Global / concurrentRestrictions` 활성화

이펙트 시스템 문제:
- Cats Effect: `import cats.effect.*` 또는 `import cats.syntax.all.*` 누락 확인
- ZIO: `ZIO.serviceWith`와 `ZIO.serviceWithZIO`로 레이어 구성 검증
- Akka: 액터 계층과 supervision 전략 검토

---

## 관련 스킬

- `do-lang-java` - JVM 상호운용성, Spring Boot 통합
- `do-domain-backend` - REST API, GraphQL, 마이크로서비스 패턴
- `do-domain-database` - Doobie, Slick, 데이터베이스 패턴

---

Last Updated: 2025-12-07
Status: Production Ready (v1.0.0)
