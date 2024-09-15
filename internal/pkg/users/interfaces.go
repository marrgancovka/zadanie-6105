package users

type UserRepository interface {
	UserIsExists(username string) (bool, error)
}
