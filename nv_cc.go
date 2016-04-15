/*
Copyright 2016 IBM

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Licensed Materials - Property of IBM
© Copyright IBM Corp. 2016
*/
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/openblockchain/obc-peer/openchain/chaincode/shim"
)

const   BANKA = "BANKA"
const   BANKB = "BANKB"
const   BANKC = "BANKC"
const 	AUDITOR = "AUDITOR"

const AUDUSD = 0.74
const USDAUD = 1.34
const EURUSD = 1.10
const USDEUR = 0.90
const AUDEUR = 0.67
const EURAUD = 1.48
const TESTCONV= 1.13

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Account struct {
	Holder    	string  `json:"holder"`
	Currency  	string  `json:"currency"`
	CashBalance float64 `json:"cashBalance"`
}

type FinancialInst struct {
	Owner     	string  `json:"owner"`
	Accounts []Account `json:"accounts"`
}

type Transaction struct {
	RefNumber   string   `json:"refNumber"`
	Date 		string   `json:"date"`
	Description string   `json:"description"`
	Type 		string   `json:"type"`
	Amount    	float64  `json:"amount"`
	To			string   `json:"to"`
	From		string   `json:"from"`
	Contract	string   `json:"contract"`
	StatusCode	int 	 `json:"statusCode"`
	StatusMsg	string   `json:"statusMsg"`
}


type User struct {
	Name   		string   `json:"Name"`
	Balance 	float64  `json:"Balance"`
	Status      string 	 `json:"status"`
	Expiration  string   `json:"expiration"`
	Join		string   `json:"join"`
	Modified	string   `json:"from"`
}

type AllTransactions struct{
	Transactions []Transaction `json:"transactions"`
}

type AllUsers struct{
	User []User `json:"users"`
}

type NVAccounts struct {
	User 		[]User `json:"user"`
	Vostro 		[]FinancialInst `json:"vostro"`
}

// ============================================================================================================================
// Init 
// ============================================================================================================================
func (t *SimpleChaincode) init(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	var err error

	// Natalie
	var natalie User
	natalie.Name = "Natalie"
	natalie.Balance = 1000
	natalie.Status  = "Platinum"
	natalie.Expiration = "June"
	natalie.Join  = "January"
	natalie.Modified = "Today"
	
	jsonAsBytes, _ := json.Marshal(natalie)
	err = stub.PutState(natalie.Name, jsonAsBytes)								
	if err != nil {
		fmt.Println("Error Creating initial user account")
		return nil, err
	}
	
	//BANK A
	var fid FinancialInst
	fid.Owner = BANKA
	
	var actAB Account
	actAB.Holder = BANKB
	actAB.Currency = "USD"
	actAB.CashBalance = 250000.00
	fid.Accounts = append(fid.Accounts, actAB)
	var actAC Account
	actAC.Holder = BANKC
	actAC.Currency = "USD"
	actAC.CashBalance = 360000.00
	fid.Accounts = append(fid.Accounts, actAC)

	jsonAsBytes, _ = json.Marshal(fid)
	err = stub.PutState("BANKA", jsonAsBytes)								
	if err != nil {
		fmt.Println("Error creating account "+BANKA)
		return nil, err
	}

	// BANK B
	var fid2 FinancialInst
	fid2.Owner = BANKB

	var actBA Account
	actBA.Holder = BANKA
	actBA.Currency = "AUD"
	actBA.CashBalance = actAB.CashBalance * USDAUD
	fid2.Accounts = append(fid2.Accounts, actBA)
	var actBC Account
	actBC.Holder = BANKC
	actBC.Currency = "AUD"
	actBC.CashBalance = 120000.00
	fid2.Accounts = append(fid2.Accounts, actBC)

	jsonAsBytes, _ = json.Marshal(fid2)
	err = stub.PutState("BANKB", jsonAsBytes)								
	if err != nil {
		fmt.Println("Error creating account "+BANKB)
		return nil, err
	}

	// BANK C
	var fid3 FinancialInst
	fid3.Owner = BANKC

	var actCA Account
	actCA.Holder = BANKA
	actCA.Currency = "EUR"
	actCA.CashBalance = actAC.CashBalance * USDEUR
	fid3.Accounts = append(fid3.Accounts, actCA)
	var actCB Account
	actCB.Holder = BANKB
	actCB.Currency = "EUR"
	actCB.CashBalance = actBC.CashBalance * AUDEUR
	fid3.Accounts = append(fid3.Accounts, actCB)

	jsonAsBytes, _ = json.Marshal(fid3)
	err = stub.PutState("BANKC", jsonAsBytes)								
	if err != nil {
		fmt.Println("Error creating account "+BANKC)
		return nil, err
	}
	
	var transactions AllTransactions
	jsonAsBytes, _ = json.Marshal(transactions)
	err = stub.PutState("allTx", jsonAsBytes)
	if err != nil {
		return nil, err
	}
	
	return nil, nil
}



// ============================================================================================================================
// Run - Our entry point
// ============================================================================================================================
func (t *SimpleChaincode) Run(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("run is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state
		return t.init(stub, args)
	} else if function == "submitTx" {											//create a transaction
		return t.submitTx(stub, args)
	} 
	fmt.Println("run did not find func: " + function)						//error

	return nil, errors.New("Received unknown function invocation")
}

// ============================================================================================================================
// Query - read a variable from chaincode state - (aka read)
// ============================================================================================================================
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	if len(args) != 2 { return nil, errors.New("Incorrect number of arguments passed") }

	if args[0] != "getFIDetails" && args[0] != "getTxs" && args[0] != "getNVAccounts"{
		return nil, errors.New("Invalid query function name.")
	}

	if args[0] == "getFIDetails" { return t.getFinInstDetails(stub, args[1]) }
	if args[0] == "getNVAccounts" { return t.getNVAccounts(stub, args[1]) }
	if args[0] == "getTxs" { return t.getTxs(stub, args[1]) }

	return nil, nil										
}


// ============================================================================================================================
// Get Financial Institution Details
// ============================================================================================================================
func (t *SimpleChaincode) getFinInstDetails(stub *shim.ChaincodeStub, finInst string)([]byte, error){
	
	fmt.Println("Start find getFinInstDetails")
	fmt.Println("Looking for " + finInst);

	//get the finInst index
	fdAsBytes, err := stub.GetState(finInst)
	if err != nil {
		return nil, errors.New("Failed to get Financial Institution")
	}

	return fdAsBytes, nil
	
}

// ============================================================================================================================
// Get Nostro/Vostro accounts for a specific Financial Institution
// ============================================================================================================================
func (t *SimpleChaincode) getNVAccounts(stub *shim.ChaincodeStub, finInst string)([]byte, error){
	
	fmt.Println("Start find getNVAccounts")
	fmt.Println("Looking for " + finInst);

	//get the User index
	fdAsBytes, err := stub.GetState("Natalie")
	if err != nil {
		return nil, errors.New("Failed to get Financial Institution")
	}

	var fd User
	json.Unmarshal(fdAsBytes, &fd)

	var res NVAccounts
	res.User = append(res.User, fd)
	
	
	//get the finInst index
	fdAsBytes, err = stub.GetState("BANKA")
	if err != nil {
		return nil, errors.New("Failed to get Financial Institution")
	}

	var fin FinancialInst
	json.Unmarshal(fdAsBytes, &fin)

	res.Vostro = append(res.Vostro, fin)

	resAsBytes, _ := json.Marshal(res)

	return resAsBytes, nil
	
}

// ============================================================================================================================
// Get Transactions for a specific Financial Institution (Inbound and Outbound)
// ============================================================================================================================
func (t *SimpleChaincode) getTxs(stub *shim.ChaincodeStub, finInst string)([]byte, error){
	
	var res AllTransactions

	fmt.Println("Start find getTransactions")
	fmt.Println("Looking for " + finInst);

	//get the AllTransactions index
	allTxAsBytes, err := stub.GetState("allTx")
	if err != nil {
		return nil, errors.New("Failed to get all Transactions")
	}

	var txs AllTransactions
	json.Unmarshal(allTxAsBytes, &txs)

	for i := range txs.Transactions{

		if txs.Transactions[i].From == finInst{
			res.Transactions = append(res.Transactions, txs.Transactions[i])
		}

		if txs.Transactions[i].To == finInst{
			res.Transactions = append(res.Transactions, txs.Transactions[i])
		}

		if(finInst == AUDITOR) {
			res.Transactions = append(res.Transactions, txs.Transactions[i])
		}
	}

	resAsBytes, _ := json.Marshal(res)

	return resAsBytes, nil
	
}



// ============================================================================================================================
// Submit Transaction
	// RefNumber   string   `json:"refNumber"`
	// OpCode 		string   `json:"opCode"`
	// VDate 		string   `json:"vDate"`
	// Currency  	string   `json:"currency"`
	// Amount    	float64  `json:"amount"`
	// From		string   `json:"From"`
	// To	string   `json:"To"`
	// OrdCust		string   `json:"ordcust"`
	// BenefCust	string   `json:"benefcust"`
	// DetCharges  string   `json:"detcharges"`
// ============================================================================================================================
func (t *SimpleChaincode) submitTx(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	var err error
	fmt.Println("Running submitTx")

	if len(args) != 10 {
		fmt.Println("Incorrect number of arguments. Expecting 10 - MT103 format")
		return nil, errors.New("Incorrect number of arguments. Expecting 10 - MT103 format")
	}
															
	var tx Transaction
	tx.RefNumber 	= args[0]
	tx.Date 		= args[1]
	tx.Description 	= args[2]
	tx.Type 	    = args[3]
	tx.To 			= args[5]
	tx.From 		= args[6]
	tx.Contract 	= args[7]
	tx.StatusCode 	= 1
	tx.StatusMsg 	= "Transaction Completed"

	amountValue, err := strconv.ParseFloat(args[4], 64)
	if err != nil {
		tx.StatusCode = 0
		tx.StatusMsg = "Invalid Amount"
	}else{
		tx.Amount = amountValue
	}

	
	
	// Get Sender account from BC
	rfidBytes, err := stub.GetState(tx.From)
	if err != nil {
		return nil, errors.New("SubmitTx Failed to get Financial Institution")
	}
	var sender FinancialInst
	fmt.Println("SubmitTx Unmarshalling Financial Institution");
	err = json.Unmarshal(rfidBytes, &sender)
	
	
	// Get Receiver account from BC
	rfidBytes, err = stub.GetState(tx.To)
	if err != nil {
		return nil, errors.New("SubmitTx Failed to get User from BC")
	}
	var receiver User
	fmt.Println("SubmitTx Unmarshalling User Struct");
	err = json.Unmarshal(rfidBytes, &receiver)
	
	
	
	receiver.Balance = receiver.Balance  + tx.Amount
	sender.Accounts[0].CashBalance   = sender.Accounts[0].CashBalance  - tx.Amount
	
	
	
	//Commit Sender to ledger
	fmt.Println("SubmitTx Commit Updated Sender To Ledger");
	txsAsBytes, _ := json.Marshal(sender)
	err = stub.PutState(tx.From, txsAsBytes)	
	if err != nil {
		return nil, err
	}
	
	//Commit Receiver to ledger
	fmt.Println("SubmitTx Commit Updated Sender To Ledger");
	txsAsBytes, _ = json.Marshal(receiver)
	err = stub.PutState(tx.To, txsAsBytes)	
	if err != nil {
		return nil, err
	}
	

	
	//get the AllTransactions index
	allTxAsBytes, err := stub.GetState("allTx")
	if err != nil {
		return nil, errors.New("SubmitTx Failed to get all Transactions")
	}
	
	//Commit transaction to ledger
	fmt.Println("SubmitTx Commit Transaction To Ledger");
	var txs AllTransactions
	json.Unmarshal(allTxAsBytes, &txs)
	txs.Transactions = append(txs.Transactions, tx)
	txsAsBytes, _ = json.Marshal(txs)
	err = stub.PutState("allTx", txsAsBytes)	
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (t *SimpleChaincode) creditVostroAccount(stub *shim.ChaincodeStub, sender string, receiver string, amount float64) ([]byte, error) {

	senderBytes, err := stub.GetState(sender)
	if err != nil {
		return nil, errors.New("Failed to get Financial Institution")
	}
	var sfid FinancialInst
	fmt.Println("CreditVostroAccount Unmarshalling Financial Institution");
	err = json.Unmarshal(senderBytes, &sfid)
	if err != nil {
		return nil, err
	}

	for i := range sfid.Accounts{
		if sfid.Accounts[i].Holder == receiver{
			sfid.Accounts[i].CashBalance = sfid.Accounts[i].CashBalance + amount
		}
	}

	sfidAsBytes, _ := json.Marshal(sfid)
	err = stub.PutState(sender, sfidAsBytes)	
	if err != nil {
		return nil, err
	}

	return nil, nil

}

func (t *SimpleChaincode) debitNostroAccount(stub *shim.ChaincodeStub, sender string, receiver string, amount float64) ([]byte, error) {

	receiverBytes, err := stub.GetState(receiver)
	if err != nil {
		return nil, errors.New("Failed to get Financial Institution")
	}
	var rfid FinancialInst
	fmt.Println("DebitNostroAccount Unmarshalling Financial Institution");
	err = json.Unmarshal(receiverBytes, &rfid)
	if err != nil {
		return nil, err
	}

	for i := range rfid.Accounts{
		if rfid.Accounts[i].Holder == sender{
			rfid.Accounts[i].CashBalance = rfid.Accounts[i].CashBalance - amount
		}
	}

	rfidAsBytes, _ := json.Marshal(rfid)
	err = stub.PutState(receiver, rfidAsBytes)	
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

func getFXRate(curS string, curR string) (float64, error){
	if(curS == "USD" && curR == "AUD"){ return USDAUD,nil }
	if(curS == "USD" && curR == "EUR"){ return USDEUR,nil }
	if(curS == "EUR" && curR == "AUD"){ return EURAUD,nil }
	if(curS == "EUR" && curR == "USD"){ return EURUSD,nil }
	if(curS == "AUD" && curR == "EUR"){ return AUDEUR,nil }
	if(curS == "AUD" && curR == "USD"){ return AUDUSD,nil }
	return 0.0, errors.New("Not matching Currency")
}

func FloatToString(input_num float64) string {
    // to convert a float number to a string
    return strconv.FormatFloat(input_num, 'f', 4, 64)
}


