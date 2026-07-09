import type { StepCard, ComparisonItem, UseCase } from "./types";

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
    code: "atomicwrite.Write(path, data, fp)",
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

export const comparisons: ComparisonItem[] = [
  {
    variant: "os.WriteFile",
    accent: false,
    pros: ["Zero dependencies"],
    cons: ["No TOCTOU protection", "Partial writes on crash", "No concurrent-write safety"],
  },
  {
    variant: "go-atomic-write",
    accent: true,
    pros: [
      "TOCTOU-safe via fingerprint",
      "Crash-durable via fsync",
      "Concurrent-safe via flock",
      "Atomic rename — no missing files",
      "Only 2 dependencies",
    ],
    cons: [],
  },
  {
    variant: "DIY",
    accent: false,
    pros: ["No external deps"],
    cons: [
      "Two-rename window (not atomic)",
      "No fingerprint verification",
      "Manual fsync + lock handling",
    ],
  },
];

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
