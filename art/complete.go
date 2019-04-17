package art

// this file contains function(s) for transitioning art from pending* to completed

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
)

//maybeComplete checks if all chunks have been uploaded and then calls complete if they have
func maybeComplete(id string) error {
	sLength, err := pendingExpectedLengths.Get(id)
	if err != nil {
		return handleInternalError(fmt.Sprintf("fetching expected length of %s", id), err)
	}
	expectedLength, err := strconv.Atoi(sLength)
	if err != nil {
		return handleInternalError(fmt.Sprintf("converting expected length of %s to int", id), err)
	}

	// XXX: One of my "extra" goals was to write this so that you could swap out xyproto/simpleredis
	// for simplebolt and make this a distributed service without corectness problems / data
	// races.
	//
	// To achieve that, absent xyproto/simple* supporting multi op transactions or adding own locks,
	// we can't guard against two simultaneous "passing" checks for an existing chunk
	// causing a running length calculation from double counting the same chunk.
	//
	// With a memory-mapped BoltDB - it ought to still be easily fast enough in practice
	//
	// If we mod the program to actually enforce chunk size, this will be O(n) in number of chunks
	// rather than the current O(n*c) for number of chunks * chunk size. Maybe. I'm not sure
	// whether bolt's Go-slice oriented storage already has lengths precomputed. If so, it's a wash.
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

//complete adds the art to the completed id -> content collection,
// then removes it from the pending data structures.
//It does not verify that the piece is fully uploaded. For that, call maybeComplete
// (which will call complete if appropriate)
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

	//Absent transactions, there really is no perfect answer about the order for removing expectedLength & ID.
	// But if we managed to delete the length and not the ID, we could, from a chunk submission, attempt
	// a complete without having a length to use to determine it. That would be a problem.
	//Leaving length without ID should just waste storage, so it seems preferable to handle it second.
	pendingExpectedLengths.Del(id)
	if err != nil {
		return handleInternalError(fmt.Sprintf("deleting expected length of %s", id), err)
	}

	return nil
}
