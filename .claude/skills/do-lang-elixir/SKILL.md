---
name: do-lang-elixir
description: Elixir 1.17+ 개발 전문가 - Phoenix 1.7, LiveView, Ecto, OTP 패턴 지원. 실시간 애플리케이션, 분산 시스템, Phoenix 프로젝트 개발 시 사용.
version: 1.0.0
updated: 2025-12-07
status: active
allowed-tools: Read, Grep, Glob, Bash, mcp__context7__resolve-library-id, mcp__context7__query-docs
---

## 빠른 참조

Elixir 1.17+ 개발 전문가 - Phoenix 1.7, LiveView, Ecto, OTP 패턴, 함수형 프로그래밍 지원.

자동 트리거: `.ex`, `.exs` 파일, `mix.exs`, `config/`, Phoenix/LiveView 관련 논의

핵심 기능:
- Elixir 1.17: 패턴 매칭, 파이프 연산자, 프로토콜, 비헤이비어, 매크로
- Phoenix 1.7: 컨트롤러, LiveView, 채널, PubSub, Verified Routes
- Ecto: 스키마, Changeset, 쿼리, 마이그레이션, Multi
- OTP: GenServer, Supervisor, Agent, Task, Registry
- ExUnit: 테스트 setup, describe, async
- Oban: 백그라운드 작업 처리

---

## 기본 패턴

### Phoenix 컨트롤러

```elixir
defmodule MyAppWeb.UserController do
  use MyAppWeb, :controller
  alias MyApp.Accounts

  def create(conn, %{"user" => user_params}) do
    case Accounts.create_user(user_params) do
      {:ok, user} ->
        conn |> put_flash(:info, "Created!") |> redirect(to: ~p"/users/#{user}")
      {:error, changeset} ->
        render(conn, :new, changeset: changeset)
    end
  end
end
```

### Ecto 스키마와 Changeset

```elixir
defmodule MyApp.Accounts.User do
  use Ecto.Schema
  import Ecto.Changeset

  schema "users" do
    field :name, :string
    field :email, :string
    field :password, :string, virtual: true
    timestamps()
  end

  def changeset(user, attrs) do
    user
    |> cast(attrs, [:name, :email, :password])
    |> validate_required([:name, :email, :password])
    |> validate_format(:email, ~r/@/)
    |> unique_constraint(:email)
  end
end
```

### GenServer 패턴

```elixir
defmodule MyApp.Counter do
  use GenServer

  def start_link(init), do: GenServer.start_link(__MODULE__, init, name: __MODULE__)
  def increment, do: GenServer.call(__MODULE__, :increment)

  @impl true
  def init(value), do: {:ok, value}

  @impl true
  def handle_call(:increment, _from, count), do: {:reply, count + 1, count + 1}
end
```

---

## Elixir 1.17 기능

### with 구문을 활용한 에러 처리

```elixir
def process_order(params) do
  with {:ok, validated} <- validate_order(params),
       {:ok, total} <- calculate_total(validated),
       {:ok, order} <- create_order(total) do
    {:ok, order}
  else
    {:error, reason} -> {:error, reason}
  end
end
```

### 프로토콜

```elixir
defprotocol Stringify do
  def to_string(data)
end

defimpl Stringify, for: Map do
  def to_string(map), do: Jason.encode!(map)
end
```

---

## Phoenix 1.7 패턴

### LiveView 컴포넌트

```elixir
defmodule MyAppWeb.CounterLive do
  use MyAppWeb, :live_view

  def mount(_params, _session, socket), do: {:ok, assign(socket, count: 0)}

  def handle_event("increment", _, socket) do
    {:noreply, update(socket, :count, &(&1 + 1))}
  end

  def render(assigns) do
    ~H"""
    <div>
      <h1>Count: <%= @count %></h1>
      <button phx-click="increment">+</button>
    </div>
    """
  end
end
```

### LiveView 폼

```elixir
defmodule MyAppWeb.UserFormLive do
  use MyAppWeb, :live_view

  def mount(_params, _session, socket) do
    {:ok, assign(socket, form: to_form(Accounts.change_user(%User{})))}
  end

  def handle_event("save", %{"user" => params}, socket) do
    case Accounts.create_user(params) do
      {:ok, user} -> {:noreply, push_navigate(socket, to: ~p"/users/#{user}")}
      {:error, cs} -> {:noreply, assign(socket, form: to_form(cs))}
    end
  end

  def render(assigns) do
    ~H"""
    <.form for={@form} phx-submit="save">
      <.input field={@form[:name]} label="Name" />
      <.button>Save</.button>
    </.form>
    """
  end
end
```

### Phoenix 채널

```elixir
defmodule MyAppWeb.RoomChannel do
  use MyAppWeb, :channel

  @impl true
  def join("room:" <> room_id, _params, socket) do
    {:ok, assign(socket, :room_id, room_id)}
  end

  @impl true
  def handle_in("new_message", %{"body" => body}, socket) do
    broadcast!(socket, "new_message", %{body: body})
    {:noreply, socket}
  end
end
```

### Verified Routes

```elixir
# router.ex
live "/users", UserLive.Index, :index
live "/users/:id", UserLive.Show, :show

# 사용
~p"/users"           # "/users"
~p"/users/#{user}"   # "/users/123"
```

---

## Ecto 패턴

### Multi를 활용한 트랜잭션

```elixir
def transfer_funds(from, to, amount) do
  Ecto.Multi.new()
  |> Ecto.Multi.update(:withdraw, withdraw_changeset(from, amount))
  |> Ecto.Multi.update(:deposit, deposit_changeset(to, amount))
  |> Repo.transaction()
end
```

### 쿼리 조합

```elixir
defmodule MyApp.UserQuery do
  import Ecto.Query
  def base, do: from(u in User)
  def active(q \\ base()), do: from u in q, where: u.active == true
  def with_posts(q \\ base()), do: from u in q, preload: [:posts]
end

# 사용
User |> UserQuery.active() |> UserQuery.with_posts() |> Repo.all()
```

---

## OTP 패턴

### Supervisor 트리

```elixir
defmodule MyApp.Application do
  use Application

  @impl true
  def start(_type, _args) do
    children = [
      MyApp.Repo,
      {Phoenix.PubSub, name: MyApp.PubSub},
      MyAppWeb.Endpoint
    ]
    Supervisor.start_link(children, strategy: :one_for_one)
  end
end
```

### Dynamic Supervisor

```elixir
defmodule MyApp.WorkerSupervisor do
  use DynamicSupervisor

  def start_link(arg), do: DynamicSupervisor.start_link(__MODULE__, arg, name: __MODULE__)

  @impl true
  def init(_), do: DynamicSupervisor.init(strategy: :one_for_one)

  def start_worker(args), do: DynamicSupervisor.start_child(__MODULE__, {MyApp.Worker, args})
end
```

---

## 테스트 패턴

### ExUnit 테스트

```elixir
defmodule MyApp.AccountsTest do
  use MyApp.DataCase, async: true

  describe "users" do
    test "create_user/1 with valid data" do
      assert {:ok, %User{}} = Accounts.create_user(%{name: "Test", email: "t@e.com"})
    end
  end
end
```

### LiveView 테스트

```elixir
defmodule MyAppWeb.CounterLiveTest do
  use MyAppWeb.ConnCase
  import Phoenix.LiveViewTest

  test "increments counter", %{conn: conn} do
    {:ok, view, _} = live(conn, ~p"/counter")
    assert view |> element("button") |> render_click() =~ "Count: 1"
  end
end
```

---

## Oban 백그라운드 작업

```elixir
defmodule MyApp.Workers.EmailWorker do
  use Oban.Worker, queue: :mailers

  @impl Oban.Worker
  def perform(%Oban.Job{args: %{"email" => email}}) do
    MyApp.Mailer.send_email(email)
  end
end

# 등록
%{email: "user@example.com"} |> MyApp.Workers.EmailWorker.new() |> Oban.insert()
```

---

## Context7 라이브러리

- `/elixir-lang/elixir` - Elixir 언어 문서
- `/phoenixframework/phoenix` - Phoenix 웹 프레임워크
- `/phoenixframework/phoenix_live_view` - LiveView 실시간 UI
- `/elixir-ecto/ecto` - 데이터베이스 래퍼 및 쿼리 언어
- `/sorentwo/oban` - 백그라운드 작업 처리

---

## 연관 스킬

- `do-domain-backend` - REST API 및 마이크로서비스 아키텍처
- `do-domain-database` - SQL 패턴 및 쿼리 최적화
- `do-workflow-testing` - TDD 및 테스트 전략

---

## 문제 해결

버전 확인:
```bash
elixir --version  # 1.17+ 필요
```

의존성:
```bash
mix deps.get      # 의존성 가져오기
mix deps.compile  # 컴파일
```

데이터베이스:
```bash
mix ecto.create   # DB 생성
mix ecto.migrate  # 마이그레이션 실행
```

Phoenix 서버:
```bash
mix phx.server           # 서버 시작
MIX_ENV=prod mix release # 릴리즈 빌드
```

LiveView 로딩 실패 시:
- 브라우저 콘솔에서 웹소켓 연결 확인
- mix.exs 의존성에 Phoenix.LiveView 포함 여부 확인

---

Last Updated: 2025-12-07
Status: Active (v1.0.0)
