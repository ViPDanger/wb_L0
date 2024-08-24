SET search_path TO wb_schema;
DROP TABLE IF EXISTS order_item;
DROP TABLE IF EXISTS delivery;
DROP TABLE IF EXISTS payment;
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS orders;

CREATE TABLE orders
(
order_uid varchar(50) PRIMARY KEY NOT NULL,
track_number varchar(50) NOT NULL,
entry varchar(50) NOT NULL,
locale varchar(10) NOT NULL,
internal_signature varchar(50) NOT NULL,
customer_id varchar(50) NOT NULL,
delivery_service varchar(50) NOT NULL,
shardkey varchar(10) NOT NULL,
sm_id int NOT NULL,
date_created timestamp NOT NULL,
oof_shard varchar(10) NOT NULL
);
CREATE TABLE delivery
(
order_uid varchar(50) NOT NULL,
Name 	varchar(50)	NOT NULL,
Phone 	varchar(50)	NOT NULL,
Zip     varchar(50)	NOT NULL,
City    varchar(50)	NOT NULL,
Address	varchar(50)	NOT NULL,
Region  varchar(50)	NOT NULL,
Email   varchar(50)	NOT NULL,
CONSTRAINT fk_order
	FOREIGN KEY(order_uid)
	REFERENCES orders(order_uid) 
	ON UPDATE CASCADE ON DELETE CASCADE
);
CREATE TABLE payment
(
order_uid varchar(50)  NOT NULL,
transaction varchar(50) NOT NULL,
request_id varchar(50),
currency varchar(10) NOT NULL,
provider varchar(50) NOT NULL,
amount int8 NOT NULL,
payment_dt int8 NOT NULL,
bank varchar(50) NOT NULL,
delivery_cost int8 NOT NULL,
goods_total int8 NOT NULL,
custom_fee int8 NOT NULL,
	CONSTRAINT fk_order
	FOREIGN KEY(order_uid)
	REFERENCES orders(order_uid) 
	ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE items
(
chrt_id int8 PRIMARY KEY NOT NULL,
track_number varchar(50) NOT NULL,
price int8 NOT NULL,
rid varchar(50) NOT NULL,
name varchar(50) NOT NULL,
sale int8 NOT NULL,
size varchar(50) NOT NULL,
total_price int8 NOT NULL,
nm_id int8 NOT NULL,
brand varchar(50) NOT NULL,
status int8 NOT NULL
);

CREATE TABLE order_item (
  order_uid  varchar(50) REFERENCES orders (order_uid) ON UPDATE CASCADE ON DELETE CASCADE,
  chrt_id int8 REFERENCES items (chrt_id) ON UPDATE CASCADE ON DELETE CASCADE
);