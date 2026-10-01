# Passive Providers

ReconHawk's passive provider layer is pluggable.

Included providers:

- crt.sh certificate transparency
- HackerTarget hostsearch
- AlienVault OTX passive DNS
- Common Crawl historical URLs

Providers are queried only for discovery. ReconHawk does not treat a discovered hostname as automatically in scope.

Provider failures are isolated: a failing public provider does not have to stop the complete passive collection.

External provider availability and terms can change; production integrations should handle rate limits, errors and provider-specific policies.
