package main

import "github.com/robfig/cron/v3"

type Scheduler struct {
	cron   *cron.Cron
	sender func()
}

func NewScheduler(sender func()) *Scheduler {
	return &Scheduler{
		cron:   cron.New(),
		sender: sender,
	}
}

func (s *Scheduler) Add(spec string) error {
	_, err := s.cron.AddFunc(spec, s.sender)
	return err
}

func (s *Scheduler) Start() {
	s.cron.Start()
}

func (s *Scheduler) Stop() error {
	s.cron.Stop()
	return nil
}
