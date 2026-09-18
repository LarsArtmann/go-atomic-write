import type { StepCard, UseCase, ComparisonMatrix, FaqItem } from "./types";

export const steps: StepCard[] = [
  {
    step: "1",
    stepColor: "accent",
    title: "Fingerprint",
    desc: "Compute an xxhash64 digest when you read the file.",
    code: "fp := atomicwrite.FingerprintFile(path)",
  },
  {
    step: "2",
    stepColor: "accent",
    title: "Stage + fsync",
    desc: "Write new content to a unique temp file, then fsync for durability.",
    code: "atomicwrite.WriteVerified(path, data, fp)",
  },
  {
    step: "3",
    stepColor: "amber",
    title: "Lock + Verify",
    desc: "Acquire exclusive lock, re-read the file, verify fingerprint still matches.",
    code: "// flock.Lock() + Fingerprint.Matches()",
  },
  {
    step: "4",
    stepColor: "amber",
    title: "Atomic Rename",
    desc: "Single rename(2) replaces the target. Directory is fsync'd (POSIX).",
    code: "// os.Rename(tmp, path) + syncDir()",
  },
];

export const comparisonMatrix: ComparisonMatrix = {
  columns: ["os.WriteFile", "DIY", "go-atomic-write"],
  rows: [
    { feature: "TOCTOU-safe", values: ["no", "no", "yes"] },
    { feature: "Crash-durable (fsync)", values: ["no", "no", "yes"] },
    { feature: "Concurrent-write safe", values: ["no", "no", "yes"] },
    { feature: "Atomic rename", values: ["no", "partial", "yes"] },
    { feature: "Fingerprint verification", values: ["no", "no", "yes"] },
    { feature: "Cross-platform locking", values: ["no", "no", "yes"] },
    { feature: "Dependencies", values: ["0", "0", "2"] },
  ],
};

export const useCases: UseCase[] = [
  {
    title: "Config Files",
    desc: "Read-modify-write cycles safe from concurrent edits",
    icon: "cog",
  },
  {
    title: "State & Checkpoints",
    desc: "Corruption-free persistence for databases and jobs",
    icon: "chart",
  },
  {
    title: "CI/CD Pipelines",
    desc: "Concurrent writers on shared artifacts and lock files",
    icon: "refresh",
  },
];

export const faqItems: FaqItem[] = [
  {
    question: "Why not just write a temp file and rename it?",
    answer:
      "That handles crashes and missing files, but not races. If another process changes the file between your read and your write, one of the two updates is silently lost. go-atomic-write locks the file, re-reads it, and verifies the fingerprint first. On a mismatch you get ErrConcurrentModification instead of a lost update.",
  },
  {
    question: "When should I not use this?",
    answer:
      "If one process owns the file, never crashes mid-write, and you do not need durability across reboots, plain os.WriteFile is simpler. This library earns its keep when a file is shared between processes or has to survive a crash.",
  },
  {
    question: "What does a write cost?",
    answer:
      "One xxhash64 pass (~27 GB/s, zero allocations), one advisory lock, and one full re-read of the file while the lock is held. Negligible for config, state, and checkpoint files; measurable for multi-gigabyte files.",
  },
  {
    question: "What happens when two writers race?",
    answer:
      "Exactly one wins. The loser gets ErrConcurrentModification and the original file is left untouched, because the temp file is discarded before any rename. Writers holding fresh data can retry; writers that do not should surface the error.",
  },
  {
    question: "Does it work on Windows?",
    answer:
      "Yes. LockFileEx replaces flock, a failed os.Rename is retried five times with exponential backoff on ERROR_ACCESS_DENIED and ERROR_SHARING_VIOLATION (antivirus, open handles), and directory fsync is a no-op because Windows has no equivalent.",
  },
  {
    question: "Do I have to change how I read files?",
    answer:
      "No. Reads stay as they are. You only change the write: compute a fingerprint when you read, then call WriteVerified. WriteIfChanged skips the write entirely when the bytes already match, which keeps file watchers and mtimes quiet on no-op runs.",
  },
];
