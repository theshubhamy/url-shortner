# ✂️ URL Shortener API

A simple and efficient URL Shortener API built with **Go**, **Fiber**, and **Redis**. This service allows users to shorten long URLs, create custom aliases, control expiry time, and handle redirections — all with IP-based rate limiting and Docker support.

---

## 🚀 Features

- 🔗 Shorten URLs with optional custom aliases
- 🕒 Set expiry duration (in hours)
- 🚦 Rate limiting per IP using Redis
- 🌐 URL validation and HTTPS enforcement
- 🐳 Full Docker support for local development
- 📈 Visit counter for shortened URLs

---

## 🧱 Tech Stack

- **Language**: Go (Golang)
- **Framework**: [Fiber](https://gofiber.io)
- **Database**: [Redis](https://redis.io)
- **Validation**: [govalidator](https://github.com/asaskevich/govalidator)
- **UUID**: [google/uuid](https://pkg.go.dev/github.com/google/uuid)
- **Containerization**: Docker + Docker Compose

---

## 📁 Project Structure

```

.
├── api/ # Go backend service
│ ├── config/ # App configuration
│ ├── database/ # Redis setup
│ ├── helpers/ # Utilities (e.g., URL validation)
│ └── routes/ # Route handlers (shorten, resolve)
├── db/ # Redis Docker setup (optional)
├── data/ # Redis persistent volume
├── docker-compose.yml # Docker Compose file
├── README.md
└── .gitignore

```

---

## 🐳 Getting Started with Docker

### 1. Clone the Project

```bash
git clone https://github.com/yourusername/urlshortener.git
cd urlshortener
```

### 2. Start with Docker Compose

```bash
docker-compose up -d
```

This will:

- Build and start the Go API (`api`) on port **3300**
- Start Redis (`db`) on port **6379**
- Map Redis data to local `./data` for persistence

---

## 📬 API Endpoints

### 1. `POST /api/v1` — Create Short URL

Create a shortened URL (with optional custom alias and expiry).

#### Request

```json
{
  "url": "https://example.com",
  "customShort": "myalias", // optional
  "expirey": 24 // optional (in hours)
}
```

#### Response

```json
{
  "url": "https://example.com",
  "customShort": "http://localhost:3300/myalias",
  "expirey": 24,
  "XRateLimitRest": 29,
  "XRateRemaining": 9
}
```

- **XRateLimitRest**: minutes left until your rate limit resets
- **XRateRemaining**: how many more requests you can make

---

### 2. `GET /:url` — Redirect to Original

Redirect to the original URL using the shortened alias.

#### Example

```http
GET http://localhost:3300/myalias
```

#### Behavior

### - ✅ Redirects to the original URL with `301 Moved Permanently`

### - ❌ If not found, responds with `404 Not Found`:

```json
{
  "error": "short url not found!"
}
```

---

## 🚦 Rate Limiting

- Each IP address gets **30 requests per 30 minutes**
- After limit is exceeded, the API returns:

```json
{
  "error": "Rate limit exceeded",
  "rate_limit_rest": 28
}
```

Where `rate_limit_rest` is the number of minutes left before your limit resets.

---

## 🛠 Configuration

Update `api/config/config.go` to change domain and quota:

```go
const (
    Domain   = "http://localhost:3300"
    ApiQuota = 30
)
```

---

## 🧪 Testing

Use [Postman](https://www.postman.com/) or `curl` to test the API:

```bash
curl -X POST http://localhost:3300/api/v1 \
     -H "Content-Type: application/json" \
     -d '{"url": "https://example.com"}'
```

```bash
curl -I http://localhost:3300/myalias
```

---

## 📄 License

MIT License. See `LICENSE` file for more info.
# url-shortner
