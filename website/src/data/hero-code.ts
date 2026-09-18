import { siteConfig } from "./config";

const importPath = siteConfig.modulePath;

interface CodeSegment {
  text: string;
  className?: string;
}

const codeSegments: CodeSegment[] = [
  { text: "package", className: "text-text-muted" },
  { text: " " },
  { text: "main", className: "text-text-secondary" },
  { text: "\n\n" },
  { text: "import", className: "text-text-muted" },
  { text: " (\n    " },
  { text: '"os"', className: "text-text-muted" },
  { text: "\n\n    atomicwrite " },
  { text: `"${importPath}"`, className: "text-accent-hover" },
  { text: "\n)\n\n" },
  { text: "func", className: "text-text-muted" },
  { text: " " },
  { text: "main", className: "text-accent" },
  { text: "() {\n    path := " },
  { text: '"/etc/app/config.json"', className: "text-amber" },
  { text: "\n\n    data, _ := os." },
  { text: "ReadFile", className: "text-accent" },
  { text: "(path)\n    fp := atomicwrite." },
  { text: "FingerprintFromBytes", className: "text-accent" },
  { text: "(data)\n\n    newData := []byte(" },
  { text: '`{"updated": true}`', className: "text-amber" },
  { text: ")\n\n    err := atomicwrite." },
  { text: "WriteVerified", className: "text-accent" },
  { text: "(path, newData, fp)\n    " },
  { text: "// err == ErrConcurrentModification if someone", className: "text-code-comment" },
  { text: "\n    " },
  { text: "// else wrote between your read and write", className: "text-code-comment" },
  { text: "\n    _ = err\n}" },
];

function escapeHtml(text: string): string {
  return text.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
}

/**
 * The one and only source of the hero snippet. Both the highlighted markup and
 * the copy-to-clipboard payload are derived from `codeSegments`, so the code a
 * visitor reads can never drift from the code they copy.
 */
export const heroCode = codeSegments.map((segment) => segment.text).join("");

export function highlightedHeroCode(): string {
  return codeSegments
    .map((segment) =>
      segment.className
        ? `<span class="${segment.className}">${escapeHtml(segment.text)}</span>`
        : escapeHtml(segment.text),
    )
    .join("");
}
