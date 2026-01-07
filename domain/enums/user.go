package enums

type UserRole string

const (
	ADMINRole UserRole = "ADMIN"
	USERRole  UserRole = "USER"
)

func (ur UserRole) String() string {
	return string(ur)
}
