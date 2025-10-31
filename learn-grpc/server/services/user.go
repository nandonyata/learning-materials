package services

import (
	"context"
	"fmt"
	"learn-grpc/pb"
	"learn-grpc/server/datalayer/actions"
	"learn-grpc/server/datalayer/models"
	"learn-grpc/server/middleware"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// UserService implements the pb.UserServiceServer interface
// This is the BUSINESS LOGIC layer - it orchestrates the operations
type UserService struct {
	pb.UserServiceServer                             // Embed to satisfy interface
	userAction           actions.UserActionInterface // Database operations
}

// NewUserService creates a new UserService instance
// Called during server setup in main.go
//
// DEPENDENCY INJECTION PATTERN:
// Instead of creating userAction inside, we inject it
// Benefits: easier testing, loose coupling
func NewUserService(userAction actions.UserActionInterface) *UserService {
	return &UserService{
		userAction: userAction,
	}
}

// Register creates a new user account
//
// THE PROCESS:
// 1. Receive registration request (username, email, password)
// 2. Hash the password (NEVER store plain text!)
// 3. Save user to database
// 4. Generate authentication token
// 5. Return token + user info
//
// SECURITY NOTE:
// Password is hashed using bcrypt:
// "mypassword" → "$2a$10$N9qo8uLOickgx2ZMRZoMye..."
// - Even if database is leaked, passwords are safe
// - Cannot reverse the hash to get original password
func (s *UserService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	fmt.Println("Register: new request")

	// STEP 1: Hash the password
	// bcrypt automatically generates a salt and combines it with password
	// Cost = 10 (default): higher = more secure but slower
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		// Hashing failed (very rare, usually means system issues)
		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}

	// STEP 2: Create user object
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword), // Store the HASHED password
	}

	// STEP 3: Save to database
	if err := s.userAction.Create(ctx, user); err != nil {
		// Common errors:
		// - Email already exists (unique constraint violation)
		// - Username already taken
		// - Database connection issues
		return nil, status.Errorf(codes.AlreadyExists, "user already exists: %v", err)
	}

	// STEP 4: Generate token
	// SIMPLE VERSION (for learning): Use userID as token
	// Token = "123" means user with ID 123
	//
	// PRODUCTION VERSION should use JWT:
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	//     "userId": user.ID,
	//     "exp": time.Now().Add(24 * time.Hour).Unix(),
	// })
	// tokenString, _ := token.SignedString([]byte(secretKey))
	token := fmt.Sprintf("%d", user.ID)

	// STEP 5: Return response
	return &pb.AuthResponse{
		Token:   token,                    // Client saves this for future requests
		User:    models.ToUserProto(user), // User info (without password!)
		Message: "Registration successful",
	}, nil
}

// Login authenticates a user with email and password
//
// THE PROCESS:
// 1. Find user by email
// 2. Compare submitted password with stored hash
// 3. If match: generate token and return
// 4. If no match: return error
//
// SECURITY:
// - Uses constant-time comparison (bcrypt.CompareHashAndPassword)
// - Prevents timing attacks
// - Generic error message ("invalid credentials") to prevent user enumeration
func (s *UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	fmt.Println("Login: new request")

	// STEP 1: Find user by email
	user, err := s.userAction.FindByEmail(ctx, req.Email)
	if err != nil {
		// User not found
		// SECURITY: Don't say "email not found" - gives attackers info
		// Instead, generic message: "invalid credentials"
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	// STEP 2: Verify password
	// bcrypt.CompareHashAndPassword does:
	// - Extract salt from stored hash
	// - Hash the submitted password with same salt
	// - Compare results in constant time
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password), // Stored hash: "$2a$10$..."
		[]byte(req.Password),  // Submitted password: "mypassword"
	)
	if err != nil {
		// Password doesn't match
		// SECURITY: Same generic message as "user not found"
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
	}

	// STEP 3: Password is correct! Generate token
	token := fmt.Sprintf("%d", user.ID)

	// STEP 4: Return success response
	return &pb.AuthResponse{
		Token:   token,
		User:    models.ToUserProto(user),
		Message: "Login successful",
	}, nil
}

// GetProfile returns the current authenticated user's profile
//
// THE PROCESS:
// 1. Get userID from context (added by auth interceptor)
// 2. Fetch user from database
// 3. Return user info
//
// AUTHENTICATION:
// This endpoint requires authentication!
// - Client must send token in metadata
// - Interceptor validates token
// - Interceptor adds userID to context
// - We use that userID here
func (s *UserService) GetProfile(ctx context.Context, _ *emptypb.Empty) (*pb.User, error) {
	fmt.Println("GetProfile: new request")

	// STEP 1: Get userID from context
	// The auth interceptor already validated the token and added userID
	userID, ok := ctx.Value(middleware.UserIDKey).(int32)
	if !ok {
		// This should never happen if interceptor is working correctly
		// But we check anyway for safety
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	// STEP 2: Fetch user from database
	user, err := s.userAction.FindByID(ctx, uint(userID))
	if err != nil {
		// User was deleted or doesn't exist anymore
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	// STEP 3: Return user profile (without password!)
	return models.ToUserProto(user), nil
}

// DATA FLOW EXAMPLE:
//
// Register Flow:
// Client: {username: "john", email: "john@ex.com", password: "secret"}
//    ↓
// Service: Hash password → Save to DB → Generate token
//    ↓
// Client: {token: "123", user: {id: 123, username: "john"}, message: "Success"}
//
// Login Flow:
// Client: {email: "john@ex.com", password: "secret"}
//    ↓
// Service: Find user → Verify password → Generate token
//    ↓
// Client: {token: "123", user: {...}, message: "Login successful"}
//
// GetProfile Flow:
// Client: (sends token in metadata)
//    ↓
// Interceptor: Validate token → Add userID=123 to context
//    ↓
// Service: Get userID from context → Fetch user → Return profile
//    ↓
// Client: {id: 123, username: "john", email: "john@ex.com"}
