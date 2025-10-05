package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"health-tech/user/models"
	"health-tech/user/services"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

	// Create user service
	userService := services.NewUserService(db)

	// Run initialization
	if err := initializeDatabase(ctx, userService, db); err != nil {
		log.Fatal("Initialization failed:", err)
	}

	fmt.Println("🎉 Database initialization completed successfully!")
}

func initializeDatabase(ctx context.Context, userService *services.UserService, db *mongo.Database) error {
	fmt.Println("🌱 Initializing database with sample data...")

	// Check if data already exists
	userCount, err := userService.CountUsers(ctx, map[string]interface{}{})
	if err != nil {
		return fmt.Errorf("failed to check existing users: %w", err)
	}

	if userCount > 0 {
		fmt.Println("⚠️  Database already contains data. Skipping initialization.")
		return nil
	}

	// Create sample users
	if err := createSampleUsers(ctx, userService); err != nil {
		return fmt.Errorf("failed to create sample users: %w", err)
	}

	// Create sample hospitals
	if err := createSampleHospitals(ctx, db); err != nil {
		return fmt.Errorf("failed to create sample hospitals: %w", err)
	}

	// Create sample pharmacies
	if err := createSamplePharmacies(ctx, db); err != nil {
		return fmt.Errorf("failed to create sample pharmacies: %w", err)
	}

	// Show summary
	showSummary(ctx, userService, db)

	return nil
}

func createSampleUsers(ctx context.Context, userService *services.UserService) error {
	fmt.Println("👥 Creating 3 sample users...")

	now := time.Now()

	// Create hospital and pharmacy IDs
	hospitalID := primitive.NewObjectID()
	pharmacyID := primitive.NewObjectID()

	users := []models.User{
		// Patient
		{
			Email:     "john.doe@example.com",
			Role:      models.RolePatient,
			FirstName: "John",
			LastName:  "Doe",
			DateOfBirth: func() *time.Time {
				dob := time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC)
				return &dob
			}(),
			Contact: models.Contact{
				Phone: "+1-555-0101",
				Address: models.Address{
					Street:  "123 Main St",
					City:    "New York",
					State:   "NY",
					ZipCode: "10001",
					Country: "US",
				},
			},
			IsActive:  true,
			CreatedAt: now,
			UpdatedAt: now,
		},
		// Doctor
		{
			Email:         "dr.smith@hospital.com",
			Role:          models.RoleDoctor,
			FirstName:     "Jane",
			LastName:      "Smith",
			EntityID:      &hospitalID,
			EntityType:    func() *models.EntityType { et := models.EntityHospital; return &et }(),
			LicenseNumber: func() *string { ln := "MD123456"; return &ln }(),
			Contact: models.Contact{
				Phone: "+1-555-0201",
				Address: models.Address{
					Street:  "456 Medical Ave",
					City:    "Boston",
					State:   "MA",
					ZipCode: "02101",
					Country: "US",
				},
			},
			IsActive:  true,
			CreatedAt: now,
			UpdatedAt: now,
		},
		// Pharmacist
		{
			Email:         "pharmacist@pharmacy.com",
			Role:          models.RolePharmacist,
			FirstName:     "Mike",
			LastName:      "Johnson",
			EntityID:      &pharmacyID,
			EntityType:    func() *models.EntityType { et := models.EntityPharmacy; return &et }(),
			LicenseNumber: func() *string { ln := "RPH789012"; return &ln }(),
			Contact: models.Contact{
				Phone: "+1-555-0301",
				Address: models.Address{
					Street:  "789 Pharmacy St",
					City:    "Chicago",
					State:   "IL",
					ZipCode: "60601",
					Country: "US",
				},
			},
			IsActive:  true,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	// Hash passwords and create users
	for i := range users {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password for user %s: %w", users[i].Email, err)
		}
		users[i].PasswordHash = string(hashedPassword)

		if err := userService.CreateUser(ctx, &users[i]); err != nil {
			return fmt.Errorf("failed to create user %s: %w", users[i].Email, err)
		}
	}

	fmt.Println("✅ Created 3 users successfully!")
	return nil
}

func createSampleHospitals(ctx context.Context, db *mongo.Database) error {
	fmt.Println("🏥 Creating 3 sample hospitals...")

	collection := db.Collection("hospitals")
	now := time.Now()

	hospitals := []interface{}{
		map[string]interface{}{
			"name": "General Hospital",
			"address": map[string]interface{}{
				"street":  "123 Medical Center Dr",
				"city":    "New York",
				"state":   "NY",
				"zipCode": "10001",
				"country": "US",
			},
			"phone":         "+1-555-0100",
			"email":         "info@generalhospital.com",
			"licenseNumber": "HOSP001",
			"createdAt":     now,
		},
		map[string]interface{}{
			"name": "City Medical Center",
			"address": map[string]interface{}{
				"street":  "456 Health Ave",
				"city":    "Boston",
				"state":   "MA",
				"zipCode": "02101",
				"country": "US",
			},
			"phone":         "+1-555-0200",
			"email":         "contact@citymedical.com",
			"licenseNumber": "HOSP002",
			"createdAt":     now,
		},
		map[string]interface{}{
			"name": "Regional Hospital",
			"address": map[string]interface{}{
				"street":  "789 Care Blvd",
				"city":    "Chicago",
				"state":   "IL",
				"zipCode": "60601",
				"country": "US",
			},
			"phone":         "+1-555-0300",
			"email":         "admin@regionalhospital.com",
			"licenseNumber": "HOSP003",
			"createdAt":     now,
		},
	}

	_, err := collection.InsertMany(ctx, hospitals)
	if err != nil {
		return fmt.Errorf("failed to insert hospitals: %w", err)
	}

	fmt.Println("✅ Created 3 hospitals successfully!")
	return nil
}

func createSamplePharmacies(ctx context.Context, db *mongo.Database) error {
	fmt.Println("💊 Creating 3 sample pharmacies...")

	collection := db.Collection("pharmacies")
	now := time.Now()

	pharmacies := []interface{}{
		map[string]interface{}{
			"name": "Downtown Pharmacy",
			"address": map[string]interface{}{
				"street":  "321 Main St",
				"city":    "New York",
				"state":   "NY",
				"zipCode": "10002",
				"country": "US",
			},
			"phone":         "+1-555-1000",
			"email":         "info@downtownpharmacy.com",
			"licenseNumber": "PHARM001",
			"createdAt":     now,
		},
		map[string]interface{}{
			"name": "Health Plus Pharmacy",
			"address": map[string]interface{}{
				"street":  "654 Wellness St",
				"city":    "Boston",
				"state":   "MA",
				"zipCode": "02102",
				"country": "US",
			},
			"phone":         "+1-555-2000",
			"email":         "contact@healthplus.com",
			"licenseNumber": "PHARM002",
			"createdAt":     now,
		},
		map[string]interface{}{
			"name": "Community Pharmacy",
			"address": map[string]interface{}{
				"street":  "987 Care Lane",
				"city":    "Chicago",
				"state":   "IL",
				"zipCode": "60602",
				"country": "US",
			},
			"phone":         "+1-555-3000",
			"email":         "info@communitypharmacy.com",
			"licenseNumber": "PHARM003",
			"createdAt":     now,
		},
	}

	_, err := collection.InsertMany(ctx, pharmacies)
	if err != nil {
		return fmt.Errorf("failed to insert pharmacies: %w", err)
	}

	fmt.Println("✅ Created 3 pharmacies successfully!")
	return nil
}

func showSummary(ctx context.Context, userService *services.UserService, db *mongo.Database) {
	fmt.Println("\n📊 Database Summary:")
	fmt.Println("===================")

	// User count
	userCount, _ := userService.CountUsers(ctx, map[string]interface{}{})
	fmt.Printf("👥 Users: %d\n", userCount)

	// Hospital count
	hospitalCount, _ := db.Collection("hospitals").CountDocuments(ctx, map[string]interface{}{})
	fmt.Printf("🏥 Hospitals: %d\n", hospitalCount)

	// Pharmacy count
	pharmacyCount, _ := db.Collection("pharmacies").CountDocuments(ctx, map[string]interface{}{})
	fmt.Printf("💊 Pharmacies: %d\n", pharmacyCount)

	fmt.Println("\n🔑 Login Credentials:")
	fmt.Println("===================")
	fmt.Println("Patient: john.doe@example.com / password123")
	fmt.Println("Doctor: dr.smith@hospital.com / password123")
	fmt.Println("Pharmacist: pharmacist@pharmacy.com / password123")

	fmt.Println("\n📡 Connection String:")
	fmt.Println("===================")
	fmt.Println("mongodb://admin:password123@localhost:27017/health-tech?authSource=admin")
}
