# 📚 Bookstore Backend API (Golang)

A high-performance, secure backend service for managing book inventory, built with Go, MongoDB, and integrated with OpenTelemetry for full observability in New Relic.


## ✨ Features

- **CRUD Operations:** Full set of RESTful CRUD operations (Create, Read, Update, Delete).
- **Data Storage:** MongoDB persistence.
- **Security:** JWT Authentication with Role-Based Access Control (RBAC) for Admin/User roles.
- **Rate Limiting:** Token-bucket based rate limiting to prevent API abuse and ensure fair usage.
- **API Documentation:** Interactive Swagger/OpenAPI 3.0 documentation with authentication support.
- **Logging:** Logrus for structured logging with OTLP correlation.
- **Observability:** Complete integration of Traces, Metrics, and Logs via OpenTelemetry Protocol (OTLP/gRPC).
- **Monitoring:** All telemetry is sent directly to New Relic for analysis and alerting.


## � API Documentation

This API includes comprehensive **Swagger/OpenAPI 3.0 documentation** with interactive testing capabilities.

### 🌐 Access Swagger UI
Once the server is running, visit: **http://localhost:4000/swagger/index.html**

### 🔐 Authentication in Swagger
1. **Generate Token**: Use the `POST /token` endpoint with your credentials
2. **Authorize**: Click the **🔒 Authorize** button in Swagger UI
3. **Enter Token**: Format: `Bearer your_jwt_token_here`
4. **Test Endpoints**: All protected endpoints will now work with your token

## 🔒 API Endpoints 

| Method    | Path         | Description                     | Access        | Auth Required |
| :-------- | :---------   | :----------------------------   | :------------ | :------------ |
| `GET`     | `/health`    | Application health check        | Public        | ❌            |
| `POST`    | `/token`     | Generate JWT bearer token       | Public        | ❌            |
| `GET`     | `/books`     | Retrieve a list of all books    | User or Admin | ✅            |
| `POST`    | `/book`      | Create a new book entry         | Admin Only    | ✅            |
| `PUT`     | `/book/{id}` | Update an existing book by ID   | Admin Only    | ✅            |
| `DELETE`  | `/book/{id}` | Delete a book by its ID         | Admin Only    | ✅            |
| `DELETE`  | `/books`     | **[CRITICAL]** Delete all books | Admin Only    | ✅            |


## 🛠️ Prerequisites

Before running this service, ensure you have:
- **Go (1.21+):** The language runtime.
- **MongoDB:** A running instance (local or cloud).
- **New Relic Account:** To receive and visualize the telemetry data.
- **JWT Key Pair:** A set of public and private keys (RSA) for signing and verifying tokens.


## 🚀 Getting Started

1. **Clone the repository** 

```bash
git clone [YOUR_REPO_URL]
cd [YOUR_REPO_NAME]
```
2. **Install Dependencies**

```bash
go mod tidy
```

3. **Configure Environment Variables**

Create a .env file in the root directory and populate it with the required configuration. This is critical for connecting to MongoDB and New Relic.
    
| Variable               | Description                                 |
| :--------------------- | :----------------------------------------   |
| `MONGO_URL`            | Connection string for your MongoDB instance.|
| `NEW_RELIC_LICENSE_KEY`| New Relic Ingest - License key.             |
| `JWT_PRIVATE_KEY_B64 ` | JWT private key Base64-encoded.             |
| `JWT_PUBLIC_KEY_B64`   | JWT private key Base64-encoded.             |

4. **Generate or Update Swagger Documentation (optional)** 

```bash
# Install swag CLI tool (if not installed)
go install github.com/swaggo/swag/cmd/swag@latest

# Generate swagger docs
swag init
```

5. **Run the application**

```bash
go run main.go
```

The server will start on `http://localhost:4000`

## 🔧 Development

### Project Structure
```
bookstore/
├── controllers/         # HTTP handlers and business logic
├── middlewares/         # Authentication, rate limiting, logging, and recovery middleware
├── models/             # Data models and validation
├── routes/             # Route definitions and middleware chaining
├── db/                 # Database connection and configuration
├── logger/             # Logging configuration
├── otel/               # OpenTelemetry setup and configuration
├── docs/               # Auto-generated Swagger documentation
├── main.go             # Application entry point
├── go.mod              # Go module dependencies
└── README.md           # This file
```

## 📜 License

This project is licensed under the [MIT](https://choosealicense.com/licenses/mit/) License.
