package handler

import (
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
)

// ChargeRequest contains charging information
type ChargeRequest struct {
	ID string `bson:"_id"`
	AccountID string
	CardNumber string
	CardSeries string
	Amount int64
	Status string
}

// ChargeData contains charging input data
type ChargeData struct {
	AccountID string `json:"account_id"`
	CardNumber string `json:"card_number"`
	CardSeries string `json:"card_series"`
	Amount int64 `json:"amount"`
}

// ErrorMessage error message to send to client
type ErrorMessage struct {
	Message string
}

var (
	dbURL string
	dbName string
	db *mgo.Database
)

// Initialize logging and database
func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	dbURL = "mongodb://localhost"
	dbName = "tomcuaca"
	s, err := mgo.Dial(dbURL)
	if err != nil {
		log.Fatal(err)
	}
	db = s.DB(dbName)
}

// ListChargeByAccountID list all charging request by account id
func ListChargeByAccountID(c echo.Context) error {
	log.Println("Enter ListChargeByAccountID")
	defer log.Println("Exit ListChargeByAccountID")
	accountID := c.Param("id")

	var charges []ChargeRequest

	collection := db.C("charge")
	err := collection.Find(bson.M{"accountid": accountID}).All(&charges)
	if err != nil {
		log.Println("ListChargeByAccountID: ", err)
		return c.JSON(http.StatusInternalServerError, ErrorMessage{"Internal Server Error"})
	}

	return c.JSON(http.StatusOK, charges)
}

// CreateChargeRequest create a charging request record
func CreateChargeRequest(c echo.Context) error {
	log.Println("Enter CreateChargeRequest")
	defer log.Println("Exit CreateChargeRequest")
	r := c.Request()
	data, err := ioutil.ReadAll(r.Body())
	if err != nil {
		log.Println("CreateChargeRequest: ", err)
		return c.JSON(http.StatusBadRequest, ErrorMessage{"Bad json input data"})
	}
	log.Println("Input params: ", string(data))

	var chargeData ChargeData
	err = json.Unmarshal(data, &chargeData)
	if err != nil {
		log.Println("CreateChargeRequest: Unmarshalling: ", err)
		return c.JSON(http.StatusBadRequest, ErrorMessage{"Bad json input data"})
	}
	log.Println("Json input params: ", chargeData)

	if chargeData.AccountID == "" || chargeData.CardNumber == "" || chargeData.CardSeries == "" || chargeData.Amount <= 0 {
		log.Println("CreateChargeRequest: Invalid Input Data")
		return c.JSON(http.StatusBadRequest, ErrorMessage{"Invalid Input Data"})
	}

	charge:= ChargeRequest{ID: uuid.NewV4().String(),
				AccountID: chargeData.AccountID,
				CardNumber: chargeData.CardNumber,
				CardSeries: chargeData.CardSeries,
				Amount: chargeData.Amount,
				Status: "pending"}

	collection := db.C("charge")
	err = collection.Insert(&charge)
	if err != nil {
		log.Println("CreateChargeRequest: ", err)
		return c.JSON(http.StatusInternalServerError, ErrorMessage{"Internal Server Error"})
	}
	return c.JSON(http.StatusCreated, charge)
}

// UpdateChargeRequest update charging request record
func UpdateChargeRequest(c echo.Context) error {
	log.Println("Enter UpdateChargeRequest")
	defer log.Println("Exit UpdateChargeRequest")
	id := c.Param("id")

	collection := db.C("charge")
	var charge ChargeRequest
	err := collection.FindId(id).One(&charge)
	if err != nil {
		log.Println("UpdateChargeRequest: ", err)
		return c.JSON(http.StatusNotFound, ErrorMessage{"Not Found"})
	}

	if charge.Status != "pending" {
		log.Println("UpdateChargeRequest: try to approve non-pending charge")
		return c.JSON(http.StatusBadRequest, ErrorMessage{"Charge status is not pending"})
	}

	charge.Status = "success"
	err = collection.UpdateId(charge.ID, &charge)
	if err != nil {
		log.Println("UpdateChargeRequest: ", err)
		return c.JSON(http.StatusInternalServerError, ErrorMessage{"Internal Server Error"})
	}
	return c.JSON(http.StatusOK, charge)
}
