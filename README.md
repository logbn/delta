# Delta

A multi value FIFO cache.

[![Go Reference](https://godoc.org/github.com/pantopic/delta?status.svg)](https://godoc.org/github.com/pantopic/delta)
[![License](https://img.shields.io/badge/License-Apache_2.0-dd6600.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/pantopic/delta?4)](https://goreportcard.com/report/github.com/pantopic/delta)
[![Go Coverage](https://github.com/pantopic/delta/wiki/coverage.svg)](https://raw.githack.com/wiki/pantopic/delta/coverage.html)

Delta was created in order to maintain an in memory cache of events across an array of topics. Events are evicted in the
order in which they were inserted. The length of the cache is bound by the total size of all items in the cache rather
than the number of items in the cache in order to more easily manage memory usage.

## Usage

Specify the maximum size of the cache in bytes.

```go
import "github.com/pantopic/delta"

c := delta.New(
    delta.WithMaxSizeBytes(2 << 30), // 1 GiB
)

c.Add("test-1", 1, []byte("test-value-a"))
c.Add("test-1", 2, []byte("test-value-b"))
c.Add("test-2", 3, []byte("test-value-c"))

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

for cursor, value := range c.IterAfter("test-1", 1) {
    println(cursor, string(value))
}
// output:
//   2 test-b
```

## License

Delta is licensed under Apache 2.0
