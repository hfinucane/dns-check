# a dns checker

Features:
* Has a deadline- will not block for a large amount of time in a crisis
* Supports monitoring multiple DNS servers so you don't wind up running a
  distributed monitoring service for Cloudflare or something.
* Checks multiple DNS servers concurrently
* Configurable warning/criticality

Examples:
```
# Alert if Cloudflare doesn't have myhost.tld
dns-check -host myhost.tld 1.1.1.1

# Alert if more than one of Cloudflare, Google, and OpenDNS don't have myhost.tld
dns-check -w 50 -c 50 -host myhost.tld 1.1.1.1 8.8.8.8 208.67.222.222

# Warn if the popular public DNS resolvers doesn't have myhost.tld, and
# go critical if more than one don't have it.
dns-check -w 10 -c 50 -host myhost.tld 1.1.1.1 8.8.8.8 208.67.222.222

# Allow for 30 seconds
dns-check -deadline 30 -host myhost.tld 1.1.1.1
```
