package middleware

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string

const userIDKey contextKey = "userID"

// AuthInterceptor is the main struct that holds our authentication logic
// Think of it as a security configuration object
type AuthInterceptor struct {
	secretKey string // Used to validate tokens (in real JWT this would verify signatures)
}

// NewAuthInterceptor creates a new instance of our auth interceptor
// This is called once when the server starts
func NewAuthInterceptor(secretKey string) *AuthInterceptor {
	return &AuthInterceptor{
		secretKey: secretKey,
	}
}

// UnaryInterceptor handles authentication for NORMAL (unary) gRPC calls
// Unary = Simple request → response (like CreateBlog, Login, etc.)
//
// HOW IT WORKS:
// 1. Request comes in
// 2. This function runs BEFORE the actual handler
// 3. We check if auth is needed
// 4. If yes, validate token
// 5. If valid, add userID to context
// 6. Call the actual handler
// 7. Return response
func (a *AuthInterceptor) UnaryInterceptor(
	ctx context.Context, // The request context (contains metadata, deadlines, etc.)
	req interface{}, // The actual request message (e.g., CreateBlogRequest)
	info *grpc.UnaryServerInfo, // Info about which method is being called
	handler grpc.UnaryHandler, // The actual handler function to call
) (interface{}, error) {

	// STEP 1: Check if this endpoint needs authentication
	// Skip auth for public endpoints (Login and Register)
	if strings.Contains(info.FullMethod, "Login") ||
		strings.Contains(info.FullMethod, "Register") {
		// Public endpoint - just call the handler directly
		return handler(ctx, req)
	}

	// STEP 2: This is a protected endpoint, validate the token
	userID, err := a.authorize(ctx)
	if err != nil {
		// Token is invalid or missing - reject the request
		return nil, err
	}

	// STEP 3: Token is valid! Add userID to context so handlers can use it
	// This is like attaching a "verified badge" with the user's ID
	ctx = context.WithValue(ctx, userIDKey, userID)

	// STEP 4: Now call the actual handler with the enriched context
	return handler(ctx, req)
}

// StreamInterceptor handles authentication for STREAMING gRPC calls
// Streaming = Continuous data flow (like GetAllBlog which returns multiple blogs)
//
// Similar to UnaryInterceptor but works with streams
func (a *AuthInterceptor) StreamInterceptor(
	srv interface{}, // The service implementation
	ss grpc.ServerStream, // The stream object
	info *grpc.StreamServerInfo, // Info about the stream method
	handler grpc.StreamHandler, // The actual stream handler
) error {

	// STEP 1: Validate token from the stream's context
	userID, err := a.authorize(ss.Context())
	if err != nil {
		// Invalid token - reject the stream
		return err
	}

	// STEP 2: Create a wrapped stream with userID in context
	// We need to wrap because we can't modify the original stream's context
	wrappedStream := &wrappedServerStream{
		ServerStream: ss,
		ctx:          context.WithValue(ss.Context(), userIDKey, userID),
	}

	// STEP 3: Call the actual stream handler with our wrapped stream
	return handler(srv, wrappedStream)
}

// authorize is the CORE authentication logic
// It extracts and validates the token from the request
//
// PROCESS:
// 1. Get metadata from context (like HTTP headers)
// 2. Extract "authorization" header
// 3. Parse token format: "Bearer <token>"
// 4. Validate the token
// 5. Return userID if valid, error if not
func (a *AuthInterceptor) authorize(ctx context.Context) (int32, error) {
	// STEP 1: Extract metadata from context
	// Metadata is like HTTP headers - key-value pairs sent with the request
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, status.Errorf(codes.Unauthenticated, "metadata not provided")
	}

	// STEP 2: Get the "authorization" header
	// In HTTP terms: Authorization: Bearer <token>
	values := md["authorization"]
	if len(values) == 0 {
		return 0, status.Errorf(codes.Unauthenticated, "authorization token not provided")
	}

	token := values[0]

	// STEP 3: Validate token format
	// Expected format: "Bearer 123" where 123 is the userID
	// In production, this would be "Bearer eyJhbGc..." (JWT token)
	parts := strings.Split(token, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return 0, status.Errorf(codes.Unauthenticated, "invalid token format")
	}

	// STEP 4: Extract and validate the token
	// SIMPLE VERSION (for learning):
	// Token IS the userID: "Bearer 123" → userID = 123
	//
	// PRODUCTION VERSION would do:
	// - Decode JWT token
	// - Verify signature with secretKey
	// - Check expiration
	// - Extract userID from claims
	var userID int32
	_, err := fmt.Sscanf(parts[1], "%d", &userID)
	if err != nil {
		return 0, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	// STEP 5: Token is valid! Return the userID
	return userID, nil
}

// wrappedServerStream wraps grpc.ServerStream to inject custom context
// This is needed because gRPC streams have immutable contexts
//
// WHY WE NEED THIS:
// - We want to add userID to the context
// - But stream.Context() is read-only
// - So we wrap the stream and override the Context() method
type wrappedServerStream struct {
	grpc.ServerStream                 // Embed the original stream
	ctx               context.Context // Our custom context with userID
}

// Context returns our custom context instead of the original
// This makes the userID available to the stream handler
func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}
