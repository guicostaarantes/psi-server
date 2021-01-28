package match

// IMatchUtil is an abstraction for a utility that checks requirements of strings
type IMatchUtil interface {
	IsPasswordStrong(password string) error
	IsEmailValid(email string) error
}
