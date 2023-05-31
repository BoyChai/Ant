package Ant

func New(validator Validator) Validator {
	v := validator
	if validator.Parity == "" {
		v.Parity = Ant
	}
	return v
}
