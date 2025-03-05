package models

type TaskResult struct {
	TaskID string  `json:"id"`
	Result float64 `json:"result"`
}

type Argument struct {
	Value        float64
	Ready        bool
	ParentTaskID uint32
}

type Task struct {
	ID            uint32
	ExpressionID  uint32
	ParentArgID   uint32
	Arg1          *Argument
	Arg2          *Argument
	Operation     string
	OperationTime int
}

type TaskResponse struct {
	ID            string  `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int     `json:"operation_time"`
}

func (task *Task) IsReady() bool {
	return task.Arg1.Ready && task.Arg2.Ready
}
