# Persistent Result Cache

ReconHawk can persist scan reports in a local SQLite database.

The cache key includes the target and scan options, so changing worker/depth/timeout settings creates a different entry.

Example integration:

    cache, err := recon.OpenCache("reconhawk.db")
    if err != nil { ... }
    defer cache.Close()

    report, hit, err := cache.Get(target, options, 30*time.Minute)
    if err != nil { ... }
    if !hit {
        report, err = recon.Scan(options)
        if err != nil { ... }
        _ = cache.Put(target, options, report)
    }

Use a bounded TTL when cached reconnaissance data should expire.

The cache stores scan results only. It does not automatically authorize newly discovered hosts.
