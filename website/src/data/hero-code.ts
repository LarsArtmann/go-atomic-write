import { siteConfig } from "./config";

const importPath = siteConfig.github.replace("https://github.com/", "github.com/");

export const heroCode = `package main

import (
    "os"

    atomicwrite "${importPath}"
)

func main() {
    path := "/etc/app/config.json"

    data, _ := os.ReadFile(path)
    fp := atomicwrite.FingerprintFromBytes(data)

    newData := []byte(\`{"updated": true}\`)

    err := atomicwrite.WriteVerified(path, newData, fp)
    // err == ErrConcurrentModification if someone
    // else wrote between your read and write
    _ = err
}`;
