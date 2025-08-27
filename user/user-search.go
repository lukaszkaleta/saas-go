package user

type UserSearch interface {
	ByPhone(phone string) (User, error)
}
