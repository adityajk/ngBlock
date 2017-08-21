package main

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
	"os/exec"
	"time"
	"log"
	"strconv"
)

var logger = shim.NewLogger("mylogger")

// This could be stored on a blockchain table or an application
var statusType = []string{"NEW", "TAX_DEDUCTED", "AMOUNT_CREDITED", "CLOSED"}

type EscrowChaincode struct {
}

type  Bank struct {
	 BankId				string		`json:"bankId"`
	 Name				string		`json:"name"`
	 AmountCredited		string		`json:"amountCredited"`
}

type TaxFinancialInfo struct { 
	TaxAuthorityName   		string		`json:"taxAuthorityName"`
    TaxPercentage        	int 		`json:"taxPercentage"`
    StartDate   			int 		`json:"startDate"`
    Frequency 				int 		`json:"frequency"`
    AmountCredited			string		`json:"amountCredited"`
    TaxCurrentBalance		int			`json:"taxCurrentBalance"`
}
 
type EscrowApplication struct {   
    ParcelId               string        		`json:"parcelId"`
    PropertyValue		   int					`json:"propertyValue"`
    CustomerId             string  		 		`json:"customerId"`
    CurrentBalance		   int 		 			`json:"currentBalance"`
   // TaxFinancialInfo       *TaxFinancialInfo 	`json:"taxFinancialInfo,omitempty"`
    Bank				   *Bank				`json:"bank,omitempty"`
    Source                 string        		`json:"source"`
    Status				   string				`json:"status"`
    LastModifiedDate       string       		`json:"lastModifiedDate"`
}

// Custom Event
type escrowEvent struct {
    Type        string `json:"type"`
    Description string `json:"description"`
}

func (t *EscrowChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    return nil, nil
}

func (t *EscrowChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
   
   switch function {
	   case "GetLastTransaction":
		   return GetLastTransaction(stub, args)
	   case "GetNoOfTransactions":
		   return GetNoOfTransactions(stub, args)
	   case "GetEscrowTransactionLog":
		   return GetEscrowTransactionLog(stub, args)
	   case "GetCurrentEscrowBalance":
		   return GetCurrentEscrowBalance(stub, args)
	   case "GetNextUpcomingPayments":
		   return GetNextUpcomingPayments(stub, args)
	   case "GetParcelPayments":
		   return GetParcelPayments(stub, args)	 
	   default:
			return nil, nil	        	   		   
   }
   
   /**
    if function == "GetEscrowApplication" {
		return GetEscrowApplication(stub, args)
	}
	return nil, nil
	**/
}
 
func (t *EscrowChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	
	switch function {
		case "CreditIntoEscrowAccount":
			return CreditIntoEscrowAccount(stub, args)
		case "PerformEscrowTaxDeduction":
			return PerformEscrowTaxDeduction(stub, args)
		case "EscrowAmountManualCredit":
			return EscrowAmountManualCredit(stub, args)
		case "ImportEscrowAccount":
			return ImportEscrowAccount(stub, args)	
		default:
			return nil, nil		
	}
	
	/**
    if function == "CreditIntoEscrowAccount" {
        username, _ := GetCertAttribute(stub, "username")
        role, _ := GetCertAttribute(stub, "role")
        if role == "Bank_Escrow_Admin" {
            return CreditIntoEscrowAccount(stub, args)
        } else {
            return nil, errors.New(username + " with role " + role + " does not have access to create a escrow application block")
        }
 
    } else if function == "CreditIntoEscrowAccount" {
    
    }
    **/
   // return nil, nil
}

/**********************************************************************************/
/************************** Invoke APIs *******************************************/
/**********************************************************************************/
func CreditIntoEscrowAccount(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    fmt.Println("Entering CreditIntoEscrowAccount")
 
    if len(args) != 7 {
        fmt.Println("Invalid input args")
        return nil, errors.New("Expected the escrow details for Escrow Block creation")
    }
 
    out, err := exec.Command("uuidgen").Output()
    if err != nil {
        log.Fatal(err)
        return nil, errors.New("Not able to generate Unique ID")
    }
    uId := string(out[:])
    var escrowApplicationId = "e_"+uId
    var bankId = "user_type1_1"
    //var escrowApplicationInput = args[1]
    //fmt.Printf("%s", out)
   
   propVal, err := strconv.Atoi(args[1])
   if err != nil {
		fmt.Println("Int conversion error: ", err)
		return nil, err
	}
   curBal, err := strconv.Atoi(args[3])
   if err != nil {
		fmt.Println("Int conversion error: ", err)
		return nil, err
	}
   mapD := &Bank{bankId, args[4], args[5]}
   var dTime = time.Now().UTC().Format("2006-01-02 15:04:05 UCT") 
    
	escrowApplicationInput := EscrowApplication {							 
							 args[0],
							 propVal,
							 args[2],
							 curBal,
							 mapD,
							 args[6],
							 statusType[0],
							 dTime,
						 }
	
	ajson, err := json.Marshal(escrowApplicationInput)
	if err != nil {
		fmt.Println("toJSON error: ", err)
		return nil, err
	}
	 
    err = stub.PutState(escrowApplicationId, []byte(ajson))
    if err != nil {
        fmt.Println("Could not save escrow application to ledger", err)
        return nil, err
    }
    
    var event = escrowEvent{"CreditIntoEscrowAccount", "Successfully created escrow application with ID " + escrowApplicationId}
    eventBytes, err := json.Marshal(&event)
    
   
    if err != nil {
        return nil, err
    }
    
    fmt.Println("Successfully saved escrow application")
    return eventBytes, nil
}

func PerformEscrowTaxDeduction(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	 fmt.Println("Entering PerformEscrowTaxDeduction")
	return nil, nil
}

func EscrowAmountManualCredit(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	 fmt.Println("Entering EscrowAmountManualCredit")
	return nil, nil
}

func ImportEscrowAccount(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	 fmt.Println("Entering ImportEscrowAccount")
	return nil, nil
}

/**********************************************************************************/
/************************** Query APIs *******************************************/
/**********************************************************************************/
func GetLastTransaction(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    fmt.Println("Entering GetLastTransaction")
 
    if len(args) < 1 {
        fmt.Println("Invalid number of arguments")
        return nil, errors.New("Missing Escrow application ID")
    }
 
    var escrowApplicationId = args[0]
    bytes, err := stub.GetState(escrowApplicationId)
    if err != nil {
        fmt.Println("Could not fetch escrow application with id "+escrowApplicationId+" from ledger", err)
        return nil, err
    }
    return bytes, nil
}

func GetNoOfTransactions(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    fmt.Println("Entering GetNoOfTransactions")
	
	var escrowApplicationId = args[0]
    bytes, err := stub.GetState(escrowApplicationId)
	 if err != nil {
        return nil, err
    }
	return bytes, nil
}  

func GetEscrowTransactionLog(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    fmt.Println("Entering GetEscrowTransactionLog")
	
	var escrowApplicationId = args[0]
    bytes, err := stub.GetState(escrowApplicationId)
	 if err != nil {
        return nil, err
    }
	return bytes, nil
}

func GetCurrentEscrowBalance(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    fmt.Println("Entering GetCurrentEscrowBalance")
	
	var escrowApplicationId = args[0]
    bytes, err := stub.GetState(escrowApplicationId)
	 if err != nil {
        return nil, err
    }
	return bytes, nil
}

func GetNextUpcomingPayments(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    fmt.Println("Entering GetNextUpcomingPayments")
	
	var escrowApplicationId = args[0]
    bytes, err := stub.GetState(escrowApplicationId)
	 if err != nil {
        return nil, err
    }
	return bytes, nil
}

func GetParcelPayments(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    fmt.Println("Entering GetParcelPayments")
	
	var escrowApplicationId = args[0]
    bytes, err := stub.GetState(escrowApplicationId)
	 if err != nil {
        return nil, err
    }
	return bytes, nil
}  

/**
Updates the status of the Escrow application

func UpdateEscrowApplication(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Debug("Entering UpdateEscrowApplication")

	if len(args) < 2 {
		logger.Error("Invalid number of args")
		return nil, errors.New("Expected atleast two arguments for escrow application update")
	}

	var escrowApplicationId = args[0]
	var status = args[1]

	eaBytes, err := stub.GetState(escrowApplicationId)
	if err != nil {
		logger.Error("Could not fetch escrow application from ledger", err)
		return nil, err
	}
	var escrowApplication EscrowApplication
	err = json.Unmarshal(eaBytes, &escrowApplication)
	escrowApplication.Status = status

	eaBytes, err = json.Marshal(&escrowApplication)
	if err != nil {
		logger.Error("Could not marshal escrow application post update", err)
		return nil, err
	}

	err = stub.PutState(escrowApplicationId, eaBytes)
	if err != nil {
		logger.Error("Could not save escrow application post update", err)
		return nil, err
	}

	var escrowEvent = "{eventType: 'escrowApplicationUpdate', description:" + escrowApplicationId + "' Successfully updated status'}"
	err = stub.SetEvent("evtSender", []byte(escrowEvent))
	if err != nil {
		return nil, err
	}
	logger.Info("Successfully updated escrow application")
	return nil, nil

}
**/

//Retrieving attributes from the transaction certificate of the caller
/**
func GetCertAttribute(stub shim.ChaincodeStubInterface, attributeName string) (string, error) {
    fmt.Println("Entering GetCertAttribute")
    attr, err := stub.ReadCertAttribute(attributeName)
    if err != nil {
        return "", errors.New("Couldn't get attribute " + attributeName + ". Error: " + err.Error())
    }
    attrString := string(attr)
    return attrString, nil
}
**/
 

 
func main() {
	lld, _ := shim.LogLevel("DEBUG")
	fmt.Println(lld)

	logger.SetLevel(lld)
	fmt.Println(logger.IsEnabledFor(lld))

	err := shim.Start(new(EscrowChaincode))
	if err != nil {
		logger.Error("Could not start SampleChaincode")
	} else {
		logger.Info("SampleChaincode successfully started")
	}
}
