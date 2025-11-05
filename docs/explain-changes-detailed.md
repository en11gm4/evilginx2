# Detailed explanation of the recent stealth-oriented changes

## Overview
The branch introduces a number of steps intended to make deployments of the framework less recognizable during phishing engagements. The changes are concentrated around three areas:

1. The bootstrap logic (`main.go`), where runtime assets such as the configuration directory are chosen.
2. The HTTP proxy pipeline (`core/http_proxy.go`), which now mutates outbound requests and session metadata to look more organic.
3. The redirector script (`core/scripts.go`), which now produces more varied network patterns on the phishing landing page.

The following sections break down the concrete modifications and their expected effects.

## `main.go`
- **Randomised configuration directory** – when `-c` is not provided, the binary no longer defaults to `~/.evilginx`. Instead it pulls a single byte of crypto-grade entropy and uses it to pick one of four common hidden folders (`~/.cache`, `~/.local`, `~/.config`, or `~/.data`). The actual configuration lives in a `session` subdirectory under the chosen base. This reduces the chance that defenders notice the presence of the framework by looking for the historical `.evilginx` folder.
- The rest of the bootstrap logic stays identical (phishlet loading, database initialisation, blacklist creation, etc.), so operator workflows are unchanged aside from the new configuration location.

## `core/http_proxy.go`
### Request mutation before forwarding upstream
- **Randomised `User-Agent` header** – every proxied request now overwrites the browser-provided user agent with a randomly selected string from a hard-coded pool of modern Chrome, Firefox, and Safari identifiers. The selection uses one byte from `crypto/rand`, guaranteeing uniform distribution without introducing PRNG state.
- **Artificial latency injection** – requests are delayed by a random 50–250 ms pause (`addRandomDelay`). This aims to mimic real-world network jitter and mask the deterministic timing of the proxy.
- **Noise traffic generation** – with a 30 % chance, the proxy spins up a goroutine that logs a faux background request to common assets such as `/favicon.ico` or `/robots.txt`. The target host is chosen from the active phishlet’s proxy hosts. Although the helper only logs the intent right now, it establishes a hook where genuine background fetches could be added later.

### Session bookkeeping adjustments
- **Removed framework-identifying header** – the header previously set via `req.Header.Set(p.getHomeDir(), o_host)` is now intentionally omitted to eliminate a predictable marker that could be matched by defenders.
- **Reworked `getHomeDir`** – the helper now returns the static string `X-Session-Config`, decoupling the rest of the code from the historical `.evilginx` path spelling.
- **Natural-looking cookie names** – `getSessionCookieName` no longer hashes the cookie name. It keeps the first four characters from the upstream cookie, appends a realistic suffix (`id`, `token`, `csrf`, etc.), and finishes with an 8–11 character random segment. This produces identifiers that resemble native session cookies while still being unique.

### Supporting helpers
- New helper functions (`getRandomUserAgent`, `addRandomDelay`, `generateNoiseRequest`) encapsulate the behaviour above and keep the main proxy loop readable.

## `core/scripts.go`
- **Dynamic redirect behaviour** – the embedded JavaScript now randomises the redirect endpoint among `/redirect/`, `/auth/`, `/callback/`, and `/return/`, adds timestamp/random query parameters, and alternates between `window.location.replace` and direct assignment. Retries use a back-off pattern that varies with the HTTP status. These changes make the lure traffic less deterministic and harder to pattern-match.

## Operational notes
- All randomness relies on the Go `crypto/rand` reader, so the patterns are not reproducible across runs.
- Operators should document the new configuration directory layout to avoid confusion when migrating existing setups.
- The newly added helper hooks (especially `generateNoiseRequest`) only log activity today; adding a real HTTP client would be the next step if background traffic needs to be observable on the wire.
