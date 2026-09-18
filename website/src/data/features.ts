import type { Feature } from "./types";

export const features: Feature[] = [
  {
    icon: "shield",
    title: "No silent overwrites",
    desc: "A fingerprint check catches concurrent modification between your read and your write, so you get an error instead of a lost update.",
  },
  {
    icon: "lightning",
    title: "Zero-cost fingerprinting",
    desc: "~11x faster than SHA-256, zero allocations, ~27 GB/s throughput. Memory-bound, not compute-bound.",
  },
  {
    icon: "lock",
    title: "Safe across processes",
    desc: "flock on Unix, LockFileEx on Windows. Protects across processes, not just goroutines.",
  },
  {
    icon: "refresh",
    title: "The file is never missing",
    desc: "Single rename(2) on POSIX — no window where the file is missing. MoveFileEx on Windows.",
  },
  {
    icon: "database",
    title: "Survives power loss",
    desc: "fsync on the temp file before rename, fsync on the directory after. Data survives power loss.",
  },
  {
    icon: "folder",
    title: "Two deliberate dependencies",
    desc: "Only xxhash for hashing and flock for locking. Both intentional, neither replaceable.",
  },
];
