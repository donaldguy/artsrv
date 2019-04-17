package art

import (
	"errors"
	"fmt"
	"log"
	"strconv"
)

var (
	ErrAlreadyExists = errors.New("Art with the given id already exists")

	ErrAlreadyRegistered = errors.New("Art with the given ID is already being uploaded")

	ErrNotRegistered = errors.New("No art with the given ID has been registered")

	ErrChunkAlreadySubmitted = errors.New("A different chunk for this art with this chunk ID was already submitted")
)

func IsPending(id string) (bool, error) {
	return pendingIDSet.Has(id)
}

func Register(id string, totalSize int) error {
	exists, err := Exists(id)
	if err != nil {
		return handleInternalError("checking prior existence", err)
	}

	if exists {
		return ErrAlreadyExists
	}

	isPending, err := IsPending(id)
	if err != nil {
		return handleInternalError("checking prior registration", err)
	}

	if isPending {
		return ErrAlreadyRegistered
	}

	err = pendingIDSet.Add(id)
	if err != nil {
		return handleInternalError(fmt.Sprintf("adding %s as pending id", id), err)
	}

	err = pendingExpectedLengths.Set(id, strconv.Itoa(totalSize))
	if err != nil {
		return handleInternalError("setting expected length", err)
	}

	return nil
}

func SubmitChunk(id, chunkID, chunk string) error {
	isPending, err := IsPending(id)
	if err != nil {
		return handleInternalError("checking registration in submit", err)
	}

	if !isPending {
		return ErrNotRegistered
	}

	oldChunk, err := pendingChunks.Get(id, chunkID)
	if err != nil && errIsNotExists(err) {
		// it doesn't exist; everything is good
	} else if err == nil {
		if oldChunk == chunk {
			log.Printf("Warning: Redundant submission of chunk")
		} else {
			return ErrChunkAlreadySubmitted
		}
	} else {
		return handleInternalError("checking for previous submission of chunk", err)
	}

	err = pendingChunks.Set(id, chunkID, chunk)
	if err != nil {
		return handleInternalError(
			fmt.Sprintf("saving chunk %s of %s", chunkID, id),
			err,
		)
	}

	return maybeComplete(id)
}
