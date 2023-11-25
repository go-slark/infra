package machinery

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func HelloWord(a, b int) error {
	fmt.Println("+++++++++++++++++:", a+b)
	return nil
}

func TestTask(t *testing.T) {
	srv := NewServer()
	err := srv.RegisterTask("ADD-TEST", HelloWord)
	if err != nil {
		fmt.Println("register task error:", err)
		return
	}
	err = srv.SendTask(TaskName("ADD-TEST"), Args([]Arg{{
		Name:  "",
		Type:  "int",
		Value: 3,
	},
		{
			Name:  "",
			Type:  "int",
			Value: 9,
		},
	}))
	if err != nil {
		fmt.Println("send task error:", err)
		return
	}
	w := srv.NewWorker("worker-test", 1)
	err = w.Start()
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	w.Stop(context.TODO())
}

func TestDelayTask(t *testing.T) {
	srv := NewServer()
	err := srv.RegisterTask("ADD-Delay", HelloWord)
	if err != nil {
		fmt.Println("register task error:", err)
		return
	}
	err = srv.SendTask(TaskName("ADD-Delay"), Delay(time.Now().UTC().Add(30*time.Second)), Args([]Arg{{
		Name:  "",
		Type:  "int",
		Value: 30,
	},
		{
			Name:  "",
			Type:  "int",
			Value: 9,
		},
	}))
	if err != nil {
		fmt.Println("send task error:", err)
		return
	}
	w := srv.NewWorker("worker-test", 1)
	err = w.Start()
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	w.Stop(context.TODO())
}

func TestPeriodTask(t *testing.T) {
	srv := NewServer()
	err := srv.RegisterPeriodTask("*/1 * * * *", TaskName("ADD"))
	if err != nil {
		fmt.Println("register task error:", err)
		return
	}
	err = srv.SendTask(TaskName("ADD"), Args([]Arg{{
		Name:  "",
		Type:  "int",
		Value: 10,
	},
		{
			Name:  "",
			Type:  "int",
			Value: 9,
		},
	}))
	if err != nil {
		fmt.Println("send task error:", err)
		return
	}
	w := srv.NewWorker("worker-test", 1)
	err = w.Start()
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	w.Stop(context.TODO())
}
