package model

type ResponseList struct {
	Page      int `json:"page"`
	PageTotal int `json:"page_total"`
	Data      any `json:"data"`
	DataTotal int `json:"data_total"`
	Error     any `json:"errors"`
}

type Response struct {
	Data  any `json:"data"`
	Error any `json:"errors"`
}
