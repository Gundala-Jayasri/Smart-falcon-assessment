package main

import (
    "encoding/json"
    "fmt"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
    "strconv"
)

// SmartContract provides functions for managing assets
type SmartContract struct {
    contractapi.Contract
}

// Asset defines the structure for a financial asset
type Asset struct {
    DealerID    string `json:"dealerID"`
    MSISDN      string `json:"msisdn"`
    MPIN        string `json:"mpin"`
    Balance     int    `json:"balance"`
    Status      string `json:"status"`
    TransAmount int    `json:"transAmount"`
    TransType   string `json:"transType"`
    Remarks     string `json:"remarks"`
}

// InitLedger initializes the ledger with some basic assets (for testing purposes)
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
    assets := []Asset{
        {DealerID: "D001", MSISDN: "9876543210", MPIN: "1234", Balance: 1000, Status: "active", TransAmount: 500, TransType: "credit", Remarks: "Initial deposit"},
        {DealerID: "D002", MSISDN: "8765432109", MPIN: "5678", Balance: 2000, Status: "inactive", TransAmount: 1000, TransType: "debit", Remarks: "Withdrawal"},
    }

    for _, asset := range assets {
        assetJSON, err := json.Marshal(asset)
        if err != nil {
            return err
        }

        err = ctx.GetStub().PutState(asset.DealerID, assetJSON)
        if err != nil {
            return fmt.Errorf("failed to put asset in world state: %v", err)
        }
    }

    return nil
}

// CreateAsset adds a new asset to the world state
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, dealerID string, msisdn string, mpin string, balance int, status string, transAmount int, transType string, remarks string) error {
    // Check if asset already exists
    exists, err := s.AssetExists(ctx, dealerID)
    if err != nil {
        return err
    }
    if exists {
        return fmt.Errorf("asset %s already exists", dealerID)
    }

    // Create the asset
    asset := Asset{
        DealerID:    dealerID,
        MSISDN:      msisdn,
        MPIN:        mpin,
        Balance:     balance,
        Status:      status,
        TransAmount: transAmount,
        TransType:   transType,
        Remarks:     remarks,
    }

    assetJSON, err := json.Marshal(asset)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(dealerID, assetJSON)
}

// QueryAsset retrieves an asset by its DealerID
func (s *SmartContract) QueryAsset(ctx contractapi.TransactionContextInterface, dealerID string) (*Asset, error) {
    assetJSON, err := ctx.GetStub().GetState(dealerID)
    if err != nil {
        return nil, fmt.Errorf("failed to read from world state: %v", err)
    }
    if assetJSON == nil {
        return nil, fmt.Errorf("asset %s does not exist", dealerID)
    }

    var asset Asset
    err = json.Unmarshal(assetJSON, &asset)
    if err != nil {
        return nil, err
    }

    return &asset, nil
}

// UpdateAsset modifies an existing asset's balance or status
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, dealerID string, newBalance int, newStatus string) error {
    asset, err := s.QueryAsset(ctx, dealerID)
    if err != nil {
        return err
    }

    asset.Balance = newBalance
    asset.Status = newStatus

    assetJSON, err := json.Marshal(asset)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(dealerID, assetJSON)
}

// AssetExists checks if an asset exists in the ledger
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, dealerID string) (bool, error) {
    assetJSON, err := ctx.GetStub().GetState(dealerID)
    if err != nil {
        return false, fmt.Errorf("failed to read from world state: %v", err)
    }

    return assetJSON != nil, nil
}

// GetTransactionHistory returns the transaction history for an asset
func (s *SmartContract) GetTransactionHistory(ctx contractapi.TransactionContextInterface, dealerID string) ([]string, error) {
    resultsIterator, err := ctx.GetStub().GetHistoryForKey(dealerID)
    if err != nil {
        return nil, fmt.Errorf("failed to get history for key %v: %v", dealerID, err)
    }
    defer resultsIterator.Close()

    var history []string
    for resultsIterator.HasNext() {
        result, err := resultsIterator.Next()
        if err != nil {
            return nil, err
        }

        history = append(history, string(result.Value))
    }

    return history, nil
}
