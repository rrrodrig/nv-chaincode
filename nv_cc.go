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
	"time"
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



// New structs for API demo
type Transaction struct {
	RefNumber   string   `json:"RefNumber"`
	Date 		string   `json:"Date"`
	Description string   `json:"description"`
	Type 		string   `json:"Type"`
	Amount    	float64  `json:"Amount"`
	To			string   `json:"ToUserid"`
	From		string   `json:"FromUserid"`
	ToName	    string   `json:"ToName"`
	FromName	string   `json:"FromName"`
	Contract	string   `json:"Contract"`
	StatusCode	int 	 `json:"StatusCode"`
	StatusMsg	string   `json:"StatusMsg"`
}


type User struct {
	UserId		string   `json:"UserId"`
	Name   		string   `json:"Name"`
	Balance 	float64  `json:"Balance"`
	Status      string 	 `json:"Status"`
	Expiration  string   `json:"ExpirationDate"`
	Join		string   `json:"JoinDate"`
	Modified	string   `json:"LastModifiedDate"`
}


// Old structs from NV demo - to be deleted 
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


type Account struct {
	Holder    	string  `json:"holder"`
	Currency  	string  `json:"currency"`
	CashBalance float64 `json:"cashBalance"`
}

type FinancialInst struct {
	Owner     	string  `json:"owner"`
	Accounts []Account `json:"accounts"`
}

// ============================================================================================================================
// Init - initiate data structures and blockchain
// ============================================================================================================================
func (t *SimpleChaincode) init(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	var err error

		
	// Create the 'Bank' user and add it to the blockchain
	var bank User
	bank.UserId = "1";
	bank.Name = "Open Financial Network"
	bank.Balance = 1000000
	bank.Status  = "Originator"
	bank.Expiration = "2099-12-31"
	bank.Join  = "2015-01-01"
	bank.Modified = "2016-05-06"
	
	jsonAsBytes, _ := json.Marshal(bank)
	err = stub.PutState(bank.UserId, jsonAsBytes)								
	if err != nil {
		fmt.Println("Error Creating Bank user account")
		return nil, err
	}
	
	
    // Create the 'Travel Agency' user and add it to the blockchain
	var travel User
	travel.UserId = "2";
	travel.Name = "Open Travel Network"
	travel.Balance = 500000
	travel.Status  = "Member"
	travel.Expiration = "2099-12-31"
	travel.Join  = "2015-01-01"
	travel.Modified = "2016-05-06"
	
	jsonAsBytes, _ = json.Marshal(travel)
	err = stub.PutState(travel.UserId, jsonAsBytes)								
	if err != nil {
		fmt.Println("Error Creating Travel user account")
		return nil, err
	}
	
	
	// Create the 'Natalie' user and add her to the blockchain
	var natalie User
	natalie.UserId = "3";
	natalie.Name = "Natalie"
	natalie.Balance = 1000
	natalie.Status  = "Platinum"
	natalie.Expiration = "2017-06-01"
	natalie.Join  = "2015-05-31"
	natalie.Modified = "2016-05-06"
	
	jsonAsBytes, _ = json.Marshal(natalie)
	err = stub.PutState(natalie.UserId, jsonAsBytes)								
	if err != nil {
		fmt.Println("Error Creating Natalie user account")
		return nil, err
	}
	
	
	// Create the 'Anthony' user and add him to the blockchain
	var anthony User
	anthony.UserId = "4";
	anthony.Name = "Anthony"
	anthony.Balance = 500
	anthony.Status  = "Silver"
	anthony.Expiration = "2017-03-15"
	anthony.Join  = "2015-08-15"
	anthony.Modified = "2016-04-17"
	
	jsonAsBytes, _ = json.Marshal(anthony)
	err = stub.PutState(anthony.UserId, jsonAsBytes)								
	if err != nil {
		fmt.Println("Error Creating Anthony user account")
		return nil, err
	}
	
	
	var transactions AllTransactions
	jsonAsBytes, _ = json.Marshal(transactions)
	err = stub.PutState("allTx", jsonAsBytes)
	if err != nil {
		return nil, err
	}
	
	// Create current reference number if necessary
	var refNumber int
	refNumberBytes, numErr := stub.GetState("refNumber")
	if numErr != nil {
	
		refNumber = 1
		jsonAsBytes, _ = json.Marshal(refNumber)
		err = stub.PutState("refNumber", jsonAsBytes)								
		if err != nil {
			fmt.Println("Error Creating reference number")
			return nil, err
		}
	} else {
		err = json.Unmarshal(refNumberBytes, &refNumber)
	}
	



	
	
	//BANK A
	var fid FinancialInst
	fid.Owner = BANKA
	
	var actAB Account
	actAB.Holder = BANKB
	actAB.Currency = "USD"
	actAB.CashBalance = 250000
	fid.Accounts = append(fid.Accounts, actAB)
	var actAC Account
	actAC.Holder = BANKC
	actAC.Currency = "USD"
	actAC.CashBalance = 300000
	fid.Accounts = append(fid.Accounts, actAC)

	jsonAsBytes, _ = json.Marshal(fid)
	err = stub.PutState("BANKA", jsonAsBytes)								
	if err != nil {
		fmt.Println("Error creating account "+BANKA)
		return nil, err
	}

	
	return nil, nil
}



// ============================================================================================================================
// Run - Our entry point
// Function names called from Node JS must be added here in order to be callable
// ============================================================================================================================
func (t *SimpleChaincode) Run(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("run is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state
		return t.init(stub, args)
	} else if function == "submitTx" {											//create a transaction
		return t.submitTx(stub, args) 
	} else if function == "updateUserAccount" {											//create a transaction
		return t.updateUserAccount(stub, args) 
	} else if function == "transferPoints" {											//create a transaction
		return t.transferPoints(stub, args)
	} 
	
	
	fmt.Println("run did not find func: " + function)						//error

	return nil, errors.New("Received unknown function invocation")
}

// ============================================================================================================================
// Query - read a variable from chaincode state - (aka read)
// ============================================================================================================================
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	if len(args) != 2 { return nil, errors.New("Incorrect number of arguments passed") }

	if args[0] != "getFIDetails" && args[0] != "getTxs" && args[0] != "getNVAccounts"&& args[0] != "getUserAccount"{
		return nil, errors.New("Invalid query function name.")
	}

	if args[0] == "getFIDetails" { return t.getFinInstDetails(stub, args[1]) }
	if args[0] == "getNVAccounts" { return t.getNVAccounts(stub, args[1]) }
	if args[0] == "getTxs" { return t.getTxs(stub, args[1]) }
	if args[0] == "getUserAccount" { return t.getUserAccount(stub, args[1]) }

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
	fdAsBytes, err := stub.GetState("1")
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
func (t *SimpleChaincode) getUserAccount(stub *shim.ChaincodeStub, userId string)([]byte, error){
	
	fmt.Println("Start getUserAccount")
	fmt.Println("Looking for user with ID " + userId);

	//get the User index
	fdAsBytes, err := stub.GetState(userId)
	if err != nil {
		return nil, errors.New("Failed to get user account from blockchain")
	}

	return fdAsBytes, nil
	
}

// ============================================================================================================================
// Get Transactions for a specific Financial Institution (Inbound and Outbound)
// ============================================================================================================================
func (t *SimpleChaincode) getTxs(stub *shim.ChaincodeStub, userId string)([]byte, error){
	
	var res AllTransactions

	fmt.Println("Start find getTransactions")
	fmt.Println("Looking for " + userId);

	//get the AllTransactions index
	allTxAsBytes, err := stub.GetState("allTx")
	if err != nil {
		return nil, errors.New("Failed to get all Transactions")
	}

	var txs AllTransactions
	json.Unmarshal(allTxAsBytes, &txs)

	for i := range txs.Transactions{

		if txs.Transactions[i].From == userId{
			res.Transactions = append(res.Transactions, txs.Transactions[i])
		}

		if txs.Transactions[i].To == userId{
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

	fmt.Println("Running submitTx")
	
	
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
	
	
	//***************************************************************
	// Get Receiver account from BC
	rfidBytes, err := stub.GetState(tx.To)
	if err != nil {
		return nil, errors.New("SubmitTx Failed to get User from BC")
	}
	var receiver User
	fmt.Println("SubmitTx Unmarshalling User Struct");
	err = json.Unmarshal(rfidBytes, &receiver)
	receiver.Balance = receiver.Balance  + tx.Amount
	
	
	//Commit Receiver to ledger
	fmt.Println("SubmitTx Commit Updated Sender To Ledger");
	txsAsBytes, _ := json.Marshal(receiver)
	err = stub.PutState(tx.To, txsAsBytes)	
	if err != nil {
		return nil, err
	}
	
	// Get Sender account from BC
	rfidBytes, err = stub.GetState(tx.From)
	if err != nil {
		return nil, errors.New("SubmitTx Failed to get Financial Institution")
	}
	var sender FinancialInst
	fmt.Println("SubmitTx Unmarshalling Financial Institution");
	err = json.Unmarshal(rfidBytes, &sender)
	sender.Accounts[0].CashBalance   = sender.Accounts[0].CashBalance  - tx.Amount
	
	//Commit Sender to ledger
	fmt.Println("SubmitTx Commit Updated Sender To Ledger");
	txsAsBytes, _ = json.Marshal(sender)
	err = stub.PutState(tx.From, txsAsBytes)	
	if err != nil {
		return nil, err
	}
	
	
	return nil, nil
	//***********************************************************************
}


// ============================================================================================================================
func (t *SimpleChaincode) transferPoints(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	fmt.Println("Running transferPoints")
	currentDateStr := time.Now().Format(time.RFC850)

	
	var tx Transaction
	//tx.RefNumber 	= "1000"
	tx.Date 		=  currentDateStr
	tx.Description 	= "PointsTransfer"
	tx.Type 	    = "Valid"
	tx.To 			= args[0]
	tx.From 		= args[1]
	tx.Contract 	= "Standard"
	tx.StatusCode 	= 1
	tx.StatusMsg 	= "Transaction Completed"
	
	
	amountValue, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		tx.StatusCode = 0
		tx.StatusMsg = "Invalid Amount"
	}else{
		tx.Amount = amountValue
	}
	
	
	// Get the current reference number and update it
	var refNumber int
	refNumberBytes, numErr := stub.GetState("refNumber")
	if numErr != nil {
		fmt.Println("Error Getting  ref number for transferring points")
		return nil, err
	}
	
	json.Unmarshal(refNumberBytes, &refNumber)
	tx.RefNumber 	= strconv.Itoa(refNumber)
	refNumber = refNumber + 1;
	refNumberBytes, _ = json.Marshal(refNumber)
	err = stub.PutState("refNumber", refNumberBytes)								
	if err != nil {
		fmt.Println("Error Creating updating ref number")
		return nil, err
	}
	
	
	//***************************************************************
	// Get Receiver account from BC
	rfidBytes, err := stub.GetState(tx.To)
	if err != nil {
		return nil, errors.New("transferPoints Failed to get Receiver from BC")
	}
	var receiver User
	fmt.Println("transferPoints Unmarshalling User Struct");
	err = json.Unmarshal(rfidBytes, &receiver)
	receiver.Balance = receiver.Balance  + tx.Amount
	receiver.Modified = currentDateStr
	tx.ToName = receiver.Name;
	
	
	//Commit Receiver to ledger
	fmt.Println("transferPoints Commit Updated receiver To Ledger");
	txsAsBytes, _ := json.Marshal(receiver)
	err = stub.PutState(tx.To, txsAsBytes)	
	if err != nil {
		return nil, err
	}
	
	// Get Sender account from BC
	rfidBytes, err = stub.GetState(tx.From)
	if err != nil {
		return nil, errors.New("transferPoints Failed to get Financial Institution")
	}
	var sender User
	fmt.Println("transferPoints Unmarshalling Sender");
	err = json.Unmarshal(rfidBytes, &sender)
	sender.Balance   = sender.Balance  - tx.Amount
	sender.Modified = currentDateStr
	tx.FromName = sender.Name;
	
	//Commit Sender to ledger
	fmt.Println("transferPoints Commit Updated Sender To Ledger");
	txsAsBytes, _ = json.Marshal(sender)
	err = stub.PutState(tx.From, txsAsBytes)	
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
	//***********************************************************************
}

// ============================================================================================================================
func (t *SimpleChaincode) updateUserAccount(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	fmt.Println("Running updateUserAccount")

	userId := args[0]
	amountValue, err := strconv.ParseFloat(args[1], 64)
	

	// Get user account from the blockchain
	rfidBytes, err := stub.GetState(userId)
	if err != nil {
		return nil, errors.New("updateUserAccount Failed to get User from BC")
	}
	
	var account User
	fmt.Println("SubmitTx Unmarshalling User Struct");
	err = json.Unmarshal(rfidBytes, &account)
	account.Balance = account.Balance  + amountValue
	
	// Commit user account to ledger
	fmt.Println("SubmitTx Commit Updated user account To Ledger");
	txsAsBytes, _ := json.Marshal(account)
	err = stub.PutState(userId, txsAsBytes)	
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


