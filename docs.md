# Sato Framework Documentation

Sato Framework is a progressive Go framework for building efficient and scalable server-side applications. It uses modern Go and combines elements of OOP (Object Oriented Programming), FP (Functional Programming), and FRP (Functional Reactive Programming).

## Table of Contents
1. [Installation](#installation)
2. [Quick Start](#quick-start)
3. [Core Features](#core-features)
   - [Modules](#modules)
   - [Controllers](#controllers)
   - [Providers](#providers)
   - [Middleware](#middleware)
   - [Interceptors](#interceptors)
   - [Pipes](#pipes)
   - [Guards](#guards)
   - [Database](#database)
   - [Authentication](#authentication)
   - [Event System](#event-system)
   - [Dependency Injection](#dependency-injection)
4. [CLI Tools](#cli-tools)

## Installation

```bash
go get github.com/your-username/sato-framework
```

## Quick Start

```go
package main

import (
    "github.com/your-username/sato-framework/core"
    "github.com/gofiber/fiber/v2"
)

@core.Module(core.ModuleOptions{
    Controllers: []interface{}{&UserController{}},
    Providers:   []interface{}{&UserService{}},
})
type AppModule struct{}

func main() {
    app := fiber.New()
    
    // Initialize framework
    framework := core.NewApp()
    
    // Register routes
    core.RegisterRoutes(app)
    
    // Start server
    app.Listen(":3000")
}
```

## Core Features

### Modules

Modules are used to organize your application into cohesive blocks of functionality.

```go
@core.Module(core.ModuleOptions{
    Imports:     []interface{}{&DatabaseModule{}},
    Controllers: []interface{}{&UserController{}},
    Providers:   []interface{}{&UserService{}, &UserRepository{}},
    Exports:     []interface{}{&UserService{}},
})
type UserModule struct{}
```

### Controllers

Controllers are responsible for handling incoming requests and returning responses to the client.

```go
@core.Controller(core.ControllerOptions{
    Path: "/users",
    Version: "v1",
})
type UserController struct {
    UserService *UserService `inject:"userService"`
}

@core.Get(core.RouteOptions{Path: "/"})
@core.UseInterceptor(core.LoggingInterceptor(), core.InterceptorOptions{})
func (c *UserController) FindAll(ctx *fiber.Ctx) error {
    users, err := c.UserService.FindAll()
    if err != nil {
        return ctx.Status(500).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    return ctx.JSON(users)
}

@core.Post(core.RouteOptions{Path: "/"})
@core.UseGuards(AuthGuard)
func (c *UserController) Create(
    @core.UsePipe(core.ValidationPipe(), core.PipeOptions{})
    user User,
    ctx *fiber.Ctx,
) error {
    created, err := c.UserService.Create(user)
    if err != nil {
        return ctx.Status(500).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    return ctx.Status(201).JSON(created)
}
```

### Providers

Providers are a fundamental concept in Sato. Many of the basic Sato classes may be treated as a provider – services, repositories, factories, helpers, and so on.

```go
type UserService struct {
    UserRepository *UserRepository `inject:"userRepository"`
}

func NewUserService(repo *UserRepository) *UserService {
    return &UserService{
        UserRepository: repo,
    }
}

func (s *UserService) FindAll() ([]User, error) {
    return s.UserRepository.FindAll()
}
```

### Interceptors

Interceptors are used to intercept requests and responses.

```go
// Logging interceptor
func LoggingInterceptor() core.Interceptor {
    return func(c *fiber.Ctx, next fiber.Handler) error {
// Logger middleware
func Logger() fiber.Handler {
    return func(c *fiber.Ctx) error {
        start := time.Now()
        
        err := c.Next()
        if err != nil {
            return err
        }
        
        duration := time.Since(start)
        log.Printf("%s %s %d %v", c.Method(), c.Path(), c.Response().StatusCode(), duration)
        
        return nil
    }
}

// Using middleware
app.Use(Logger())
```

### Database

Sato provides database integration with various databases.

#### MongoDB

```go
// Create MongoDB provider
mongoProvider := core.NewMongoDBProvider("mongodb://localhost:27017", "mydb")

// Connect to database
if err := mongoProvider.Connect(); err != nil {
    panic(err)
}

// Get collection
collection := mongoProvider.GetCollection("users")
```

### Authentication

Sato provides a complete authentication system.

```go
// JWT Strategy
type JWTStrategy struct {
    Secret string
}

func (s *JWTStrategy) Validate(token string) (interface{}, error) {
    // Validate JWT token
    return nil, nil
}

// Using JWT Auth
jwtStrategy := &JWTStrategy{Secret: "your-secret-key"}
authMiddleware := core.AuthMiddleware(jwtStrategy)
app.Use(authMiddleware)
```

### Event System

Sato provides an event system for decoupled communication.

```go
// Create event bus
eventBus := core.NewEventBus()

// Subscribe to event
eventBus.Subscribe("user.created", func(event core.Event) error {
    // Handle event
    return nil
})

// Publish event
eventBus.Publish(&UserCreatedEvent{
    UserID: "123",
})
```

### Dependency Injection

Sato has a built-in Dependency Injection (DI) container.

```go
// Create container
container := core.NewContainer()

// Register providers
container.Register("userService", NewUserService())
container.Register("userRepository", NewUserRepository())

// Inject dependencies
type UserController struct {
    UserService *UserService `inject:"userService"`
}
```

## CLI Tools

Sato provides a set of CLI tools to help you develop your application.

### Generate Module

```bash
sato g module users
```

This will generate:
- `module/users/controller/users_controller.go`
- `module/users/service/users_service.go`
- `module/users/repository/users_repository.go`
- `module/users/dto/user_dto.go`
- `module/users/entity/user_entity.go`

### Generate CRUD

```bash
sato g crud users
```

This will generate a complete CRUD implementation for the users module.

## Configuration

Create a `config.json` file in your project root:

```json
{
    "app": {
        "port": 3000,
        "env": "development",
        "logLevel": "debug"
    },
    "database": {
        "driver": "mongodb",
        "host": "localhost",
        "port": 27017,
        "database": "mydb"
    },
    "auth": {
        "secret": "your-secret-key"
    }
}
```

## Best Practices

1. Use dependency injection for better testability
2. Implement proper error handling
3. Use middleware for cross-cutting concerns
4. Follow the module structure
5. Use the CLI tools for code generation
6. Implement proper authentication and authorization
7. Use the event system for decoupled communication
8. Follow the configuration management guidelines 