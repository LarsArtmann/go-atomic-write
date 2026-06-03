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

| Term        | Definition                                                          | Context                                                                  |
| ----------- | ------------------------------------------------------------------- | ------------------------------------------------------------------------ |
| Fingerprint | An `[8]byte` xxhash64 digest representing file content at read time | Zero-value means "no prior file existed"                                 |
| `.tmp` file | Staging file written before atomic rename                           | Cleaned up on failure                                                    |
| `.bak` file | Backup of previous file content before overwrite                    | Created on successful non-first writes, overwritten on subsequent writes |

## Commands

| Term                 | Definition                                               | Context                                        |
| -------------------- | -------------------------------------------------------- | ---------------------------------------------- |
| Write                | Stage `.tmp` → lock → verify fingerprint → atomic rename | The main operation                             |
| FingerprintFile      | Compute an xxhash64 digest of a file's current content   | Returns zero-value for nonexistent files       |
| FingerprintFromBytes | Compute an xxhash64 digest from raw bytes                | Used when caller already has content in memory |

## Events

| Term                      | Definition                                                       | Context                                 |
| ------------------------- | ---------------------------------------------------------------- | --------------------------------------- |
| ErrConcurrentModification | The file was modified between fingerprint read and write attempt | Caller should re-read, merge, and retry |
| Successful Write          | `.tmp` renamed to target, `.bak` of previous content created     | Target file is always consistent        |

## Bounded Contexts

Single context — this is a focused library with one bounded context: **safe file writes**.

---

> **How to use this file:**
>
> - Keep terms concise — one clear sentence per definition
> - Update when new domain concepts emerge
> - Use these terms consistently in code, docs, and conversations
> - When in doubt about a word's meaning, check here first
