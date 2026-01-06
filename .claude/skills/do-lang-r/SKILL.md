---
name: do-lang-r
description: R 4.4+ development specialist covering tidyverse, ggplot2, Shiny, and data science patterns. Use when developing data analysis pipelines, visualizations, or Shiny applications.
version: 1.0.0
updated: 2025-12-07
status: active
allowed-tools: Read, Grep, Glob, Bash, mcp__context7__resolve-library-id, mcp__context7__query-docs
---

## 빠른 참조

R 4.4+ 개발 전문가 - tidyverse, ggplot2, Shiny, renv, 모던 R 패턴 지원.

자동 트리거: `.R` 파일, `.Rmd`, `.qmd`, `DESCRIPTION`, `renv.lock`, Shiny/ggplot2 관련 작업

핵심 역량:
- R 4.4 기능: Native pipe |>, lambda 문법 \(x), 개선된 오류 메시지
- 데이터 조작: dplyr, tidyr, purrr, stringr, forcats
- 시각화: ggplot2, plotly, scales, patchwork
- 웹 애플리케이션: Shiny, reactivity, modules, bslib
- 테스팅: testthat 3.0, 스냅샷 테스팅, 모킹
- 패키지 관리: renv, pak, DESCRIPTION
- 재현 가능 리포트: R Markdown, Quarto
- 데이터베이스: DBI, dbplyr, pool

---

## R 4.4 모던 기능

### Native Pipe Operator |>

```r
result <- data |>
  filter(!is.na(value)) |>
  mutate(log_value = log(value)) |>
  summarise(mean_log = mean(log_value))

# Placeholder _ 사용 (첫 번째 인자가 아닌 경우)
data |>
  lm(formula = y ~ x, data = _)
```

### Lambda 문법

```r
map(data, \(x) x^2)
map2(list1, list2, \(x, y) x + y)

# dplyr 컨텍스트에서
data |>
  mutate(across(where(is.numeric), \(x) scale(x)[,1]))
```

---

## tidyverse 데이터 조작

### dplyr 핵심 동사

```r
library(dplyr)

processed <- raw_data |>
  filter(status == "active", amount > 0) |>
  select(id, date, amount, category) |>
  mutate(
    month = floor_date(date, "month"),
    amount_scaled = amount / max(amount)
  ) |>
  arrange(desc(date))

# group_by와 summarise
summary <- processed |>
  group_by(category, month) |>
  summarise(
    n = n(),
    total = sum(amount),
    avg = mean(amount),
    .groups = "drop"
  )

# across로 여러 열 처리
data |>
  mutate(across(starts_with("price"), \(x) round(x, 2)))
```

### tidyr 재형성

```r
library(tidyr)

# pivot_longer (넓은 형식 > 긴 형식)
wide_data |>
  pivot_longer(
    cols = starts_with("year_"),
    names_to = "year",
    names_prefix = "year_",
    values_to = "value"
  )

# pivot_wider (긴 형식 > 넓은 형식)
long_data |>
  pivot_wider(
    names_from = category,
    values_from = value,
    values_fill = 0
  )
```

### purrr 함수형 프로그래밍

```r
library(purrr)

files |> map(\(f) read_csv(f))
files |> map_dfr(\(f) read_csv(f), .id = "source")
values |> map_dbl(\(x) mean(x, na.rm = TRUE))

# safely로 오류 처리
safe_read <- safely(read_csv)
results <- files |> map(safe_read)
successes <- results |> map("result") |> compact()
```

---

## ggplot2 시각화

### 완전한 플롯 구조

```r
library(ggplot2)
library(scales)

p <- ggplot(data, aes(x = x, y = y, color = group)) +
  geom_point(alpha = 0.7, size = 3) +
  geom_smooth(method = "lm", se = TRUE) +
  scale_x_continuous(labels = comma) +
  scale_y_log10(labels = dollar) +
  scale_color_brewer(palette = "Set2") +
  facet_wrap(~ category, scales = "free_y") +
  labs(
    title = "Analysis Title",
    subtitle = "Descriptive subtitle",
    x = "X Axis Label",
    y = "Y Axis Label"
  ) +
  theme_minimal(base_size = 12) +
  theme(legend.position = "bottom")

ggsave("output.png", p, width = 10, height = 6, dpi = 300)
```

### patchwork로 다중 플롯

```r
library(patchwork)

p1 <- ggplot(data, aes(x)) + geom_histogram()
p2 <- ggplot(data, aes(x, y)) + geom_point()
p3 <- ggplot(data, aes(group, y)) + geom_boxplot()

combined <- (p1 | p2) / p3 +
  plot_annotation(title = "Combined Analysis", tag_levels = "A")
```

### 커스텀 테마

```r
theme_custom <- function(base_size = 12, base_family = "") {
  theme_minimal(base_size = base_size, base_family = base_family) %+replace%
    theme(
      plot.title = element_text(face = "bold", size = rel(1.2), hjust = 0),
      plot.subtitle = element_text(color = "grey40", hjust = 0),
      panel.grid.major = element_line(color = "grey90"),
      panel.grid.minor = element_blank(),
      axis.title = element_text(face = "bold", size = rel(0.9)),
      legend.position = "bottom",
      strip.background = element_rect(fill = "grey95", color = NA),
      strip.text = element_text(face = "bold", size = rel(0.9))
    )
}
```

---

## Shiny 애플리케이션

### 기본 앱 구조

```r
library(shiny)

ui <- fluidPage(
  selectInput("var", "Variable:", choices = names(mtcars)),
  plotOutput("plot")
)

server <- function(input, output, session) {
  output$plot <- renderPlot({
    ggplot(mtcars, aes(.data[[input$var]])) +
      geom_histogram()
  })
}

shinyApp(ui, server)
```

### 모듈 패턴

```r
dataFilterUI <- function(id) {
  ns <- NS(id)
  tagList(
    selectInput(ns("category"), "Category:", choices = NULL),
    sliderInput(ns("range"), "Range:", min = 0, max = 100, value = c(0, 100))
  )
}

dataFilterServer <- function(id, data) {
  moduleServer(id, function(input, output, session) {
    observe({
      categories <- unique(data()$category)
      updateSelectInput(session, "category", choices = categories)
    })

    reactive({
      req(input$category)
      data() |>
        filter(
          category == input$category,
          value >= input$range[1],
          value <= input$range[2]
        )
    })
  })
}
```

### Reactive 패턴

```r
server <- function(input, output, session) {
  # reactive: 캐시된 계산
  processed_data <- reactive({
    raw_data() |>
      filter(year == input$year)
  })

  # reactiveVal: 변경 가능한 상태
  counter <- reactiveVal(0)
  observeEvent(input$increment, {
    counter(counter() + 1)
  })

  # eventReactive: 특정 이벤트에 트리거
  analysis <- eventReactive(input$run_analysis, {
    expensive_computation(processed_data())
  })

  # debounce: 빠른 입력 처리
  search_term <- reactive(input$search) |> debounce(300)
}
```

### 비동기 처리

```r
library(shiny)
library(promises)
library(future)
plan(multisession)

server <- function(input, output, session) {
  output$result <- renderText({
    future({
      Sys.sleep(5)
      expensive_calculation()
    }) %...>%
      as.character()
  })
}
```

### 캐싱

```r
library(memoise)

expensive_function <- memoise(function(data) {
  heavy_processing(data)
})

# Shiny에서
server <- function(input, output, session) {
  cached_data <- reactive({
    expensive_function(input$data_source)
  }) |> bindCache(input$data_source)
}
```

---

## testthat 테스팅

### 테스트 구조

```r
library(testthat)

test_that("calculate_growth returns correct values", {
  data <- tibble(year = 2020:2022, value = c(100, 110, 121))

  result <- calculate_growth(data)

  expect_equal(nrow(result), 3)
  expect_equal(result$growth[2], 0.1, tolerance = 0.001)
  expect_true(is.na(result$growth[1]))
})

test_that("calculate_growth handles edge cases", {
  expect_error(calculate_growth(NULL), "data cannot be NULL")
})
```

### 테스트 픽스처

```r
# tests/testthat/helper.R
setup_test_db <- function() {
  conn <- DBI::dbConnect(RSQLite::SQLite(), ":memory:")
  DBI::dbExecute(conn, "CREATE TABLE users (id INTEGER, name TEXT)")
  conn
}

teardown_test_db <- function(conn) {
  DBI::dbDisconnect(conn)
}

# tests/testthat/test-db.R
test_that("database operations work correctly", {
  conn <- setup_test_db()
  on.exit(teardown_test_db(conn))

  DBI::dbExecute(conn, "INSERT INTO users VALUES (1, 'Test')")
  result <- DBI::dbGetQuery(conn, "SELECT * FROM users")

  expect_equal(nrow(result), 1)
  expect_equal(result$name, "Test")
})
```

---

## 데이터베이스 통합

### dbplyr와 pool

```r
library(DBI)
library(pool)
library(dbplyr)

# 연결 풀 생성
pool <- dbPool(
  drv = RPostgres::Postgres(),
  dbname = "mydb",
  host = "localhost",
  user = Sys.getenv("DB_USER"),
  password = Sys.getenv("DB_PASSWORD"),
  minSize = 1,
  maxSize = 5
)

# dbplyr로 쿼리
users_db <- tbl(pool, "users")

result <- users_db |>
  filter(active == TRUE) |>
  group_by(department) |>
  summarise(
    count = n(),
    avg_salary = mean(salary, na.rm = TRUE)
  ) |>
  arrange(desc(count)) |>
  collect()

# 종료 시 풀 닫기
onStop(function() {
  poolClose(pool)
})
```

### 트랜잭션 처리

```r
with_transaction <- function(pool, expr) {
  conn <- poolCheckout(pool)
  on.exit(poolReturn(conn))

  dbBegin(conn)

  tryCatch({
    result <- expr
    dbCommit(conn)
    result
  }, error = function(e) {
    dbRollback(conn)
    stop(e)
  })
}
```

---

## renv 의존성 관리

```r
renv::init()
renv::install("tidyverse")
renv::install("shiny")
renv::snapshot()
renv::restore()
```

---

## 성능 최적화

### data.table 대용량 데이터

```r
library(data.table)

dt <- as.data.table(large_df)

# 빠른 그룹화
result <- dt[, .(
  count = .N,
  mean_value = mean(value, na.rm = TRUE),
  max_value = max(value, na.rm = TRUE)
), by = .(category, year)]

# 인플레이스 업데이트
dt[, new_col := value * 2]

# 효율적 조인
dt1[dt2, on = .(key_col)]

# 롤링 연산
dt[, rolling_mean := frollmean(value, n = 7), by = category]
```

### 병렬 처리

```r
library(future)
library(furrr)

plan(multisession, workers = 4)

results <- future_map(file_list, \(f) {
  read_csv(f) |>
    process_data()
}, .progress = TRUE)

# 오류 처리와 함께
safe_process <- possibly(process_data, otherwise = NULL)
results <- future_map(data_list, safe_process)

plan(sequential)
```

---

## 패키지 개발

### Roxygen 문서화

```r
#' Calculate Growth Rate
#'
#' Calculates the period-over-period growth rate for a numeric vector.
#'
#' @param x A numeric vector of values.
#' @param periods Number of periods for growth calculation. Default is 1.
#' @param na.rm Logical. Should NA values be removed? Default is TRUE.
#'
#' @return A numeric vector of growth rates.
#'
#' @examples
#' calculate_growth(c(100, 110, 121))
#' calculate_growth(c(100, 110, 121), periods = 2)
#'
#' @export
calculate_growth <- function(x, periods = 1, na.rm = TRUE) {
  if (!is.numeric(x)) {
    stop("x must be numeric")
  }

  growth <- (x - dplyr::lag(x, n = periods)) / dplyr::lag(x, n = periods)

  if (na.rm) {
    growth[is.na(growth)] <- NA_real_
  }

  growth
}
```

---

## 프로덕션 배포

### Docker

```dockerfile
FROM rocker/shiny:4.4.0

RUN apt-get update && apt-get install -y \
    libcurl4-gnutls-dev \
    libssl-dev \
    libpq-dev \
    && rm -rf /var/lib/apt/lists/*

COPY renv.lock renv.lock
RUN R -e "install.packages('renv'); renv::restore()"

COPY . /srv/shiny-server/
RUN chown -R shiny:shiny /srv/shiny-server

EXPOSE 3838

CMD ["/usr/bin/shiny-server"]
```

### Posit Connect 배포

```r
library(rsconnect)

rsconnect::setAccountInfo(
  name = "your-account",
  token = Sys.getenv("CONNECT_TOKEN"),
  secret = Sys.getenv("CONNECT_SECRET")
)

rsconnect::deployApp(
  appDir = ".",
  appName = "my-shiny-app",
  appTitle = "My Shiny Application",
  forceUpdate = TRUE
)
```

---

## 오류 처리

### Condition 시스템

```r
# 커스텀 조건 정의
validation_error <- function(message, field = NULL) {
  rlang::abort(
    message,
    class = "validation_error",
    field = field
  )
}

# 조건 처리
process_input <- function(data) {
  tryCatch(
    {
      validate_data(data)
      transform_data(data)
    },
    validation_error = function(e) {
      cli::cli_alert_danger("Validation failed: {e$message}")
      cli::cli_alert_info("Field: {e$field}")
      NULL
    },
    error = function(e) {
      cli::cli_alert_danger("Unexpected error: {e$message}")
      rlang::abort("Processing failed", parent = e)
    }
  )
}
```

---

## 문제 해결

R 버전 확인:
```r
R.version.string  # 4.4+ 권장
packageVersion("dplyr")
```

Native Pipe 동작 안함:
- R 버전 4.1+ 확인
- RStudio 설정: Tools > Global Options > Code > Use native pipe

renv 문제:
```r
renv::clean()
renv::rebuild()
renv::snapshot(force = TRUE)
```

Shiny Reactivity 디버그:
```r
options(shiny.reactlog = TRUE)
reactlog::reactlog_enable()
shiny::reactlogShow()
```

ggplot2 폰트 문제:
```r
library(showtext)
font_add_google("Roboto", "roboto")
showtext_auto()
```

---

## Context7 라이브러리

```
/tidyverse/dplyr - 데이터 조작 동사
/tidyverse/ggplot2 - Grammar of Graphics 시각화
/tidyverse/purrr - 함수형 프로그래밍 도구
/tidyverse/tidyr - 데이터 정리 함수
/rstudio/shiny - 웹 애플리케이션 프레임워크
/r-lib/testthat - 단위 테스팅 프레임워크
/rstudio/renv - 의존성 관리
```

---

## 연결 스킬

- `do-lang-python` - Python/R 상호운용 (reticulate)
- `do-domain-database` - SQL 패턴 및 데이터베이스 최적화
- `do-workflow-testing` - TDD 및 테스팅 전략
- `do-foundation-quality` - TRUST 5 품질 원칙

---

Last Updated: 2025-12-07
Status: Active (v1.0.0)
