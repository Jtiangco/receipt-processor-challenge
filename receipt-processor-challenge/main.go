package main

import (
	"errors"
	"math"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type receipt struct {
	ID           string `json:"id"`
	Retailer     string `json:"retailer"`
	PurchaseTime string `json:"purchaseTime"`
	PurchaseDate string `json:"purchaseDate"`
	Items        []struct {
		ShortDescription string  `json:"shortDescription"`
		Price            float64 `json:"price,string"`
	} `json:"items"`
	Total float64 `json:"total,string"`
}

type id struct {
	ID string `json:"id"`
}

type points struct {
	Points int `json:"points"`
}

var receipts = []receipt{}

func processReceipt(c *gin.Context) {
	var savedReceipt receipt
	var returnIdJson id

	if err := c.BindJSON(&savedReceipt); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "BadRequest"})
		return
	}

	savedReceipt.ID = uuid.NewString()

	receipts = append(receipts, savedReceipt)

	returnIdJson.ID = savedReceipt.ID

	c.IndentedJSON(http.StatusOK, returnIdJson)
}

func recieptById(c *gin.Context) {
	id := c.Param("id")
	reciept, err := getRecieptById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "No receipt found for that id"})
		return
	}

	c.IndentedJSON(http.StatusOK, reciept)
}

func getRecieptById(id string) (*receipt, error) {
	for i, r := range receipts {
		if r.ID == id {
			return &receipts[i], nil
		}
	}

	return nil, errors.New("not found")
}

func getPoints(c *gin.Context) {
	var Points points
	var runningTotal float64 = 0

	id := c.Param("id")
	receipt, err := getRecieptById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "No receipt found for that id"})
		return
	}

	dateLayout := "2006-01-02"
	timeLayout := "15:04"

	date, err := time.Parse(dateLayout, receipt.PurchaseDate)
	if err != nil {
		panic(err)
	}

	startTime, err := time.Parse(timeLayout, "14:00")
	if err != nil {
		panic(err)
	}
	endTime, err := time.Parse(timeLayout, "16:00")
	if err != nil {
		panic(err)
	}

	time, err := time.Parse(timeLayout, receipt.PurchaseTime)
	if err != nil {
		panic(err)
	}

	//Business logic and addition of points

	//One point for every alphanumeric character in the retailer name.
	var nonAlphanumeric = regexp.MustCompile(`[^a-zA-Z0-9]+`)
	var retailerAlphanumeric = nonAlphanumeric.ReplaceAllString(receipt.Retailer, "")
	runningTotal += float64(len(retailerAlphanumeric))

	//50 points if the total is a round dollar amount with no cents.
	if math.Mod(receipt.Total*100, 100) == 0 {
		runningTotal += 50
	}

	//25 points if the total is a multiple of 0.25.
	if math.Mod(receipt.Total, .25) == 0 {
		runningTotal += 25
	}

	//5 points for every two items on the receipt.
	runningTotal += float64((len(receipt.Items) / 2)) * 5

	for i := range receipt.Items {
		var trimmedDesc int = len(strings.TrimSpace(receipt.Items[i].ShortDescription))

		//If the trimmed length of the item description is a multiple of 3,
		if math.Mod(float64(trimmedDesc), 3) == 0 {
			//multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
			runningTotal += math.Ceil(receipt.Items[i].Price * .2)
		}
	}

	//6 points if the day in the purchase date is odd.
	if math.Mod(float64(date.Day()), 2) != 0 {
		runningTotal += 6
	}

	//10 points if the time of purchase is after 2:00pm and before 4:00pm.
	if time.After(startTime) && time.Before(endTime) {
		runningTotal += 10
	}

	Points.Points = int(runningTotal)
	c.IndentedJSON(http.StatusOK, Points)
}

func main() {
	router := gin.Default()
	router.POST("/receipts/process", processReceipt)
	router.GET("/receipt/:id", recieptById) //API created for testing purposes
	router.GET("/receipts/:id/points", getPoints)
	router.Run("localhost:8080")
}
