# Health-Tech Database Setup

Simple Docker Compose setup that creates a MongoDB database and populates it with 3 sample records using Go.

## 🚀 Quick Start

```bash
# Navigate to migrate directory
cd migrate

# Start MongoDB with sample data (Go-based initialization)
docker-compose up

# Or run in background
docker-compose up -d
```

## 📊 What Gets Created

### Database: `health-tech`
- **3 Users**: Patient, Doctor, Pharmacist
- **3 Hospitals**: General Hospital, City Medical Center, Regional Hospital  
- **3 Pharmacies**: Downtown Pharmacy, Health Plus Pharmacy, Community Pharmacy

### Sample Data
- **Patient**: john.doe@example.com / password123
- **Doctor**: dr.smith@hospital.com / password123
- **Pharmacist**: pharmacist@pharmacy.com / password123

## 🛠 Management

### Start Database
```bash
docker-compose up
```

### Stop Database
```bash
docker-compose down
```

### Reset Database
```bash
docker-compose down -v
docker-compose up
```

### View Initialization Logs
```bash
docker-compose logs init
```

### Connect to Database
```bash
# MongoDB shell
docker exec -it health-tech-mongodb mongosh

# Or from host
mongosh mongodb://admin:password123@localhost:27017/health-tech?authSource=admin
```

## 📡 Connection Details

- **Host**: localhost:27017
- **Database**: health-tech
- **Username**: admin
- **Password**: password123
- **Connection String**: `mongodb://admin:password123@localhost:27017/health-tech?authSource=admin`

## 🔍 Verify Data

```bash
# Connect to MongoDB
docker exec -it health-tech-mongodb mongosh

# Switch to database
use health-tech

# Count records
db.users.countDocuments()
db.hospitals.countDocuments()
db.pharmacies.countDocuments()

# View users
db.users.find().pretty()

# View hospitals
db.hospitals.find().pretty()

# View pharmacies
db.pharmacies.find().pretty()
```

## 🏗 How It Works

1. **MongoDB Container**: Starts MongoDB 7.0 with authentication
2. **Go Initialization**: Runs `init.go` to create collections and populate data
3. **Automatic Setup**: Creates 3 users, 3 hospitals, and 3 pharmacies
4. **Ready to Use**: Database is immediately available for your application

That's it! Simple Go-based initialization. 🎉