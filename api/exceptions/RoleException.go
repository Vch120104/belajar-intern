package exceptions

type RoleError struct {
	Error string
}

func NewRoleError(error string) RoleError {
	return RoleError{Error: error}
}
