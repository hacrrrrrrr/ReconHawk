package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/hacrrrrrrr/ReconHawk/internal/recon"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] != "scan" {
		usage()
		os.Exit(2)
	}
	fs := flag.NewFlagSet("scan", flag.ExitOnError)
	workers := fs.Int("workers", 8, "maximum concurrent requests")
	timeout := fs.Duration("timeout", 8*time.Second, "per-request timeout")
	depth := fs.Int("depth", 1, "same-origin crawl depth")
	output := fs.String("output", "", "write JSON report to a file")
	jsonl := fs.String("jsonl", "", "write observations as JSONL")
	fs.Parse(os.Args[2:])
	if fs.NArg() != 1 { fmt.Fprintln(os.Stderr, "error: scan requires one target"); usage(); os.Exit(2) }

	report, err := recon.Scan(recon.Options{Target: fs.Arg(0), Workers: *workers, Timeout: *timeout, Depth: *depth})
	if err != nil { fmt.Fprintln(os.Stderr, "error:", err); os.Exit(1) }
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
	data = append(data, '\n')
	if *output != "" {
		if err := os.WriteFile(*output, data, 0o644); err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
	} else { fmt.Print(string(data)) }
	if *jsonl != "" {
		f, err := os.Create(*jsonl)
		if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
		defer f.Close()
		enc := json.NewEncoder(f)
		for _, o := range report.Observations {
			if err := enc.Encode(o); err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
		}
	}
}

func usage() {
	fmt.Println("ReconHawk - automated reconnaissance and attack-surface discovery")
	fmt.Println("Usage: reconhawk scan <target> [flags]")
	fmt.Println("  --workers N   maximum concurrent requests (default 8)")
	fmt.Println("  --timeout D   per-request timeout (default 8s)")
	fmt.Println("  --depth N     same-origin crawl depth (default 1)")
	fmt.Println("  --output F    write complete JSON report")
	fmt.Println("  --jsonl F     write observations as JSONL")
	fmt.Println("Only scan assets you are authorized to test.")
}
