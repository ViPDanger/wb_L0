package structures

type Get_Page struct {
	Title            string
	SearchAllAuthor  string
	SearchAllContent string
	SearchAuthorID   string
	SearchContentID  string
	InputAuthor      string
	InputContent     string
	DeleteAuthorID   string
	DeleteContentID  string
}

type items struct {
	chrt_id      int    `json:"chrt_id"`
	track_number string `json:"track_number"`
	price        int    `json:"price"`
	rid          string `json:"rid"`
	name         string `json:"name"`
	sale         int    `json:"sale"`
	size         int    `json:"size"`
	total_price  int    `json:"total_price"`
	nm_id        int    `json:"nm_id"`
	brand        string `json:"brand"`
	status       int    `json:"status"`
}

type payment struct {
	chrt_id int `json:"chrt_id"`

	transaction   string `json:"transaction"`
	request_id    string `json:"request_id"`
	currency      string `json:"currency"`
	provider      string `json:"provider"`
	amount        int    `json:"amount"`
	payment_dt    int    `json:"payment_dt"`
	bank          string `json:"bank"`
	delivery_cost int    `json:"delivery_cost"`
	goods_total   int    `json:"goods_total"`
	custom_fee    int    `json:"custom_fee"`
}

type delivery struct {
	name    string `json:"name"`
	phone   string `json:"phone"`
	zip     string `json:"zip"`
	city    string `json:"city"`
	address string `json:"address"`
	region  string `json:"region"`
	email   string `json:"email"`
}

type order struct {
	order_uid          string   `json:"order_uid"`
	track_number       string   `json:"track_number"`
	entry              string   `json:"entry"`
	delivery           delivery `json:"delivery"`
	payment            payment  `json:"payment"`
	items              []items  `json:"items"`
	locale             string   `json:"locale"`
	internal_signature string   `json:"internal_signature"`
	customer_id        string   `json:"customer_id"`
	delivery_service   string   `json:"delivery_service"`
	shardkey           string   `json:"shardkey"`
	sm_id              int      `json:"sm_id"`
	date_created       string   `json:"date_created"`
	oof_shard          string   `json:"entry"`
}

type Result_Page struct {
	Title      string
	First_Line []string
	Data       any
	BackButton string
}
