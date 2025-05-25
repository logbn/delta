package mvfifo

type Option func(*Cache)

func WithMaxSizeBytes(s int) Option {
	return func(c *Cache) {
		c.maxSize = max(0, s)
	}
}
