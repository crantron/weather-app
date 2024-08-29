package param_validate

type ValidateInterface interface {
	validate(query string) (bool, error)
}
