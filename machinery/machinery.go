package machinery

//
//import (
//	"context"
//	"fmt"
//	"github.com/RichardKnop/machinery/v2"
//	ifaceBackend "github.com/RichardKnop/machinery/v2/backends/iface"
//	redisbackend "github.com/RichardKnop/machinery/v2/backends/redis"
//	ifaceBroker "github.com/RichardKnop/machinery/v2/brokers/iface"
//	redisbroker "github.com/RichardKnop/machinery/v2/brokers/redis"
//	"github.com/RichardKnop/machinery/v2/config"
//	eagerlock "github.com/RichardKnop/machinery/v2/locks/eager"
//	ifaceLock "github.com/RichardKnop/machinery/v2/locks/iface"
//	"github.com/RichardKnop/machinery/v2/tasks"
//)
//
//type RedisOpt struct {
//	Broker                 []string
//	DB                     int
//	MaxIdle                int
//	IdleTimeout            int
//	ReadTimeout            int
//	WriteTimeout           int
//	ConnectTimeout         int
//	NormalTasksPollPeriod  int
//	DelayedTasksPollPeriod int
//}
//
//const (
//	BrokerRedis = iota
//	BrokerAMQP
//	BrokerSQS
//)
//
//const (
//	BackendRedis = iota
//	BackendQMQP
//	BackendMongo
//)
//
//type Config struct {
//	Broker      string
//	BrokerType  int
//	Backend     string
//	BackendType int
//	ExpireIn    int
//	Queue       string
//	Concurrency int
//	Redis       RedisOpt
//}
//
//type Option func(*Config)
//
//// 生成并发布任务到broker
//
//type Server struct {
//	queue string
//	*machinery.Server
//}
//
//func NewServer(opts ...Option) *Server {
//	cfg := &Config{
//		Broker:   "",
//		Backend:  "",
//		ExpireIn: 3600,
//		Queue:    "machinery",
//		Redis: RedisOpt{
//			Broker:                 []string{"redis://CtHHQNbFkXpw33ew@192.168.3.13:2379"},
//			DB:                     0,
//			MaxIdle:                3,
//			IdleTimeout:            240,
//			ReadTimeout:            15,
//			WriteTimeout:           15,
//			ConnectTimeout:         15,
//			NormalTasksPollPeriod:  1000,
//			DelayedTasksPollPeriod: 500,
//		},
//	}
//	for _, opt := range opts {
//		opt(cfg)
//	}
//
//	c := &config.Config{
//		// amqp
//		Broker:       cfg.Broker,
//		DefaultQueue: cfg.Queue,
//		// amqp
//		ResultBackend:   cfg.Backend,
//		ResultsExpireIn: cfg.ExpireIn, //s
//		Redis: &config.RedisConfig{
//			MaxIdle:                cfg.Redis.MaxIdle,
//			IdleTimeout:            cfg.Redis.IdleTimeout,
//			ReadTimeout:            cfg.Redis.ReadTimeout,
//			WriteTimeout:           cfg.Redis.WriteTimeout,
//			ConnectTimeout:         cfg.Redis.ConnectTimeout,
//			NormalTasksPollPeriod:  cfg.Redis.NormalTasksPollPeriod,
//			DelayedTasksPollPeriod: cfg.Redis.DelayedTasksPollPeriod,
//		},
//	}
//	var (
//		broker  ifaceBroker.Broker
//		backend ifaceBackend.Backend
//		lock    ifaceLock.Lock
//	)
//	// default redis
//	broker = redisbroker.NewGR(c, cfg.Redis.Broker, cfg.Redis.DB)
//	backend = redisbackend.NewGR(c, cfg.Redis.Broker, cfg.Redis.DB)
//	lock = eagerlock.New()
//	srv := machinery.NewServer(c, broker, backend, lock)
//	return &Server{Server: srv, queue: c.DefaultQueue}
//}
//
//func (s *Server) RegisterTask(n string, f interface{}) error {
//	return s.Server.RegisterTask(n, f)
//}
//
//func (s *Server) RegisterTasks(fs map[string]interface{}) error {
//	return s.Server.RegisterTasks(fs)
//}
//
//// spec : "0 6 * * ?"
//
//func (s *Server) RegisterPeriodTask(spec string, opts ...TaskOption) error {
//	sig := &tasks.Signature{}
//	for _, opt := range opts {
//		opt(sig)
//	}
//	return s.Server.RegisterPeriodicTask(spec, "", sig)
//}
//
//func (s *Server) SendTask(opts ...TaskOption) error {
//	sig := &tasks.Signature{}
//	for _, opt := range opts {
//		opt(sig)
//	}
//	_, err := s.Server.SendTask(sig)
//	return err
//
//	//result, err := s.Server.SendTask(sig)
//	//if err != nil {
//	//	return err
//	//}
//
//	// 异步阻塞轮询backend获取任务执行结果
//	//value, err := result.Get(5 * time.Millisecond)
//	//if err != nil {
//	//	return err
//	//}
//	//tasks.HumanReadableResults(value)
//	//return err
//}
//
//// 执行任务保存状态结果到backend
//
//type Worker struct {
//	*machinery.Worker
//}
//
//func (s *Server) NewWorker(tag string, concurrency int) *Worker {
//	w := s.Server.NewCustomQueueWorker(tag, concurrency, s.queue)
//	// callback
//	w.SetPreTaskHandler(func(s *tasks.Signature) {
//		fmt.Println("start ------------")
//	})
//	w.SetPostTaskHandler(func(s *tasks.Signature) {
//		fmt.Println("end --------------")
//	})
//	return &Worker{Worker: w}
//}
//
//func (w *Worker) Start() error {
//	return w.Launch()
//}
//
//func (w *Worker) Stop(ctx context.Context) error {
//	w.Quit()
//	return nil
//}
//
//// 任务编排 group chain chord
