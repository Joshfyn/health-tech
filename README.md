# Health-Tech Go Schema

A comprehensive MongoDB schema implementation in Go for a health-tech application with user management, role-based access control, and database operations.

## 🚀 Quick Start

### 1. Start Database with Sample Data

```bash
# Navigate to migrate directory
cd migrate

# Start MongoDB with sample data (Go-based initialization)
docker-compose up

# Or run in background
docker-compose up -d
```

# Connect to mongo db with the following
    MONGO_INITDB_ROOT_USERNAME: admin
    MONGO_INITDB_ROOT_PASSWORD: password123
    MONGO_INITDB_DATABASE: health-tech

# To run the user API
Copy the envSample.txt to .env file and provide the env vars values
```bash
export $(cat .env | xargs)
```
The use the command to build and run the API
```bash
go run main.go
```