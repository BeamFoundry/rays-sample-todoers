# Todoers &mdash; Rays Framework Sample Application

[![Go Version](https://img.shields.io/badge/Go-1.27+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Framework](https://img.shields.io/badge/Framework-BeamFoundry%2FRays-6f42c1?style=flat)](https://github.com/BeamFoundry/rays)
[![Web](https://img.shields.io/badge/Web-Fiber%20v3-00ACD7?style=flat)](https://github.com/gofiber/fiber)
[![Frontend](https://img.shields.io/badge/Frontend-HTMX%20%2B%20Bootstrap%205-336699?style=flat)](https://htmx.org)
[![Database](https://img.shields.io/badge/ORM-GORM%20%2B%20SQLite-E9573F?style=flat)](https://gorm.io)

**Todoers** is a full-featured, multi-list task management application built to demonstrate the capabilities and design paradigms of the **[Rays framework](https://github.com/BeamFoundry/rays)** by BeamFoundry.

Rays brings proven enterprise architectural patterns &mdash; reminiscent of Spring Boot and Java EE &mdash; natively into Go. It provides an Inversion of Control (IoC) container, declarative Stereotype scanning, compile-time validation, automatic dependency injection, profile-based configuration, Spring Data-style ORM repositories, and declarative web security.

---

## Table of Contents

- [Features](#features)
- [Architecture & Tech Stack](#architecture--tech-stack)
- [Core Rays Concepts Demonstrated](#core-rays-concepts-demonstrated)
  - [1. Stereotypes & Beams](#1-stereotypes--beams)
  - [2. The `stereotyper` Code Generator](#2-the-stereotyper-code-generator)
  - [3. Dependency Injection](#3-dependency-injection)
  - [4. Conditionals & Profiles](#4-conditionals--profiles)
  - [5. Declarative ORM Repositories (`rays-orm`)](#5-declarative-orm-repositories-rays-orm)
  - [6. Declarative Web Security (`rays-web-security`)](#6-declarative-web-security-rays-web-security)
  - [7. Web Controllers & HTMX (`rays-fiber`)](#7-web-controllers--htmx-rays-fiber)
  - [8. Internationalization (`i18n`)](#8-internationalization-i18n)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Building](#building)
  - [Running the Application](#running-the-application)
- [Configuration & CLI Overrides](#configuration--cli-overrides)
  - [Command-line Flags](#command-line-flags)
  - [Environment Variables](#environment-variables)
  - [Active Profiles](#active-profiles)
- [Code Walkthrough](#code-walkthrough)
  - [Defining a Service with Transactions](#defining-a-service-with-transactions)
  - [Defining a Derived Query Repository](#defining-a-derived-query-repository)
  - [Configuring Web Security](#configuring-web-security)
  - [Handling Routes with FiberController](#handling-routes-with-fibercontroller)
- [License](#license)

---

## Features

- **Multi-List Kanban Board**: Create, rename, delete, and manage multiple todo lists.
- **Drag-and-Drop Reordering**: Interactive reordering of lists and moving items across lists powered by [Sortable.js](https://sortablejs.github.io/Sortable/) and [HTMX](https://htmx.org).
- **Inline Editing & Instant Updates**: Dynamic CRUD operations without full page reloads.
- **Authentication & Authorization**:
  - Secure form-based authentication with CSRF protection.
  - Persistent "Remember Me" token authentication.
  - Role-Based Access Control (`ROLE_USER`, `ROLE_ADMIN`).
  - Automatic admin promotion: the first registered user is automatically granted `ROLE_ADMIN`.
- **User Profile Management**: Update display names and change passwords with server-side validation.
- **Admin Dashboard**: Live search of registered users with pagination (`ROLE_ADMIN` only).
- **Internationalization (i18n)**: Instant switching between English (`en`) and French (`fr`) with localized messages and validation feedback.
- **Cache-Busted Static Assets**: Asset hashing via `hashfs` for seamless client-side caching.
- **Pure Go & CGO Dual Support**: Configured to run either completely without CGO (`glebarez/sqlite`) or with native CGO SQLite (`mattn/go-sqlite3`).

---

## Architecture & Tech Stack

```
                     ┌──────────────────────────────────────────────┐
                     │             Browser / Client                 │
                     │  (Bootstrap 5 + HTMX + Sortable.js + CSS)    │
                     └──────────────────────┬───────────────────────┘
                                            │ HTTP / REST / HTML Fragments
                                            ▼
                     ┌──────────────────────────────────────────────┐
                     │       Fiber v3 Web Engine (rays-fiber)       │
                     │   - FiberControllers & RequestMapping        │
                     │   - HttpSecurity Filters & SecurityContext   │
                     │   - Error Handler & Locale Resolver          │
                     └──────────────────────┬───────────────────────┘
                                            │ Dependency Injection (@:"Inject")
                                            ▼
                     ┌──────────────────────────────────────────────┐
                     │          Services & Business Logic           │
                     │      (TodoService, UserService, etc.)        │
                     │   - Declarative GormTransactional Tx         │
                     └──────────────────────┬───────────────────────┘
                                            │
                                            ▼
                     ┌──────────────────────────────────────────────┐
                     │        ORM Layer (rays-orm + GORM)           │
                     │   - Generic Repositories: OrmRepository[T,ID]│
                     │   - Method-name derived queries              │
                     └──────────────────────┬───────────────────────┘
                                            │
                                            ▼
                     ┌──────────────────────────────────────────────┐
                     │            SQLite Database File              │
                     │   (Auto-migrated by DbMigrator at startup)   │
                     └──────────────────────────────────────────────┘
```

- **Core Framework**: [BeamFoundry/rays](https://github.com/BeamFoundry/rays)
- **Web Layer**: [BeamFoundry/rays-fiber](https://github.com/BeamFoundry/rays-fiber) on top of [Fiber v3](https://github.com/gofiber/fiber)
- **Data Layer**: [BeamFoundry/rays-orm](https://github.com/BeamFoundry/rays-orm) on top of [GORM](https://gorm.io)
- **Security**: [BeamFoundry/rays-web-security](https://github.com/BeamFoundry/rays-web-security)
- **Templates**: Go `html/template` with [Sprig](https://github.com/Masterminds/sprig) functions
- **Static Asset Serving**: [HashFS](https://github.com/benbjohnson/hashfs) for content-hashed asset URLs

---

## Core Rays Concepts Demonstrated

### 1. Stereotypes & Beams

In Rays, components managed by the IoC container are called **Beams** (analogous to Beans in Spring). Structs declare their role by embedding one of the framework's **Stereotypes**:

| Stereotype | Equivalent | Purpose |
| :--- | :--- | :--- |
| `Component` | `@Component` | General-purpose managed component |
| `Service` | `@Service` | Service layer holding business logic |
| `Repository` | `@Repository` | Data access object interface / implementation |
| `Controller` | `@Controller` | Web request handler / endpoint controller |
| `Configuration` | `@Configuration` | Producer of Beams through factory methods |
| `ConfigurationProperties` | `@ConfigurationProperties` | Binds external YAML/environment configuration to typed structs |

### 2. The `stereotyper` Code Generator

Rays provides a dedicated code generation tool:

```bash
go tool stereotyper
```

The tool scans the packages and generates `zz_stereotype.go` in each directory. These files contain compile-time interface assertions (e.g. `var _, _ = any(&DbConfig{}).(_rays_core.Stereotype)`), ensuring that all annotated structs correctly satisfy framework constraints at build time.

### 3. Dependency Injection

Dependencies are injected into struct fields or factory method arguments using the `@:"Inject"` struct tag:

```go
type TodoServiceImpl struct {
    Service
    transactional  GormTransactional `@:"Inject"`
    todoRepository TodoRepository    `@:"Inject"`
    itemRepository ItemRepository    `@:"Inject"`
    messageSource  MessageSource     `@:"Inject"`
}
```

Rays resolves and wires references automatically during `ctx.Initialize()`.

### 4. Conditionals & Profiles

Rays supports conditional Beam creation:

- `ConditionalOnMissingBeam`: Instantiates a fallback Beam only if no other Beam of that type is registered.
- `ConditionalOnProfile`: Activates a Beam only when a specific profile (such as `test`) is enabled.
- `ConditionalOnNoCGO` / `ConditionalOnCGO`: Selects different implementations depending on whether the binary was compiled with CGO enabled or disabled.

For instance, in `pkg/configurations/database.go`, the application switches between native SQLite and pure-Go SQLite based on the build target:

```go
func (this *DialectorConfig) DBOpener(_ struct { ConditionalOnMissingBeam }) DBOpener {
    return &SQLiteDB{}       // CGO SQLite
}

func (this *DialectorConfig) DBOpenerNoCGO(_ struct { ConditionalOnNoCGO }) DBOpener {
    return &SQLiteDBNoCGO{}  // Pure-Go SQLite
}
```

### 5. Declarative ORM Repositories (`rays-orm`)

Similar to Spring Data JPA, `rays-orm` provides generic repositories (`OrmRepository[T, ID]`) where query functions are derived automatically from the method names declared on the interface:

```go
type ItemRepository interface {
    OrmRepositoryInterface[Item, uint]
    FindAllByTodoID(todoId uint) []Item
    FindAllByTodoIDOrderByRank(todoId uint) []Item
    FindAllByTodoIDAndRankGreaterThanOrEqual(uint, uint) []Item
    FindOneByTodoIDAndItemID(uint, uint) Optional[*Item]
}
```

No boilerplate SQL or GORM query chains need to be written by hand for standard queries.

### 6. Declarative Web Security (`rays-web-security`)

Web security is configured in a fluent, expressive DSL using `HttpSecurityBuilder`:

```go
func (this *SecurityConfig) HttpSecurity(that struct { rememberme RememberMeService `@:"Inject"` }) *HttpSecurity {
    return (&HttpSecurityBuilder{}).
        AuthorizedHTTPRequests().
        Paths("/profile").Authenticated().
        Paths("/todolists**").HasAuthority("ROLE_USER").
        Paths("/admin", "/admin/**").HasAuthority("ROLE_ADMIN").
        Paths("/static/**").PermitAll().
        Paths("/*").PermitAll().
        And().
        AnyRequest().DenyAll().
        FormLogin().
        ViewName("views/login").
        And().
        RememberMe(that.rememberme).
        Build()
}
```

### 7. Web Controllers & HTMX (`rays-fiber`)

Controllers embed `FiberController` and use `RequestMapping` struct tags to define route patterns and HTTP methods:

```go
type TodolistsController struct {
    FiberController
    RequestMapping `Path:"/todolists"`
    AppSecure
    todoService TodoService `@:"Inject"`
}

func (this *TodolistsController) ListAll(
    _ struct { RequestMapping `Method:"GET"` },
    c fiber.Ctx,
) error {
    user := fiber.Locals[*User](c, "User")
    return c.Render("views/fragments/todolist", fiber.Map{
        "Todos": this.todoService.GetUserTodoLists(user),
    }, "views/fragments/todoitem")
}
```

Handlers return HTML template fragments that HTMX swaps directly into the DOM, giving a single-page app feel with zero JavaScript framework overhead.

### 8. Internationalization (`i18n`)

Application messages and validation error translations are loaded from embedded YAML resource files (`resources/config/translation/`):

- Default locale: `en` (English)
- Alternate locale: `fr` (French)
- Switched via query parameter (e.g. `?lang=fr` or `?lang=en`) handled by `LocaleResolver`.

---

## Project Structure

```
rays-sample-todoers/
├── cmd/
│   └── main.go                  # Bootstrap, CLI flags, configuration loading, container execution
├── internal/
│   ├── configurations/          # Configuration beams (DB, Fiber, Security, i18n, Sprig, Static FS)
│   ├── repositories/            # rays-orm derived repositories (TodoRepository, ItemRepository, UserRepository)
│   └── services/                # Business logic & transactions (TodoService, UserService)
├── pkg/
│   ├── controllers/             # FiberControllers (Todolists, Page, Register, Admin, Static)
│   │   └── handlers/            # Security & locale route middleware (AppSecure, AppSecureTodo, etc.)
│   ├── core/                    # Services and Repository interfaces.
│   └── domains/                 # Domain entities (User, Todo, Item)
├── resources/
│   ├── config/
│   │   ├── rays.yaml            # Default application configuration (logging, server port, etc.)
│   │   └── translation/         # i18n translation bundles for en and fr
│   ├── static/                  # Static assets: CSS, JavaScript, favicons
│   ├── views/                   # HTML templates, layouts, and HTMX fragments
│   └── resources.go             # Go embed.FS wrappers for resources
├── Makefile                     # Build targets (all, tool, cgo, nocgo)
├── todolist.db                  # Default SQLite database
├── go.mod & go.sum              # Module definitions and dependencies
└── README.md                    # Project documentation
```

---

## Getting Started

### Prerequisites

- **Go**: Version `1.26` or later (or appropriate Go toolchain)
- **Make** (optional, for running Makefile targets)
- **GCC / C compiler** (optional, only if building with CGO SQLite)

### Building

The project provides a `Makefile` with targets for both pure Go and CGO builds:

```bash
# Pure Go build (no CGO required, uses glebarez/sqlite)
make nocgo

# Alternatively, build with CGO enabled (uses mattn/go-sqlite3)
make cgo
```

The Makefile automatically runs `go tool stereotyper` before compiling the binary to ensure stereotype assertions are up to date.

To run tests:

```bash
go test ./...
```

### Running the Application

Start the compiled binary:

```bash
./todolist
```

By default, the server listens on **`http://localhost:3000`**.

1. Navigate to `http://localhost:3000` in your web browser.
2. Click **Register** to create an account.
   > **Note**: The first registered account is automatically assigned the `ROLE_ADMIN` authority and has access to the `/admin` user management panel.
3. Start creating todo lists, adding items, editing them, and dragging them around to reorder!

---

## Configuration & CLI Overrides

Configuration is driven by `resources/config/rays.yaml` embedded directly inside the binary. You can inspect or customize settings via command-line flags or environment variables without recompiling.

### Command-line Flags

```bash
Usage of ./todolist:
  -c, --config strings            Additional configuration directory paths (io/fs.FS)
  -p, --profile strings           Profiles to activate (comma-separated or multiple flags)
      --override stringToString   Override configuration keys inline (key=value)
```

#### Examples

Change the listening port to `:8080`:

```bash
./todolist --override rays.fiber.servers.primary.address=":8080"
```

Activate the `test` profile (uses an in-memory SQLite database):

```bash
./todolist -p test
```

Specify additional configuration directories:

```bash
./todolist -c /etc/todoers/config
```

### Environment Variables

Rays automatically scans environment variables starting with `RAYS_` or `APP_`. Variables are converted to lowercase and underscores are translated to dots:

```bash
# Overrides rays.fiber.servers.primary.address=":9000"
export RAYS_FIBER_SERVERS_PRIMARY_ADDRESS=":9000"

# Overrides rays.logging.level="DEBUG"
export RAYS_LOGGING_LEVEL="DEBUG"

./todolist
```

### Active Profiles

Configurations follow Spring Boot-style naming conventions:
- `rays.yaml` &mdash; loaded by default.
- `rays-<profile>.yaml` &mdash; loaded when `<profile>` is active (e.g. `--profile test`).

---

## Code Walkthrough

### Defining a Service with Transactions

Services embed `Service` and use `@:"Inject"` to declare dependencies. Declarative transactions are managed via `GormTransactional`:

```go
package services

import (
    "context"
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays-orm/pkg/core"
    . "github.com/BeamFoundry/rays-sample-todoers/pkg/domains"
    . "github.com/BeamFoundry/rays-sample-todoers/internal/repositories"
)

type TodoServiceImpl struct {
    Service
    transactional  GormTransactional `@:"Inject"`
    todoRepository TodoRepository    `@:"Inject"`
    itemRepository ItemRepository    `@:"Inject"`
}

func (this *TodoServiceImpl) SaveList(todo *Todo) (*Todo, error) {
    tx := this.transactional.Start(context.Background())
    defer this.transactional.End()

    t, err := this.todoRepository.Save(todo)
    if err != nil {
        tx.Rollback()
        return nil, err
    }

    tx.Commit()
    return t, nil
}
```

### Defining a Derived Query Repository

Define your repository interface by embedding `OrmRepositoryInterface[Model, IDType]` and declaring query method signatures:

```go
package repositories

import (
    . "github.com/BeamFoundry/rays-orm/pkg/core"
    . "github.com/BeamFoundry/rays-sample-todoers/pkg/domains"
    . "github.com/BeamFoundry/rays/pkg/lang"
)

type UserRepository interface {
    OrmRepositoryInterface[User, uint]
    FindByUsername(string) Optional[*User]
    FindAllUsersByUsername(string, ...OrmQueryFunc) []User
}

type UserRepositoryImpl struct {
    OrmRepository[User, uint]
    findByUsername                      func(string) Optional[*User]
    findAllUsersByUsernameLikeOrderByID func(string, ...OrmQueryFunc) []User
}

func (this *UserRepositoryImpl) FindByUsername(username string) Optional[*User] {
    return this.findByUsername(username)
}
```

### Configuring Web Security

Use `HttpSecurityBuilder` inside a `Configuration` beam to secure URL patterns:

```go
package configurations

import (
    . "github.com/BeamFoundry/rays-web-security/pkg/config"
    . "github.com/BeamFoundry/rays-web-security/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/core"
)

type SecurityConfig struct {
    Configuration
}

func (this *SecurityConfig) HttpSecurity(that struct {
    rememberme RememberMeService `@:"Inject"`
}) *HttpSecurity {
    return (&HttpSecurityBuilder{}).
        AuthorizedHTTPRequests().
        Paths("/profile").Authenticated().
        Paths("/todolists**").HasAuthority("ROLE_USER").
        Paths("/admin", "/admin/**").HasAuthority("ROLE_ADMIN").
        Paths("/static/**", "/*").PermitAll().
        And().
        AnyRequest().DenyAll().
        FormLogin().ViewName("views/login").
        And().
        RememberMe(that.rememberme).
        Build()
}
```

### Handling Routes with FiberController

Controller endpoints receive a typed struct with `RequestMapping` tags and the Fiber context:

```go
type AdminController struct {
    FiberController
    RequestMapping `Path:"/admin"`
    AppSecure
    userService UserService `@:"Inject"`
}

func (this *AdminController) AdminListUsers(
    _ struct { RequestMapping `Path:"/users"` },
    ctx fiber.Ctx,
) error {
    page := fiber.Query(ctx, "page", 1)
    query := ctx.Query("q", "")
    users := this.userService.FindAllUsersByUsername("%" + query + "%", page, 10)

    return ctx.Render("views/fragments/admin-user-list", fiber.Map{
        "Users": users,
        "CurrentPage": page,
    })
}
```

---

## License

This sample project is distributed under the terms of the project's repository license. Refer to the repository root for details.
