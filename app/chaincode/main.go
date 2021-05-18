package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	//"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Record provides functions for managing an Asset
type Record struct {
	contractapi.Contract
}

// medical records should contain metadata of a patient-provider encounter (visit date/time, location, etc.),

// Asset describes basic details of what makes up a simple asset
type Asset struct {
	UUID       string `json:"uuid"`
	Timestamp  string `json:"dateTime"`
	Owner      string `json:"owner"`
	Expiration string `json:"expiration"`
}

// InitLedger adds a base set of assets to the ledger
func (s *Record) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []Asset{
		{UUID: "0a970bbe-b436-4601-87d2-becd3bf84054", Timestamp: "1511382400000", Owner: "Jane Doe", Expiration: "1655510400000"}, // add an attribute for Immunisation, representing the disesase + variant?
		{UUID: "f57d594d-3a88-4fcc-80a2-4147d23923e3", Timestamp: "1721382400000", Owner: "Ajay Singh", Expiration: "1655510400000"},
		{UUID: "6d8670a9-8550-40b9-9c57-30bb33a4d6f4", Timestamp: "1621382402100", Owner: "Zhang San", Expiration: "1655510400000"},
		{UUID: "3519bdc7-eb13-4ac0-bd63-137241af313b", Timestamp: "1621382448000", Owner: "Max Mustermann", Expiration: "1655510400000"},
		{UUID: "183ce652-0788-4ea7-896b-93d6036db83e", Timestamp: "1611382409500", Owner: "Pierre Paul", Expiration: "1655510400000"},
		{UUID: "8e0c5ba3-4dfe-41fc-b454-26a3a276ac92", Timestamp: "1511382400120", Owner: "Wang Wu", Expiration: "1655510400000"},
	}

	for _, asset := range assets {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(asset.UUID, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *Record) CreateAsset(ctx contractapi.TransactionContextInterface, uuid string, timestamp string, owner string, expiration string) error {
	exists, err := s.AssetExists(ctx, uuid)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", uuid)
	}

	asset := Asset{
		UUID:       uuid,
		Timestamp:  timestamp,
		Owner:      owner,
		Expiration: expiration,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(uuid, assetJSON)
}

// ReadAsset returns the asset stored in the world state with given uuid.
func (s *Record) ReadAsset(ctx contractapi.TransactionContextInterface, uuid string) (*Asset, error) {
	assetJSON, err := ctx.GetStub().GetState(uuid)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", uuid)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// DeleteAsset deletes an given asset from the world state.
func (s *Record) DeleteAsset(ctx contractapi.TransactionContextInterface, uuid string) error {
	exists, err := s.AssetExists(ctx, uuid)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", uuid)
	}

	return ctx.GetStub().DelState(uuid)
}

// AssetExists returns true when asset with given UUID exists in world state
func (s *Record) AssetExists(ctx contractapi.TransactionContextInterface, uuid string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(uuid)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// AssetValid returns true when an asset's timestamp is less than the expiration
func (s *Record) AssetValid(ctx contractapi.TransactionContextInterface, uuid string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(uuid)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	var asset Asset
	json.Unmarshal(assetJSON, &asset)

	var stamp, error = strconv.ParseUint(asset.Timestamp, 10, 64)
	if error != nil {
		return false, fmt.Errorf("failed to Validate asset %v", error)
	}

	var exp, anotherErr = strconv.ParseUint(asset.Expiration, 10, 64)
	if anotherErr != nil {
		return false, fmt.Errorf("failed to Validate asset %v", anotherErr)
	}

	return stamp < exp, nil
}

// TransferAsset updates the owner field of asset with given uuid in world state.
func (s *Record) TransferAsset(ctx contractapi.TransactionContextInterface, uuid string, newOwner string) error {
	asset, err := s.ReadAsset(ctx, uuid)
	if err != nil {
		return err
	}

	asset.Owner = newOwner
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(uuid, assetJSON)
}

// GetAllAssets returns all assets found in world state
func (s *Record) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

func main() {
	assetChaincode, err := contractapi.NewChaincode(&Record{})
	if err != nil {
		log.Panicf("Error creating record chaincode: %v", err)
	}

	if err := assetChaincode.Start(); err != nil {
		log.Panicf("Error starting record chaincode: %v", err)
	}
}
