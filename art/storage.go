package art

import (
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/xyproto/pinterface"
	"github.com/xyproto/simplebolt"
)

// ErrInternal is the error returned to a user when an unexpected error occurs
// that we perhaps should not surface to a user (i.e. a 500)
var ErrInternal = errors.New("An unexpected error occured and has been logged")

// pendingIDSet is a set containing the sha256 ids of images which have been
// registered but are not yet uploaded
var pendingIDSet pinterface.ISet

// pendingChunks is a map (sha256 id) -> (chunk id) -> (chunk content)
var pendingChunks pinterface.IHashMap

// pendingExpectedLengths is a map (sha256 id) -> (expected length when all chunks submitted)
var pendingExpectedLengths pinterface.IKeyValue

// completed is the query interface for art which has been succesfully uploaded and is ready for download
var completed pinterface.IKeyValue

// InitBoltStorage creates or opens a BoltDB database file at the given path
// as well as initializing the higher-level package-level storage handles
func InitBoltStorage(filename string) {
	db, err := simplebolt.New(filename)
	if err != nil {
		log.Fatal().Msgf("Could not open database: %s", err)
	}

	c := simplebolt.NewCreator(db)

	pendingIDSet, err = c.NewSet("pending_art_ids")
	if err != nil {
		log.Fatal().Msgf("Could not create or load pendingChunks id set: %s", err)
	}

	pendingChunks, err = c.NewHashMap("pending_art")
	if err != nil {
		log.Fatal().Msgf("Could not create or load pendingChunks map: %s", err)
	}

	pendingExpectedLengths, err = c.NewKeyValue("pending_total_lengths")
	if err != nil {
		log.Fatal().Msgf("Could not create or load pending expected lengths storage: %s", err)
	}

	completed, err = c.NewKeyValue("art")
	if err != nil {
		log.Fatal().Msgf("Could not create or load complete art storage: %s", err)
	}
}

func handleInternalError(msg string, err error) error {
	log.Error().Msgf("%s: %v", msg, err)
	return ErrInternal
}

func errIsNotExists(err error) bool {
	switch err {
	case simplebolt.ErrBucketNotFound:
		fallthrough
	case simplebolt.ErrKeyNotFound:
		fallthrough
	case simplebolt.ErrDoesNotExist:
		return true
	default:
		return false
	}
}
