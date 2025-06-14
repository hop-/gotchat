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
```

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
