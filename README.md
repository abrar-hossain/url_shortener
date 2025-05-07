# 📎 Go URL Shortener

A high-performance URL shortening service written in Go. This system generates short aliases for long URLS, stores them efficiently, and retrieves them with optimised access speed using Redis.

---

## 🚀 Features

- ✅ Shorten any valid URL with a unique key
- 🔁 Redirect users from short links to original URLS
- ⚡ In-memory storage using **Redis** for fast access
- 🧪 Unit-tested core functionality
- 🐳 Docker Compose support for quick setup

---

## 🛠️ Tech Stack

- **Language**: Go
- **Storage**: Redis
- **Web Framework**: net/http
- **Containerization**: Docker, Docker Compose
- **Testing**: `testing` package

---

## 📁 Project Structure

```bash
.
├── cmd/
│   └── main.go          # Entry point
├── internal/
│   ├── handler.go       # HTTP handlers
│   ├── service.go       # URL generation and lookup logic
│   └── storage.go       # Redis logic
├── config/
│   └── config.go        # Redis & app config
├── test/
│   └── url_test.go      # Unit tests
├── Dockerfile
├── docker-compose.yml
└── README.md
