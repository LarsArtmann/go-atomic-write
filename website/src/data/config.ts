const modulePath = "github.com/larsartmann/go-atomic-write";

export const siteConfig = {
  name: "go-atomic-write",
  title: "go-atomic-write — Crash-Safe, Race-Free File Writes for Go",
  description:
    "Crash-safe, race-free file writes for Go. Detect concurrent modification instead of silently overwriting data.",
  siteUrl: "https://atomicwrite.lars.software",
  modulePath,
  github: `https://${modulePath}`,
  installCommand: `go get ${modulePath}@latest`,
  author: {
    name: "LarsArtmann",
    url: "https://larsartmann.com/",
  },
  pkgGoDev: `https://pkg.go.dev/${modulePath}`,
} as const;
