# 🤖 Package aiyou

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.16-blue.svg)](https://golang.org/dl/)

A minimalist and efficient Go SDK for interacting with the AI.You API. Designed to be simple to use while offering maximum flexibility.

## ✨ Features

- 🎯 Simple and intuitive interface
- 🤖 Two main functions:
  - List available models
  - Send messages to a model or assistant, with temperature control
- 🔄 Streaming support for real-time responses
- ⚡ Automatic retry handling
- 🛠️ Flexible configuration options
- 🔍 Debug mode for development
- 🌡️ Response temperature control

## 📦 Installation

```bash
go get github.com/n1neT10ne/aiyou
```

## 🚀 Usage

```go
package main

import (
    "fmt"
    "github.com/n1neT10ne/aiyou"
)

func main() {
    // List available models
    models, err := aiyou.ListModels(
        "your-token",
        aiyou.WithDebug(true),
    )
    if err != nil {
        panic(err)
    }
    for _, model := range models {
        fmt.Printf("- %s\n", model.Name)
    }

    // Basic configuration
    response, err := aiyou.Completion(
        "model-name",
        "your-token",
        "your message",
    )
    if err != nil {
        panic(err)
    }
    fmt.Println(response)

    // With options
    response, err = aiyou.Completion(
        "model-name",
        "your-token",
        "your message",
        aiyou.WithDebug(true),
        aiyou.WithTemperature(0.7),
        aiyou.WithSystemPrompt("system prompt"),
        aiyou.WithStream(true),
    )
}
```

## ⚙️ Options

The package supports several configuration options:

```go
// Enable debug mode
WithDebug(debug bool)

// Set generation temperature (0.0-2.0)
WithTemperature(temp float64)

// Configure retry behavior
WithRetry(maxRetries int, delay time.Duration)

// Set system prompt
WithSystemPrompt(prompt string)

// Enable streaming mode
WithStream(stream bool)
```

## 🔄 Streaming Mode

Streaming mode allows receiving the response as it's being generated. It's particularly useful for long responses or to display the response progressively.

```go
response, err := aiyou.Completion(
    "model-name",
    "your-token",
    "your message",
    aiyou.WithStream(true),
)
```

In streaming mode:
- The response is built progressively from received chunks
- Each chunk contains a part of the final response
- Debug mode displays received chunks and their content

## ⚠️ Error Handling

The package defines several error types:

```go
var (
    ErrEmptyToken      = errors.New("token is required")
    ErrEmptyMessage    = errors.New("message is required")
    ErrInvalidToken    = errors.New("invalid token")
    ErrRateLimit       = errors.New("rate limit exceeded")
    ErrStreamCorrupted = errors.New("stream response corrupted")
)
```

## 🧪 Unit Tests

The package includes a complete suite of unit tests. To run them, you need to set your AI.You token in the `AIYOU_TEST_TOKEN` environment variable:

```bash
export AIYOU_TEST_TOKEN="your-token"
go test ./...
```

If the environment variable is not set, tests requiring authentication will be automatically skipped with an explanatory message. This approach:
- Avoids storing tokens directly in the code
- Facilitates continuous integration without exposing sensitive data
- Allows each developer to use their own test token

## 🔗 CLI

A command-line tool is available in a separate project: [aiyou-cli](https://github.com/n1neT10ne/aiyou-cli). This CLI provides a quick and simple way to interact with the AI.You API directly from your terminal.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👤 Author

**Cyrille BARTHELEMY**

* Github: [@n1neT10ne](https://github.com/n1neT10ne)
