package disco

const surveyPort uint = 5677

// NewSurveyor constructs a new Surveyor.
func NewSurveyor(addr string) *Surveyor {
	return &Surveyor{}
}

// Surveyor surveys the network for a service.
type Surveyor struct {
}

// Survey looks for a service at the specified address.
func (s *Surveyor) Survey(service ...string) ([]*Endpoint, error) {
}
