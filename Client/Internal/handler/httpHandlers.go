package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	n "github.com/ViPDanger/L0/Client/Internal/nats"
	"github.com/ViPDanger/L0/Client/Internal/structures"
	"github.com/bluele/gcache"
	"github.com/nats-io/nats.go"
)

type ShutdownHandler struct {
	cancel_on_http context.CancelFunc
}

func (sht ShutdownHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Shutdown by http")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Server is going got Killed. Murderer."))
	sht.cancel_on_http()
}

func SetupHttpHandlers(mux *http.ServeMux, cancel_on_http context.CancelFunc, natsHandler NatsHanlder, cache *gcache.Cache) {
	mux.Handle("/Shutdown", ShutdownHandler{cancel_on_http: cancel_on_http})
	mux.HandleFunc("/", DefaultPage)
	mux.HandleFunc("/response/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		path := r.URL.Path
		path = path[strings.LastIndexAny(path[1:], "/")+2:]
		switch path {
		case "Delete":
			body, err := io.ReadAll(r.Body)
			body = []byte(strings.TrimPrefix(string(body), "request_uid="))
			if err != nil || len(body) < 2 {
				errorPage(w, r, errors.New("request_uid не может быть пустым"))
				return
			}

			// Удаление Order при его наличии в кэше
			(*cache).Remove(string(body))

			body, _ = json.Marshal(&structures.Order_uid{
				Order_uid: string(body),
			})
			n.PublishMsg(*natsHandler.JetstreamContext, "DelOrder."+natsHandler.ClientName, body)
			var msg *nats.Msg
			msg = <-natsHandler.DelChannel
			delResultPage(w, r, msg.Data)
		case "Get":
			body, err := io.ReadAll(r.Body)
			body = []byte(strings.TrimPrefix(string(body), "request_uid="))
			if err != nil || len(body) < 1 {
				errorPage(w, r, errors.New("request_uid не может быть пустым"))
				return
			}
			// Проверка Кэша на наличие Order
			if (*cache).Has(string(body)) {
				cacheOrder, _ := (*cache).Get(string(body))
				getResultPage(w, r, cacheOrder.(structures.Order))
				return
			}

			body, _ = json.Marshal(&structures.Order_uid{
				Order_uid: string(body),
			})

			n.PublishMsg(*natsHandler.JetstreamContext, "GetOrder."+natsHandler.ClientName, body)
			var msg *nats.Msg
			msg = <-natsHandler.GetChannel
			// Форматирование полученного значения
			var order structures.Order
			err = json.Unmarshal(msg.Data, &order)
			if err != nil {
				errorPage(w, r, errors.New(string(msg.Data)))
			} else {
				getResultPage(w, r, order)
			}
		case "Put":
			body, err := io.ReadAll(r.Body)
			if err != nil || len(body) < 1 {
				errorPage(w, r, errors.New("Запрос не может быть пустым"))
				return
			}
			order, err := putOrderDataTest(w, r, body)
			if err != nil {
				errorPage(w, r, err)
				return
			}
			if (*cache).Has(string(body)) {
				errorPage(w, r, errors.New("Order with order_uid "+string(body)+" is already in the base (and cache)"))
			}

			body, _ = json.Marshal(order)
			n.PublishMsg(*natsHandler.JetstreamContext, "PutOrder."+natsHandler.ClientName, body)
			var msg *nats.Msg
			msg = <-natsHandler.PutChannel
			// Добавление значения в кэш
			if string(msg.Data) == "Order "+order.Order_uid+" Добавлен успешно" {
				err := (*cache).Set(order.Order_uid, order)
				println(err)
			}
<<<<<<< HEAD
=======

>>>>>>> b8b5f62acf310e5886c10669814987e1c23c205d
			putResultPage(w, r, msg.Data)
		}
	})
}

func DefaultPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	t, err := template.ParseFiles("./../Internal/templates/default.html")
	if err != nil {
		log.Println("Template:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to parse files"))
		return
	}
	// Изменения кол-ва Items
	int := 1
	if true { // OK
		body, _ := io.ReadAll(r.Body)
		if strings.HasPrefix(string(body), "ItemCount=") {
			int, err = strconv.Atoi(strings.TrimPrefix(string(body), "ItemCount="))
			if err != nil || int < 0 {
				log.Println("Defaultpage: ", err)
				int = 1
				errorPage(w, r, errors.New("Число Items должно быть натуральным"))
			}
		}
	}
	items := make([]structures.Items, int)
	if len(items) > 0 {
		items[0] = structures.Items{
			Chrt_id:      9934930,
			Track_number: "WBILMTESTTRACK",
			Price:        453,
			Rid:          "ab4219087a764ae0btest",
			Name:         "Mascaras",
			Sale:         30,
			Size:         "0",
			Total_price:  317,
			Nm_id:        2389212,
			Brand:        "Vivienne Sabo",
			Status:       202,
		}
	}
<<<<<<< HEAD
=======

>>>>>>> b8b5f62acf310e5886c10669814987e1c23c205d
	page := structures.DefaultPage{
		Title:          "HTTP-Postgress API",
		GetOrderButton: "GetOrderByOrder_uid",
		PutOrderButton: "PutNewOrder",
		DelOrderButton: "DeleteOrderByOrder_uid",
		Order: structures.Order{
			Order_uid:    "b563feb7b2b84b6test",
			Track_number: "WBILMTESTTRACK",
			Entry:        "WBIL",
			Delivery: structures.Delivery{
				Name:    "Test Testov",
				Phone:   "+9720000000",
				Zip:     "2639809",
				City:    "Kiryat Mozkin",
				Address: "Ploshad Mira 15",
				Region:  "Kraiot",
				Email:   "test@gmail.com",
			},
			Payment: structures.Payment{
				Transaction:   "b563feb7b2b84b6test",
				Request_id:    "",
				Currency:      "USD",
				Provider:      "wbpay",
				Amount:        1817,
				Payment_dt:    1637907727,
				Bank:          "alpha",
				Delivery_cost: 1500,
				Goods_total:   317,
				Custom_fee:    0,
			},
			Items:              items,
			Locale:             "en",
			Internal_signature: "",
			Customer_id:        "test",
			Delivery_service:   "meest",
			Shardkey:           "9",
			Sm_id:              99,
			Date_created:       "2021-11-26T06:22:19Z",
			Oof_shard:          "1",
		},
	}
	t.Execute(w, &page)
}
func errorPage(w http.ResponseWriter, r *http.Request, errors error) {
	w.Header().Set("Content-type", "text/html")
	t, err := template.ParseFiles("./../Internal/templates/error.html")
	if err != nil {
		log.Println("Template:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to parse files"))
		return
	}
	page := structures.ErrorPage{
		Title:        "Error",
		ErrorMessage: errors.Error(),
	}
	t.Execute(w, &page)
}
<<<<<<< HEAD
func putResultPage(w http.ResponseWriter, r *http.Request, data []byte) {
=======

func putOrderDataTest(w http.ResponseWriter, r *http.Request, data []byte) (structures.Order, error) {
	var order structures.Order
	arr := make([]int, 6)
	for j, s := range []string{"Payment.Amount", "Payment.Payment_dt", "Payment.Delivery_cost", "Payment.Goods_total", "Payment.Custom_fee", "Sm_id"} {
		i, err := strconv.Atoi(GetfromHTMLData(&data, []byte(s)))
		//log.Println(s, ":", i)
		if err != nil || i < 0 {
			return order, errors.New(s + " Должна являтся натуральным числом")
		}
		arr[j] = i
	}
	order = structures.Order{
		Order_uid:    GetfromHTMLData(&data, []byte("Order_uid")),
		Track_number: GetfromHTMLData(&data, []byte("Track_number")),
		Entry:        GetfromHTMLData(&data, []byte("Entry")),
		Delivery: structures.Delivery{
			Name:    GetfromHTMLData(&data, []byte("Delivery.Name")),
			Phone:   GetfromHTMLData(&data, []byte("Delivery.Phone")),
			Zip:     GetfromHTMLData(&data, []byte("Delivery.Zip")),
			City:    GetfromHTMLData(&data, []byte("Delivery.City")),
			Address: GetfromHTMLData(&data, []byte("Delivery.Address")),
			Region:  GetfromHTMLData(&data, []byte("Delivery.Region")),
			Email:   GetfromHTMLData(&data, []byte("Delivery.Email")),
		},
		Payment: structures.Payment{
			Transaction:   GetfromHTMLData(&data, []byte("Payment.Transaction")),
			Request_id:    GetfromHTMLData(&data, []byte("Payment.Request_id")),
			Currency:      GetfromHTMLData(&data, []byte("Payment.Currency")),
			Provider:      GetfromHTMLData(&data, []byte("Payment.Provider")),
			Amount:        arr[0],
			Payment_dt:    arr[1],
			Bank:          GetfromHTMLData(&data, []byte("Payment.Bank")),
			Delivery_cost: arr[2],
			Goods_total:   arr[3],
			Custom_fee:    arr[4],
		},
		Locale:             GetfromHTMLData(&data, []byte("Locale")),
		Internal_signature: GetfromHTMLData(&data, []byte("Internal_signature")),
		Customer_id:        GetfromHTMLData(&data, []byte("Customer_id")),
		Delivery_service:   GetfromHTMLData(&data, []byte("Delivery_service")),
		Shardkey:           GetfromHTMLData(&data, []byte("Shardkey")),
		Sm_id:              arr[5],
		Date_created:       GetfromHTMLData(&data, []byte("Date_created")),
		Oof_shard:          GetfromHTMLData(&data, []byte("Oof_shard")),
	}
	// ТЕСТЫ ORDER
	if order.Order_uid == "" {
		return order, errors.New("Order_uid не должно быть пустым")
	}
	if order.Track_number == "" {
		return order, errors.New("Track_number не должно быть пустым")
	}
	if order.Entry == "" {
		return order, errors.New("Entry не должно быть пустым")
	}
	_, err := strconv.Atoi(strings.TrimPrefix(order.Delivery.Phone, "+"))
	if len(order.Delivery.Phone) != 11 || !strings.HasPrefix(order.Delivery.Phone, "+") || err != nil {
		return order, errors.New("Delivery.Phone должен иметь формат +12345667890")
	}
	_, err = time.Parse("2006-01-02T15:04:05Z", order.Date_created)
	if err != nil {
		return order, errors.New("Date_created должно быть установлено в формате yyyy-mm-dd'T'HH:mm:ss'Z'")
	}
	// ITEMS
	items := make([]structures.Items, 0)
	for len(data) != 0 {
		for j, s := range []string{"Items.Chrt_id", "Items.Price", "Items.Sale", "Items.Total_price", "Items.Nm_id", "Items.Status"} {
			i, err := strconv.Atoi(GetfromHTMLData(&data, []byte(s)))
			//log.Println(s, ":", i)
			if err != nil || i < 0 {
				return order, errors.New("Items" + strconv.Itoa(len(items)) + "." + s + " Должна являтся натуральным числом")
			}
			arr[j] = i
		}
		items = append(items, structures.Items{
			Chrt_id:      arr[0],
			Track_number: GetfromHTMLData(&data, []byte("Items.Track_number")),
			Price:        arr[1],
			Rid:          GetfromHTMLData(&data, []byte("Items.Rid")),
			Name:         GetfromHTMLData(&data, []byte("Items.Name")),
			Sale:         arr[2],
			Size:         GetfromHTMLData(&data, []byte("Items.Size")),
			Total_price:  arr[3],
			Nm_id:        arr[4],
			Brand:        GetfromHTMLData(&data, []byte("Items.Brand")),
			Status:       arr[5],
		})

		// ТЕСТЫ ITEMS
	}
	order.Items = items
	return order, nil
}
func GetfromHTMLData(data *[]byte, nameInHTML []byte) string {
	var firstId, lastId, endmarkId, counter int
	firstId = -1
	lastId = -1
	endmarkId = -1
	lenData := len(*data)
	for i := 0; i < lenData; i++ {
		if lastId == -1 {
			if (*data)[i] == nameInHTML[counter] {
				if counter == 0 {
					firstId = i
				}
				counter++
				if counter == len(nameInHTML) {
					lastId = i
				}
			} else {
				firstId = -1
				counter = 0
			}
		} else {
			// Перевод
			if (*data)[i] == '%' {
				if (*data)[i+1] == '2' && (*data)[i+2] == 'B' {
					(*data)[i] = '+'
					(*data) = append((*data)[:i+1], (*data)[i+3:]...)
					lenData = len(*data)
				}
				if (*data)[i+1] == '3' && (*data)[i+2] == 'A' {
					(*data)[i] = ':'
					(*data) = append((*data)[:i+1], (*data)[i+3:]...)
					lenData = len(*data)
				}
			}
			// Завершение
			if (*data)[i] == '&' {
				endmarkId = i
				break
			}
		}
	}
	if firstId == -1 || lastId == -1 {
		return ""
	}
	var str string
	if endmarkId == -1 {
		endmarkId = len(*data)
		str = string((*data)[lastId+2 : endmarkId])
		(*data) = (*data)[:firstId]
	} else {
		str = string((*data)[lastId+2 : endmarkId])
		(*data) = append((*data)[:firstId], (*data)[endmarkId+1:]...)
	}
	return str
}

func delResultPage(w http.ResponseWriter, r *http.Request, data []byte) {
>>>>>>> b8b5f62acf310e5886c10669814987e1c23c205d
	w.Header().Set("Content-type", "text/html")
	t, err := template.ParseFiles("./../Internal/templates/delresult.html")
	if err != nil {
		log.Println("Template:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to parse files"))
		return
	}

	page := structures.DefaultPage{
		Title: "HTTP-Postgress API",
		Data:  string(data),
	}
	t.Execute(w, &page)
}
func getResultPage(w http.ResponseWriter, r *http.Request, order structures.Order) {
	w.Header().Set("Content-type", "text/html")
	t, err := template.ParseFiles("./../Internal/templates/getresult.html")
	if err != nil {
		log.Println("Template:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to parse files"))
		return
	}

	page := structures.DefaultPage{
		Title: "HTTP-Postgress API",
		Order: order}
	t.Execute(w, &page)
}
func delResultPage(w http.ResponseWriter, r *http.Request, data []byte) {
	w.Header().Set("Content-type", "text/html")
	t, err := template.ParseFiles("./../Internal/templates/delresult.html")
	if err != nil {
		log.Println("Template:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to parse files"))
		return
	}

	page := structures.DefaultPage{
		Title: "HTTP-Postgress API",
		Data:  string(data),
	}
	t.Execute(w, &page)
}

<<<<<<< HEAD
func putOrderDataTest(w http.ResponseWriter, r *http.Request, data []byte) (structures.Order, error) {
	var order structures.Order
	arr := make([]int, 6)
	for j, s := range []string{"Payment.Amount", "Payment.Payment_dt", "Payment.Delivery_cost", "Payment.Goods_total", "Payment.Custom_fee", "Sm_id"} {
		i, err := strconv.Atoi(GetfromHTMLData(&data, []byte(s)))
		//log.Println(s, ":", i)
		if err != nil || i < 0 {
			return order, errors.New(s + " Должна являтся натуральным числом")
		}
		arr[j] = i
=======
func putResultPage(w http.ResponseWriter, r *http.Request, data []byte) {
	w.Header().Set("Content-type", "text/html")
	t, err := template.ParseFiles("./Client/templates/delresult.html")
	if err != nil {
		log.Println("Template:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to parse files"))
		return
	}

	page := structures.DefaultPage{
		Title: "HTTP-Postgress API",
		Data:  string(data),
	}
	t.Execute(w, &page)
}
func getResultPage(w http.ResponseWriter, r *http.Request, order structures.Order) {
	w.Header().Set("Content-type", "text/html")
	t, err := template.ParseFiles("./Client/templates/getresult.html")
	if err != nil {
		log.Println("Template:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to parse files"))
		return
>>>>>>> b8b5f62acf310e5886c10669814987e1c23c205d
	}
	order = structures.Order{
		Order_uid:    GetfromHTMLData(&data, []byte("Order_uid")),
		Track_number: GetfromHTMLData(&data, []byte("Track_number")),
		Entry:        GetfromHTMLData(&data, []byte("Entry")),
		Delivery: structures.Delivery{
			Name:    GetfromHTMLData(&data, []byte("Delivery.Name")),
			Phone:   GetfromHTMLData(&data, []byte("Delivery.Phone")),
			Zip:     GetfromHTMLData(&data, []byte("Delivery.Zip")),
			City:    GetfromHTMLData(&data, []byte("Delivery.City")),
			Address: GetfromHTMLData(&data, []byte("Delivery.Address")),
			Region:  GetfromHTMLData(&data, []byte("Delivery.Region")),
			Email:   GetfromHTMLData(&data, []byte("Delivery.Email")),
		},
		Payment: structures.Payment{
			Transaction:   GetfromHTMLData(&data, []byte("Payment.Transaction")),
			Request_id:    GetfromHTMLData(&data, []byte("Payment.Request_id")),
			Currency:      GetfromHTMLData(&data, []byte("Payment.Currency")),
			Provider:      GetfromHTMLData(&data, []byte("Payment.Provider")),
			Amount:        arr[0],
			Payment_dt:    arr[1],
			Bank:          GetfromHTMLData(&data, []byte("Payment.Bank")),
			Delivery_cost: arr[2],
			Goods_total:   arr[3],
			Custom_fee:    arr[4],
		},
		Locale:             GetfromHTMLData(&data, []byte("Locale")),
		Internal_signature: GetfromHTMLData(&data, []byte("Internal_signature")),
		Customer_id:        GetfromHTMLData(&data, []byte("Customer_id")),
		Delivery_service:   GetfromHTMLData(&data, []byte("Delivery_service")),
		Shardkey:           GetfromHTMLData(&data, []byte("Shardkey")),
		Sm_id:              arr[5],
		Date_created:       GetfromHTMLData(&data, []byte("Date_created")),
		Oof_shard:          GetfromHTMLData(&data, []byte("Oof_shard")),
	}
	// ТЕСТЫ ORDER
	if order.Order_uid == "" {
		return order, errors.New("Order_uid не должно быть пустым")
	}
	if order.Track_number == "" {
		return order, errors.New("Track_number не должно быть пустым")
	}
	if order.Entry == "" {
		return order, errors.New("Entry не должно быть пустым")
	}
	_, err := strconv.Atoi(strings.TrimPrefix(order.Delivery.Phone, "+"))
	if len(order.Delivery.Phone) != 11 || !strings.HasPrefix(order.Delivery.Phone, "+") || err != nil {
		return order, errors.New("Delivery.Phone должен иметь формат +12345667890")
	}
	_, err = time.Parse("2006-01-02T15:04:05Z", order.Date_created)
	if err != nil {
		return order, errors.New("Date_created должно быть установлено в формате yyyy-mm-dd'T'HH:mm:ss'Z'")
	}
	// ITEMS
	items := make([]structures.Items, 0)
	for len(data) != 0 {
		for j, s := range []string{"Items.Chrt_id", "Items.Price", "Items.Sale", "Items.Total_price", "Items.Nm_id", "Items.Status"} {
			i, err := strconv.Atoi(GetfromHTMLData(&data, []byte(s)))
			//log.Println(s, ":", i)
			if err != nil || i < 0 {
				return order, errors.New("Items" + strconv.Itoa(len(items)) + "." + s + " Должна являтся натуральным числом")
			}
			arr[j] = i
		}
		items = append(items, structures.Items{
			Chrt_id:      arr[0],
			Track_number: GetfromHTMLData(&data, []byte("Items.Track_number")),
			Price:        arr[1],
			Rid:          GetfromHTMLData(&data, []byte("Items.Rid")),
			Name:         GetfromHTMLData(&data, []byte("Items.Name")),
			Sale:         arr[2],
			Size:         GetfromHTMLData(&data, []byte("Items.Size")),
			Total_price:  arr[3],
			Nm_id:        arr[4],
			Brand:        GetfromHTMLData(&data, []byte("Items.Brand")),
			Status:       arr[5],
		})

		// ТЕСТЫ ITEMS
	}
	order.Items = items
	return order, nil
}
func GetfromHTMLData(data *[]byte, nameInHTML []byte) string {
	var firstId, lastId, endmarkId, counter int
	firstId = -1
	lastId = -1
	endmarkId = -1
	lenData := len(*data)
	for i := 0; i < lenData; i++ {
		if lastId == -1 {
			if (*data)[i] == nameInHTML[counter] {
				if counter == 0 {
					firstId = i
				}
				counter++
				if counter == len(nameInHTML) {
					lastId = i
				}
			} else {
				firstId = -1
				counter = 0
			}
		} else {
			// Перевод
			if (*data)[i] == '%' {
				if (*data)[i+1] == '2' && (*data)[i+2] == 'B' {
					(*data)[i] = '+'
					(*data) = append((*data)[:i+1], (*data)[i+3:]...)
					lenData = len(*data)
				}
				if (*data)[i+1] == '3' && (*data)[i+2] == 'A' {
					(*data)[i] = ':'
					(*data) = append((*data)[:i+1], (*data)[i+3:]...)
					lenData = len(*data)
				}
			}
			// Завершение
			if (*data)[i] == '&' {
				endmarkId = i
				break
			}
		}
	}
	if firstId == -1 || lastId == -1 {
		return ""
	}
	var str string
	if endmarkId == -1 {
		endmarkId = len(*data)
		str = string((*data)[lastId+2 : endmarkId])
		(*data) = (*data)[:firstId]
	} else {
		str = string((*data)[lastId+2 : endmarkId])
		(*data) = append((*data)[:firstId], (*data)[endmarkId+1:]...)
	}
	return str
}
