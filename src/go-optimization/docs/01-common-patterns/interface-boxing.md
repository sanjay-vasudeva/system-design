# Avoiding Interface Boxing

Go’s interfaces make it easy to write flexible, decoupled code. But behind that convenience is a detail that can trip up performance: when a concrete value is assigned to an interface, Go wraps it in a hidden structure—a process called interface boxing.

In many cases, boxing is harmless. But in performance-sensitive code—like tight loops, hot paths, or high-throughput services—it can introduce hidden heap allocations, extra memory copying, and added pressure on the garbage collector. These effects often go unnoticed during development, only showing up later as latency spikes or memory bloat.

## What is Interface Boxing?

Interface boxing refers to the process of converting a concrete value to an interface type. In Go, an interface value consists of two parts internally:

- A **type descriptor**, which describes the concrete type.
- A **data pointer**, which points to the actual value.

When you assign a value to an interface variable, Go creates this two-word structure under the hood. If the value is a non-pointer (a struct or a primitive) and is not already heap-allocated, Go **may** allocate a copy of it on the heap. This is especially relevant when the value is large or when it's stored in a slice of interfaces.

Here’s a simple example:

```go
var i interface{}
i = 42
```

In this case, the integer `42` is boxed into an interface: Go stores the type information (`int`) and a copy of the value `42`. This is inexpensive for small values like `int`, but for large structs, the cost becomes non-trivial.

Another example:

```go
type Shape interface {
    Area() float64
}

type Square struct {
    Size float64
}

func (s Square) Area() float64 { return s.Size * s.Size }

func main() {
    var shapes []Shape
    for i := 0; i < 1000; i++ {
        s := Square{Size: float64(i)}
        shapes = append(shapes, s) // boxing occurs here
    }
}
```

!!! warning
    **Pay attention to this code!** In this example, even though `shapes` is a slice of interfaces, each `Square` value is copied into an interface when appended to `shapes`. If `Square` were a large struct, this would introduce 1000 allocations and large memory copying.

To avoid that, you could pass pointers:

```go
        shapes = append(shapes, &s) // avoids large struct copy
```

This way, only an 8-byte pointer is stored in the interface, reducing both allocation size and copying overhead.

## Why It Matters

In tight loops or high-throughput paths, such as unmarshalling JSON, rendering templates, or processing large collections, interface boxing can degrade performance by triggering unnecessary heap allocations and increasing GC pressure. This overhead is especially costly in systems with high concurrency or real-time responsiveness constraints.

Boxing can also make profiling and benchmarking misleading, since allocations attributed to innocuous-looking lines may actually stem from implicit conversions to interfaces.


## Benchmarking Impact

For the benchmarking we will define an interface and a struct with a significant payload that implements the interface.

```go
{%
    include-markdown "01-common-patterns/src/interface-boxing_test.go"
    start="// interface-start"
    end="// interface-end"
%}
```

### Boxing Large Structs

To demonstrate the real impact of boxing large values vs. pointers, we benchmarked the cost of assigning 1,000 large structs to an interface slice:

```go
{%
    include-markdown "01-common-patterns/src/interface-boxing_test.go"
    start="// bench-slice-start"
    end="// bench-slice-end"
%}
```

Benchmark Results

| Benchmark                | Time per op (ns) | Bytes per op | Allocs per op |
|--------------------------------|---------|-----------|-----------|
| BoxedLargeSliceGrowth          | 404,649 | ~4.13 MB  | 1011      |
| PointerLargeSliceGrowth        | 340,549 | ~4.13 MB  | 1011      |

Boxing large values is significantly slower—about 19% in this case—due to the cost of copying the entire 4KB struct for each interface assignment. Boxing a pointer, however, avoids that cost and keeps the copy small (just 8 bytes). While both approaches allocate the same overall memory (since all values escape to the heap), pointer boxing has clear performance advantages under pressure.

### Passing to a Function That Accepts an Interface

Another common source of boxing is when a large value is passed directly to a function that accepts an interface. Even without storing to a slice, boxing will occur at the call site.

```go
{%
    include-markdown "01-common-patterns/src/interface-boxing_test.go"
    start="// bench-call-start"
    end="// bench-call-end"
%}
```

Benchmark Results

| Benchmark                | ns/op   | B/op  | allocs/op |
|--------------------------|---------|--------|-----------|
| CallWithValue            | 422.5   | 4096   | 1         |
| CallWithPointer          | 379.9   | 4096   | 1         |

Passing a value to a function expecting an interface causes boxing, copying the full struct and allocating it on the heap. In our benchmark, this results in approximately 11% higher CPU cost compared to using a pointer. Passing a pointer avoids copying the struct, reduces memory movement, and results in smaller, more cache-friendly interface values, making it the more efficient choice in performance-sensitive scenarios.

??? example "Show the complete benchmark file"
    ```go
    {% include "01-common-patterns/src/interface-boxing_test.go" %}
    ```

## When Interface Boxing Is Acceptable

Despite its performance implications in some contexts, interface boxing is often perfectly reasonable—and sometimes preferred.

### When abstraction is more important than performance
Interfaces enable decoupling and modularity. If you're designing a clean, testable API, the cost of boxing is negligible compared to the benefit of abstraction.

```go
type Storage interface {
    Save([]byte) error
}
func Process(s Storage) { /* ... */ }
```

### When values are small and boxing is allocation-free
Boxing small, copyable values like `int`, `float64`, or small structs typically causes no allocations.

```go
var i interface{}
i = 123 // safe and cheap
```

### When values are short-lived
If the boxed value is used briefly (e.g. for logging or interface-based sorting), the overhead is minimal.

```go
fmt.Println("value:", someStruct) // implicit boxing is fine
```

### When dynamic behavior is required
Interfaces allow runtime polymorphism. If you need different types to implement the same behavior, boxing is necessary and idiomatic.

```go
for _, s := range []Shape{Circle{}, Square{}} {
    fmt.Println(s.Area())
}
```

Use boxing when it supports clarity, reusability, or design goals—and avoid it only in performance-critical code paths.

## How to Avoid Interface Boxing

- Use pointers when assigning to interfaces. If the method set requires a pointer receiver or the value is large, explicitly pass a pointer to avoid repeated copying and heap allocation.
    ```go
    for i := range tasks {
       result = append(result, &tasks[i]) // Avoids boxing copies
    }
    ```
- Avoid interfaces in hot paths. If the concrete type is known and stable, avoid interface indirection entirely—especially in compute-intensive or allocation-sensitive functions.
- Use type-specific containers. Instead of `[]interface{}`, prefer generic slices or typed collections where feasible. This preserves static typing and reduces unnecessary allocations.
- Benchmark and inspect with pprof. Use `go test -bench` and `pprof` to observe where allocations occur. If the allocation site is in `runtime.convT2E` (convert T to interface), you're likely boxing.