import { siteConfig } from "./config";

const importPath = siteConfig.modulePath;

interface CodeSegment {
  text: string;
  className?: string;
}

const codeSegments: CodeSegment[] = [
  { text: "package", className: "text-editor-muted" },
  { text: " " },
  { text: "main", className: "text-editor-text" },
  { text: "\n\n" },
  { text: "import", className: "text-editor-muted" },
  { text: " (\n    " },
  { text: '"os"', className: "text-editor-muted" },
  { text: "\n\n    atomicwrite " },
  { text: `"${importPath}"`, className: "text-editor-path" },
  { text: "\n)\n\n" },
  { text: "func", className: "text-editor-muted" },
  { text: " " },
  { text: "main", className: "text-editor-fn" },
  { text: "() {\n    path := " },
  { text: '"/etc/app/config.json"', className: "text-editor-string" },
  { text: "\n\n    data, _ := os." },
  { text: "ReadFile", className: "text-editor-fn" },
  { text: "(path)\n    fp := atomicwrite." },
  { text: "FingerprintFromBytes", className: "text-editor-fn" },
  { text: "(data)\n\n    newData := []byte(" },
  { text: '`{"updated": true}`', className: "text-editor-string" },
  { text: ")\n\n    err := atomicwrite." },
  { text: "WriteVerified", className: "text-editor-fn" },
  { text: "(path, newData, fp)\n    " },
  { text: "// err == ErrConcurrentModification if someone", className: "text-editor-comment" },
  { text: "\n    " },
  { text: "// else wrote between your read and write", className: "text-editor-comment" },
  { text: "\n    _ = err\n}" },
];

function escapeHtml(text: string): string {
  return text.replaceAll("&", "&amp;").replaceAll("<", "&lt;").replaceAll(">", "&gt;");
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
