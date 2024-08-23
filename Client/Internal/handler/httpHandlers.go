package handler

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	n "github.com/ViPDanger/L0/Client/Internal/nats"
	"github.com/ViPDanger/L0/Client/Internal/structures"
	"github.com/nats-io/nats.go"
)

func SetupHttpHandlers(mux *http.ServeMux, cancel_on_http context.CancelFunc, natsHandler NatsHanlder) {
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
				ErrorPage(w, r,[]byte("Тело не может быть пустым"))
				return
			}
			body, _ = json.Marshal(&structures.Order_uid{
				Order_uid: string(body),
			})
			n.PublishMsg(*natsHandler.JetstreamContext, "DelOrder."+natsHandler.ClientName, body)
			var msg *nats.Msg
			msg = <-natsHandler.DelChannel
			DelResultPage(w, r, msg.Data)

		case "Get":

			body, err := io.ReadAll(r.Body)
			body = []byte(strings.TrimPrefix(string(body), "request_uid="))
			if err != nil || len(body) < 1 {
				ErrorPage(w, r,[]byte("Тело не может быть пустым"))
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
			err = json.Unmarshal(msg.Data, order)
			if err != nil {
				ErrorPage(w, r,[]byte("Нет Order с "+string(body)))
			} else {
				GetResultPage(w, r, order)
			}
		}
	})
}

type ShutdownHandler struct {
	cancel_on_http context.CancelFunc
}

func (sht ShutdownHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Shutdown by http")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Server is going got Killed. Murderer."))
	sht.cancel_on_http()
}

func DefaultPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	t, err := template.ParseFiles("./Client/templates/default.html")
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
				ErrorPage(w, r,[]byte("Число Items должно быть натуральным"))

			}
		}
	}
	items := make([]structures.Items, int)
	page := structures.DefaultPage{
		Title:          "HTTP-Postgress API",
		GetOrderButton: "GetOrderByOrder_uid",
		PutOrderButton: "PutNewOrder",
		DelOrderButton: "DeleteOrderByOrder_uid",
		Order: structures.Order{
			Items: items,
		},
	}
	t.Execute(w, &page)
}
func ErrorPage(w http.ResponseWriter, r *http.Request,data []byte) {
	w.Header().Set("Content-type", "text/html")
	t, err := template.ParseFiles("./Client/templates/error.html")
	if err != nil {
		log.Println("Template:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to parse files"))
		return
	}
	page := structures.ErrorPage{
		Title:        "Error",
		ErrorMessage: string(data),
	}
	t.Execute(w, &page)
}

func PutOrder(w http.ResponseWriter, r *http.Request) {

}
func DelResultPage(w http.ResponseWriter, r *http.Request, data []byte) {
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

func GetResultPage(w http.ResponseWriter, r *http.Request, order structures.Order) {
	w.Header().Set("Content-type", "text/html")
	t, err := template.ParseFiles("./Client/templates/getresult.html")
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
