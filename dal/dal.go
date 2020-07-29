package dal

import (
	"gowebsite/app"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
	"context"
	
	"errors"
	"fmt"
	"log"
	"strconv"
	"os"
)

// MsSQL represents data for connection to Data base
type MsSQL struct {
	Host     string
	DataBase *sql.DB
}


// NewMsSQL constructs object of MsSQL
func NewMsSQL(host string, port int) (*MsSQL, error) {

	e := godotenv.Load() //Загрузить файл .env
	if e != nil {
		fmt.Print(e)
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")

	//connString := "user=postgres password=mypass dbname=web_pages sslmode=disable"
	connString := fmt.Sprintf("user=%s dbname=%s sslmode=disable password=%s", username, dbName, password) 
	//fmt.Println(connString)
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
	fmt.Printf("Connected to db!\n")

	res := &MsSQL{
		Host:     host,
		DataBase: db}

	return res, nil
}



// Create inserts new Event into DB
func (t *MsSQL) Create(current *app.Event) error {
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

	tsql := fmt.Sprintf("INSERT INTO comments (page, name, content) VALUES ('%s','%s','%s');", current.Page, current.Title, current.Description)

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

func (t *MsSQL) Read1(currentID int) (*app.ContentPage, error) {
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

func (t *MsSQL) Read2(current *app.Event) (*app.AllEvents, error) {
	events := app.GetEvents()

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
	i1, err := strconv.Atoi(current.Pag)
	if err != nil {
		fmt.Printf("pagERROR i=%d, type: %T\n", current.Pag, current.Pag)
	}

	tsql := fmt.Sprintf("SELECT page, name, content FROM comments where page=%d", i1)

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
		var id, page, title, description string

		// Get values from row.
		err := rows.Scan(&page, &title, &description)
		if err != nil {
			log.Fatal("Error reading rows: " + err.Error())
			return nil, err
		}

		ev := *app.NewEvent()
		ev.ID = id
		ev.Page = page
		ev.Title = title
		ev.Description = description
		*events = append(*events, ev)

		fmt.Printf("ID: %d, Name: %s, Description: %s\n", id, title, description)
		count++
	}

	return events, nil
}

func (t *MsSQL) Read3() (*app.AllTitles, error) {
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
