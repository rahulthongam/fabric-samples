package chaincode

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing bank accounts
type SmartContract struct {
	contractapi.Contract
}

// Account describes the basic details of a bank account
type Account struct {
	ID        string  `json:"ID"`
	Owner     string  `json:"Owner"`
	Balance   float64 `json:"Balance"`
}

// InitLedger adds a base set of accounts to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	accounts := []Account{
		{ID: "account1", Owner: "Tomoko", Balance: 1000.0},
		{ID: "account2", Owner: "Brad", Balance: 2000.0},
		{ID: "account3", Owner: "Jin Soo", Balance: 3000.0},
		{ID: "account4", Owner: "Max", Balance: 4000.0},
		{ID: "account5", Owner: "Adriana", Balance: 5000.0},
		{ID: "account6", Owner: "Michel", Balance: 6000.0},
	}

	for _, account := range accounts {
		accountJSON, err := json.Marshal(account)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(account.ID, accountJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state: %v", err)
		}
	}

	return nil
}

// CreateAccount creates a new bank account with the given details.
func (s *SmartContract) CreateAccount(ctx contractapi.TransactionContextInterface, id string, owner string, balance float64) error {
	exists, err := s.AccountExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the account %s already exists", id)
	}

	account := Account{
		ID:      id,
		Owner:   owner,
		Balance: balance,
	}
	accountJSON, err := json.Marshal(account)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, accountJSON)
}

// ReadAccount returns the account stored in the world state with the given id.
func (s *SmartContract) ReadAccount(ctx contractapi.TransactionContextInterface, id string) (*Account, error) {
	accountJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if accountJSON == nil {
		return nil, fmt.Errorf("the account %s does not exist", id)
	}

	var account Account
	err = json.Unmarshal(accountJSON, &account)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

// UpdateAccount updates an existing account in the world state with the provided parameters.
func (s *SmartContract) UpdateAccount(ctx contractapi.TransactionContextInterface, id string, owner string, balance float64) error {
	exists, err := s.AccountExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the account %s does not exist", id)
	}

	account := Account{
		ID:      id,
		Owner:   owner,
		Balance: balance,
	}
	accountJSON, err := json.Marshal(account)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, accountJSON)
}

// DeleteAccount deletes the given account from the world state.
func (s *SmartContract) DeleteAccount(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.AccountExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the account %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// AccountExists returns true when an account with the given ID exists in the world state.
func (s *SmartContract) AccountExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	accountJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return accountJSON != nil, nil
}

// GetAllAccounts returns all accounts found in the world state.
func (s *SmartContract) GetAllAccounts(ctx contractapi.TransactionContextInterface) ([]*Account, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var accounts []*Account
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var account Account
		err = json.Unmarshal(queryResponse.Value, &account)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, &account)
	}

	return accounts, nil
}

// TransferFunds transfers funds from one account to another.
func (s *SmartContract) TransferFunds(ctx contractapi.TransactionContextInterface, fromID string, toID string, amount float64) error {
	fromAccount, err := s.ReadAccount(ctx, fromID)
	if err != nil {
		return err
	}
	toAccount, err := s.ReadAccount(ctx, toID)
	if err != nil {
		return err
	}

	if fromAccount.Balance < amount {
		return fmt.Errorf("insufficient funds in the account %s", fromID)
	}

	fromAccount.Balance -= amount
	toAccount.Balance += amount

	fromAccountJSON, err := json.Marshal(fromAccount)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(fromID, fromAccountJSON)
	if err != nil {
		return err
	}

	toAccountJSON, err := json.Marshal(toAccount)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(toID, toAccountJSON)
	if err != nil {
		return err
	}

	return nil
}
