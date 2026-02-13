package tempuscache

import (
	"fmt"
	"testing"
	"time"
)

/*
BenchmarkSet measures the performance of the Set() operation
under repeated overwrites of the same key.

================================================================================
OBJECTIVE
================================================================================

This benchmark evaluates the steady-state cost of:

- Expiration timestamp calculation (time.Now + ttl)
- Mutex Lock()/Unlock() overhead
- Map update (existing key path)
- LRU move-to-front operation

================================================================================
BENCHMARK MODEL
================================================================================

- The same key is repeatedly overwritten.
- Map size remains constant.
- No eviction pressure.
- Minimal memory growth.

This isolates the hot write path without structural expansion effects.

================================================================================
GO BENCHMARK MECHANICS
================================================================================

The testing framework dynamically determines b.N
to achieve stable timing results.

Use:

    go test -bench=. -benchmem

to observe:
- ns/op
- B/op
- allocs/op

================================================================================
INTERPRETATION
================================================================================

This benchmark reflects ideal write throughput
in a stable cache state.
*/

func BenchmarkSet(b *testing.B) {
	cache := New()

	for i := 0; i < b.N; i++ {
		cache.Set("key", "value", 5*time.Second)
	}
}

/*
BenchmarkSetUnique measures write performance
when inserting unique keys.

================================================================================
OBJECTIVE
================================================================================

Unlike BenchmarkSet, this benchmark:

- Forces map growth.
- Exercises LRU insert path.
- Tests capacity-bound behavior.

With WithMaxEntries(b.N + 1),
evictions are avoided to measure pure growth cost.

================================================================================
WHAT IT CAPTURES
================================================================================

- Map allocation behavior
- Linked list node allocation
- Memory growth impact
- Lock overhead under expanding state

This benchmark represents a more realistic
write-heavy workload scenario.
*/

func BenchmarkSetUnique(b *testing.B) {
	cache := New(WithMaxEntries(b.N + 1))

	for i := 0; i < b.N; i++ {
		cache.Set(fmt.Sprintf("key%d", i), i, 0)
	}
}

/*
BenchmarkGet measures the performance of the read path.

================================================================================
OBJECTIVE
================================================================================

Evaluates the cost of:

- Map lookup
- Expiration check
- LRU move-to-front operation
- Mutex overhead
- Stats increment

================================================================================
SCENARIO
================================================================================

- A single key is preloaded.
- Repeated Get() calls simulate high read locality.
- No expiration or eviction occurs.

This benchmark reflects hot-cache read performance.
*/

func BenchmarkGet(b *testing.B) {
	cache := New()

	cache.Set("key", "value", 0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get("key")
	}
}

/*
BenchmarkParallelGet measures read performance
under concurrent access.

================================================================================
OBJECTIVE
================================================================================

Simulates high-concurrency read workloads
using b.RunParallel.

Evaluates:

- Lock contention behavior
- Throughput scaling across CPU cores
- Stability under parallel goroutines

================================================================================
WHY THIS MATTERS
================================================================================

Caches are typically read-heavy systems.
Understanding concurrent read scalability
is critical for backend performance analysis.

Run with:

    go test -bench=. -cpu=1,2,4,8

to observe scaling behavior.
*/

func BenchmarkParallelGet(b *testing.B) {
	cache := New()

	cache.Set("key", "value", 0)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cache.Get("key")
		}
	})
}

/*
BenchmarkEviction measures write performance
under constant eviction pressure.

================================================================================
SCENARIO
================================================================================

- Cache capacity is limited (WithMaxEntries(100)).
- Continuous unique inserts exceed capacity.
- LRU eviction is triggered repeatedly.

================================================================================
WHAT IT EVALUATES
================================================================================

- O(1) eviction behavior
- removeElement() efficiency
- Map delete performance
- List removal cost
- Stats increment overhead

================================================================================
SYSTEM INSIGHT
================================================================================

This benchmark reflects memory-bounded
production scenarios where eviction
is part of the steady-state workload.
*/

func BenchmarkEviction(b *testing.B) {
	cache := New(WithMaxEntries(100))

	for i := 0; i < b.N; i++ {
		cache.Set(fmt.Sprintf("key%d", i), i, 0)
	}
}
