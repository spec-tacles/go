package config

// Broker represents all broker-related configuration
type Broker struct {
	URL    string
	Groups BrokerGroups
}

// BrokerGroups represents all the group names for different services on the broker
type BrokerGroups struct {
	Gateway string
}
