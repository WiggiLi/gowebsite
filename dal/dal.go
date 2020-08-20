package dal

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/WiggiLi/gowebsite/app"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// PostgreSQL represents data for connection to Data base
type PSQL struct {
	Host     string
	DataBase *sql.DB
}

// NewPSQL constructs object of PostgreSQL
func NewPSQL(host string, port int) (*PSQL, error) {

	e := godotenv.Load() //load file .env with data for db connecting string
	if e != nil {
		log.Print("Not exist .env with data for db connecting string.", e)
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	//dbHost := os.Getenv("db_host")

	//connString := "user=postgres password=mypass dbname=web_pages sslmode=disable"
	connString := fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable", host, username, dbName, password)
	var err error
	var db *sql.DB

	db, err = sql.Open("postgres", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}

	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Print("Connected to db!\n")

	res := &PSQL{
		Host:     host,
		DataBase: db}

	return res, nil
}

// CreareComment inserts new Commentinto DB
func (t *PSQL) CreareComment(current *app.Comment) error {
	ctx := context.Background()
	var err error
	if t.DataBase == nil {
		err = errors.New("DB is null")
		log.Println("Null DB")
		return err
	}

	// Check if database is alive.
	err = t.DataBase.PingContext(ctx)
	if err != nil {
		log.Fatal("Error pinging database: " + err.Error())
	}
	log.Printf("('%s','%s','%s');", current.Page, current.Name, current.Content)
	tsql := fmt.Sprintf("INSERT INTO comments (page, name, content) VALUES ('%s','%s','%s');", current.Page, current.Name, current.Content)

	stmt, err := t.DataBase.Prepare(tsql)
	if err != nil {
		log.Println("Prepare error")
		return err
	}
	defer stmt.Close()
	fmt.Printf(current.Page + "\n")

	row, err := t.DataBase.ExecContext(ctx, tsql)
	_ = row
	//fmt.Printf("Inserted successfully222.\n")
	if err != nil {
		log.Fatal("Error inserting new row: " + err.Error())
		return err
	}
	//fmt.Printf("Inserted successfully222.\n")
	return nil
}

func (t *PSQL) ReadContent(currentID int) (*app.ContentPage, error) {
	events := app.NewContentPage()

	ctx := context.Background()
	var err error

	if t.DataBase == nil {
		err = errors.New("DB is null")
		log.Println("Null DB")
		return nil, err
	}

	// Check if database is alive.
	err = t.DataBase.PingContext(ctx)
	if err != nil {
		log.Fatal("Error pinging database: " + err.Error())
	}

	tsql := fmt.Sprintf("SELECT page, title, content FROM content_page where page=%d", currentID)

	stmt, err := t.DataBase.Prepare(tsql)
	if err != nil {
		log.Println("Prepare error")
		return nil, err
	}
	defer stmt.Close()

	rows, err := t.DataBase.QueryContext(ctx, tsql)
	if err != nil {
		log.Fatal("Error reading rows: " + err.Error())
		return nil, err
	}

	defer rows.Close()

	var count int = 0

	for rows.Next() {
		var page, title, description string

		// Get values from row.
		err := rows.Scan(&page, &title, &description)
		if err != nil {
			log.Fatal("Error reading rows: " + err.Error())
			return nil, err
		}

		events.Page = page
		events.Title = title
		events.Content = description

		//fmt.Printf("ID: %d, Name: %s, Description: %s\n", page, title, description)
		count++
	}

	return events, nil
}

func (t *PSQL) ReadComments(currentPage int) (*app.AllComments, error) {
	events := app.GetComments()

	ctx := context.Background()
	var err error

	if t.DataBase == nil {
		err = errors.New("DB is null")
		log.Println("Null DB")
		return nil, err
	}

	// Check if database is alive.
	err = t.DataBase.PingContext(ctx)
	if err != nil {
		log.Fatal("Error pinging database: " + err.Error())
	}
	//fmt.Printf("pag2 before i=%d, type: %T\n", current.Pag, current.Pag)
	//i1, err := strconv.Atoi(currentPage)
	//if err != nil {
	//	log.Printf(err)
	//}

	tsql := fmt.Sprintf("SELECT page, name, content FROM comments where page=%d", currentPage)

	stmt, err := t.DataBase.Prepare(tsql)
	if err != nil {
		log.Println("Prepare error")
		return nil, err
	}
	defer stmt.Close()

	rows, err := t.DataBase.QueryContext(ctx, tsql)
	if err != nil {
		log.Fatal("Error reading rows: " + err.Error())
		return nil, err
	}

	defer rows.Close()

	var count int = 0

	for rows.Next() {
		var page, title, description string

		// Get values from row.
		err := rows.Scan(&page, &title, &description)
		if err != nil {
			log.Fatal("Error reading rows: " + err.Error())
			return nil, err
		}

		ev := *app.NewComment()
		ev.Page = page
		ev.Name = title
		ev.Content = description
		*events = append(*events, ev)

		//fmt.Printf("ID: %d, Name: %s, Description: %s\n", id, title, description)
		count++
	}

	return events, nil
}

func (t *PSQL) ReadTitles() (*app.AllTitles, error) {
	events := app.GetTitles()

	ctx := context.Background()
	var err error

	if t.DataBase == nil {
		err = errors.New("DB is null")
		log.Println("Null DB")
		return nil, err
	}

	// Check if database is alive.
	err = t.DataBase.PingContext(ctx)
	if err != nil {
		log.Fatal("Error pinging database: " + err.Error())
	}

	tsql := fmt.Sprintf("SELECT title FROM content_page")

	stmt, err := t.DataBase.Prepare(tsql)
	if err != nil {
		log.Println("Prepare error")
		return nil, err
	}
	defer stmt.Close()

	rows, err := t.DataBase.QueryContext(ctx, tsql)
	if err != nil {
		log.Fatal("Error reading rows: " + err.Error())
		return nil, err
	}

	defer rows.Close()

	var count int = 0

	for rows.Next() {
		var title string

		// Get values from row.
		err := rows.Scan(&title)
		if err != nil {
			log.Fatal("Error reading rows: " + err.Error())
			return nil, err
		}

		ev := *app.NewTitle()
		ev.Title = title
		*events = append(*events, ev)

		//fmt.Printf("Name: %s\n", title)
		count++
	}

	return events, nil
}
