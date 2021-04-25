package main

import (
  "encoding/json"
  "fmt"
  "log"
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
      ID             string `json:"ID"`
      Timestamp          string `json:"dateTime"`
      Owner          string `json:"owner"`
      Expiration     uint64    `json:"expiration"`
    }

// InitLedger adds a base set of assets to the ledger
   func (s *Record) InitLedger(ctx contractapi.TransactionContextInterface) error {
    assets := []Asset{
      {ID: "asset1", Timestamp: "2012-10-31 15:50:13.793654 +0000 UTC", Owner: "Tomoko", Expiration: 104467440737095516},
      {ID: "asset2", Timestamp: "2011-10-21 12:07:13.793654 +0000 UTC", Owner: "Brad", Expiration: 438579485745748},
      {ID: "asset3", Timestamp: "2009-01-11 13:12:13.263674 +0000 UTC", Owner: "Jin Soo", Expiration: 89273489237483},
      {ID: "asset4", Timestamp: "2020-03-13 11:03:13.795664 +0000 UTC", Owner: "Max", Expiration: 239048203894},
      {ID: "asset5", Timestamp: "2019-11-05 14:12:11.798454 +0000 UTC", Owner: "Adriana", Expiration: 9023423948333},
      {ID: "asset6", Timestamp: "2013-05-21 10:01:13.274683 +0000 UTC", Owner: "Michel", Expiration: 290378489237},
    }

    for _, asset := range assets {
      assetJSON, err := json.Marshal(asset)
      if err != nil {
        return err
      }

      err = ctx.GetStub().PutState(asset.ID, assetJSON)
      if err != nil {
        return fmt.Errorf("failed to put to world state. %v", err)
      }
    }

    return nil
  }

// CreateAsset issues a new asset to the world state with given details.
   func (s *Record) CreateAsset(ctx contractapi.TransactionContextInterface, id string, timestamp string, owner string, expiration uint64) error {
    exists, err := s.AssetExists(ctx, id)
    if err != nil {
      return err
    }
    if exists {
      return fmt.Errorf("the asset %s already exists", id)
    }

    asset := Asset{
      ID: id,
      Timestamp:  timestamp,
      Owner:  owner,
      Expiration: expiration,
    }
    assetJSON, err := json.Marshal(asset)
    if err != nil {
      return err
    }

    return ctx.GetStub().PutState(id, assetJSON)
  }

// ReadAsset returns the asset stored in the world state with given id.
   func (s *Record) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
    assetJSON, err := ctx.GetStub().GetState(id)
    if err != nil {
      return nil, fmt.Errorf("failed to read from world state: %v", err)
    }
    if assetJSON == nil {
      return nil, fmt.Errorf("the asset %s does not exist", id)
    }

    var asset Asset
    err = json.Unmarshal(assetJSON, &asset)
    if err != nil {
      return nil, err
    }

    return &asset, nil
  }

// UpdateAsset updates an existing asset in the world state with provided parameters.
   func (s *Record) UpdateAsset(ctx contractapi.TransactionContextInterface, id string, timestamp string, owner string, expiration uint64) error {
    exists, err := s.AssetExists(ctx, id)
    if err != nil {
      return err
    }
    if !exists {
      return fmt.Errorf("the asset %s does not exist", id)
    }

    // overwriting original asset with new asset
    asset := Asset{
      ID: id,
      Timestamp:  timestamp,
      Owner:  owner,
      Expiration: expiration,
    }
    assetJSON, err := json.Marshal(asset)
    if err != nil {
      return err
    }

    return ctx.GetStub().PutState(id, assetJSON)
  }

  // DeleteAsset deletes an given asset from the world state.
  func (s *Record) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
    exists, err := s.AssetExists(ctx, id)
    if err != nil {
      return err
    }
    if !exists {
      return fmt.Errorf("the asset %s does not exist", id)
    }

    return ctx.GetStub().DelState(id)
  }

// AssetExists returns true when asset with given ID exists in world state
   func (s *Record) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
    assetJSON, err := ctx.GetStub().GetState(id)
    if err != nil {
      return false, fmt.Errorf("failed to read from world state: %v", err)
    }

    return assetJSON != nil, nil
  }

// TransferAsset updates the owner field of asset with given id in world state.
   func (s *Record) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string) error {
    asset, err := s.ReadAsset(ctx, id)
    if err != nil {
      return err
    }

    asset.Owner = newOwner
    assetJSON, err := json.Marshal(asset)
    if err != nil {
      return err
    }

    return ctx.GetStub().PutState(id, assetJSON)
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