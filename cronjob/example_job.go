package cronjob

import "fmt"

type ExampleJob struct {
}

func (g ExampleJob) Run() {
	fmt.Println("hello world")
}
