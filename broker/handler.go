package broker

// EventHandler represents a function that handles an event
type EventHandler = func(string, []byte)
