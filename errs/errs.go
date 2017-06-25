package errs

import "fmt"

func Recovered(e interface{}) error {
	if e == nil {
		return nil
	}

	switch v := e.(type) {
	case error:
		return v
	default:
		return fmt.Errorf("%v", e)
	}
}
