package art

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
)

func maybeComplete(id string) error {
	sLength, err := pendingExpectedLengths.Get(id)
	if err != nil {
		return handleInternalError(fmt.Sprintf("fetching expected length of %s", id), err)
	}
	expectedLength, err := strconv.Atoi(sLength)
	if err != nil {
		return handleInternalError(fmt.Sprintf("converting expected length of %s to int", id), err)
	}

	// XXX: One of my "extra" goals was to write this so that you could swap out simpleredis
	// for simplebolt and make this a distributed service without corectness problems or data
	// races.
	//
	// To achieve that, absent xyproto/simple* supporting multi op transactions,
	// we can't guard against two simultaneous checks for existing chunks fouling a running length
	// calculation. So I'mma gonna do this less-than-ideally-performant re-total.
	//
	// With a memory-mapped boltDB - it ought to still be easily fast enough
	//
	// We can do a bit better if we actually enforce chunk size
	totalLength := 0
	keys, err := pendingChunks.Keys(id)
	if err != nil {
		return handleInternalError(fmt.Sprintf("fetching keys for '%s' in maybeComplete", id), err)
	}

	for _, k := range keys {
		chunk, err := pendingChunks.Get(id, k)
		if err != nil {
			return handleInternalError(
				fmt.Sprintf("fetching chunk '%s' of '%s' in maybeComplete", id, k),
				err,
			)
		}

		totalLength += len(chunk)
	}

	if totalLength >= expectedLength {
		return complete(id)
	}

	return nil
}

func complete(id string) error {
	var art bytes.Buffer

	sKeys, err := pendingChunks.Keys(id)
	if err != nil {
		return handleInternalError(fmt.Sprintf("fetching keys for '%s' in complete", id), err)
	}

	iKeys := make([]int, len(sKeys))
	for i, v := range sKeys {
		iKeys[i], err = strconv.Atoi(v)
		if err != nil {
			return handleInternalError("converting keys to int for sort", err)
		}
	}

	sort.Ints(iKeys)

	for _, k := range iKeys {
		chunk, err := pendingChunks.Get(id, strconv.Itoa(k))
		if err != nil {
			return handleInternalError(
				fmt.Sprintf("fetching chunk '%d' of '%s' in complete", k, id),
				err,
			)
		}

		art.WriteString(chunk)
	}

	err = completed.Set(id, art.String())
	if err != nil {
		return handleInternalError(fmt.Sprintf("Setting completed art %s", id), err)
	}
	err = pendingChunks.Del(id)
	if err != nil {
		return handleInternalError(fmt.Sprintf("deleting pending chunks for %s", id), err)
	}
	err = pendingIDSet.Del(id)
	if err != nil {
		return handleInternalError(fmt.Sprintf("deleting registration for %s", id), err)
	}
	pendingExpectedLengths.Del(id)
	if err != nil {
		return handleInternalError(fmt.Sprintf("deleting expected length of %s", id), err)
	}

	return nil
}
