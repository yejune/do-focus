---
name: do-lang-csharp
description: C# 12 / .NET 8 개발 전문가 - ASP.NET Core, Entity Framework, Blazor 및 현대적 C# 패턴 지원. .NET API, 웹 애플리케이션, 엔터프라이즈 솔루션 개발에 사용
version: 1.0.0
category: language
tags:
  - csharp
  - dotnet
  - aspnetcore
  - efcore
  - blazor
updated: 2025-01-06
status: active
allowed-tools:
  - Read
  - Glob
  - Grep
  - Edit
  - Write
  - Bash
  - mcp__context7__resolve-library-id
  - mcp__context7__query-docs
user-invocable: false
---

## 빠른 참조

C# 12 / .NET 8 개발 전문가 - ASP.NET Core, Entity Framework Core, Blazor, 엔터프라이즈 패턴 기반의 현대적 C# 개발 지원

자동 트리거: `.cs`, `.csproj`, `.sln` 파일, C# 프로젝트, .NET 솔루션, ASP.NET Core 애플리케이션

핵심 스택:
- C# 12: Primary constructors, Collection expressions, Alias any type, Default lambda parameters
- .NET 8: Minimal APIs, Native AOT, 향상된 성능, WebSockets
- ASP.NET Core 8: Controllers, Endpoints, Middleware, Authentication
- Entity Framework Core 8: DbContext, Migrations, LINQ, 쿼리 최적화
- Blazor: Server/WASM 컴포넌트, InteractiveServer, InteractiveWebAssembly
- Testing: xUnit, NUnit, FluentAssertions, Moq

빠른 명령어:
```bash
# .NET 8 Web API 프로젝트 생성
dotnet new webapi -n MyApi --framework net8.0

# Blazor Web App 생성
dotnet new blazor -n MyBlazor --interactivity Auto

# Entity Framework Core 추가
dotnet add package Microsoft.EntityFrameworkCore.SqlServer
dotnet add package Microsoft.EntityFrameworkCore.Design

# FluentValidation, MediatR 추가
dotnet add package FluentValidation.AspNetCore
dotnet add package MediatR
```

---

## 구현 가이드

### C# 12 주요 기능

Primary Constructors - 클래스 수준 생성자 매개변수:
```csharp
// 의존성 주입과 함께 사용하는 Primary constructor
public class UserService(IUserRepository repository, ILogger<UserService> logger)
{
    public async Task<User?> GetByIdAsync(Guid id)
    {
        logger.LogInformation("Fetching user {UserId}", id);
        return await repository.FindByIdAsync(id);
    }
}

// Record와 Primary constructor
public record CreateUserCommand(string Name, string Email);
```

Collection Expressions - 통합 컬렉션 문법:
```csharp
// 배열, 리스트, Span 통합 문법
int[] numbers = [1, 2, 3, 4, 5];
List<string> names = ["Alice", "Bob", "Charlie"];
Span<int> span = [10, 20, 30];

// Spread 연산자
int[] combined = [..numbers, 6, 7, 8];
List<string> allNames = [..names, "David", "Eve"];
```

Alias Any Type - 복잡한 타입에 별칭 지정:
```csharp
// 튜플 별칭
using Point = (int X, int Y);

// 복잡한 제네릭 별칭
using UserCache = System.Collections.Generic.Dictionary<Guid, User>;

public class LocationService
{
    public Point GetLocation() => (10, 20);
    private readonly UserCache _cache = [];
}
```

Default Lambda Parameters:
```csharp
// 기본 매개변수가 있는 람다
var greet = (string name, string greeting = "Hello") => $"{greeting}, {name}!";

Console.WriteLine(greet("Alice"));           // "Hello, Alice!"
Console.WriteLine(greet("Bob", "Hi"));       // "Hi, Bob!"
```

### ASP.NET Core 8 패턴

Minimal API with Endpoints:
```csharp
var builder = WebApplication.CreateBuilder(args);

// 서비스 등록
builder.Services.AddDbContext<AppDbContext>(options =>
    options.UseSqlServer(builder.Configuration.GetConnectionString("Default")));
builder.Services.AddScoped<IUserService, UserService>();

var app = builder.Build();

// 타입화된 결과와 함께 Endpoint 라우팅
app.MapGet("/api/users/{id:guid}", async (Guid id, IUserService service) =>
{
    var user = await service.GetByIdAsync(id);
    return user is not null ? Results.Ok(user) : Results.NotFound();
})
.WithName("GetUser")
.WithOpenApi()
.Produces<User>(200)
.Produces(404);

app.MapPost("/api/users", async (CreateUserRequest request, IUserService service) =>
{
    var user = await service.CreateAsync(request);
    return Results.Created($"/api/users/{user.Id}", user);
})
.WithValidation<CreateUserRequest>();

app.Run();
```

Controller 기반 API:
```csharp
[ApiController]
[Route("api/[controller]")]
public class UsersController(IUserService userService, ILogger<UsersController> logger)
    : ControllerBase
{
    [HttpGet("{id:guid}")]
    [ProducesResponseType<User>(StatusCodes.Status200OK)]
    [ProducesResponseType(StatusCodes.Status404NotFound)]
    public async Task<ActionResult<User>> GetById(Guid id)
    {
        var user = await userService.GetByIdAsync(id);
        if (user is null)
        {
            logger.LogWarning("User {UserId} not found", id);
            return NotFound();
        }
        return user;
    }

    [HttpPost]
    [ProducesResponseType<User>(StatusCodes.Status201Created)]
    [ProducesResponseType<ValidationProblemDetails>(StatusCodes.Status400BadRequest)]
    public async Task<ActionResult<User>> Create([FromBody] CreateUserRequest request)
    {
        var user = await userService.CreateAsync(request);
        return CreatedAtAction(nameof(GetById), new { id = user.Id }, user);
    }
}
```

### Entity Framework Core 8 패턴

DbContext 구성:
```csharp
public class AppDbContext(DbContextOptions<AppDbContext> options) : DbContext(options)
{
    public DbSet<User> Users => Set<User>();
    public DbSet<Post> Posts => Set<Post>();
    public DbSet<Tag> Tags => Set<Tag>();

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.ApplyConfigurationsFromAssembly(typeof(AppDbContext).Assembly);
    }
}

// 엔티티 구성
public class UserConfiguration : IEntityTypeConfiguration<User>
{
    public void Configure(EntityTypeBuilder<User> builder)
    {
        builder.HasKey(u => u.Id);
        builder.Property(u => u.Email).HasMaxLength(256).IsRequired();
        builder.HasIndex(u => u.Email).IsUnique();
        builder.HasMany(u => u.Posts).WithOne(p => p.Author).HasForeignKey(p => p.AuthorId);
    }
}
```

Repository 패턴과 Specification:
```csharp
public interface IRepository<T> where T : class
{
    Task<T?> GetByIdAsync(Guid id, CancellationToken ct = default);
    Task<IReadOnlyList<T>> ListAsync(CancellationToken ct = default);
    Task<T> AddAsync(T entity, CancellationToken ct = default);
    Task UpdateAsync(T entity, CancellationToken ct = default);
    Task DeleteAsync(T entity, CancellationToken ct = default);
}

public class EfRepository<T>(AppDbContext context) : IRepository<T> where T : class
{
    private readonly DbSet<T> _dbSet = context.Set<T>();

    public async Task<T?> GetByIdAsync(Guid id, CancellationToken ct = default)
        => await _dbSet.FindAsync([id], ct);

    public async Task<IReadOnlyList<T>> ListAsync(CancellationToken ct = default)
        => await _dbSet.ToListAsync(ct);

    public async Task<T> AddAsync(T entity, CancellationToken ct = default)
    {
        await _dbSet.AddAsync(entity, ct);
        await context.SaveChangesAsync(ct);
        return entity;
    }

    public async Task UpdateAsync(T entity, CancellationToken ct = default)
    {
        _dbSet.Update(entity);
        await context.SaveChangesAsync(ct);
    }

    public async Task DeleteAsync(T entity, CancellationToken ct = default)
    {
        _dbSet.Remove(entity);
        await context.SaveChangesAsync(ct);
    }
}
```

### FluentValidation 패턴

요청 검증:
```csharp
public record CreateUserRequest(string Name, string Email, string Password);

public class CreateUserRequestValidator : AbstractValidator<CreateUserRequest>
{
    public CreateUserRequestValidator(IUserRepository userRepository)
    {
        RuleFor(x => x.Name)
            .NotEmpty().WithMessage("Name is required")
            .MaximumLength(100).WithMessage("Name cannot exceed 100 characters");

        RuleFor(x => x.Email)
            .NotEmpty().WithMessage("Email is required")
            .EmailAddress().WithMessage("Invalid email format")
            .MustAsync(async (email, ct) => !await userRepository.EmailExistsAsync(email, ct))
            .WithMessage("Email already exists");

        RuleFor(x => x.Password)
            .NotEmpty().WithMessage("Password is required")
            .MinimumLength(8).WithMessage("Password must be at least 8 characters")
            .Matches(@"[A-Z]").WithMessage("Password must contain uppercase letter")
            .Matches(@"[a-z]").WithMessage("Password must contain lowercase letter")
            .Matches(@"[0-9]").WithMessage("Password must contain digit");
    }
}

// Program.cs에서 등록
builder.Services.AddValidatorsFromAssemblyContaining<CreateUserRequestValidator>();
```

### MediatR CQRS 패턴

Command와 Query 분리:
```csharp
// Query
public record GetUserByIdQuery(Guid Id) : IRequest<User?>;

public class GetUserByIdQueryHandler(AppDbContext context)
    : IRequestHandler<GetUserByIdQuery, User?>
{
    public async Task<User?> Handle(GetUserByIdQuery request, CancellationToken ct)
        => await context.Users
            .AsNoTracking()
            .Include(u => u.Posts)
            .FirstOrDefaultAsync(u => u.Id == request.Id, ct);
}

// Command
public record CreateUserCommand(string Name, string Email, string Password) : IRequest<User>;

public class CreateUserCommandHandler(
    AppDbContext context,
    IPasswordHasher passwordHasher,
    IValidator<CreateUserCommand> validator)
    : IRequestHandler<CreateUserCommand, User>
{
    public async Task<User> Handle(CreateUserCommand request, CancellationToken ct)
    {
        await validator.ValidateAndThrowAsync(request, ct);

        var user = new User
        {
            Id = Guid.NewGuid(),
            Name = request.Name,
            Email = request.Email,
            PasswordHash = passwordHasher.Hash(request.Password),
            CreatedAt = DateTime.UtcNow
        };

        context.Users.Add(user);
        await context.SaveChangesAsync(ct);
        return user;
    }
}

// Controller에서 사용
[HttpPost]
public async Task<ActionResult<User>> Create(
    [FromBody] CreateUserCommand command,
    [FromServices] IMediator mediator)
{
    var user = await mediator.Send(command);
    return CreatedAtAction(nameof(GetById), new { id = user.Id }, user);
}
```

### Blazor 패턴

Interactive Server 컴포넌트:
```csharp
@page "/users"
@rendermode InteractiveServer
@inject IUserService UserService

<h1>Users</h1>

@if (_loading)
{
    <p>Loading...</p>
}
else if (_users is null)
{
    <p>No users found.</p>
}
else
{
    <table class="table">
        <thead>
            <tr><th>Name</th><th>Email</th><th>Actions</th></tr>
        </thead>
        <tbody>
            @foreach (var user in _users)
            {
                <tr>
                    <td>@user.Name</td>
                    <td>@user.Email</td>
                    <td>
                        <button @onclick="() => DeleteUser(user.Id)" class="btn btn-danger btn-sm">
                            Delete
                        </button>
                    </td>
                </tr>
            }
        </tbody>
    </table>
}

@code {
    private List<User>? _users;
    private bool _loading = true;

    protected override async Task OnInitializedAsync()
    {
        _users = await UserService.GetAllAsync();
        _loading = false;
    }

    private async Task DeleteUser(Guid id)
    {
        await UserService.DeleteAsync(id);
        _users = await UserService.GetAllAsync();
    }
}
```

### 인증과 권한 부여

JWT Authentication 설정:
```csharp
builder.Services.AddAuthentication(JwtBearerDefaults.AuthenticationScheme)
    .AddJwtBearer(options =>
    {
        options.TokenValidationParameters = new TokenValidationParameters
        {
            ValidateIssuer = true,
            ValidateAudience = true,
            ValidateLifetime = true,
            ValidateIssuerSigningKey = true,
            ValidIssuer = builder.Configuration["Jwt:Issuer"],
            ValidAudience = builder.Configuration["Jwt:Audience"],
            IssuerSigningKey = new SymmetricSecurityKey(
                Encoding.UTF8.GetBytes(builder.Configuration["Jwt:Key"]!))
        };
    });

builder.Services.AddAuthorization(options =>
{
    options.AddPolicy("Admin", policy => policy.RequireRole("Admin"));
    options.AddPolicy("CanEdit", policy => policy.RequireClaim("permission", "edit"));
});
```

---

## 고급 패턴

### Context7 통합

최신 문서 조회 방법:
```csharp
// ASP.NET Core - mcp__context7__query-docs("/dotnet/aspnetcore", "minimal-apis middleware")
// EF Core - mcp__context7__query-docs("/dotnet/efcore", "dbcontext migrations")
// .NET Runtime - mcp__context7__query-docs("/dotnet/runtime", "collections threading")
// Blazor - mcp__context7__query-docs("/dotnet/aspnetcore", "blazor components")
```

---

## 관련 스킬

- `do-domain-backend` - API 설계, 데이터베이스 통합 패턴
- `do-workflow-testing` - 테스팅 전략과 패턴
- `do-foundation-quality` - 코드 품질 표준

---

## 빠른 문제 해결

빌드와 런타임:
```bash
dotnet build --verbosity detailed    # 상세 빌드 출력
dotnet run --launch-profile https    # HTTPS 프로필로 실행
dotnet ef database update            # EF 마이그레이션 적용
dotnet ef migrations add Initial     # 새 마이그레이션 생성
```

일반적인 이슈:
```csharp
// Null 참조 처리
var user = await context.Users.FindAsync(id);
ArgumentNullException.ThrowIfNull(user, nameof(user));

// 스트리밍을 위한 비동기 열거형
public async IAsyncEnumerable<User> StreamUsersAsync(
    [EnumeratorCancellation] CancellationToken ct = default)
{
    await foreach (var user in context.Users.AsAsyncEnumerable().WithCancellation(ct))
    {
        yield return user;
    }
}
```

---

Version: 1.0.0
Last Updated: 2025-01-06
Status: Production Ready
