// Package fixture is a tiny module for testing the docgen action.
package fixture

import "io"

// Greeter writes greetings.
type Greeter struct{}

// Greet writes a greeting to w.
func (Greeter) Greet(w io.Writer) error {
	_, err := io.WriteString(w, "hello")
	return err
}
