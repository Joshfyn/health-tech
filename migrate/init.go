package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"health-tech/user/model"

	"github.com/gofrs/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Get MongoDB URI from environment variable or use default
	mongoURI := "mongodb://admin:password123@localhost:27017/health-tech?authSource=admin"
	if envURI := os.Getenv("MONGODB_URI"); envURI != "" {
		mongoURI = envURI
	}

	fmt.Println("🚀 Health-Tech Database Initialization")
	fmt.Println("=====================================")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer client.Disconnect(context.TODO())

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	fmt.Println("✅ Connected to MongoDB successfully!")

	// Get database
	db := client.Database("health-tech")

	// Run initialization
	if err := initializeDatabase(ctx, db); err != nil {
		log.Fatal("Initialization failed:", err)
	}

	fmt.Println("🎉 Database initialization completed successfully!")
}

func initializeDatabase(ctx context.Context, db *mongo.Database) error {
	fmt.Println("🌱 Initializing database with sample data...")

	// Create sample users
	if err := createSampleUsers(ctx, db); err != nil {
		return fmt.Errorf("failed to create sample users: %w", err)
	}

	return nil
}

func createSampleUsers(ctx context.Context, db *mongo.Database) error {
	fmt.Println("👥 Creating 3 sample users...")
	collection := db.Collection("users")

	now := time.Now()

	// user ids
	patientID, _ := uuid.NewV4()
	doctorID, _ := uuid.NewV4()
	pharmacistID, _ := uuid.NewV4()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password %w", err)
	}

	users := []interface{}{
		// Patient
		map[string]interface{}{
			"id":           patientID,
			"email":        "john.doe@example.com",
			"passwordHash": string(hashedPassword),
			"role":         model.RoleSuperAdmin,
			"firstName":    "John",
			"lastName":     "Doe",
			"dateOfBirth": func() *time.Time {
				dob := time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC)
				return &dob
			}(),
			"contact": model.Contact{
				Phone: "+1-555-0101",
				Address: model.Address{
					Street:  "123 Main St",
					City:    "New York",
					State:   "NY",
					ZipCode: "10001",
					Country: "US",
				},
			},
			"isActive":  true,
			"createdAt": now,
			"updatedAt": now,
		},
		// Doctor
		map[string]interface{}{
			"id":            doctorID,
			"email":         "dr.smith@hospital.com",
			"passwordHash":  string(hashedPassword),
			"role":          model.RoleHospitalAdmin,
			"firstName":     "Jane",
			"lastName":      "Smith",
			"entityType":    func() *model.EntityType { et := model.EntityHospital; return &et }(),
			"licenseNumber": func() *string { ln := "MD123456"; return &ln }(),
			"contact": model.Contact{
				Phone: "+1-555-0201",
				Address: model.Address{
					Street:  "456 Medical Ave",
					City:    "Boston",
					State:   "MA",
					ZipCode: "02101",
					Country: "US",
				},
			},
			"isActive":  true,
			"createdAt": now,
			"updatedAt": now,
		},
		// Pharmacist
		map[string]interface{}{
			"id":            pharmacistID,
			"email":         "pharmacist@pharmacy.com",
			"passwordHash":  string(hashedPassword),
			"role":          model.RolePharmacyAdmin,
			"firstName":     "Mike",
			"lastName":      "Johnson",
			"entityType":    func() *model.EntityType { et := model.EntityPharmacy; return &et }(),
			"licenseNumber": func() *string { ln := "RPH789012"; return &ln }(),
			"contact": model.Contact{
				Phone: "+1-555-0301",
				Address: model.Address{
					Street:  "789 Pharmacy St",
					City:    "Chicago",
					State:   "IL",
					ZipCode: "60601",
					Country: "US",
				},
			},
			"isActive":  true,
			"createdAt": now,
			"updatedAt": now,
		},
	}

	_, err = collection.InsertMany(ctx, users)
	if err != nil {
		return fmt.Errorf("failed to insert users: %w", err)
	}

	fmt.Println("✅ Created 3 users successfully!")
	return nil
}
