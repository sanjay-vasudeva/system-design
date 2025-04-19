## Lazy Initialization for Performance in Go

In Go, some resources are expensive to initialize, or simply unnecessary unless certain code paths are triggered. That’s where lazy initialization becomes useful: it defers the construction of a value until the moment it’s actually needed. This pattern can improve performance, reduce startup overhead, and avoid unnecessary work—especially in high-concurrency applications.

### Why Lazy Initialization Matters

Initializing complex resources—such as database connections, caches, or large data structures—at application startup can significantly delay launch time and unnecessarily consume memory. Lazy initialization ensures these resources are only created when needed, optimizing resource usage and performance.

Additionally, lazy initialization is crucial when you have code that might be executed multiple times, but you need a resource or logic executed precisely once. This pattern helps ensure idempotency and avoids redundant processing.

### Using `sync.Once` for Thread-Safe Initialization

Go provides the `sync.Once` type to implement lazy initialization safely in concurrent environments:

```go
var (
	resource *MyResource
	once     sync.Once
)

func getResource() *MyResource {
	once.Do(func() {
		resource = expensiveInit()
	})
	return resource
}
```

In this example, the function `expensiveInit()` executes exactly once, no matter how many goroutines invoke `getResource()` concurrently. This ensures thread-safe initialization without additional synchronization overhead.

### Using `sync.OnceValue` and `sync.OnceValues` for Initialization with Output Values

Since Go 1.21, if your initialization logic returns a value, you might prefer using `sync.OnceValue` (single value) or `sync.OnceValues` (multiple values) for simpler, more expressive code:

```go
var getResource = sync.OnceValue(func() *MyResource {
	return expensiveInit()
})

func processData() {
	res := getResource()
	// use res
}
```

Here, `sync.OnceValue` neatly encapsulates initialization logic and directly returns the initialized value, eliminating explicit state management.

For scenarios where your initialization returns multiple values, `sync.OnceValues` offers a clean and efficient solution:

```go
var getConfig = sync.OnceValues(func() (*Config, error) {
	return loadConfig("config.yml")
})

func processData() {
	config, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}
	// use config
}
```

Choosing `sync.OnceValue` or `sync.OnceValues` helps you clearly express initialization logic with direct value returns, whereas `sync.Once` remains best suited for general scenarios requiring flexible initialization logic without immediate value returns.

### Custom Lazy Initialization with Atomic Operations

Yes, it’s technically possible to replace `sync.Once`, `sync.OnceValue`, or `sync.OnceFunc` with custom logic using low-level atomic operations. This approach may offer slightly finer control or avoid allocations in extremely performance-critical code paths.

**That said, it’s rarely worth the tradeoff.**

Manual atomic-based initialization is more error-prone, harder to read, and easier to get wrong—especially when concurrency and memory visibility guarantees are involved. For the vast majority of cases, `sync.Once*` is safer, clearer, and performant enough.

!!! info
	If you’re convinced that atomic-based lazy initialization is justified in your case, this blog post walks through the details and caveats:  
	:material-hand-pointing-right: [Lazy initialization in Go using atomics](https://goperf.dev/blog/2025/04/03/lazy-initialization-in-go-using-atomics/)

### Performance Considerations

While lazy initialization can offer clear benefits, it also brings added complexity. It’s important to handle initialization carefully to avoid subtle issues like race conditions or concurrency bugs. Using built-in tools like `sync.Once` or `atomic` operations typically ensures thread-safety without much hassle. Still, it’s always a good idea to measure actual improvements through profiling, confirming lazy initialization truly enhances startup speed, reduces memory usage, or boosts your application's responsiveness.

## Benchmarking Impact

There is typically nothing specific to benchmark with lazy initialization itself, as the main benefit is deferring expensive resource creation. The performance gains are inherently tied to the avoided cost of unnecessary initialization, startup speed improvements, and reduced memory consumption, rather than direct runtime throughput differences.

## When to Choose Lazy Initialization

- When resource initialization is costly or involves I/O. Delaying construction avoids paying the cost of setup—like opening files, querying databases, or loading large structures—unless it’s actually needed.
- To improve startup performance and memory efficiency. Deferring work until first use allows your application to start faster and avoid allocating memory for resources that may never be used.
- When not all resources are needed immediately or at all during runtime. Lazy initialization helps you avoid initializing fields or services that only apply in specific code paths.
- To guarantee a block of code executes exactly once despite repeated calls. Using tools like `sync.Once` ensures thread-safe, one-time setup in concurrent environments.