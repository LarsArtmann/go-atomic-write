import type { StepCard, UseCase, ComparisonMatrix } from "./types";

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
