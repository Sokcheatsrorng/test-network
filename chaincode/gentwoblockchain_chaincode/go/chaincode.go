package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing products
type SmartContract struct {
	contractapi.Contract
}

// Product represents a product in the world state
type Product struct {
	Brand  string `json:"brand"`
	Price  int    `json:"price"`
	Count  int    `json:"count"`
}

// QueryResult structure used for handling the result of queries
type QueryResult struct {
	Key    string   `json:"key"`
	Record *Product `json:"record"`
}

// InitLedger initializes the ledger with some default products
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	products := []Product{
		{Brand: "Samsung TV", Price: 250, Count: 20},
		{Brand: "Apple TV", Price: 250, Count: 30},
		{Brand: "Xiaomi Mi Phone", Price: 150, Count: 50},
		{Brand: "Toshiba Laptop", Price: 200, Count: 40},
		{Brand: "Huawei Watch", Price: 150, Count: 60},
	}

	for i, product := range products {
		productAsBytes, err := json.Marshal(product)
		if err != nil {
			return fmt.Errorf("failed to marshal product: %v", err)
		}
		err = ctx.GetStub().PutState("PRODUCT"+strconv.Itoa(i), productAsBytes)
		if err != nil {
			return fmt.Errorf("failed to put to world state: %v", err)
		}
	}

	return nil
}

// QueryAllProducts gets all products from the world state
func (s *SmartContract) QueryAllProducts(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := "PRODUCT0"
	endKey := "PRODUCT99"

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

		product := new(Product)
		err = json.Unmarshal(queryResponse.Value, product)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal product: %v", err)
		}

		queryResult := QueryResult{Key: queryResponse.Key, Record: product}
		results = append(results, queryResult)
	}

	return results, nil
}

// CreateProduct adds a new product to the world state with given details
func (s *SmartContract) CreateProduct(ctx contractapi.TransactionContextInterface, productNumber string, brand string, price int, count int) error {
	product := Product{
		Brand: brand,
		Price: price,
		Count: count,
	}

	productAsBytes, err := json.Marshal(product)
	if err != nil {
		return fmt.Errorf("failed to marshal product: %v", err)
	}

	return ctx.GetStub().PutState(productNumber, productAsBytes)
}

// ChangeProductPrice updates the price of a product in the world state
func (s *SmartContract) ChangeProductPrice(ctx contractapi.TransactionContextInterface, productNumber string, newPrice int) error {
	productAsBytes, err := ctx.GetStub().GetState(productNumber)
	if err != nil {
		return fmt.Errorf("failed to read product from world state: %v", err)
	}
	if productAsBytes == nil {
		return fmt.Errorf("product does not exist: %s", productNumber)
	}

	product := new(Product)
	err = json.Unmarshal(productAsBytes, product)
	if err != nil {
		return fmt.Errorf("failed to unmarshal product: %v", err)
	}

	product.Price = newPrice

	productAsBytesToPut, err := json.Marshal(product)
	if err != nil {
		return fmt.Errorf("failed to marshal product: %v", err)
	}

	return ctx.GetStub().PutState(productNumber, productAsBytesToPut)
}

// main function starts up the chaincode in the container during instantiation
func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error creating chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting chaincode: %s", err.Error())
	}
}
