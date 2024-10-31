package main

import (
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

func main() {
	bankContract := new(BankContract)

	cc, err := contractapi.NewChaincode(bankContract)

	if err != nil {
		panic(err.Error())
	}

	if err := cc.Start(); err != nil {
		panic(err.Error())
	}
}
