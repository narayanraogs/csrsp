// Package parameters handles the validation of header parameters.
package parameters

import "csrspServer/session"

// Validator defines the interface for any parameter validation logic.
type Validator interface {
	ParamName() string
	Process(input <-chan []byte, store *session.Store, key string)
}
