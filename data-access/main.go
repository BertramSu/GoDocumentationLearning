package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

type Album struct {
	ID     int64   `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float32 `json:"price"`
}

type Config struct {
	Username     string
	Password     string
	Address      string
	DatabaseName string
}

var cfg = Config{
	Username:     "admin",
	Password:     "passwordHere",
	Address:      "localhost:5432",
	DatabaseName: "album_db",
}

var connStr = fmt.Sprintf("postgres://%s:%s@%s/%s", cfg.Username, cfg.Password, cfg.Address, cfg.DatabaseName)

func main() {
	albums, err := getAllAlbums()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Albums found: %v", albums)

	album, err := getAlbumById(2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album found: %v", album)

	/*
		newAlbum := Album{
			Title:  "Pretending",
			Artist: "Fletcher",
			Price:  0.99,
		}
		insertInt, err := insertAlbum(newAlbum)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Album inserted with id: %v", insertInt)

		albumss, err := getAllAlbums()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Albums found: %v", albumss)
	*/
}

func getAllAlbums() ([]Album, error) {
	// Connect to the database
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	// Query example
	query := "SELECT id, title, artist, price FROM album"
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		log.Fatalf("Query failed getAllAlbums: %v", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[Album])
}

func getAlbumById(id int) (Album, error) {
	// Connect to the database
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	query := "SELECT id, title, artist, price FROM album WHERE id = @id"

	args := pgx.NamedArgs{
		"id": id,
	}
	// Query example
	rows, err := conn.Query(context.Background(), query, args)
	if err != nil {
		log.Fatalf("Query failed getAlbumById): %v", err)
	}
	defer rows.Close()

	return pgx.CollectOneRow(rows, pgx.RowToStructByName[Album])
}

func insertAlbum(album Album) (int, error) {
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
		return 0, err
	}
	defer conn.Close(context.Background())

	query := "INSERT INTO album (title, artist, price) VALUES (@title, @artist, @price) RETURNING id"

	args := pgx.NamedArgs{
		"title":  album.Title,
		"artist": album.Artist,
		"price":  album.Price,
	}

	var id int
	qErr := conn.QueryRow(context.Background(), query, args).Scan(&id)
	if qErr != nil {
		log.Fatalf("Query failed insertAlbum): %v", err)
		return 0, qErr
	}

	return id, qErr
}

func deleteAlbum(id int) error {
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
		return err
	}
	defer conn.Close(context.Background())

	query := "DELETE FROM album WHERE id = $1"
	commandTag, err := conn.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return fmt.Errorf("No album found with ID %d", id)
	}

	return nil
}

func updateAlbum(id int, album Album) (int, error) {
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
		return 0, err
	}
	defer conn.Close(context.Background())

	query := `
        UPDATE album
        SET title = @title, artist = @artist, price = @price
        WHERE id = @id
    `

	args := pgx.NamedArgs{
		"title":  album.Title,
		"artist": album.Artist,
		"price":  album.Price,
		"id":     id,
	}

	var rowId int
	qErr := conn.QueryRow(context.Background(), query, args).Scan(&rowId)
	if qErr != nil {
		log.Fatalf("Query failed insertAlbum): %v", err)
		return 0, qErr
	}

	return rowId, qErr
}
