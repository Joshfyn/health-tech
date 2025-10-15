# Health-Tech Go Schema

A comprehensive MongoDB schema implementation in Go for a health-tech application with user management, role-based access control, and database operations.

## 🚀 Quick Start

### 1. Start Database with Sample Data

```bash
# Navigate to migrate directory
cd migrate

# Start MongoDB with sample data (Go-based initialization)
docker-compose up -d
```

---

## 🐳 Running Services with Docker

### 2. Setup Environment Variables

Create a `docker.env` file in the project root for Docker containers:

```bash
MONGO_URI=mongodb://mongodb:27017/health-tech?replicaSet=rs0
USER_SERVICE_PORT=8080
AUTH_SERVICE_PORT=8081
JWT_SECRET=your-jwt-secret-here
SMTP_PASS=your-smtp-password
SMTP_USER=your-email@example.com
```

> **Note:** The `docker.env` file is formatted for Docker (no `export` statements). This is different from `env.sh` which is used for local development.

**Important MongoDB Host Settings:**
- If MongoDB is running locally (not in Docker): Set Mongo host as `localhost:27017/?replicaSet=rs0`
- 

### 3. Build and Run User Service

```bash
# Build the Docker image
docker build -t user-service -f user/Dockerfile .

# Run the container
docker run -d \
 --name user-service-container \
 -p 8080:80 \
 --env-file docker.env \
 --network migrate_health-tech-network \
 user-service ./user-service
```

**Access User Service:** `http://localhost:8080`

### 4. Build and Run Auth Service

```bash
# Build the Docker image
docker build -t auth-service -f auth/Dockerfile .

# Run the container
docker run -d \
 --name auth-service-container \
 -p 8081:80 \
 --env-file docker.env \
 --network migrate_health-tech-network \
 auth-service ./auth-service
```

## Local Development (Without Docker)

### Setup Local Environment

1. Copy `envSample.txt` to create `env.sh`:
```bash
cp envSample.txt env.sh
```

2. Update `env.sh` with your values (should have `export` statements)

3. Source the environment variables:
```bash
source env.sh
# or
export $(cat env.sh | xargs)
```

### Run User Service Locally

```bash
go run user/main.go
```

### Run Auth Service Locally

```bash
go run auth/main.go
```