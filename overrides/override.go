package overrides

// Overrider is an interface that has a Send method.
// It is used to override the printing of logs in the line or JSON logger.
type Overrider interface {
	Send(string)
}
