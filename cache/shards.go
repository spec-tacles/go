package cache

// ShardInfo represents a cache of sharding information
type ShardInfo interface {
	Count() int
	Gateway() string
}
