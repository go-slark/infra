package asynq

import (
	"time"
)

type Option func(*Opt)

func Address(addr string) Option {
	return func(o *Opt) {
		o.redisOpt.Addr = addr
	}
}

func Password(password string) Option {
	return func(o *Opt) {
		o.redisOpt.Password = password
	}
}

func DialTimeout(tm time.Duration) Option {
	return func(o *Opt) {
		o.redisOpt.DialTimeout = tm
	}
}

func ReadTimeout(tm time.Duration) Option {
	return func(o *Opt) {
		o.redisOpt.ReadTimeout = tm
	}
}

func WriteTimeout(tm time.Duration) Option {
	return func(o *Opt) {
		o.redisOpt.WriteTimeout = tm
	}
}

func DB(db int) Option {
	return func(o *Opt) {
		o.redisOpt.DB = db
	}
}

func Concurrency(c int) Option {
	return func(o *Opt) {
		o.config.Concurrency = c
	}
}

func Queues(queues map[string]int) Option {
	return func(o *Opt) {
		o.config.Queues = queues
	}
}

func HandlerFunc(hfs []HandleFunc) Option {
	return func(o *Opt) {
		o.hfs = hfs
	}
}

func SetHandler(hs []Handler) Option {
	return func(o *Opt) {
		o.hs = hs
	}
}

func Location(location *time.Location) Option {
	return func(o *Opt) {
		o.schedulerOpt.Location = location
	}
}

func SetScheduler(scheduler []Scheduler) Option {
	return func(o *Opt) {
		o.schedulers = scheduler
	}
}
