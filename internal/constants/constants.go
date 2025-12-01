package constants

// DefaultPage is the default page number for pagination.
const DefaultPage = 1

// DefaultPageSize is the default page size for pagination.
const DefaultPageSize = 10

// DefaultLimit is the default limit for pagination.
const DefaultLimit = 10

// MaxLimit is the maximum limit for pagination.
const MaxLimit = 100

// Context keys
type contextKey string

const (
	// UserClaimsKey is the context key for storing the UserClaims.
	UserClaimsKey contextKey = "userClaims"
)
const (
	// Topic for user authentication and lifecycle events (registration, login)
	TopicUserEvents = "user-auth-events"

	// Add other topics here as features grow (e.g., TopicCourseEnrollments = "course-enroll-events")
)
