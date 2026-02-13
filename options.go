package tempuscache

import (
	"time"
)

/*
Option represents a functional configuration modifier for Cache.

================================================================================
DESIGN PATTERN: FUNCTIONAL OPTIONS
================================================================================

TempusCache uses the Functional Options Pattern — an idiomatic Go
approach for flexible and future-proof configuration.

Instead of passing multiple constructor parameters, New() accepts
a variadic list of Option functions:

    cache := New(
        WithCleanupInterval(10 * time.Second),
    )

Each Option is a function that mutates the Cache instance
during initialization.

================================================================================
WHY THIS PATTERN?
================================================================================

1. API STABILITY
   - Adding new configuration fields does not change the New() signature.
   - Prevents breaking changes.

2. READABILITY
   - Configuration is explicit and self-documenting.
   - Avoids confusing positional constructor arguments.

3. EXTENSIBILITY
   - New features (e.g., capacity limits, eviction strategies,
     metrics hooks, logging integrations) can be added seamlessly.

4. COMPOSABILITY
   - Multiple options can be combined in a clear and modular way.

================================================================================
ENGINEERING PHILOSOPHY
================================================================================

The constructor remains minimal and stable,
while configuration logic remains modular and isolated.

This pattern is widely used in production Go libraries
for long-term maintainability.
*/

type Option func(*Cache)

/*
WithCleanupInterval configures the active expiration frequency.

================================================================================
PARAMETER
================================================================================

d (time.Duration):
    Interval at which the background janitor scans
    and removes expired entries.

================================================================================
BEHAVIOR
================================================================================

If d > 0:
    - A background janitor goroutine is started.
    - Expired entries are periodically removed.
    - Enables active expiration strategy.

If d <= 0:
    - The janitor is disabled.
    - The cache relies solely on lazy expiration during Get().

================================================================================
PERFORMANCE TRADE-OFFS
================================================================================

Short intervals:
    - Faster cleanup of expired entries
    - Increased CPU usage due to frequent scans

Long intervals:
    - Lower CPU overhead
    - Expired items may occupy memory longer

================================================================================
SYSTEM DESIGN CONSIDERATION
================================================================================

Choosing an appropriate cleanup interval depends on:

- Cache size
- TTL distribution
- Memory sensitivity
- Throughput requirements

This option provides operational control over
the balance between performance and memory efficiency.
*/

func WithCleanupInterval(d time.Duration) Option {
	return func(c *Cache) {
		c.interval = d
	}
}

/*
WithMaxEntries configures the maximum number of entries
allowed in the cache before LRU eviction is triggered.

================================================================================
PARAMETER
================================================================================

n (int):
    Maximum number of entries permitted in the cache.

================================================================================
BEHAVIOR
================================================================================

If n > 0:
    - The cache enforces a hard capacity limit.
    - When inserting a new key and the limit is reached:
        → The least recently used (LRU) entry is evicted.
        → Eviction statistics are incremented.

If n <= 0:
    - The cache operates without capacity restriction.
    - No LRU-based eviction will occur due to size constraints.

================================================================================
EVICTION STRATEGY
================================================================================

TempusCache uses a strict LRU (Least Recently Used) policy:

- Most recently accessed entries move to the front.
- Least recently used entries remain at the back.
- The back element is evicted first when capacity is exceeded.

Eviction occurs in O(1) time due to the doubly linked list design.

================================================================================
SYSTEM DESIGN CONSIDERATION
================================================================================

Setting a maximum entry limit allows:

- Predictable memory usage
- Protection against unbounded growth
- Stable performance under heavy load

Capacity tuning should consider:

- Available memory
- Average item size
- Expected workload characteristics
- TTL distribution

This option enables bounded, production-ready cache behavior.
*/

func WithMaxEntries(n int) Option {
	return func(c *Cache) {
		c.maxEntries = n
	}
}
