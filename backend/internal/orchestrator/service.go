package orchestrator

import (
	"calc-website/config"
	"calc-website/internal/models"
	"calc-website/pkg/calc"
	"errors"
	"github.com/google/uuid"
	"strconv"
)

var ErrIDTaskNotExists = errors.New("task with this ID does not exist")

var tasksQueue = make(chan *models.TaskResponse, 1024)
var allTasks = make(map[uint32]*models.Task)
var taskArgs = make(map[uint32]*models.Argument)
var allExpressions = make(map[uint32]*models.Expression)

type APIService struct {
	TimeAdditionMs        int
	TimeSubtractionMs     int
	TimeMultiplicationsMs int
	TimeDivisionsMs       int
}

func NewAPIService(cfg *config.Config) *APIService {
	return &APIService{
		TimeAdditionMs:        cfg.TimeAdditionMs,
		TimeSubtractionMs:     cfg.TimeSubtractionMs,
		TimeMultiplicationsMs: cfg.TimeMultiplicationsMs,
		TimeDivisionsMs:       cfg.TimeDivisionsMs,
	}
}

func getOperationTime(s *APIService, operator string) int {
	switch operator {
	case "+":
		return s.TimeAdditionMs
	case "-":
		return s.TimeSubtractionMs
	case "*":
		return s.TimeMultiplicationsMs
	case "/":
		return s.TimeDivisionsMs
	default:
		return 0
	}
}

func (s *APIService) addTasks(node *calc.Node, parentArgID uint32, expressionID uint32) {
	left, right := node.Left, node.Right
	if left == nil && right == nil {
		value, _ := strconv.ParseFloat(node.Value, 64)
		taskArgs[parentArgID].Value = value
		taskArgs[parentArgID].Ready = true
		return
	}
	taskID := uuid.New().ID()
	arg1ID := uuid.New().ID()
	arg2ID := uuid.New().ID()

	taskArgs[arg1ID] = &models.Argument{ParentTaskID: taskID}
	taskArgs[arg2ID] = &models.Argument{ParentTaskID: taskID}

	s.addTasks(left, arg1ID, 0)
	s.addTasks(right, arg2ID, 0)

	task := &models.Task{
		ID:            taskID,
		ParentArgID:   parentArgID,
		Arg1:          taskArgs[arg1ID],
		Arg2:          taskArgs[arg2ID],
		ExpressionID:  expressionID,
		Operation:     node.Value,
		OperationTime: getOperationTime(s, node.Value),
	}
	if task.IsReady() {
		tasksQueue <- &models.TaskResponse{
			ID:            strconv.Itoa(int(task.ID)),
			Arg1:          task.Arg1.Value,
			Arg2:          task.Arg2.Value,
			Operation:     task.Operation,
			OperationTime: task.OperationTime,
		}
	}
	allTasks[taskID] = task
}

func (s *APIService) CreateTasks(expression string) (uint32, error) {
	expressionTree, err := calc.ToTree(expression)
	if err != nil {
		return 0, err
	}

	expressionID := uuid.New().ID()
	allExpressions[expressionID] = &models.Expression{ID: expressionID, Status: "pending"}

	s.addTasks(&expressionTree, 0, expressionID)

	return expressionID, nil
}

func (s *APIService) GetTask() *models.TaskResponse {
	select {
	case task := <-tasksQueue:
		return task
	default:
		return nil
	}
}

func (s *APIService) GetExpressionByID(expressionID uint32) *models.Expression {
	expression, exists := allExpressions[expressionID]
	if exists {
		return expression
	}
	return nil
}

func (s *APIService) GetAllExpressions() []*models.Expression {
	expressions := []*models.Expression{}
	for _, expression := range allExpressions {
		expressions = append(expressions, expression)
	}
	return expressions
}

func (s *APIService) ConfirmTask(taskID uint32, result float64) error {
	task, taskExists := allTasks[taskID]
	if taskExists {
		expression, expressionExists := allExpressions[task.ExpressionID]
		if expressionExists {
			expression.Result = result
			expression.Status = "confirmed"
		} else {
			taskArgs[task.ParentArgID].Value = result
			taskArgs[task.ParentArgID].Ready = true

			task = allTasks[taskArgs[task.ParentArgID].ParentTaskID]
			if task.IsReady() {
				tasksQueue <- &models.TaskResponse{
					ID:            strconv.Itoa(int(task.ID)),
					Arg1:          task.Arg1.Value,
					Arg2:          task.Arg2.Value,
					Operation:     task.Operation,
					OperationTime: task.OperationTime,
				}
			}
		}
		return nil
	}
	return ErrIDTaskNotExists
}

// 5 | 38
// 38 55682538
