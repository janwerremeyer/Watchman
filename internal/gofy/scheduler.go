package gofy

import "time"

type Scheduler struct {
	ticker      *time.Ticker
	stopChannel chan struct{}
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		ticker:      nil,
		stopChannel: nil,
	}

}

func (s *Scheduler) Start(interval time.Duration, task func()) {
	s.ticker = time.NewTicker(interval)
	s.stopChannel = make(chan struct{})

	go func() {
		for {
			select {
			case <-s.ticker.C:
				task()
			case <-s.stopChannel:
				s.ticker.Stop()
				return
			}
		}
	}()
}

func (s *Scheduler) Stop() {
	close(s.stopChannel)
}
