package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// For handling writing and reading from the world state
type BankContract struct {
	contractapi.Contract
}

type BankAccount struct {
	AccountID string  `json:"AccountID"`
	Owner     string  `json:"Owner"`
	Balance   float64 `json:"Balance"`
}

func (ba *BankAccount) Deposit(amount float64) {
	ba.Balance += amount
}
func (ba *BankAccount) Withdraw(amount float64) {
	ba.Balance -= amount
}

// InitLedger adds a base set of assets to the Ledger.
func (bc *BankContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	accounts := []BankAccount{
		{AccountID: "123456789", Owner: "Steven", Balance: 20000},
		{AccountID: "987654321", Owner: "NSLSC", Balance: 0},
		{AccountID: "098765432", Owner: "SaskEnergy", Balance: 0},
	}

	for _, account := range accounts {
		accountJSON, err := json.Marshal(account)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(account.AccountID, accountJSON)
		if err != nil {
			return fmt.Errorf("failed to put the world state: %v", err)
		}
	}

	return nil
}

func (bc *BankContract) GetAllBankAccounts(ctx contractapi.TransactionContextInterface) ([]*BankAccount, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var bankAccounts []*BankAccount
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var bankAccount BankAccount
		err = json.Unmarshal(queryResponse.Value, &bankAccount)

		if err != nil {
			return nil, err
		}
		bankAccounts = append(bankAccounts, &bankAccount)
	}

	return bankAccounts, nil
}

func (bc *BankContract) CreateBankAccount(ctx contractapi.TransactionContextInterface, accountID string, name string) error {
	exists, err := ctx.GetStub().GetState(accountID)

	if err != nil {
		return fmt.Errorf("failed to get the world state: %v", err)
	}

	if exists != nil {
		return fmt.Errorf("cannot create the provided account number %s. Already exists", accountID)
	}

	bankAccount := BankAccount{
		AccountID: accountID,
		Owner:     name,
		Balance:   0,
	}
	bankAccountJSON, err := json.Marshal(bankAccount)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(accountID, bankAccountJSON)
}

func (bc *BankContract) ReadBankAccount(ctx contractapi.TransactionContextInterface, accountID string) (*BankAccount, error) {
	bankAccountJSON, err := ctx.GetStub().GetState(accountID)

	if err != nil {
		return nil, fmt.Errorf("failed to get the world state: %v", err)
	}
	if bankAccountJSON == nil {
		return nil, fmt.Errorf("cannot find the provided account number %s. Does not exist", accountID)
	}

	var bankAccount BankAccount
	err = json.Unmarshal(bankAccountJSON, &bankAccount)
	if err != nil {
		return nil, err
	}

	return &bankAccount, nil
}

func (bc *BankContract) TransactPayment(ctx contractapi.TransactionContextInterface, from string, to string, amount float64) (string, error) {
	fromAccountJSON, err := ctx.GetStub().GetState(from)

	toAccountJSON, err := ctx.GetStub().GetState(to)

	if err != nil {
		return "", fmt.Errorf("failed to get the world state: %v", err)
	}
	if fromAccountJSON == nil {
		return "", fmt.Errorf("cannot find the provided account number %s. Does not exist", from)
	}
	if toAccountJSON == nil {
		return "", fmt.Errorf("cannot find the provided from number %s. Does not exist", to)
	}

	var bankAccount BankAccount
	err = json.Unmarshal(fromAccountJSON, &bankAccount)
	if err != nil {
		return "", err
	}

	var toAccount BankAccount
	err = json.Unmarshal(toAccountJSON, &toAccount)
	if err != nil {
		return "", err
	}

	bankAccount.Withdraw(amount)
	toAccount.Deposit(amount)

	bankAccountJSONUpdate, err := json.Marshal(bankAccount)
	if err != nil {
		return "", err
	}

	toAccountJSONUpdate, err := json.Marshal(toAccount)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(from, bankAccountJSONUpdate)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(to, toAccountJSONUpdate)
	if err != nil {
		return "", err
	}

	// Retrieve the transaction ID
	transactionID := ctx.GetStub().GetTxID()

	return transactionID, nil
}

func (bc *BankContract) DepositBalance(ctx contractapi.TransactionContextInterface, accountID string, amount float64) (string, error) {
	bankAccountJSON, err := ctx.GetStub().GetState(accountID)

	if err != nil {
		return "", fmt.Errorf("failed to get the world state: %v", err)
	}
	if bankAccountJSON == nil {
		return "", fmt.Errorf("cannot find the provided account number %s. Already exists", accountID)
	}

	var bankAccount BankAccount
	err = json.Unmarshal(bankAccountJSON, &bankAccount)
	if err != nil {
		return "", err
	}

	bankAccount.Deposit(amount)

	bankAccountJSONUpdate, err := json.Marshal(bankAccount)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(accountID, bankAccountJSONUpdate)
	if err != nil {
		return "", err
	}

	// Retrieve the transaction ID
	transactionID := ctx.GetStub().GetTxID()

	return transactionID, nil
}

func (bc *BankContract) WithdrawBalance(ctx contractapi.TransactionContextInterface, accountID string, amount float64) (string, error) {
	bankAccountJSON, err := ctx.GetStub().GetState(accountID)

	if err != nil {
		return "", fmt.Errorf("failed to get the world state: %v", err)
	}
	if bankAccountJSON == nil {
		return "", fmt.Errorf("cannot find the provided account number %s. Already exists", accountID)
	}

	var bankAccount BankAccount
	err = json.Unmarshal(bankAccountJSON, &bankAccount)
	if err != nil {
		return "", err
	}

	bankAccount.Withdraw(amount)

	bankAccountJSONUpdate, err := json.Marshal(bankAccount)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(accountID, bankAccountJSONUpdate)
	if err != nil {
		return "", err
	}

	// Retrieve the transaction ID
	transactionID := ctx.GetStub().GetTxID()

	return transactionID, nil
}
