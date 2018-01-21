package disco

// ProtocolSignature is used to identify some messages belonging to the Pub-Sub protocol.
var ProtocolSignature = [...]byte{0x02, 0x01}

// Survey is a message sent by a Surveyor when searching for a service by name.
type Survey struct {
	Signature   [2]byte
	Major       uint8
	Minor       uint8
	servicesLen uint8
	services    map[string]string
}
