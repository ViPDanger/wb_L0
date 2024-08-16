package main

import (
	"context"

	c "github.com/ViPDanger/L0/Server/Internal/config"
	pg "github.com/ViPDanger/L0/Server/Internal/postgres"
)

func main() {

	config := c.Read_Config()
	client, _ := pg.NewClient(context.Background(), config)
	rep := pg.NewRepository(client)
	if rep != nil {
	}
}
