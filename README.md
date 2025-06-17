# 🧮 TCP Math Expression Server (Go)

A simple TCP server in Go that listens on a port, accepts expressions like `5+2`, evaluates them, and returns the result.

It supports basic arithmetic operations like addition (+) and multiplication (*).

---

## 🚀 Features

- Accepts TCP connections from clients
- Parses simple arithmetic expressions (`5+2`, `7 * 3`)
- Returns evaluated result
- Logs all input/output and handles timeouts

---

## 🛠 How It Works

- The server listens on port `1234`
- Accepts expressions like `5 + 2` or `5+2`
- Returns the result, e.g., `7`

---

## 📦 Requirements

- Go 1.18 or higher

---

## 🧪 Usage

### 1. Run the server

```bash
go run server.go
