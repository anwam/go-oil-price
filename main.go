package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type OilPriceResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Text    string   `xml:",chardata"`
	Soap    string   `xml:"soap,attr"`
	Xsi     string   `xml:"xsi,attr"`
	Xsd     string   `xml:"xsd,attr"`
	Body    struct {
		Text                    string `xml:",chardata"`
		CurrentOilPriceResponse struct {
			Text                  string `xml:",chardata"`
			Xmlns                 string `xml:"xmlns,attr"`
			CurrentOilPriceResult string `xml:"CurrentOilPriceResult"`
		} `xml:"CurrentOilPriceResponse"`
	} `xml:"Body"`
}

type PTTORDS struct {
	XMLName xml.Name `xml:"PTTOR_DS" json:"-"`
	Text    string   `xml:",chardata" json:"-"`
	FUEL    []struct {
		Text      string `xml:",chardata" json:"-"`
		PRICEDATE string `xml:"PRICE_DATE" json:"priceDate"`
		PRODUCT   string `xml:"PRODUCT" json:"product"`
		PRICE     string `xml:"PRICE" json:"price"`
	} `xml:"FUEL" json:"fuel"`
}

func main() {
	url := "https://orapiweb.pttor.com/oilservice/OilPrice.asmx"
	method := "POST"

	payload := strings.NewReader(`<?xml version="1.0" encoding="utf-8"?>
<soap12:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap12="http://www.w3.org/2003/05/soap-envelope">
  <soap12:Body>
    <CurrentOilPrice xmlns="http://www.pttor.com">
      <Language>th</Language>
    </CurrentOilPrice>
  </soap12:Body>
</soap12:Envelope>`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/soap+xml; charset=utf-8")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp := OilPriceResponse{}
	xml.Unmarshal(body, &resp)
	detail := PTTORDS{}
	xml.Unmarshal([]byte(resp.Body.CurrentOilPriceResponse.CurrentOilPriceResult), &detail)
	for _, fuel := range detail.FUEL {
		fmt.Println(fuel.PRICEDATE, fuel.PRODUCT, fuel.PRICE)
	}
	jsonString, _ := json.Marshal(&detail)
	fmt.Println(string(jsonString))
}
