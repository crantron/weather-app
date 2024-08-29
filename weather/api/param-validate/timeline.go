package param_validate

type Timeline struct {
	latitude  float64
	longitude float64
	city      string
}

func (timeline *Timeline) validate(query string) (bool, error) {

	return true, nil
}
