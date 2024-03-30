package models

import "github.com/go-playground/validator/v10"

type Model struct {
	OrderUid          string `json:"order_uid" validate:"required"`
	Tracknumber       string `json:"track_number" validate:"required"`
	Entry             string `json:"entry"`
	Delivery          `json:"delivery"`
	Payment           `json:"payment"`
	Items             []Item `json:"items"`
	Locale            string `json:"locale"`
	Internalsignature string `json:"internal_signature"`
	Customerid        string `json:"customer_id"`
	Deliveryservice   string `json:"delivery_service"`
	Shardkey          string `json:"shardkey"`
	Smid              int    `json:"sm_id"`
	Datecreated       string `json:"date_created"`
	Oofshard          string `json:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction" validate:"required"`
	Requestid    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	Paymentdt    int    `json:"payment_dt"`
	Bank         string `json:"bank"`
	Deliverycost int    `json:"delivery_cost"`
	Goodstotal   int    `json:"goods_total"`
	Customfee    int    `json:"custom_fee"`
}

type Item struct {
	Chrtid      int    `json:"chrt_id" validate:"required"`
	Tracknumber string `json:"track_number" validate:"required"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	Totalprice  int    `json:"total_price"`
	Nmid        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

func NewModel() *Model {
	return &Model{}
}

func (m *Model) Validate() error {
	validate := validator.New()
	return validate.Struct(m)
}
