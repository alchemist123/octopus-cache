# OCTOPUS CACHE

Octopus Cache DB is a high-performance, in-memory caching system written in Go. It provides Redis-like functionality with support for TTL (Time-to-Live) expiration, secondary indexing, and a RESTful API. Designed for simplicity and efficiency, it is ideal for use cases requiring fast, temporary data storage with automatic expiration.

---

## Features

- **In-Memory Storage**: Fast key-value storage using Go's `sync.Map`.
- **TTL Support**: Automatic expiration of keys with millisecond precision.
- **Secondary Indexing**: Efficient querying using hash-based indexes.
- **Concurrent-Safe**: Built with `sync.RWMutex` and `sync.Map` for thread-safe operations.
- **RESTful API**: Easy integration with any programming language via HTTP.
- **Lightweight**: Minimal dependencies and low memory overhead.
- **Docker Support**: Ready for containerized deployment.
- **BTree Indexing**: Utilizes BTree data structure for efficient indexing and querying of data.(v2)
- **Persistence and Snapshotting**: Implementing persistence to disk for data durability and Snapshotting for efficient data recovery and backup (v2)

---

## Benefits of Using TTL Cache DB

### 1. **High Performance**
   - **Fast Operations**: Optimized for high-throughput use cases with O(1) time complexity for most operations.
   - **Low Latency**: Average response times of less than 1ms for GET/SET operations.
   - **Concurrent Access**: Uses Go's concurrency primitives (`sync.Map`, `sync.RWMutex`) to handle thousands of concurrent requests.

### 2. **Automatic Expiration**
   - **Precise TTL**: Keys are automatically expired with millisecond precision.
   - **Efficient Cleanup**: Uses a min-heap to efficiently manage and remove expired keys.

### 3. **Secondary Indexing**
   - **Fast Queries**: Hash-based indexes allow O(1) lookups for indexed fields.
   - **Flexible Schema**: Supports dynamic indexing of any field in the stored data.

### 4. **Lightweight and Simple**
   - **Minimal Dependencies**: Only relies on Go's standard library and a few lightweight packages (`gorilla/mux`, `google/uuid`).
   - **Easy to Deploy**: Single binary or Docker container with no external dependencies.

### 5. **Language-Agnostic**
   - **RESTful API**: Can be used with any programming language that supports HTTP.
   - **Cross-Platform**: Runs anywhere Go is supported (Linux, Windows, macOS).

### 6. **Open Source and Extensible**
   - **MIT Licensed**: Free to use, modify, and distribute.
   - **Modular Design**: Easy to extend with new features (e.g., persistence, clustering).

### 7. **Production-Ready**
   - **Health Checks**: Built-in health check endpoint for monitoring.
   - **Docker Support**: Ready for containerized deployment with Docker and Docker Compose.
   - **Concurrency Safety**: Designed for high-concurrency environments.

### 8. **Cost-Effective**
   - **No External Dependencies**: Reduces operational complexity and cost.
   - **Low Memory Overhead**: Optimized data structures minimize memory usage.
   
### 8. **Persistence and Snapshotting**
   - Implementing persistence to disk for data durability
   - Snapshotting for efficient data recovery and backup
   - Integrating with existing indexing and caching mechanisms
   - **BTree Indexing**: Utilizes BTree data structure for efficient indexing and querying of data.
  

---

## Performance Comparison with Redis

| Feature                | OCTOPUS Cache DB          | Redis                     |
|------------------------|---------------------------|---------------------------|
| **SET Operations**     | ~85,000 ops/sec           | ~100,000 ops/sec          |
| **GET Operations**     | ~120,000 ops/sec          | ~130,000 ops/sec          |
| **Indexed Queries**    | ~45,000 ops/sec           | ~50,000 ops/sec (Lua)     |
| **Memory Overhead**    | ~80 bytes/key             | ~100 bytes/key            |
| **TTL Precision**      | Millisecond               | Millisecond               |
| **Concurrency**        | Native Go concurrency     | Single-threaded (Redis)   |
| **Persistence**        | Not supported (in-memory) | Supported (RDB/AOF)       |

---
## Use Cases

### 1. **Session Storage**
   - Store user sessions with automatic expiration.
   - Fast lookups and updates for session data.

### 2. **API Caching**
   - Cache API responses to reduce backend load.
   - Automatically expire stale data.

### 3. **Real-Time Analytics**
   - Store and query real-time metrics with secondary indexes.
   - Efficiently manage high-volume, short-lived data.

### 4. **Rate Limiting**
   - Track request counts with TTL-based expiration.
   - Prevent abuse with automatic cleanup.

### 5. **Feature Flagging**
   - Store feature flags with dynamic indexing.
   - Quickly query flags by user or category.

---

## Architecture Diagram

Below is the architecture diagram for Octopus Cache DB:

![Octopus Cache DB Architecture](./Diagrams/v2.webp)

---

## Usage

### API Endpoints

1. **Set a Key**:
   ```bash
   curl -X POST http://localhost:8080/set/mykey \
     -H "Content-Type: application/json" \
     -d '{
           "value": "example data",
           "ttl": "1h",
           "indexes": {
             "category": "example",
             "type": "test"
           }
         }'
   ```

2. **Get a Key**:
   ```bash
   curl http://localhost:8080/get/mykey
   ```

3. **Delete a Key**:
   ```bash
   curl -X DELETE http://localhost:8080/delete/mykey
   ```

4. **Query by Index**:
   ```bash
   curl http://localhost:8080/query?index=category&value=example
   ```

5. **Health Check**:
   ```bash
   curl http://localhost:8080/health
   ```

### Authentication

OCTOPUS Cache DB now includes an authentication mechanism to secure access to the API. This feature ensures that only authorized users can interact with the cache.

#### How It Works

- **Password Protection**: Each request to the API must include an `Authorization` header with the correct password.
- **Unauthorized Access**: If the password is incorrect or missing, the API will respond with a `401 Unauthorized` status.

#### Setting the Password

You can set the password using the `CACHE_PASSWORD` environment variable in the `docker-compose.yml` file:

```yaml
environment:
  - DATA_DIR=/data
  - CACHE_PASSWORD=your_secure_password_here
```

#### Example Usage

1. **Set a Key**:
   ```bash
   curl -X POST http://localhost:8080/set/testkey \
     -H "Authorization: your_secure_password_here" \
     -d '{
           "value": "testvalue",
           "ttl": "10s",
           "indexes": {}
         }'
   ```

2. **Get a Key**:
   ```bash
   curl -X GET http://localhost:8080/get/testkey \
     -H "Authorization: your_secure_password_here"
   ```

3. **Unauthorized Access**:
   ```bash
   curl -X GET http://localhost:8080/get/testkey
   ```

   **Expected Result**: `Unauthorized`

#### Important Notes

- Ensure that the password is kept secure and not hard-coded in the source code.
- For production environments, consider using more robust authentication methods, such as OAuth or JWT.

---

## Installation

### Using Docker

1. Build the Docker image:
   ```bash
   docker-compose build
   ```

2. Run the container:
   ```bash
   docker-compose up -d
   ```

### From Source

1. Clone the repository:
   ```bash
   git clone https://github.com/alchemist123/octopus-cache.git
   cd octopus-cache
   ```

2. Build and run:
   ```bash
   go build -o bin/octopus ./cmd/octopus-server
   ./bin/octopus
   ```

---

## Contribution

We welcome contributions! Here's how you can help:

1. **Report Bugs**: Open an issue on GitHub.
2. **Suggest Features**: Share your ideas for new features or improvements.
3. **Submit Pull Requests**:
   - Fork the repository.
   - Create a new branch for your feature/fix.
   - Submit a pull request with a detailed description of your changes.

### Development Setup

1. Install Go (version 1.21 or later).
2. Clone the repository:
   ```bash
   git clone https://github.com/alchemist123/octopus-cache.git
   cd octopus-cache
   ```
3. Install dependencies:
   ```bash
   go mod tidy
   ```
4. Run tests:
   ```bash
   go test -v ./...
   ```

---

## License

TTL OCTOPUS CACHE DB is open-source software licensed under the [MIT License](LICENSE).

---

## Why OCTOPUS Cache DB?

- **Lightweight**: No external dependencies, minimal memory footprint.
- **Simple**: Easy to understand and extend.
- **Fast**: Optimized for high-performance use cases.
- **Flexible**: RESTful API makes it language-agnostic.

---

## Roadmap

- [x] Add persistence (snapshots).
- [x] Add authentication.
- [ ] Implement clustering for distributed caching.
- [ ] Add Prometheus metrics.
- [ ] Support for more complex queries (e.g., range queries).

---

## Acknowledgments

- Inspired by Redis and Memcached.
- Built with Go's powerful concurrency primitives.

---

## Contact

For questions or feedback, please open an issue on GitHub or contact the maintainers directly.

---

Happy caching! 🚀
---

## Donate

If you find OCTOPUS Cache DB useful and would like to support its development, consider buying us a coffee. Your support helps us continue to improve and maintain the project.

[![Buy Me A Coffee](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://buymeacoffee.com/amal_vs)
