package art

import (
	"errors"
)

var ErrDoesNotExist = errors.New("No Art with the given ID is available for download")

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

func Get(id string) (string, error) {
	art, err := completed.Get(id)
	if err != nil {
		if errIsNotExists(err) {
			return "", ErrDoesNotExist
		}

		return "", handleInternalError("while getting art", err)
	}

	return art, nil
}
