
package perf

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "os"
    "strings"
    "testing"
)

var lines = make([]string, 10000)

func init() {
    for i := range lines {
        lines[i] = fmt.Sprintf("log entry %d %s", i, strings.Repeat("x", 100))
    }
}

// --- 1. No I/O ---

func BenchmarkUnbatchedProcessing(b *testing.B) {
    for b.Loop() {
        for _, line := range lines {
            strings.ToUpper(line)
        }
    }
}

func BenchmarkBatchedProcessing(b *testing.B) {
    batchSize := 100
    for b.Loop() {
        for i := 0; i < len(lines); i += batchSize {
            end := i + batchSize
            if end > len(lines) {
                end = len(lines)
            }
            batch := strings.Join(lines[i:end], "|")
            strings.ToUpper(batch)
        }
    }
}

// --- 2. With I/O ---

func BenchmarkUnbatchedIO(b *testing.B) {
    for b.Loop() {
        f, err := os.CreateTemp("", "unbatched")
        if err != nil {
            b.Fatal(err)
        }
        for _, line := range lines {
            _, _ = f.WriteString(line + "\n")
        }
        f.Close()
        os.Remove(f.Name())
    }
}

func BenchmarkBatchedIO(b *testing.B) {
    batchSize := 100
    for b.Loop() {
        f, err := os.CreateTemp("", "batched")
        if err != nil {
            b.Fatal(err)
        }
        for i := 0; i < len(lines); i += batchSize {
            end := i + batchSize
            if end > len(lines) {
                end = len(lines)
            }
            batch := strings.Join(lines[i:end], "\n") + "\n"
            _, _ = f.WriteString(batch)
        }
        f.Close()
        os.Remove(f.Name())
    }
}

// --- 3. With Crypto ---

func hash(s string) string {
    h := sha256.Sum256([]byte(s))
    return hex.EncodeToString(h[:])
}

func BenchmarkUnbatchedCrypto(b *testing.B) {
    for b.Loop() {
        for _, line := range lines {
            hash(line)
        }
    }
}

func BenchmarkBatchedCrypto(b *testing.B) {
    batchSize := 100
    for b.Loop() {
        for i := 0; i < len(lines); i += batchSize {
            end := i + batchSize
            if end > len(lines) {
                end = len(lines)
            }
            joined := strings.Join(lines[i:end], "")
            hash(joined)
        }
    }
}
