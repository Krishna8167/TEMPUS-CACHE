##  TempusCache

<p align="center">
  <img 
    src="https://github.com/user-attachments/assets/0a2442e4-f640-43fb-a2a9-b8d219d1ae51" 
    alt="TempusCache Logo" 
    width="380"
  />
</p>

> A high-performance, concurrent in-memory cache for Go with TTL, LRU eviction, capacity limits, and dual expiration strategy.

TempusCache is a lightweight, thread-safe caching engine designed for backend systems that require:

- Low-latency access  
- Predictable memory bounds  
- Time-based lifecycle control  
- Safe concurrent usage  

It combines **O(1) lookup**, **LRU eviction**, and **per-key TTL** in a clean, extensible architecture.

---

#  Features

-  Thread-safe read/write operations  
-  Per-entry TTL support  
-  LRU eviction policy  
-  Configurable maximum capacity  
-  Never returns stale data (lazy expiration)  
-  Optional background cleanup worker  
-  Runtime statistics (hits, misses, evictions)  
-  Functional configuration model  
-  Race-condition safe (`go test -race`)  
-  Benchmark-tested performance  

---

 Installation
===============

`go get github.com/yourusername/tempuscache`

 Quick Start
==============

Import
------

`import "github.com/yourusername/tempuscache"`

> Only use `/v2` if you have tagged a v2 release.

* * * * *

Create Cache
------------

`cache := tempuscache.New(
    tempuscache.WithCleanupInterval(10 * time.Second),
    tempuscache.WithMaxEntries(1000),
)`

-   If no cleanup interval is provided → lazy expiration only
-   If no max entries are provided → unbounded growth

* * * * *

Set Value
---------

`cache.Set("user:1", "Krishna", 5*time.Second)`

-   `ttl > 0` → expires after duration
-   `ttl == 0` → never expires

* * * * *

Get Value
---------

`value, found := cache.Get("user:1")
if found {
    fmt.Println(value)
}`

Expired entries are automatically removed and treated as cache misses.

* * * * *

Delete Value
------------

`cache.Delete("user:1")`

* * * * *

Retrieve Stats
--------------

`stats := cache.Stats()
fmt.Println(stats.Hits, stats.Misses, stats.Evictions)`

* * * * *

Graceful Shutdown
-----------------

`cache.Stop()`

Stops the background cleanup goroutine (if configured).

* * * * *

 Architecture
===============

TempusCache combines two core data structures:

|     Component             |       Purpose          |
| ------------------------- | ---------------------- |
| `map[string]*list.Element`| O(1) key lookup        |
| `*list.List`              | Maintains LRU ordering |
| `sync.RWMutex`            | Concurrency control    |
| Background Janitor        | Active expiration      |

### Storage Model

`Map (key → list element)
        ↓
Doubly Linked List (LRU ordering)
        ↓
Item { key, value, expiration }`

This hybrid structure ensures:

-   O(1) lookup
-   O(1) eviction
-   O(1) recency updates

* * * * *

 Expiration Strategy
=====================

TempusCache implements a **dual expiration model**.

Lazy Expiration (Always Enabled)
--------------------------------

-   Checked during `Get()`
-   Expired entries are removed immediately
-   Guarantees stale data is never returned

Active Expiration (Optional)
----------------------------

-   Configurable cleanup interval
-   Background goroutine scans and deletes expired entries
-   Prevents memory retention of unused keys

If cleanup is disabled, lazy expiration alone guarantees correctness.

* * * * *

 Concurrency Model
====================

The cache uses `sync.RWMutex`:

-   `Lock()` → writes & internal mutations
-   `RLock()` → read-only access

Guarantees:

-   No concurrent map write panic
-   No race conditions
-   Safe multi-goroutine access

Verified using:

`go test -race ./...`

* * * * *

 Configuration (Functional Options)
=====================================

`cache := tempuscache.New(
    tempuscache.WithCleanupInterval(5 * time.Second),
    tempuscache.WithMaxEntries(500),
)`

### Why Functional Options?

-   Stable constructor signature
-   Explicit configuration
-   Extensible design
-   Clean API surface

* * * * *

 Performance
==============

| Operation    | Complexity  |
| ------------ | ------------|
| Set          |        O(1) |
| Get          |        O(1) |
| Delete       |        O(1) |
| LRU Eviction |        O(1) |
| Cleanup Cycle|        O(n) |

Cleanup performs a full scan during active expiration.

For moderate workloads, performance remains predictable and efficient.

* * * * *

 Testing & Benchmarking
=========================

Run tests:

`go test ./...`

Run race detection:

`go test -race ./...`

Run benchmarks:

`go test -bench=. -benchmem`

Benchmarks report:

-   `ns/op`

-   `B/op`

-   `allocs/op`

Parallel benchmarks simulate concurrent read workloads.

* * * * *

 Use Cases
============

-   API response caching
-   Session/token storage
-   Rate limiting metadata
-   In-memory microservice cache layer
-   Temporary computation results

* * * * *

 Roadmap
==========

Future enhancements:

-   Sharded locking (reduced contention)
-   Pluggable eviction strategies
-   Prometheus metrics exporter
-   Context-aware operations
-   Distributed cache mode

* * * * *

 Design Philosophy
====================

TempusCache is intentionally:

-   Minimal
-   Predictable
-   Idiomatic Go
-   Concurrency-safe
-   Performance-aware

It is built to demonstrate **backend systems engineering principles**, not just simple key-value storage.

* * * * *

 License
==========

MIT License

 `---
If you want next level polish, I can now:

- Add GitHub badges (Go version, CI, License)
- Add benchmark result table
- Write a strong "Resume Project Description"
- Help you tag v1.0.0 properly

Now this project actually looks serious.`

from installation to graceful shutdown in markdown file format

Here is the section from **Installation → Graceful Shutdown**, properly formatted in Markdown so you can directly paste it into your `README.md`.

* * * * *

`#  Installation

```bash
go get github.com/yourusername/tempuscache `

* * * * *

 Quick Start
==============

Import
------

`import "github.com/yourusername/tempuscache"`

> Only use `/v2` if you have tagged a v2 release.

* * * * *

Create Cache
------------

`cache := tempuscache.New(
    tempuscache.WithCleanupInterval(10 * time.Second),
    tempuscache.WithMaxEntries(1000),
)`

-   If no cleanup interval is provided → lazy expiration only

-   If no max entries are provided → unbounded growth

* * * * *

Set Value
---------

`cache.Set("user:1", "Krish", 5*time.Second)`

-   `ttl > 0` → expires after duration

-   `ttl == 0` → never expires

* * * * *

Get Value
---------

`value, found := cache.Get("user:1")
if found {
    fmt.Println(value)
}`

Expired entries are automatically removed and treated as cache misses.

* * * * *

Delete Value
------------

`cache.Delete("user:1")`

* * * * *

Retrieve Stats
--------------

`stats := cache.Stats()
fmt.Println(stats.Hits, stats.Misses, stats.Evictions)`

* * * * *

Graceful Shutdown
-----------------

`cache.Stop()`

Stops the background cleanup goroutine (if configured).