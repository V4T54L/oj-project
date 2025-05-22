package services

import (
	"algo-arena-be/internals/models"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type executionService struct{}

func NewExecutionService() ExecutionService {
	return &executionService{}
}

func (e *executionService) Execute(
	payload models.ExecutionPayload,
) (models.ExecutionResponse, error) {

	empty := make(models.ExecutionResponse, 0)
	json_data, err := json.Marshal(payload)
	if err != nil {
		log.Println("error marshaling payload in execute: ", err)
		return empty, err
	}

	resp, err := http.Post("http://localhost:8081/run", "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		log.Println("error sending request in execute: ", err)
		return empty, err
	}
	defer resp.Body.Close()

	log.Println("\n\nResponse:\n",resp.Body)

	var res models.ExecutionResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		log.Println("error unmarshalling response in execute: ", err)
		return empty, err
	}

	return res, nil
}
