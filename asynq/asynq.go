package asynq

import (
	"context"
	"github.com/hibiken/asynq"
	"time"
)

type Opt struct {
	redisOpt asynq.RedisClientOpt
	// server
	config asynq.Config
	hfs    []HandleFunc
	hs     []Handler
	// client
	schedulerOpt *asynq.SchedulerOpts
	schedulers   []Scheduler
}

type Server struct {
	*asynq.Server
	*asynq.ServeMux
}

type HandleFunc struct {
	Pattern string
	Handler func(ctx context.Context, task *asynq.Task) error
}

type Handler struct {
	Pattern string
	Handler asynq.Handler
}

func NewServer(opts ...Option) *Server {
	option := &Opt{
		redisOpt: asynq.RedisClientOpt{
			DB:   0,
			Addr: "127.0.0.1:6379",
		},
		config: asynq.Config{},
	}
	for _, opt := range opts {
		opt(option)
	}
	srv := asynq.NewServer(option.redisOpt, option.config)
	mux := asynq.NewServeMux()
	for _, hf := range option.hfs {
		mux.HandleFunc(hf.Pattern, hf.Handler)
	}
	for _, h := range option.hs {
		mux.Handle(h.Pattern, h.Handler)
	}
	return &Server{
		Server:   srv,
		ServeMux: mux,
	}
}

func (s *Server) Start() error {
	return s.Server.Start(s.ServeMux)
}

func (s *Server) Stop(ctx context.Context) error {
	s.Server.Shutdown()
	return nil
}

type Client struct {
	*asynq.Client
}

func NewClient(opts ...Option) *asynq.Client {
	option := &Opt{
		redisOpt: asynq.RedisClientOpt{
			DB:   0,
			Addr: "127.0.0.1:6379",
		},
	}
	for _, opt := range opts {
		opt(option)
	}
	return asynq.NewClient(option.redisOpt)
}

// enqueue

type Scheduler struct {
	TypeName string
	Payload  []byte
	Cron     string // 00 02 * * * 每天凌晨2点
}

func NewScheduler(opts ...Option) *asynq.Scheduler {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	option := &Opt{
		redisOpt: asynq.RedisClientOpt{
			DB:   0,
			Addr: "127.0.0.1:6379",
		},
		schedulerOpt: &asynq.SchedulerOpts{
			Location: location,
		},
	}
	for _, opt := range opts {
		opt(option)
	}
	ns := asynq.NewScheduler(option.redisOpt, option.schedulerOpt)
	for _, scheduler := range option.schedulers {
		_, err = ns.Register(scheduler.Cron, asynq.NewTask(scheduler.TypeName, scheduler.Payload))
		if err != nil {
			panic(err)
		}
	}
	err = ns.Start()
	if err != nil {
		panic(err)
	}
	return ns
}
