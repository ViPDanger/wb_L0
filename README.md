# Тестовый сервис и приложение клиента на REST API для выполнения задания L0 Горутиновый Голэнг Wildberries

Основная суть задания - создать цепочу запросов Клиент -> NATSJetstream -> Сервер -> PostgreSQL -> Сервер -> NATSJetstream -> Клиент

Конфиги с основными настройками NATs jetstream и прослушеваемых портов лежат по пути
Server/cmd/config.cfg
Client/cmd/config.cfg

Основные запросы:
  GET localhost:8080/  запрос со стартовой html страницей с формой заполнения заказа
  
Go version 1.23.0
