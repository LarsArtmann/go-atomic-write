import type { Feature } from "./types";

export const features: Feature[] = [
  {
    icon: "shield",
    title: "TOCTOU-Safe Writes",
    desc: "Fingerprint verification detects concurrent modification between read and write. No silent overwrites.",
  },
  {
    icon: "lightning",
    title: "xxhash64 Fingerprinting",
    desc: "~11x faster than SHA-256, zero allocations, ~27 GB/s throughput. Memory-bound, not compute-bound.",
  },
  {
    icon: "lock",
    title: "Cross-Platform Locking",
    desc: "flock on Unix, LockFileEx on Windows. Protects across processes, not just goroutines.",
  },
  {
    icon: "refresh",
    title: "Atomic Rename",
    desc: "Single rename(2) on POSIX — no window where the file is missing. MoveFileEx on Windows.",
  },
  {
    icon: "database",
    title: "Crash Durability",
    desc: "fsync on the temp file before rename, fsync on the directory after. Data survives power loss.",
  },
  {
    icon: "folder",
    title: "Minimal Dependencies",
    desc: "Only xxhash for hashing and flock for locking. Both intentional, neither replaceable.",
  },
];
