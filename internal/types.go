package internal

type UserStatus string
type Role string
type OrderStatus string

const (
	Unverified UserStatus = "unverified"
	Verified   UserStatus = "verified"
	Active     UserStatus = "active"
)

const (
	UserRole       Role = "user"
	SupervisorRole Role = "supervisor"
	AdminRole      Role = "administrator"
)

const (
	Created  UserStatus = "created"
	Approved UserStatus = "approved"
	Declined UserStatus = "Decline"
)
