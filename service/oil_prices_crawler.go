package service

import (
	_ "fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	// OilPriceSrc - URL for looking up oil prices
	OilPriceSrc string = "http://www.qiyoujiage.com/"
)

// crawlOilPrices visits OilPriceSrc and returns oil price data
func crawlOilPrices() (map[string][]float64, error) {
	var res map[string][]float64
	response, err := http.Get(OilPriceSrc)
	if err != nil {
		return res, err
	}
	defer response.Body.Close()

	res, err = analyzeOilPrices(response)
	if err != nil {
		return res, err
	}
	return res, err
}

// analyzeOilPrices processes oil price data for each province
func analyzeOilPrices(resp *http.Response) (map[string][]float64, error) {
	res := make(map[string][]float64)
	body, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return res, err
	}

	// use goquery to select the oil price table 
	html, err := body.Find(".ylist").Html()
	if err != nil {
		return res, err
	}

	// get oil price data
	reg := regexp.MustCompile("<[^>]+>")
	priceList := strings.Fields(reg.ReplaceAllString(html, " "))
	//fmt.Println("Data found: ", priceList)
	for i := 5; i < len(priceList); i += 5 {
		var temp []float64 = make([]float64, 0)
		start := i + 1
		for j := 0; j < len(oilIndex); j++ {
			pos := start + j
			if pos >= len(priceList) {
				return res, oilPricesHandleError{"Wrong webpage format"}
			}
			price, err := strconv.ParseFloat(priceList[pos], 64)
			if err != nil {
				return res, err
			}
			temp = append(temp, price)
		}
		res[priceList[i]] = temp
	}

	return res, nil
}
