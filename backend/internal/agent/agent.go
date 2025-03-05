package agent

import (
	"bytes"
	"calc-website/config"
	"calc-website/internal/models"
	"calc-website/pkg/calc"
	"calc-website/pkg/utils"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

func ProcessTask(orchestratorUrl string) error {
	taskUrl := orchestratorUrl + "/internal/task"
	resp, err := http.Get(taskUrl)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	defer utils.CloseResponseBody(resp.Body)
	if err != nil {
		return err
	}

	var task models.TaskResponse
	err = json.Unmarshal(body, &task)
	if err != nil {
		return err
	}
	result, err := calc.Compute(task.Arg1, task.Arg2, task.Operation)
	if err != nil {
		return err
	}
	time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)

	taskResult := models.TaskResult{
		TaskID: task.ID,
		Result: result,
	}
	taskBytes, err := json.Marshal(taskResult)
	if err != nil {
		return err
	}
	resp, err = http.Post(taskUrl, "application/json", bytes.NewBuffer(taskBytes))
	return nil
}

func StartAgents(cfg *config.Config) {
	for i := 0; i < cfg.ComputingPower; i++ {
		go func() {
			for {
				err := ProcessTask(cfg.OrchestratorUrl)
				if err != nil {
					log.Printf("error by process task: %v", err.Error())
				}
				time.Sleep(time.Millisecond * 100)
			}
		}()
	}
}
