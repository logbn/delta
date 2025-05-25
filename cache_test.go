package mvfifo

import (
	"bytes"
	"fmt"
	"testing"
)

func TestCache(t *testing.T) {
	c := New(WithMaxSizeBytes(3 * (20 + overhead)))
	t.Run(`iter`, func(t *testing.T) {
		c.Add("test-key", 1, []byte("test-value-1"))
		c.Add("test-key", 2, []byte("test-value-2"))
		c.Add("test-key", 3, []byte("test-value-3"))
		var n uint64
		for cur, val := range c.Iter("test-key") {
			n++
			if cur != n {
				t.Fatalf("Item %d returned cur %d", n, cur)
			}
			expected := fmt.Appendf(nil, "test-value-%d", n)
			if bytes.Compare(val, expected) != 0 {
				t.Fatalf("Item %d returned val %s rather than %s", n, val, expected)
			}
		}
		if n != 3 {
			t.Fatalf("Iter returned %d of 3 results", n)
		}
		t.Run(`missing`, func(t *testing.T) {
			var n uint64
			for range c.Iter("test-key-2") {
				n++
			}
			if n != 0 {
				t.Fatalf("Iter returned %d of 0 results", n)
			}
		})
		t.Run(`early`, func(t *testing.T) {
			var n uint64
			for range c.Iter("test-key") {
				n++
				break
			}
			if n != 1 {
				t.Fatalf("Iter returned %d of 1 results", n)
			}
		})
	})
	t.Run(`evict`, func(t *testing.T) {
		c.Add("test-key", 4, []byte("test-value-4"))
		var n uint64 = 0
		for cur, val := range c.Iter("test-key") {
			n++
			if cur != n+1 {
				t.Fatalf("Item %d returned cur %d", n, cur)
			}
			expected := fmt.Appendf(nil, "test-value-%d", n+1)
			if bytes.Compare(val, expected) != 0 {
				t.Fatalf("Item %d returned val %s rather than %s", n, val, expected)
			}
		}
		if n != 3 {
			t.Fatalf("Iter returned %d of 3 results", n)
		}
	})
	t.Run(`after`, func(t *testing.T) {
		var n uint64 = 0
		for cur, val := range c.IterAfter("test-key", 2) {
			n++
			if cur != n+2 {
				t.Fatalf("Item %d returned cur %d", n, cur)
			}
			expected := fmt.Appendf(nil, "test-value-%d", n+2)
			if bytes.Compare(val, expected) != 0 {
				t.Fatalf("Item %d returned val %s rather than %s", n, val, expected)
			}
		}
		if n != 2 {
			t.Fatalf("Iter returned %d of 2 results", n)
		}
		t.Run(`missing`, func(t *testing.T) {
			var n int
			for range c.IterAfter("test-key-2", 2) {
				n++
			}
			if n != 0 {
				t.Fatalf("Iter returned %d of 0 results", n)
			}
		})
		t.Run(`early`, func(t *testing.T) {
			var n int
			for range c.IterAfter("test-key", 0) {
				n++
				break
			}
			if n != 1 {
				t.Fatalf("Iter returned %d of 1 results", n)
			}
		})
	})
	t.Run(`first`, func(t *testing.T) {
		cur, val := c.First()
		if cur != 2 {
			t.Fatalf("Iter returned cursor of %d rather than %d", cur, 2)
		}
		expected := fmt.Appendf(nil, "test-value-%d", 2)
		if bytes.Compare(val, expected) != 0 {
			t.Fatalf("Iter returned value of %s rather than %s", val, expected)
		}
	})
	t.Run(`last`, func(t *testing.T) {
		cur, val := c.Last()
		if cur != 4 {
			t.Fatalf("Iter returned cursor of %d rather than %d", cur, 4)
		}
		expected := fmt.Appendf(nil, "test-value-%d", 4)
		if bytes.Compare(val, expected) != 0 {
			t.Fatalf("Iter returned value of %s rather than %s", val, expected)
		}
	})
	t.Run(`len`, func(t *testing.T) {
		n := c.Len()
		if n != 3 {
			t.Fatalf("Iter returned length of %d rather than 3", n)
		}
	})
	t.Run(`size`, func(t *testing.T) {
		n := c.Size()
		if n != 3*(20+overhead) {
			t.Fatalf("Iter returned size of %d rather than 3", 3*(20+overhead))
		}
	})
	t.Run(`resize`, func(t *testing.T) {
		c.Resize(2 * (20 + overhead))
		var n uint64 = 0
		for cur, val := range c.Iter("test-key") {
			n++
			if cur != n+2 {
				t.Fatalf("Item %d returned cur %d", n, cur)
			}
			expected := fmt.Appendf(nil, "test-value-%d", n+2)
			if bytes.Compare(val, expected) != 0 {
				t.Fatalf("Item %d returned val %s rather than %s", n, val, expected)
			}
		}
		if n != 2 {
			t.Fatalf("Iter returned %d of 2 results", n)
		}
	})
	t.Run(`empty`, func(t *testing.T) {
		c.Resize(0)
		n := c.Len()
		if n != 0 {
			t.Fatalf("Iter returned %d of 0 results", n)
		}
	})
}
