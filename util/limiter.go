package util

import "time"

type Limiter struct {
	C   chan time.Time
	Max uint

	closeChan chan struct{}
	refilled  chan struct{}
	remaining uint
}

func NewLimiter(max uint, per time.Duration) *Limiter {
	l := &Limiter{
		C:         make(chan time.Time),
		Max:       max,
		closeChan: make(chan struct{}),
		refilled:  make(chan struct{}),
		remaining: max,
	}

	go l.send()
	go l.tick(per)
	return l
}

func (l *Limiter) Close() error {
	l.closeChan <- struct{}{}
	return nil
}

func (l *Limiter) send() {
	for {
		select {
		case <-l.closeChan:
			break
		default:
			if l.remaining >= 0 {
				l.C <- time.Now()
				l.remaining--
			} else {
				<-l.refilled
			}
		}
	}
}

func (l *Limiter) tick(per time.Duration) {
	ticker := time.NewTicker(per)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			l.remaining = l.Max
			l.refilled <- struct{}{}
		case <-l.closeChan:
			break
		}
	}
}
