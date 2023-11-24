package machinery

//
//import (
//	"fmt"
//	"testing"
//)
//
//func HelloWord(a, b int) error {
//	fmt.Println("+++++++++++++++++:", a+b)
//	return nil
//}
//
//func TestMachinery(t *testing.T) {
//	srv := NewServer()
//	srv.RegisterTask("ADD-TEST", HelloWord)
//	srv.SendTask(Args([]Arg{{
//		Name:  "",
//		Type:  "int",
//		Value: 3,
//	},
//		{
//			Name:  "",
//			Type:  "int",
//			Value: 9,
//		},
//	}))
//	w := srv.NewWorker("worker-test", 1)
//	err := w.Start()
//	if err != nil {
//		fmt.Println("error:", err)
//		return
//	}
//}
