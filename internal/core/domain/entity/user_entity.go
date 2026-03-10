package entity

type UserEntity struct {
	ID       int64
	Name     string
	Email    string
	Password string
	Role     string
	Status   string
	TenantID *int64
}
