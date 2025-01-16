package asana

import "fmt"

type TooManyRequestsError struct {
	RetryAfter int
}

func (e TooManyRequestsError) Error() string {
	return fmt.Sprintf("too many requests, retry after %d", e.RetryAfter)
}
