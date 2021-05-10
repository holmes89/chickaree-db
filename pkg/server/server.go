package server

type Runner interface {
	Run() <-chan error
	Close() error
}
