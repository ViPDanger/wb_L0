package postgres

import (
	"context"
	"log"

	"github.com/ViPDanger/L0/Server/Internal/structures"
)

// Описание репозитория клиента
type Repository struct {
	client Client
}

func NewRepository(client Client) *Repository {
	return &Repository{
		client: client,
	}
}

// Функции репозитория
func (r *Repository) Insert_Author(author structures.AuthorDTO) (int, error) {
	var i int
	ctx := context.Background()
	tx, err := r.client.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)
	query := "insert into author (name) VALUES ($1)"
	_, err = tx.Exec(ctx, query, author.Name)

	if err != nil {
		return 0, err
	}
	query = "SELECT author.id FROM author WHERE author.name = $1"
	err = tx.QueryRow(context.Background(), query, author.Name).Scan(&i)
	if err != nil {
		return 0, err
	}
	tx.Commit(ctx)
	return i, err
}

func (r *Repository) All_Authors() ([]structures.Author, error) {
	query := "SELECT * FROM author"
	authors := make([]structures.Author, 0)
	rows, err := r.client.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for i := 0; rows.Next(); i++ {
		authors = append(authors, structures.Author{})
		err = rows.Scan(&authors[i].ID, &authors[i].Name)
		if err != nil {
			return nil, err
		}
	}

	return authors, err
}

func (r *Repository) Find_Author(id string) (structures.Author, error) {
	var author structures.Author
	query := "SELECT * FROM author WHERE id = $1"
	err := r.client.QueryRow(context.Background(), query, id).Scan(&author.ID, &author.Name)

	return author, err
}

func (r *Repository) Delete_Author(id string) error {
	ctx := context.Background()
	tx, err := r.client.Begin(ctx)
	if err != nil {
		log.Fatalln("")
		return err
	}
	defer tx.Rollback(ctx)
	query := "DELETE FROM author WHERE author.id = $1"
	_, err = tx.Exec(ctx, query, id)

	if err != nil {
		return err
	}
	tx.Commit(ctx)
	return err
}
