package art

// this file contains functions for operating on the `completed` id -> content KV collection

import (
	"errors"
)

// ErrDoesNotExist is the error thrown when no fully uploaded art is avaliable for a given ID
// it can be thrown even if the art has been registered and is in the process of being uploaded
var ErrDoesNotExist = errors.New("No Art with the given ID is available for download")

// Exists returns whether or not art with the given ID had been fully uploaded.
// In the event of error determining it returns false and the error
func Exists(id string) (bool, error) {
	_, err := completed.Get(id)

	if err != nil {
		if errIsNotExists(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// Get returns fully succefully uploaded art given its ID.
func Get(id string) (string, error) {
	art, err := completed.Get(id)
	if err != nil {
		if errIsNotExists(err) {
			return "", ErrDoesNotExist
		}

		return "", handleInternalError("while getting art", err)
	}

	//TODO?: Special case error for Art which is "pending" (registered but not complete?)
	// HTTP 423? 425?

	return art, nil
}
