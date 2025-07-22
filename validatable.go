package querybm

// Validatable is an interface for components that can validate themselves.
type Validatable interface {
	// Validate checks if the component is in a valid state.
	// It returns an error if validation fails.
	Validate() error
}
