package app

import (
	"log"
)

// IncomeRegistration is an interface for accepting income requesrs for neccassery operations  from Web Server
type IncomeRegistration interface {
	RegisterComment(*Comment)
	GiveContent(int) *ContentPage
	GiveComments(int) *AllComments
	GiveTitles() *AllTitles
	CreateAcc(*Account) (bool, *Account)
	LoginAcc(email, password string) map[string]interface{}
}

// DataAccessLayer is an interface for DAL usage from Application
type DataAccessLayer interface {
	CreareComment(*Comment) error
	ReadContent(int) (*ContentPage, error)
	ReadComments(int) (*AllComments, error)
	ReadTitles() (*AllTitles, error)
	CreateAccDB(*Account) (bool, *Account, error)
	LoginAccDB(email, password string) (map[string]interface{}, error)
}

// Application is responsible for all logics and communicates with other layers
type Application struct {
	DB   DataAccessLayer
	errc chan<- error
}

// RegisterComment sends comment to DAL for its saving
func (app *Application) RegisterComment(currentData *Comment) {
	err := app.DB.CreareComment(currentData)

	if err != nil {
		app.errc <- err
		return
	}

	log.Print("New comment added to PostgreSQL...")
}

// GiveContent sends content for page with currentID DAL
func (app *Application) GiveContent(currentID int) *ContentPage {
	allEv := NewContentPage()
	allEv, err := app.DB.ReadContent(currentID)

	if err != nil {
		app.errc <- err
		return nil
	}

	log.Printf("Content readed for %d page", currentID)
	return allEv
}

// GiveComments sends all comments from DAL
func (app *Application) GiveComments(currentPage int) *AllComments {
	allEv := GetComments()
	allEv, err := app.DB.ReadComments(currentPage)

	if err != nil {
		app.errc <- err
		return nil
	}

	log.Printf("Comments readed for %d page", currentPage)
	return allEv
}

func (app *Application) GiveTitles() *AllTitles {
	allEv := GetTitles()
	allEv, err := app.DB.ReadTitles()

	if err != nil {
		app.errc <- err
		return nil
	}

	log.Print("Titles readed for page")
	return allEv
}

func (app *Application) CreateAcc(acc *Account) (bool, *Account) {
	flag, acc, err := app.DB.CreateAccDB(acc)

	if err != nil {
		app.errc <- err
		return false, nil
	}

	log.Print("Created new user...")
	return flag, acc
}

func (app *Application) LoginAcc(email, password string) map[string]interface{} {
	flag, err := app.DB.LoginAccDB(email, password)

	if err != nil {
		app.errc <- err
		return map[string]interface{}{"status": false, "account": nil}
	}

	log.Print("Checked user...")
	return flag
}

// NewApplication constructs Application
func NewApplication(db DataAccessLayer, errchannel chan<- error) *Application {
	res := &Application{}
	res.DB = db
	res.errc = errchannel

	return res
}
