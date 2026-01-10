---
name: do-lang-ruby
description: Ruby 3.3+ 개발 전문가 - Rails 7.2, ActiveRecord, Hotwire/Turbo, 현대적 Ruby 패턴 커버. Ruby API, 웹 애플리케이션, Rails 프로젝트 개발 시 사용.
version: 1.0.0
updated: 2026-01-06
status: active
allowed-tools: Read, Grep, Glob, Bash, mcp__context7__resolve-library-id, mcp__context7__query-docs
user-invocable: false
---

## 빠른 참고

Ruby 3.3+ 개발 전문가 - Rails 7.2, ActiveRecord, Hotwire/Turbo, RSpec, 현대적 Ruby 패턴.

자동 트리거: `.rb` 파일, `Gemfile`, `Rakefile`, `config.ru`, Rails/Ruby 관련 논의

핵심 역량:
- Ruby 3.3 기능: YJIT 프로덕션 지원, 패턴 매칭, Data 클래스, 엔드리스 메서드
- 웹 프레임워크: Rails 7.2 + Turbo, Stimulus, ActiveRecord
- 프론트엔드: Hotwire (Turbo + Stimulus)로 SPA 유사 경험 구현
- 테스팅: RSpec + 팩토리, 요청 스펙, 시스템 스펙
- 백그라운드 작업: Sidekiq + ActiveJob
- 패키지 관리: Bundler + Gemfile
- 코드 품질: RuboCop + Rails cops
- 데이터베이스: ActiveRecord + 마이그레이션, 관계, Scopes

---

## 빠른 패턴

Rails 컨트롤러:
```ruby
class UsersController < ApplicationController
  before_action :set_user, only: %i[show edit update destroy]

  def create
    @user = User.new(user_params)
    respond_to do |format|
      if @user.save
        format.html { redirect_to @user, notice: "Created." }
        format.turbo_stream
      else
        format.html { render :new, status: :unprocessable_entity }
      end
    end
  end

  private
  def set_user = @user = User.find(params[:id])
  def user_params = params.require(:user).permit(:name, :email)
end
```

ActiveRecord 모델:
```ruby
class User < ApplicationRecord
  has_many :posts, dependent: :destroy
  validates :email, presence: true, uniqueness: true
  scope :active, -> { where(active: true) }
  def full_name = "#{first_name} #{last_name}".strip
end
```

RSpec 테스트:
```ruby
RSpec.describe User, type: :model do
  it { is_expected.to validate_presence_of(:email) }
  describe "#full_name" do
    let(:user) { build(:user, first_name: "John", last_name: "Doe") }
    it { expect(user.full_name).to eq("John Doe") }
  end
end
```

---

## Ruby 3.3 신기능

YJIT (프로덕션 지원):
- Ruby 3.3에서 기본 활성화, Rails 앱 15-20% 성능 향상
- 활성화: `ruby --yjit` 또는 `RUBY_YJIT_ENABLE=1`

패턴 매칭:
```ruby
case response
in { status: "ok", data: }
  puts "Success: #{data}"
in { status: "error", message: }
  puts "Error: #{message}"
end
```

Data 클래스 (불변 구조체):
```ruby
User = Data.define(:name, :email) do
  def greeting = "Hello, #{name}!"
end
```

엔드리스 메서드:
```ruby
def add(a, b) = a + b
def positive?(n) = n > 0
```

---

## Rails 7.2 패턴

Gemfile 기본 설정:
```ruby
gem "rails", "~> 7.2.0"
gem "turbo-rails"
gem "stimulus-rails"
gem "sidekiq", "~> 7.0"

group :development, :test do
  gem "rspec-rails", "~> 7.0"
  gem "factory_bot_rails"
end
```

Concerns 활용:
```ruby
module Sluggable
  extend ActiveSupport::Concern
  included do
    before_validation :generate_slug, on: :create
    validates :slug, uniqueness: true
  end
  def to_param = slug
end
```

서비스 객체:
```ruby
class UserRegistrationService
  Result = Data.define(:success, :user, :errors)

  def call(params)
    user = User.new(params)
    ActiveRecord::Base.transaction do
      user.save!
      user.create_profile!
      UserMailer.welcome(user).deliver_later
    end
    Result.new(true, user, nil)
  rescue => e
    Result.new(false, nil, e.message)
  end
end
```

---

## Hotwire (Turbo + Stimulus)

Turbo 프레임:
```erb
<%= turbo_frame_tag "posts" do %>
  <% @posts.each { |post| render post } %>
<% end %>
```

Turbo 스트림:
```erb
<%= turbo_stream.prepend "posts", @post %>
<%= turbo_stream.update "form", partial: "form" %>
```

Stimulus 컨트롤러:
```javascript
import { Controller } from "@hotwired/stimulus"
export default class extends Controller {
  static targets = ["input", "submit"]
  validate() {
    this.submitTarget.disabled = !this.inputTargets.every(i => i.value.length > 0)
  }
}
```

---

## RSpec 테스팅

Factory Bot:
```ruby
FactoryBot.define do
  factory :user do
    sequence(:email) { |n| "user#{n}@example.com" }
    trait :with_posts do
      after(:create) { |u| create_list(:post, 3, user: u) }
    end
  end
end
```

모델 스펙:
```ruby
RSpec.describe User, type: :model do
  it { is_expected.to have_many(:posts).dependent(:destroy) }
  it { is_expected.to validate_uniqueness_of(:email).case_insensitive }
end
```

요청 스펙:
```ruby
describe "POST /posts" do
  it "creates a post" do
    expect { post posts_path, params: { post: attributes_for(:post) } }
      .to change(Post, :count).by(1)
  end
end
```

---

## 고급 패턴

Sidekiq 작업:
```ruby
class ProcessOrderJob < ApplicationJob
  queue_as :default
  retry_on ActiveRecord::Deadlocked, attempts: 3

  def perform(order_id)
    Order.find(order_id).tap do |order|
      order.process!
      OrderMailer.confirmation(order).deliver_later
    end
  end
end
```

Query 객체:
```ruby
class PostSearchQuery
  def initialize(relation = Post.all) = @relation = relation

  def call(params)
    @relation
      .then { filter_by_status(_1, params[:status]) }
      .then { search_by_title(_1, params[:query]) }
  end
end
```

다형성 관계:
```ruby
class Comment < ApplicationRecord
  belongs_to :commentable, polymorphic: true
end

class Post < ApplicationRecord
  has_many :comments, as: :commentable
end
```

---

## 성능 최적화

N+1 쿼리 방지:
```ruby
# Bad
User.all.each { |u| puts u.posts.count }

# Good
User.includes(:posts).each { |u| puts u.posts.size }
```

Counter Cache:
```ruby
belongs_to :user, counter_cache: true
```

인덱스 추가:
```ruby
add_index :posts, [:status, :published_at]
```

---

## 보안 모범 사례

Strong Parameters:
```ruby
params.require(:user).permit(:name, :email)
```

SQL 인젝션 방지:
```ruby
# Bad
User.where("email = '#{params[:email]}'")

# Good
User.where(email: params[:email])
```

---

## 트러블슈팅

버전 확인:
```bash
ruby --version       # 3.3+ 필요
rails --version      # 7.2+ 필요
ruby -e "puts RubyVM::YJIT.enabled?"
```

데이터베이스 문제:
- `config/database.yml` 확인
- `rails db:create db:migrate` 실행

RSpec 설정:
```bash
rails generate rspec:install
bundle exec rspec --format documentation
```

---

## Context7 라이브러리

```
/rails/rails - Ruby on Rails
/rspec/rspec - RSpec 테스팅
/hotwired/turbo-rails - Turbo
/hotwired/stimulus-rails - Stimulus
/sidekiq/sidekiq - 백그라운드 작업
/rubocop/rubocop - 스타일 가이드
/thoughtbot/factory_bot - 테스트 팩토리
```

---

## 관련 Skills

- `do-domain-backend` - REST API 및 웹 애플리케이션 아키텍처
- `do-domain-database` - SQL 패턴 및 ActiveRecord 최적화
- `do-workflow-testing` - TDD 및 테스팅 전략
- `do-foundation-quality` - 품질 원칙

---

Last Updated: 2026-01-06
Status: Active (v1.0.0)
