package art

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/rs/zerolog/log"
)

var (
	//ErrAlreadyExists fires when attempting to register art that was previously succesfully uploaded
	ErrAlreadyExists = errors.New("Art with the given id already exists")

	//ErrAlreadyRegistered fires when attempting to register art that was previously registed and is thus currently being uploaded
	ErrAlreadyRegistered = errors.New("Art with the given ID is already being uploaded")

	//ErrNotRegistered fires when an unexpected chunk upload is attempted for an unkown art ID
	ErrNotRegistered = errors.New("No art with the given ID has been registered")

	//ErrChunkAlreadySubmitted fires when a chunk with the same ID but different content is uploaded.
	//We choose to log and ignore if an identical chunk is resubmitted.
	ErrChunkAlreadySubmitted = errors.New("A different chunk for this art with this chunk ID was already submitted")
)

//IsPending returns whether or not art with a given ID has been registered but not yet completed
func IsPending(id string) (bool, error) {
	return pendingIDSet.Has(id)
}

//Register prepares storage for chunk uploads for the given ID and expected finalSize
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

	err = pendingExpectedLengths.Set(id, strconv.Itoa(totalSize))
	if err != nil {
		return handleInternalError("setting expected length", err)
	}

	// We register length before ID last as it determines rejection of a duplicate request.
	// This should prevent a mid-request halt from leading to e.g. registered images of unknown length
	err = pendingIDSet.Add(id)
	if err != nil {
		return handleInternalError(fmt.Sprintf("adding %s as pending id", id), err)
	}

	return nil
}

//SubmitChunk stores the content of `chunk` in pending storage for art of `id`.
// If the provided chunk causes the calculated length to reach the totalLength declared at registration,
// this chunk and those previously submitted are assembled sorted by chunkID and inserted into the completed
// collection such that `Get(id)` should return it.
func SubmitChunk(id string, chunkID int, chunk string) error {
	isPending, err := IsPending(id)
	if err != nil {
		return handleInternalError("checking registration in submit", err)
	}

	if !isPending {
		return ErrNotRegistered
	}

	sChunkID := strconv.Itoa(chunkID)
	oldChunk, err := pendingChunks.Get(id, sChunkID)
	if err != nil && errIsNotExists(err) {
		// it doesn't exist; everything is good
	} else if err == nil {
		if oldChunk == chunk {
			log.Warn().Msg("Redundant submission of chunk")
		} else {
			return ErrChunkAlreadySubmitted
		}
	} else {
		return handleInternalError("checking for previous submission of chunk", err)
	}

	err = pendingChunks.Set(id, sChunkID, chunk)
	if err != nil {
		return handleInternalError(
			fmt.Sprintf("saving chunk %d of '%s'", chunkID, id),
			err,
		)
	}

	return maybeComplete(id)
}
