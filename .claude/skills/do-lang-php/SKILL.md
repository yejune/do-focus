---
name: do-lang-php
description: PHP 8.3+ 개발 전문가 - Laravel 11, Symfony 7, Eloquent ORM, 최신 PHP 패턴. PHP API, 웹 애플리케이션, Laravel/Symfony 프로젝트 개발 시 활용.
version: 1.0.0
updated: 2025-01-06
status: active
allowed-tools: Read, Grep, Glob, Bash, mcp__context7__resolve-library-id, mcp__context7__query-docs
---

## 빠른 참조

PHP 8.3+ 개발 전문가 - Laravel 11, Symfony 7, Eloquent, Doctrine, 최신 PHP 패턴 지원.

자동 트리거: `.php` 파일, `composer.json`, `artisan`, `symfony.yaml`, Laravel/Symfony 관련 논의

핵심 역량:
- PHP 8.3 기능: readonly 클래스, typed 속성, Attributes, enums, named arguments
- Laravel 11: Controllers, Models, Migrations, Form Requests, API Resources, Eloquent
- Symfony 7: Attribute 기반 라우팅, Doctrine ORM, Services, 의존성 주입
- ORM: Eloquent (Laravel), Doctrine (Symfony)
- 테스트: PHPUnit, Pest, feature/unit 테스트 패턴
- 패키지 관리: Composer 자동로딩, PSR-12, Laravel Pint

---

## 빠른 패턴

Laravel Controller:
```php
<?php
namespace App\Http\Controllers\Api;

class UserController extends Controller
{
    public function store(StoreUserRequest $request): JsonResponse
    {
        $user = User::create($request->validated());
        return response()->json(new UserResource($user), 201);
    }
}
```

Laravel Form Request:
```php
<?php
namespace App\Http\Requests;

class StoreUserRequest extends FormRequest
{
    public function rules(): array
    {
        return [
            'name' => ['required', 'string', 'max:255'],
            'email' => ['required', 'email', 'unique:users,email'],
            'password' => ['required', 'min:8', 'confirmed'],
        ];
    }
}
```

Symfony Controller:
```php
<?php
namespace App\Controller;

#[Route('/api/users')]
class UserController extends AbstractController
{
    #[Route('', methods: ['POST'])]
    public function create(EntityManagerInterface $em): JsonResponse
    {
        $user = new User();
        $em->persist($user);
        $em->flush();
        return $this->json($user, 201);
    }
}
```

---

## PHP 8.3 최신 기능

Readonly 클래스:
```php
readonly class UserDTO
{
    public function __construct(
        public int $id,
        public string $name,
        public string $email,
    ) {}
}
```

메서드가 있는 Enum:
```php
enum OrderStatus: string
{
    case Pending = 'pending';
    case Completed = 'completed';

    public function label(): string
    {
        return match($this) {
            self::Pending => 'Pending',
            self::Completed => 'Completed',
        };
    }
}
```

---

## Laravel 11 패턴

Eloquent Model:
```php
<?php
namespace App\Models;

class Post extends Model
{
    protected $fillable = ['title', 'content', 'user_id', 'status'];
    protected $casts = ['status' => PostStatus::class, 'published_at' => 'datetime'];

    public function user(): BelongsTo
    {
        return $this->belongsTo(User::class);
    }

    public function scopePublished($query)
    {
        return $query->where('status', PostStatus::Published);
    }
}
```

API Resource:
```php
<?php
namespace App\Http\Resources;

class PostResource extends JsonResource
{
    public function toArray(Request $request): array
    {
        return [
            'id' => $this->id,
            'title' => $this->title,
            'author' => new UserResource($this->whenLoaded('user')),
            'created_at' => $this->created_at->toIso8601String(),
        ];
    }
}
```

Migration:
```php
<?php
return new class extends Migration
{
    public function up(): void
    {
        Schema::create('posts', function (Blueprint $table) {
            $table->id();
            $table->foreignId('user_id')->constrained()->cascadeOnDelete();
            $table->string('title');
            $table->text('content');
            $table->timestamps();
        });
    }
};
```

Service Layer:
```php
<?php
namespace App\Services;

class UserService
{
    public function create(UserDTO $dto): User
    {
        return DB::transaction(function () use ($dto) {
            $user = User::create([
                'name' => $dto->name,
                'email' => $dto->email,
                'password' => Hash::make($dto->password),
            ]);
            return $user->load('profile');
        });
    }
}
```

---

## Symfony 7 패턴

Doctrine Entity:
```php
<?php
namespace App\Entity;

#[ORM\Entity(repositoryClass: UserRepository::class)]
class User
{
    #[ORM\Id]
    #[ORM\GeneratedValue]
    #[ORM\Column]
    private ?int $id = null;

    #[ORM\Column(length: 255)]
    #[Assert\NotBlank]
    private ?string $name = null;
}
```

Service with DI:
```php
<?php
namespace App\Service;

class UserService
{
    public function __construct(
        private readonly EntityManagerInterface $entityManager,
        private readonly UserPasswordHasherInterface $passwordHasher,
    ) {}

    public function createUser(string $email, string $password): User
    {
        $user = new User();
        $user->setEmail($email);
        $user->setPassword($this->passwordHasher->hashPassword($user, $password));
        $this->entityManager->persist($user);
        $this->entityManager->flush();
        return $user;
    }
}
```

---

## 테스트 패턴

PHPUnit (Laravel):
```php
<?php
class UserApiTest extends TestCase
{
    use RefreshDatabase;

    public function test_can_create_user(): void
    {
        $response = $this->postJson('/api/users', [
            'name' => 'John Doe',
            'email' => 'john@example.com',
            'password' => 'password123',
            'password_confirmation' => 'password123',
        ]);

        $response->assertStatus(201);
        $this->assertDatabaseHas('users', ['email' => 'john@example.com']);
    }
}
```

Pest (Laravel):
```php
<?php
it('can create a post', function () {
    $user = User::factory()->create();
    $response = $this->actingAs($user)
        ->postJson('/api/posts', ['title' => 'My Post', 'content' => 'Content']);
    $response->assertStatus(201);
    expect(Post::count())->toBe(1);
});
```

---

## 고급 패턴

Observers:
```php
<?php
class UserObserver
{
    public function creating(User $user): void { $user->uuid = Str::uuid(); }
    public function created(User $user): void { event(new UserCreated($user)); }
}
```

Query Scopes:
```php
public function scopePublished(Builder $query): Builder
{
    return $query->where('status', 'published')->whereNotNull('published_at');
}
// 사용: Post::published()->byAuthor($user)->get();
```

Queue Job:
```php
<?php
class ProcessUserData implements ShouldQueue
{
    use Dispatchable, InteractsWithQueue, Queueable, SerializesModels;
    public int $tries = 3;
    public int $backoff = 60;

    public function __construct(public User $user) {}
    public function handle(): void { $this->user->processData(); }
}
// Dispatch: ProcessUserData::dispatch($user)->onQueue('high');
```

Redis 캐싱:
```php
<?php
class ProductService
{
    public function getProduct(int $id): ?Product
    {
        return Cache::tags(['products'])->remember(
            "product:{$id}", now()->addHours(24), fn () => Product::find($id)
        );
    }
}
```

보안 미들웨어:
```php
<?php
class SecurityHeaders
{
    public function handle(Request $request, Closure $next)
    {
        $response = $next($request);
        $response->headers->set('X-Content-Type-Options', 'nosniff');
        $response->headers->set('X-Frame-Options', 'DENY');
        return $response;
    }
}
```

---

## Context7 라이브러리 매핑

```
/laravel/framework - Laravel 웹 프레임워크
/symfony/symfony - Symfony 컴포넌트
/doctrine/orm - Doctrine ORM
/phpunit/phpunit - PHP 테스트 프레임워크
/pestphp/pest - Pest 테스트 프레임워크
/laravel/sanctum - Laravel API 인증
```

---

## 관련 Skill

- `do-domain-backend` - REST API 및 마이크로서비스 아키텍처
- `do-domain-database` - SQL 패턴 및 ORM 최적화
- `do-workflow-testing` - TDD 및 테스트 전략

---

## 문제 해결

PHP 버전 확인:
```bash
php --version  # 8.3+ 필요
php -m | grep -E 'pdo|mbstring|openssl'
```

Composer 자동로드:
```bash
composer dump-autoload -o
composer clear-cache
```

Laravel 캐시:
```bash
php artisan config:clear && php artisan cache:clear && php artisan route:clear
```

Symfony 캐시:
```bash
php bin/console cache:clear && php bin/console cache:warmup
```

DB 연결 확인:
```php
try { DB::connection()->getPdo(); } catch (\Exception $e) { echo $e->getMessage(); }
```

Migration 롤백:
```bash
php artisan migrate:rollback --step=1
php bin/console doctrine:migrations:migrate prev
```

---

Last Updated: 2025-01-06
Status: Active (v1.0.0)
