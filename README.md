# Multi Value FIFO

A multi value FIFO cache for lists of byte slices.

[![Go Reference](https://godoc.org/github.com/logbn/mvfifo?status.svg)](https://godoc.org/github.com/logbn/mvfifo)
[![License](https://img.shields.io/badge/License-Apache_2.0-dd6600.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/logbn/mvfifo?4)](https://goreportcard.com/report/github.com/logbn/mvfifo)
[![Go Coverage](https://github.com/logbn/mvfifo/wiki/coverage.svg)](https://raw.githack.com/wiki/logbn/mvfifo/coverage.html)

This multi value FIFO cache was created to maintain an in memory segment of an event stream distributed across an array
of topics. Items are evicted in the order in which they are inserted. The length of the cache is bound by the total size
of all items in the cache rather than the total number of items in the cache to more easily manage memory usage.

## Usage

The maximum size of the cache can be specified in bytes (default 64 MiB).

```go
import "github.com/logbn/mvfifo"

c := mvfifo.NewCache(
    mvfifo.WithMaxSizeBytes(1 << 30), // 1 GiB
)
```

Values can be inserted with cursor data

```go
c.Add("test-1", 1, []byte("test-value-a"))
c.Add("test-1", 2, []byte("test-value-b"))
c.Add("test-2", 3, []byte("test-value-c"))
```

Values can be iterated for any key

```go
for cursor, value := range c.Iter("test-1") {
    println(cursor, string(value))
}
// output:
//   1 test-a
//   2 test-b

for cursor, value := range c.Iter("test-2") {
    println(cursor, string(value))
}
// output:
//   1 test-c
```

Values can be iterated for any key after a specified cursor value (non-inclusive)

```go
for cursor, value := range c.IterAfter("test-1", 1) {
    println(cursor, string(value))
}
// output:
//   2 test-b
```

The calling code is responsible for ensuring that cursors increase over time and never decrease.  
The cache makes no attempt to reorder values based on the value of inserted cursors.  
The cursor is used only for range iteration, not for sorting.

## License

Multi Value FIFO is licensed under Apache 2.0
