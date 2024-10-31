package main

import (
	"crypto/x509"
	"encoding/json"
	"strconv"

	"fmt"
	"log"
	"net/http"

	"os"
	"path"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/hash"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	mspID        = "Org1MSP"
	cryptoPath   = "../organizations/peerOrganizations/org1.example.com"
	certPath     = cryptoPath + "/users/User1@org1.example.com/msp/signcerts"
	keyPath      = cryptoPath + "/users/User1@org1.example.com/msp/keystore"
	tlsCertPath  = cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt"
	peerEndpoint = "dns:///host.docker.internal:7051"
	gatewayPeer  = "peer0.org1.example.com"
)

type BankAccount struct {
	AccountID string `json:"AccountID"`
	Owner     string `json:"Owner"`
	Balance   int64  `json:"Balance"`
}

type Account struct {
	AccountID string `json:"AccountID"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	Sin       string `json:"sin"`
}

type TransactionHash struct {
	TransactionHash string    `json:"TransactionHash"`
	Activity        string    `json:"Activity"`
	Amount          float64   `json:"Amount"`
	AccountID       int       `json:"AccountID" gorm:"primaryKey"`
	Timestamp       time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;<-:create"`
}

type Response struct {
	From   string `json:"From"`
	Amount string `json:"Amount"`
}
type ResponsePayment struct {
	From   string `json:"From"`
	To	   string `json:"To"`
	Amount string `json:"Amount"`
}

// Declare 'contract' as a global variable
var contract *client.Contract

var db *gorm.DB

func init() {
	var err error
	// Database connection details
	dsn := "host=postgres user=testdev password=test123 dbname=BankAccountDB port=5432 sslmode=disable"

	// Connect to the database
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Unable to connect to database: ", err)
	}

	log.Println("Connected to PostgreSQL!")
}

func main() {
	server := gin.Default()

	// The gRPC client connection should be shared by all Gateway connections to this endpoint
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	id := newIdentity()
	sign := newSign()

	// Create a Gateway connection for a specific client identity
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithHash(hash.SHA256),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(30*time.Second), // Increase to 30 seconds or more
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gw.Close()

	// Override default values for chaincode and channel name as they may differ in testing contexts.
	chaincodeName := "basic"
	if ccname := os.Getenv("CHAINCODE_NAME"); ccname != "" {
		chaincodeName = ccname
	}

	channelName := "mychannel"
	if cname := os.Getenv("CHANNEL_NAME"); cname != "" {
		channelName = cname
	}

	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	initLedger(contract)

	// CORS middleware to allow requests from your React app (running on port 3000, for example)
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // React app address
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	server.GET("/account/:account_id", func(ctx *gin.Context) {
		account_id := ctx.Param("account_id")
		data := readAssetByID(contract, account_id)

		var account BankAccount
		if err := json.Unmarshal(data, &account); err != nil {
			ctx.JSON(500, gin.H{"error": "failed to parse account data"})
			return
		}

		ctx.JSON(200, gin.H{
			"data": account,
		})
	})

	server.GET("/cheq/:account_id", func(ctx *gin.Context) {
		accountId := ctx.Param("account_id")

		var transactionsHash []TransactionHash // This should be a slice to store multiple records

		if err := db.Where("account_id = ?", accountId).Order("timestamp DESC").Find(&transactionsHash).Error; err != nil {
			ctx.JSON(500, gin.H{"error": "Error fetching transaction hashes"})
		} else {
			ctx.JSON(200, gin.H{
				"data": transactionsHash,
			})
		}
	})

	server.POST("/deposit", func(ctx *gin.Context) {
		var res Response

		if err := ctx.ShouldBindJSON(&res); err != nil {
			panic(fmt.Errorf("failed to bind json: %w", err))
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		floatVal, errr := strconv.ParseFloat(res.Amount+".00", 64)
		if errr != nil {
			fmt.Println("Error parsing float:", err)
			return
		}

		err := depositMoneyByID(contract, res.From, floatVal)
		if err != nil {
			ctx.JSON(500, gin.H{"error": "failed to deposit"})
			return
		}

		ctx.JSON(200, gin.H{
			"data": "Deposit successful!",
		})
	})

	server.POST("/withdraw", func(ctx *gin.Context) {
		var res Response

		if err := ctx.ShouldBindJSON(&res); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		floatVal, errr := strconv.ParseFloat(res.Amount+".00", 64)
		if errr != nil {
			fmt.Println("Error parsing float:", err)
			return
		}

		err := withdrawMoneyByID(contract, res.From, floatVal)
		if err != nil {
			ctx.JSON(500, gin.H{"error": "failed to withdraw"})
			return
		}

		ctx.JSON(200, gin.H{
			"data": "Withdraw successful!",
		})
	})

	server.POST("/payment", func(ctx *gin.Context) {
		var res ResponsePayment

		if err := ctx.ShouldBindJSON(&res); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		floatVal, errr := strconv.ParseFloat(res.Amount+".00", 64)
		if errr != nil {
			fmt.Println("Error parsing float:", err)
			return
		}

		err := paymentMoneybyID(contract, res.From, res.To, floatVal)
		if err != nil {
			ctx.JSON(500, gin.H{"error": "failed to withdraw"})
			return
		}

		ctx.JSON(200, gin.H{
			"data": "Payment successful!",
		})
	})

	server.Run() // listen and serve on 0.0.0.0:8080
}

// newGrpcConnection creates a gRPC connection to the Gateway server.
func newGrpcConnection() *grpc.ClientConn {
	certificatePEM, err := os.ReadFile(tlsCertPath)
	if err != nil {
		panic(fmt.Errorf("failed to read TLS certifcate file: %w", err))
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)

	connection, err := grpc.NewClient(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection: %w", err))
	}

	return connection
}

// newIdentity creates a client identity for this Gateway connection using an X.509 certificate.
func newIdentity() *identity.X509Identity {
	certificatePEM, err := readFirstFile(certPath)
	if err != nil {
		panic(fmt.Errorf("failed to read certificate file: %w", err))
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(err)
	}

	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		panic(err)
	}

	return id
}

// newSign creates a function that generates a digital signature from a message digest using a private key.
func newSign() identity.Sign {
	privateKeyPEM, err := readFirstFile(keyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read private key file: %w", err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}

	return sign
}

func readFirstFile(dirPath string) ([]byte, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}

	fileNames, err := dir.Readdirnames(1)
	if err != nil {
		return nil, err
	}

	return os.ReadFile(path.Join(dirPath, fileNames[0]))
}

// This type of transaction would typically only be run once by an application the first time it was started after its
// initial deployment. A new version of the chaincode deployed later would likely not need to run an "init" function.
func initLedger(contract *client.Contract) {
	fmt.Printf("\n--> Submit Transaction: InitLedger, function creates the initial set of assets on the ledger \n")

	_, err := contract.SubmitTransaction("InitLedger")
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

// Evaluate a transaction by account_id to query ledger state.
func readAssetByID(contract *client.Contract, account_id string) []byte {
	fmt.Printf("\n--> Evaluate Transaction: ReadAsset, function returns asset attributes\n")

	evaluateResult, err := contract.EvaluateTransaction("ReadBankAccount", account_id)
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}

	return evaluateResult
}

// Submit a transaction synchronously, blocking until it has been committed to the ledger.
func depositMoneyByID(contract *client.Contract, account_id string, amount float64) error {
	fmt.Printf("\n--> Submit Transaction: DepositBalance, deposit money on to the client's account.  \n")

	strFloat := strconv.FormatFloat(amount, 'f', -1, 64)

	hash, err := contract.SubmitTransaction("DepositBalance", account_id, strFloat)
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")

	intAccountID, errr := strconv.Atoi(account_id)
	if errr != nil {
		// Handle the error, for example:
		fmt.Println("Error converting string to int:", errr)
	} else {
		// Use intAccountID as an integer
		fmt.Println("Converted Account ID:", intAccountID)
	}

	// Add hash to the transaction DB
	transactionHash := TransactionHash{
		TransactionHash: string(hash),
		Activity:        "deposit",
		Amount:          amount,
		AccountID:       int(intAccountID),
	}
	if err := db.Create(&transactionHash).Error; err != nil {
		return err
	}
	return nil
}

func withdrawMoneyByID(contract *client.Contract, account_id string, amount float64) error {
	fmt.Printf("\n--> Submit Transaction: WithdrawBalance, withdraw money on to the client's account.  \n")

	strFloat := strconv.FormatFloat(amount, 'f', -1, 64)

	hash, err := contract.SubmitTransaction("WithdrawBalance", account_id, strFloat)
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")

	intAccountID, errr := strconv.Atoi(account_id)
	if errr != nil {
		// Handle the error, for example:
		fmt.Println("Error converting string to int:", errr)
	} else {
		// Use intAccountID as an integer
		fmt.Println("Converted Account ID:", intAccountID)
	}

	// Add hash to the transaction DB
	transactionHash := TransactionHash{
		TransactionHash: string(hash),
		Activity:        "withdraw",
		Amount:          amount,
		AccountID:       int(intAccountID),
	}
	if err := db.Create(&transactionHash).Error; err != nil {
		return err
	}
	return nil
}

func paymentMoneybyID(contract *client.Contract, from string, to string, amount float64) error {
	fmt.Printf("\n--> Submit Transaction: TransactPayment, make a payment to another account.  \n")

	strFloat := strconv.FormatFloat(amount, 'f', -1, 64)

	hash, err := contract.SubmitTransaction("TransactPayment", from, to, strFloat)
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")

	intFromID, errr := strconv.Atoi(from)
	if errr != nil {
		// Handle the error, for example:
		fmt.Println("Error converting string to int:", errr)
	} else {
		// Use intFromID as an integer
		fmt.Println("Converted Account ID:", intFromID)
	}

	intToID, errrr := strconv.Atoi(to)
	if errrr != nil {
		// Handle the error, for example:
		fmt.Println("Error converting string to int:", errrr)
	} else {
		// Use intToID as an integer
		fmt.Println("Converted Account ID:", intToID)
	}

	// Add hash to the transaction DB
	transactionHashFrom := TransactionHash{
		TransactionHash: string(hash),
		Activity:        "payment",
		Amount:          amount,
		AccountID:       int(intFromID),
	}

	// Add hash to the transaction DB
	transactionHashTo := TransactionHash{
		TransactionHash: string(hash),
		Activity:        "received",
		Amount:          amount,
		AccountID:       int(intToID),
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&transactionHashFrom).Error; err != nil {
			return err
		}
		if err := tx.Create(&transactionHashTo).Error; err != nil {
			return err
		}
		return nil
	})
	
	if err != nil {
		return err // Handle the error appropriately
	}

	return nil
}