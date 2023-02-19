/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a car
type SmartContract struct {
	contractapi.Contract
}


type Owner struct {
	Id      int     `json:"id"`
	Name    string  `json:"name"`
	Surname string  `json:"surname"`
	Email   string  `json:"email"`
	Money   float64 `json:"money"`
}

type Malfunction struct {
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

// Car describes basic details of what makes up a car
type Car struct {
	Id	   int     `json:"id"`
	Make   string  `json:"make"`
	Model  string  `json:"model"`
	Colour string  `json:"colour"`
	Year   string  `json:"year"`
	Owner  string  `json:"owner"`
	Price  float64 `json:"price"`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Car
}

// InitLedger adds a base set of cars to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	owners := []Owner{
		{Id: 1, Name: "Marko", Surname: "NIkolic", Email: "markonikolic@gmail.com", Money: 8000},
		{Id: 2, Name: "NIkola", Surname: "Kajtes", Email: "nikolakajtes@gmail.com", Money: 11000},
	}

	cars := []Car{
		{Id: 1, Make: "Toyota", Model: "Prius", Color: "blue",Year: "2020", Owner: "1",
			Malfunctions: []Malfunction{
				{Description: "Otpo tocak", Price: 250},
			},
			Price: 5000,
		},
		{Id: 2, Make: "Ford", Model: "Mustang", Color: "blue",Year: "2010", Owner: "1", Malfunctions: []Malfunction{}, Price: 3000},
		{Id: 3, Make: "Hyundai", Model: "Tucson", Color: "green",Year: "2021", Owner: "2",
			Malfunctions: []Malfunction{
				{Description: "Otisla lamela", Price: 1000},
			},
			Price: 2500,
		},
		{Id: 4, Make: "Volkswagen", Model: "Passat", Color: "blue",Year: "2020", Owner: "2", Malfunctions: []Malfunction{}, Price: 7000},
		{Id: 5, Make: "Tesla", Model: "S", Color: "blue",Year: "2023", Owner: "1", Malfunctions: []Malfunction{}, Price: 20000},
	}

	for i, car := range cars {
		carAsBytes, _ := json.Marshal(car)
		err := ctx.GetStub().PutState("CAR"+strconv.Itoa(i), carAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// CreateCar adds a new car to the world state with given details
func (s *SmartContract) CreateCar(ctx contractapi.TransactionContextInterface, carNumber string, make string, model string, colour string, owner string) error {
	car := Car{
		Make:   make,
		Model:  model,
		Colour: colour,
		Owner:  owner,
	}

	carAsBytes, _ := json.Marshal(car)

	return ctx.GetStub().PutState(carNumber, carAsBytes)
}

// QueryCar returns the car stored in the world state with given id
func (s *SmartContract) QueryCar(ctx contractapi.TransactionContextInterface, carNumber string) (*Car, error) {
	carAsBytes, err := ctx.GetStub().GetState(carNumber)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if carAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", carNumber)
	}

	car := new(Car)
	_ = json.Unmarshal(carAsBytes, car)

	return car, nil
}

// QueryAllCars returns all cars found in world state
func (s *SmartContract) QueryAllCars(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		car := new(Car)
		_ = json.Unmarshal(queryResponse.Value, car)

		queryResult := QueryResult{Key: queryResponse.Key, Record: car}
		results = append(results, queryResult)
	}

	return results, nil
}

// ChangeCarOwner updates the owner field of car with given id in world state
func (s *SmartContract) ChangeCarOwner(ctx contractapi.TransactionContextInterface, carNumber string, newOwner string) error {
	car, err := s.QueryCar(ctx, carNumber)

	if err != nil {
		return err
	}

	car.Owner = newOwner

	carAsBytes, _ := json.Marshal(car)

	return ctx.GetStub().PutState(carNumber, carAsBytes)
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	}
}
