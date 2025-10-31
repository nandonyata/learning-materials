package models

import (
	"learn-grpc/pb"

	"gorm.io/gorm"
)

// User represents a user in our database
// This is the DATABASE model (how data is stored)
type User struct {
	// gorm.Model adds these fields automatically:
	// - ID        uint (primary key)
	// - CreatedAt time.Time
	// - UpdatedAt time.Time
	// - DeletedAt gorm.DeletedAt (for soft deletes)
	gorm.Model

	// Custom fields:
	Username string `gorm:"unique;not null"`
	// ↑ unique: No two users can have same username
	// ↑ not null: Username is required

	Email string `gorm:"unique;not null"`
	// ↑ unique: Email must be unique
	// ↑ not null: Email is required

	Password string `gorm:"not null"`
	// ↑ This stores the HASHED password (never plain text!)
	// ↑ Example: "$2a$10$N9qo8uLOickgx2ZMRZoMye..."
}

// TIP: Table name will be "users" (GORM pluralizes automatically)
// If you want custom name: func (User) TableName() string { return "my_users" }

// ToUserProto converts database User to protobuf User
// WHY WE NEED THIS:
// - Database model has gorm.Model (ID, CreatedAt, etc.)
// - Proto model has int32 id, string created_at
// - They're different types, so we need conversion
//
// This function:
// 1. Takes database User (with uint ID, time.Time, etc.)
// 2. Returns proto User (with int32 id, string created_at)
func ToUserProto(u *User) *pb.User {
	return &pb.User{
		Id:        int32(u.ID),                               // Convert uint → int32
		Username:  u.Username,                                // String → String (no conversion needed)
		Email:     u.Email,                                   // String → String
		CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"), // time.Time → formatted string
		// Format "2006-01-02 15:04:05" is Go's reference time
		// Output example: "2024-10-31 14:30:00"
	}
}

// DESIGN NOTE: We don't include Password in proto.User
// Why? Passwords should NEVER be sent to clients!
// Even hashed passwords should stay on the server
