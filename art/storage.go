package art

import (
	"errors"
	"log"

	"github.com/xyproto/pinterface"
	"github.com/xyproto/simplebolt"
)

var ErrInternal = errors.New("An unexpected error occured and has been logged")

var pendingIDSet pinterface.ISet
var pendingChunks pinterface.IHashMap
var pendingExpectedLengths pinterface.IKeyValue

var completed pinterface.IKeyValue

func InitBoltStorage(filename string) {
	db, err := simplebolt.New(filename)
	if err != nil {
		log.Fatalf("Could not open database: %s", err)
	}

	c := simplebolt.NewCreator(db)

	pendingIDSet, err = c.NewSet("pending_art_ids")
	if err != nil {
		log.Fatalf("Could not create or load pendingChunks id set: %s", err)
	}

	pendingChunks, err = c.NewHashMap("pending_art")
	if err != nil {
		log.Fatalf("Could not create or load pendingChunks map: %s", err)
	}

	pendingExpectedLengths, err = c.NewKeyValue("pending_total_lengths")
	if err != nil {
		log.Fatalf("Could not create or load pending expected lengths storage: %s", err)
	}

	completed, err = c.NewKeyValue("art")
	if err != nil {
		log.Fatalf("Could not create or load complete art storage: %s", err)
	}
}

func handleInternalError(msg string, err error) error {
	log.Printf("Error: %s: %v", msg, err)
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
