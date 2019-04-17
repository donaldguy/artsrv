# artsrv

A service for multipart uploads and all-at-once downloads of ASCII art

## Usage

`make` will build and run with docker as described in instructions

`make test` will run the test binary



## Design Choices

### Scope of durability, recoverability

In the absence of specification about preferred trade-offs between performance and reliability guarantees, I aimed for a server which provided correctness and durability when restarted between any two requests—congruent with a system that would undergo blue-green deployment with a drain operation. 

This meant keeping no information scoped greater than a single request in memory (vs e.g. buffering chunks in memory into the correct order as received, which a fixed chunk size would definitely allow). BoltDB operates as a memory mapped file, so with OS filesystem caching, there is likely not significant overhead from the mere use of disk over memory. There is, however, the need created (absent transactions) to make some linear scans of the chunk collection to decide on completeness, and key sorting (and string ↔️ int conversion) to establish after the fact ordering. 

If we were explicitly willing to throw away any incompletely-uploaded art on server restart, I'd go the other way. Some sort of reloadable, write-through, in-memory→disk buffer would provide a best of both worlds approach, but that would—in my opinion—constitute overkill for this exercise.

For most requests, there is also little danger that mid-request interruption could lead to an inconsistent state. The exception is in the completion process following upload of the last chunk. This can last ~10s of ms or more depending on total art size. Thus, it is the most likely culprit/victim— a mid-request halt could result in art with all the chunks uploaded, but the art not available for download. With the current code, resubmission of any valid chunk after start would generate an in log warning, but also be sufficient to cause completion. 

Alternatively, a loop like

```go
pendingArt, _ := pendingChunks.All() // error handling excluded
for _, id := range pendingArt {
    go maybeComplete(id)
}
```

could be added to initialization to clean these up.

Both this limited inconsistency danger & the poor completeness-check-performance could be eliminated by proper use of batch transactions (which would e.g. record a chunk and its length to a running total). These are supported by raw BoltDB, but not simplebolt. For the sake of this assignment, I preferred the simpler API, and was more enthralled—in theory—by the idea of being able to make this code distributable/load-balance-able just by swapping out [simplebolt](https://github.com/xyproto/simplebolt) with [simpleredis](https://github.com/xyproto/simpleredis) than the performance gain from such bookkeeping. 

### Permissive implementation of API

From the given spec, I decided to ignore:

 1. The fact that IDs were sha256 sums of the content; any string can be used (though I'd change the "sha256" field on the registration API if I were making this official)

 2. All handling of chunk sizes. Chunk size doesn't need to be declared at registration or submission time. Data in a chunk upload can be any length and can vary freely between chunks.

Neither change affected correctness, nor was keeping them necessary for a performant implementation in Go. Chunk Size would be helpful for more performant in-memory assembly of uploaded chunks, but I sacrificed this in favor of more free restart-ability, as described above. 

Though I do not believe the test binary provided manifested any such behaviors, it was unclear how one should handle:

- data greater than declared chunk size (should it be rejected? truncated?)
- data shorter than declared chunk size (should it be padded?)
- chunks, other than the final, different than chunk size at registration time

One downside of ignoring chunk size is that it removes a potentially useful lever of rate-limiting. As it stands, the limit of submittable chunk size and total size is determined by Go http / gin internals and/or memory/disk constraints. (I googled it and am surprised to find there is no RFC limit on allowable POST body size)

### Functional Style

It is notable that this program is all `func Verb(noun, noun)` style rather than containing any more OO `func (*subject) Verb (object)` style Go. 

In general, this was more a result of what I happened to type than a specific choice, but I believe I went that way based on the relatively "stateless" role this code is playing between BoltDB and the HTTP clients. If I had done more in-memory operations, I think I would have favored a style more heavily around structs and "methods".

I also think that if this server were destined to have more than one object type in its REST hierarchy it would also probably be better to make it more OO style to clear up some namespace pollution.