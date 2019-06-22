package config

// Discord represents all discord-related config
type Discord struct {
	Token  string
	Shards DiscordShards
}

// DiscordShards represents all discord sharding config
type DiscordShards struct {
	Count int
	IDs   []int
}
