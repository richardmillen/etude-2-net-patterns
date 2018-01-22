package patterns

import "errors"

// ErrInvalidSig occurs when an invalid protocol signature is received.
var ErrInvalidSig = errors.New("invalid protocol signature")
