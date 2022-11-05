package api

// Action is the common interface for all actions
type Action interface {
	Execute() error
}
