# 📚 Bookstore Backend API (Golang)

A high-performance, secure backend service for managing book inventory, built with Go, MongoDB, and integrated with OpenTelemetry for full observability in New Relic.


## ✨ Features

- **CRUD Operations:** Full set of RESTful CRUD operations (Create, Read, Update, Delete).
- **Data Storage:** MongoDB persistence.
- **Security:** JWT Authentication with Role-Based Access Control (RBAC) for Admin/User roles.
- **Logging:** Logrus for structured logging with OTLP correlation.
- **Observability:** Complete integration of Traces, Metrics, and Logs via OpenTelemetry Protocol (OTLP/gRPC).
- **Monitoring:** All telemetry is sent directly to New Relic for analysis and alerting.


## 🔒 API Endpoints 

All endpoints require a valid JWT in the Authorization: Bearer <token> header.

| Method    | Path         | Description                     | Access        |
| :-------- | :---------   | :----------------------------   | :------------ |
| `GET`     | `/health`    | Application health check        | User or Admin |
| `GET`     | `/books`     | Retrieve a list of all books    | User or Admin |
| `POST`    | `/book`      | Create a new book entry         | Admin         |
| `PUT`     | `/book/{id}` | Update an existing book by ID   | Admin         |
| `DELETE`  | `/book/{id}` | Delete a book by its ID         | Admin         |
| `DELETE`  | `/books`     | **[CRITICAL]** Delete all books | Admin         |
| `POST`    | `/token`     | Generate bearer token           | User or Admin |


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

4. **Run the application**

```bash
go run main.go
```


## 📜 License

This project is licensed under the [MIT](https://choosealicense.com/licenses/mit/) License.
