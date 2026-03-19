package crontab

import (
	"log"
	"math/rand"
	"time"

	"github.com/robfig/cron/v3"
)

func NewCrontabManager(skipJitterOnStart bool) *CrontabManager {
	return &CrontabManager{
		schedule:          cron.New(),
		onStartJobs:       make([]func(), 0),
		skipJitterOnStart: skipJitterOnStart,
	}
}

type CrontabManager struct {
	schedule          *cron.Cron
	onStartJobs       []func()
	skipJitterOnStart bool
}

func (c *CrontabManager) Start() {
	for _, j := range c.onStartJobs {
		go j()
	}
	c.schedule.Start()
}

func (c *CrontabManager) Stop() {
	c.schedule.Stop()
}

func (c *CrontabManager) RegisterCrontab(desc string, runOnStart bool, d time.Duration, f func()) {
	wrapper := func(skipJitter bool) func() {
		return func() {
			if !skipJitter {
				jitter := time.Duration(rand.Int63n(int64(d / 10)))
				log.Printf("%s: jitter %s\n", desc, jitter)
				time.Sleep(jitter)
			}
			log.Printf("%s: begin\n", desc)
			f()
			log.Printf("%s: end\n", desc)
		}
	}
	if runOnStart {
		c.onStartJobs = append(c.onStartJobs, wrapper(c.skipJitterOnStart))
	}
	c.schedule.Schedule(cron.Every(d), cron.FuncJob(wrapper(false)))
}
