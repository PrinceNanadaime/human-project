package main

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"strings"
)

type HumanFactService interface {
	CreateHumanFact() (*HumanFact, error)
}

func NewHumanFactService() HumanFactService {
	return HumanFact{}
}

func (h HumanFact) CreateHumanFact() (*HumanFact, error) {
	newFact := GetNewFact()
	factToUpdate, err := json.Marshal(newFact)
	if err != nil {
		return nil, err
	}

	updateRs, err := http.Post("http://localhost:8080/calculate-length", "application/json", bytes.NewBuffer(factToUpdate))
	if err != nil {
		return nil, err
	}
	defer updateRs.Body.Close()

	updatedFact := &HumanFact{}
	if err := json.NewDecoder(updateRs.Body).Decode(updatedFact); err != nil {
		return nil, err
	}

	factToDb, err := json.Marshal(updatedFact)
	if err != nil {
		return nil, err
	}

	InsertFactToDatabase(factToDb)
	return updatedFact, err
}

func GetNewFact() HumanFact {
	facts, _ := os.ReadFile("main/HumanFacts.txt")
	randomFactNumber := rand.Intn(9)
	newFact := strings.Split(string(facts), "\n")[randomFactNumber]
	preparedFact := strings.Replace(newFact, "\r", "", -1)
	return HumanFact{
		Fact: preparedFact,
	}
}

func InsertFactToDatabase(fact []byte) {
	dbRs, err := http.Post("http://localhost:8080/db/insert", "application/json", bytes.NewBuffer(fact))
	if err != nil {
		panic(err)
	}
	defer dbRs.Body.Close()
}
