# Domain Language — go-atomic-write

A **Unified Language** for `go-atomic-write` — shared across contributors, users, and AI.
Inspired by Domain-Driven Design (DDD) Ubiquitous Language.

Every term below should mean the **same thing** to everyone who reads it.

## Glossary

| Term                    | Definition                                                                           | Context                                |
| ----------------------- | ------------------------------------------------------------------------------------ | -------------------------------------- |
| Fingerprint             | An xxhash64 digest (`[8]byte`) of file content at a point in time                    | Used to detect concurrent modification |
| TOCTOU                  | Time-of-check-to-time-of-use — the vulnerability class this library addresses        | Security concept                       |
| Atomic Write            | A write that either completes fully or leaves the original file unchanged            | Crash safety guarantee                 |
| Concurrent Modification | A file changed by another process between the fingerprint read and the write attempt | Error condition                        |

## Value Objects

| Term        | Definition                                                          | Context                                                    |
| ----------- | ------------------------------------------------------------------- | ---------------------------------------------------------- |
| Fingerprint | An `[8]byte` xxhash64 digest representing file content at read time | Zero-value means "no prior file existed"                   |
| `.tmp` file | Staging file written before atomic rename                           | Uniquely named (`path + "." + randomHex + ".tmp"`); cleaned up on failure |

## Commands

The API is **split by intent**. A plain write is crash-durable but performs no
TOCTOU check; a `*Verified` write adds fingerprint verification under an exclusive lock.

| Term                | Definition                                                                   | Context                                                                  |
| ------------------- | ---------------------------------------------------------------------------- | ------------------------------------------------------------------------ |
| Write               | Stage `.tmp` → fsync → atomic rename. No TOCTOU check.                       | Common case: temp files, first-time creation, single-writer scenarios    |
| WriteVerified       | Stage `.tmp` → fsync → lock → re-read → verify fingerprint → atomic rename   | TOCTOU protection. Zero-value fingerprint means "file must not yet exist" |
| WriteIfChanged      | Fingerprint disk → skip if unchanged, else `Write`/`WriteVerified`           | Idempotent writes for config files and code generators                   |
| WriteFunc           | Like `Write` but streams content via a callback (64KB buffer)                | Large or incrementally-produced content                                  |
| WriteFuncVerified   | Like `WriteFunc` but with TOCTOU fingerprint verification                    | Streaming content that must be race-safe                                  |
| FingerprintFile     | Compute an xxhash64 digest of a file's current content                       | Returns zero-value for nonexistent files                                  |
| FingerprintFromBytes | Compute an xxhash64 digest from raw bytes                                   | Used when caller already has content in memory                            |

## Events

| Term                      | Definition                                                  | Context                                 |
| ------------------------- | ----------------------------------------------------------- | --------------------------------------- |
| ErrConcurrentModification | The file was modified between fingerprint read and write    | Caller should re-read, merge, and retry |
| Successful Write          | `.tmp` renamed atomically to target; original overwritten   | Target file is always consistent        |

## Bounded Contexts

Single context — this is a focused library with one bounded context: **safe file writes**.

---

> **How to use this file:**
>
> - Keep terms concise — one clear sentence per definition
> - Update when new domain concepts emerge
> - Use these terms consistently in code, docs, and conversations
> - When in doubt about a word's meaning, check here first
