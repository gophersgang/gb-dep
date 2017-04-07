package semaphore

// semType controls access to a finite number of resources.
type semType chan struct{}

// New creates a Semaphore that controls access to `n` resources.
func New(n int) semType {
	return semType(make(chan struct{}, n))
}

// Acquire `n` resources.
func (s semType) Acquire(n int) {
	for i := 0; i < n; i++ {
		s <- struct{}{}
	}
}

// Release `n` resources.
func (s semType) Release(n int) {
	for i := 0; i < n; i++ {
		<-s
	}
}
