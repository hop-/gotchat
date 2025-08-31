# GoTChat

GoTChat is a chat application built in Go, designed to provide a robust and scalable chat experience. It is terminal-based and leverages an event-driven architecture to handle real-time messaging efficiently.

## Features

- User signin and signup

- Real-time messaging

- Event-driven real-time architecture

- Terminal-based user interface (TUI)

## Project Structure

```plaintext
internal/
  app/        # Application entry point and lifecycle management
  core/       # Core entities and services
  logic/      # Business logic
  ui/         # Terminal-based user interface (TUI)

pkg/
  log/        # Logging utilities and structured logging
  network/    # Network transport layer and connection management
```

## Package Details

### pkg/log

The logging package provides structured logging capabilities with different log levels and configurable output formats. It includes:

- Thread-safe logging operations
- Multiple log levels (Debug, Info, Warn, Error)
- Configurable output destinations
- Builder pattern for logger configuration

### pkg/network

The network package handles all network-related operations including:

- TCP connection management
- Secure connection handling with encryption
- Network message serialization and deserialization
- Transport layer abstraction
- Connection listeners and acceptors

## Prerequisites

- Go 1.24.2 or later

## Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/hop-/gotchat.git
   cd gotchat
   ```

2. Build the application:

   Ensure you have Go installed and set up. Then, run the following command to build the application:

   ```bash
   go build
   ```

## Installation Unsing Go

To install GoTChat using Go, you can use the following command:

```bash
go install github.com/hop-/gotchat@latest
```

## Testing

Run the unit tests:

```bash
go test ./...
```

## Usage

To run the application, execute the following command:

```bash
./gotchat
```

or if installed via Go:

```bash
gotchat
```

## License

This project is licensed under the MIT License.
