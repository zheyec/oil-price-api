package service

import (
	"fmt"
	"net/http"
)

const (
	// DBURL - URL for database
	DBURL string = "mongodb://127.0.0.1:27017"
)

var (
	oilIndex = map[string]int{
		"#92 Gas": 0,
		"#95 Gas": 1,
		"#98 Gas": 2,
		"#0 Diesel":  3,
	}

	oilName = []string{"#92 Gas", "#95 Gas", "#98 Gas", "#0 Diesel"}
)

// oilPricesHandleError - custom error class
type oilPricesHandleError struct {
	Msg string
}

func (err oilPricesHandleError) Error() string {
	return err.Msg
}

// returnErr is called when an error occurs; the default response is empty
func returnErr(e error, w http.ResponseWriter) {
	fmt.Printf("Error: %s\n", e.Error())
	resp := &Response{}
	resp.ErrorNo = e.Error()
	fmt.Printf("[Error]%s \n", resp.ToBytes())
	w.Write(resp.ToBytes())
}

// formatData - formats oil prices into reply
func formatData(data []float64, province string, oilType string) string {

	// Handles NA specially
	var text []string = make([]string, len(data))
	for i, u := range data {
		if u != 0 {
			text[i] = fmt.Sprintf("%v yuan/liter", u)
		} else {
			text[i] = "NA"
		}
	}

	// generate response
	reply := ""
	if oilType == "" {
		reply = fmt.Sprintf("Oil price today in %s: ", province)
		for idx, name := range oilName {
			reply += fmt.Sprintf("\n%s: %s", name, text[idx])
		}
	} else {
		reply = fmt.Sprintf("Price of %s today in %s: %s", oilType, province, text[0])
	}
	return reply
}

// main function
func oilPricesHandle(w http.ResponseWriter, r *http.Request) {

	province := r.URL.Query().Get("prov")
	oilType := r.URL.Query().Get("oil")

	// init
	db := OilPriceDB{}
	defer db.Close()
	err := db.Init(DBURL)
	if err != nil {
		returnErr(err, w)
		return
	}

	// read oil prices
	data, err := db.Read(province, oilType)
	if err != nil {
		returnErr(err, w)
		return
	}
	fmt.Println("Found oil prices: ", data)

	// generate response
	reply := formatData(data, province, oilType)
	resp := &Response{}
	resp.Slots = append(resp.Slots, Slots{"result", reply})
	fmt.Printf("[Response]%s \n", resp.ToBytes())
	w.Write(resp.ToBytes())

}
