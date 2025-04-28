# Sato Framework

A progressive Go framework for building efficient and scalable server-side applications. Inspired by NestJS, Sato Framework combines the best practices of modern web development with Go's performance and simplicity.

## Features

- 🚀 **Progressive Architecture**: Built with modern Go and combines elements of OOP, FP, and FRP
- 🔄 **Dependency Injection**: Built-in DI container for better testability and maintainability
- 🛡️ **Guards & Interceptors**: Protect your routes and intercept requests/responses
- 🔌 **Pipes & Adapters**: Transform and validate data with ease
- 📦 **Module System**: Organize your application into cohesive blocks
- 🗄️ **Database Support**: Built-in support for MongoDB and SQL databases
- 🔐 **Authentication**: Complete authentication system with JWT support
- 📝 **CLI Tools**: Generate modules and CRUD operations with ease

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

func main() {
    app := fiber.New()
    
    // Initialize framework
    framework := core.NewApp()
    
    // Create MongoDB provider
    mongoProvider := core.NewMongoDBProvider("mongodb://localhost:27017", "example")
    if err := mongoProvider.Connect(); err != nil {
        log.Fatal(err)
    }
    
    // Register routes
    core.RegisterRoutes(app)
    
    app.Listen(":3000")
}
```

## CLI Usage

```bash
# Generate a new module
sato g module user

# Generate CRUD operations
sato g crud user
```

## Documentation

For detailed documentation, please visit our [Documentation](docs.md).

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 

