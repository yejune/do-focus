---
name: do-lang-python
description: Python 3.13+ 개발 전문가로서 FastAPI, Django, async 패턴, 데이터 과학, pytest를 사용한 테스팅, 그리고 현대적 Python 기능을 다룸. Python API, 웹 애플리케이션, 데이터 파이프라인 개발 또는 테스트 작성 시 사용.
version: 1.0.0
updated: 2026-01-06
status: active
allowed-tools: Read, Grep, Glob, Bash, mcp__context7__resolve-library-id, mcp__context7__query-docs
user-invocable: false
---

## 빠른 참조

Python 3.13+ 개발 전문가 - FastAPI, Django, async 패턴, pytest, 현대적 Python 기능.

자동 트리거: `.py` 파일, `pyproject.toml`, `requirements.txt`, `pytest.ini`, FastAPI/Django 논의

핵심 역량:
- Python 3.13 기능: JIT 컴파일러 (PEP 744), GIL-free 모드 (PEP 703), 패턴 매칭
- 웹 프레임워크: FastAPI 0.115+, Django 5.2 LTS
- 데이터 검증: Pydantic v2.9 (model_validate 패턴)
- ORM: SQLAlchemy 2.0 async 패턴
- 테스팅: pytest, fixture, async 테스팅, parametrize
- 패키지 관리: poetry, uv, pip (pyproject.toml)
- 타입 힌트: Protocol, TypeVar, ParamSpec, 현대적 typing 패턴
- Async: asyncio, async 제너레이터, task 그룹
- 데이터 과학: numpy, pandas, polars 기초

---

## 핵심 패턴

### FastAPI 엔드포인트

```python
from fastapi import FastAPI, Depends
from pydantic import BaseModel

app = FastAPI()

class UserCreate(BaseModel):
    name: str
    email: str

@app.post("/users/")
async def create_user(user: UserCreate) -> User:
    return await UserService.create(user)
```

### Pydantic v2.9 검증

```python
from pydantic import BaseModel, ConfigDict

class User(BaseModel):
    model_config = ConfigDict(from_attributes=True, str_strip_whitespace=True)
    id: int
    name: str
    email: str

user = User.model_validate(orm_obj)  # ORM 객체에서
user = User.model_validate_json(json_data)  # JSON에서
```

### pytest Async 테스트

```python
import pytest

@pytest.mark.asyncio
async def test_create_user(async_client):
    response = await async_client.post("/users/", json={"name": "Test"})
    assert response.status_code == 201
```

---

## 구현 가이드

### Python 3.13 신기능

JIT 컴파일러 (PEP 744):
- 실험적 기능, 기본 비활성화
- 활성화: `PYTHON_JIT=1` 환경 변수
- CPU 집약적 코드 성능 향상 제공

GIL-Free 모드 (PEP 703):
- 실험적 free-threaded 빌드 (python3.13t)
- 실제 병렬 스레드 실행 허용
- 프로덕션 권장 아님

패턴 매칭 (match/case):

```python
def process_response(response: dict) -> str:
    match response:
        case {"status": "ok", "data": data}:
            return f"Success: {data}"
        case {"status": "error", "message": msg}:
            return f"Error: {msg}"
        case _:
            return "Unknown response"
```

### FastAPI 0.115+ 패턴

Async 의존성 주입:

```python
from fastapi import FastAPI, Depends
from sqlalchemy.ext.asyncio import AsyncSession
from contextlib import asynccontextmanager

@asynccontextmanager
async def lifespan(app: FastAPI):
    await init_db()
    yield
    await cleanup()

app = FastAPI(lifespan=lifespan)

async def get_db() -> AsyncGenerator[AsyncSession, None]:
    async with async_session() as session:
        yield session

@app.get("/users/{user_id}")
async def get_user(user_id: int, db: AsyncSession = Depends(get_db)) -> UserResponse:
    user = await get_user_by_id(db, user_id)
    return UserResponse.model_validate(user)
```

클래스 기반 의존성:

```python
class Paginator:
    def __init__(self, page: int = 1, size: int = 20):
        self.page = max(1, page)
        self.size = min(100, max(1, size))
        self.offset = (self.page - 1) * self.size

@app.get("/items/")
async def list_items(pagination: Paginator = Depends()) -> list[Item]:
    return await Item.get_page(pagination.offset, pagination.size)
```

### Django 5.2 LTS 기능

복합 기본 키:

```python
class OrderItem(models.Model):
    order = models.ForeignKey(Order, on_delete=models.CASCADE)
    product = models.ForeignKey(Product, on_delete=models.CASCADE)

    class Meta:
        pk = models.CompositePrimaryKey("order", "product")
```

### Pydantic v2.9 심화 패턴

Annotated를 사용한 재사용 가능한 검증기:

```python
from typing import Annotated
from pydantic import AfterValidator, BaseModel

def validate_positive(v: int) -> int:
    if v <= 0:
        raise ValueError("Must be positive")
    return v

PositiveInt = Annotated[int, AfterValidator(validate_positive)]

class Product(BaseModel):
    price: PositiveInt
    quantity: PositiveInt
```

교차 필드 검증을 위한 Model Validator:

```python
from pydantic import BaseModel, model_validator
from typing import Self

class DateRange(BaseModel):
    start_date: date
    end_date: date

    @model_validator(mode="after")
    def validate_dates(self) -> Self:
        if self.end_date < self.start_date:
            raise ValueError("end_date must be after start_date")
        return self
```

### SQLAlchemy 2.0 Async 패턴

엔진 및 세션 설정:

```python
from sqlalchemy.ext.asyncio import create_async_engine, async_sessionmaker, AsyncSession

engine = create_async_engine(
    "postgresql+asyncpg://user:pass@localhost/db",
    pool_pre_ping=True,
)

async_session = async_sessionmaker(engine, expire_on_commit=False)
```

리포지토리 패턴:

```python
class UserRepository:
    def __init__(self, session: AsyncSession):
        self.session = session

    async def get_by_id(self, user_id: int) -> User | None:
        result = await self.session.execute(select(User).where(User.id == user_id))
        return result.scalar_one_or_none()

    async def create(self, user: UserCreate) -> User:
        db_user = User(**user.model_dump())
        self.session.add(db_user)
        await self.session.commit()
        await self.session.refresh(db_user)
        return db_user
```

### pytest 고급 패턴

pytest-asyncio를 사용한 Async Fixtures:

```python
import pytest_asyncio
from httpx import AsyncClient

@pytest_asyncio.fixture
async def async_client(app) -> AsyncGenerator[AsyncClient, None]:
    async with AsyncClient(app=app, base_url="http://test") as client:
        yield client

@pytest_asyncio.fixture
async def db_session() -> AsyncGenerator[AsyncSession, None]:
    async with async_session() as session:
        async with session.begin():
            yield session
            await session.rollback()
```

매개변수화된 테스트:

```python
@pytest.mark.parametrize(
    "input_data,expected_status",
    [
        ({"name": "Valid"}, 201),
        ({"name": ""}, 422),
        ({}, 422),
    ],
    ids=["valid", "empty_name", "missing_name"],
)
async def test_create_user(async_client, input_data, expected_status):
    response = await async_client.post("/users/", json=input_data)
    assert response.status_code == expected_status
```

### 타입 힌트 현대적 패턴

구조적 타이핑을 위한 Protocol:

```python
from typing import Protocol, runtime_checkable

@runtime_checkable
class Repository(Protocol[T]):
    async def get(self, id: int) -> T | None: ...
    async def create(self, data: dict) -> T: ...
    async def delete(self, id: int) -> bool: ...
```

데코레이터를 위한 ParamSpec:

```python
from typing import ParamSpec, TypeVar, Callable
from functools import wraps

P = ParamSpec("P")
R = TypeVar("R")

def retry(times: int = 3) -> Callable[[Callable[P, R]], Callable[P, R]]:
    def decorator(func: Callable[P, R]) -> Callable[P, R]:
        @wraps(func)
        async def wrapper(*args: P.args, **kwargs: P.kwargs) -> R:
            for attempt in range(times):
                try:
                    return await func(*args, **kwargs)
                except Exception:
                    if attempt == times - 1:
                        raise
        return wrapper
    return decorator
```

### 패키지 관리

pyproject.toml (Poetry):

```toml
[tool.poetry]
name = "my-project"
version = "1.0.0"
python = "^3.13"

[tool.poetry.dependencies]
fastapi = "^0.115.0"
pydantic = "^2.9.0"
sqlalchemy = {extras = ["asyncio"], version = "^2.0.0"}

[tool.poetry.group.dev.dependencies]
pytest = "^8.0"
pytest-asyncio = "^0.24"
ruff = "^0.8"

[tool.ruff]
line-length = 100
target-version = "py313"

[tool.pytest.ini_options]
asyncio_mode = "auto"
```

uv (빠른 패키지 관리자):

```bash
curl -LsSf https://astral.sh/uv/install.sh | sh
uv venv
uv pip install -r requirements.txt
uv add fastapi
```

---

## Context7 라이브러리 매핑

```
/tiangolo/fastapi - FastAPI async 웹 프레임워크
/django/django - Django 웹 프레임워크
/pydantic/pydantic - 타입 어노테이션 기반 데이터 검증
/sqlalchemy/sqlalchemy - SQL 툴킷 및 ORM
/pytest-dev/pytest - 테스트 프레임워크
/numpy/numpy - 수치 계산
/pandas-dev/pandas - 데이터 분석 라이브러리
/pola-rs/polars - 빠른 DataFrame 라이브러리
```

---

## 관련 스킬

- `do-domain-backend` - REST API 및 마이크로서비스 아키텍처
- `do-domain-database` - SQL 패턴 및 ORM 최적화
- `do-workflow-testing` - TDD 및 테스트 전략
- `do-domain-frontend` - 프론트엔드 통합

---

## 트러블슈팅

Python 버전 확인:
```bash
python --version  # 3.13+ 이어야 함
```

Async 세션 분리 오류:
- 해결: 세션 설정에서 `expire_on_commit=False` 설정
- 또는: 커밋 후 `await session.refresh(obj)` 사용

pytest asyncio 모드 경고:
```toml
[tool.pytest.ini_options]
asyncio_mode = "auto"
asyncio_default_fixture_loop_scope = "function"
```

Pydantic v2 마이그레이션:
- `parse_obj()` -> `model_validate()`
- `parse_raw()` -> `model_validate_json()`
- `from_orm()` -> ConfigDict에서 `from_attributes=True` 필요

---

마지막 업데이트: 2026-01-06
상태: 활성 (v1.0.0)
