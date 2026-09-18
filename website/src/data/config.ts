export const siteConfig = {
  name: "go-atomic-write",
  title: "go-atomic-write — Crash-Safe, Race-Free File Writes for Go",
  description:
    "Crash-safe, race-free file writes for Go. Detect concurrent modification instead of silently overwriting data.",
  siteUrl: "https://atomicwrite.lars.software",
  github: "https://github.com/larsartmann/go-atomic-write",
  author: {
    name: "LarsArtmann",
    url: "https://larsartmann.com/",
  },
  pkgGoDev: "https://pkg.go.dev/github.com/larsartmann/go-atomic-write",
} as const;
