#!/usr/bin/env python3
import json
import sys

target = sys.argv[1] if len(sys.argv) > 1 else ""
print(json.dumps({"plugin":"python-example","target":target,"type":"enrichment"}))
