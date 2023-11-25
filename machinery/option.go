package machinery

import (
	"github.com/RichardKnop/machinery/v2/tasks"
	"time"
)

type TaskOption func(signature *tasks.Signature)

func Retry(times int) TaskOption {
	return func(s *tasks.Signature) {
		s.RetryCount = times
	}
}

func RetryTimeout(tm int) TaskOption {
	return func(s *tasks.Signature) {
		s.RetryTimeout = tm
	}
}

// time.Now().Add(time.Second * 5): 延迟任务

func Delay(tm time.Time) TaskOption {
	return func(s *tasks.Signature) {
		s.ETA = &tm
	}
}

func TaskName(name string) TaskOption {
	return func(s *tasks.Signature) {
		s.Name = name
	}
}

type Arg struct {
	Name  string
	Type  string
	Value interface{}
}

func Args(args []Arg) TaskOption {
	return func(s *tasks.Signature) {
		for _, arg := range args {
			s.Args = append(s.Args, tasks.Arg{
				Name:  arg.Name,
				Type:  arg.Type,
				Value: arg.Value,
			})
		}
	}
}
