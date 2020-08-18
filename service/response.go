package service

import (
	"encoding/json"

	"fmt"
)

type Response struct {
	Slots     []Slots       `json:"slots"`
	Responses []ResponseMsg `json:"responses"`
	ErrorNo   string        `json:"error_no"`
}

type Slots struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ResponseMsg struct {
	Content string `json:"content"`
	Type    int    `json:"type"`
}

func NewResponse(code string) *Response {
	return &Response{
		ErrorNo: code,
	}
}

//ToBytes convert to bytes
func (x *Response) ToBytes() []byte {
	bytes, err := json.Marshal(x)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	return bytes
}

type NLGAnswer struct {
	KeyRes  []string `json:"key_res"`  
	NlgList []NLG    `json:"nlg_list"` 
	Pattern struct {
	} `json:"pattern"`
}

type NLG struct {
	Content string `json:"content"`
	Type    int    `json:"type"`
}

//ToBytes convert to bytes
func (x *NLGAnswer) ToBytes() []byte {
	bytes, err := json.Marshal(x)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	return bytes
}
