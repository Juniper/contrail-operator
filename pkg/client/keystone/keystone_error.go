package keystone

import "net/http"

type keystoneError struct {
	msg        string
	statusCode int
}

func newUnauthorized() keystoneError {
	return keystoneError{
		msg:        "not authorized",
		statusCode: http.StatusUnauthorized,
	}
}

func (err keystoneError) Error() string {
	return err.msg
}

func IsUnauthorized(err error) bool {
	kerr, ok := err.(keystoneError)
	if ok {
		return kerr.statusCode == http.StatusUnauthorized
	}
	return false
}
