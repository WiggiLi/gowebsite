package app

import (
	"log"
	//"fmt"
)

// IncomeRegistration is an interface for accepting income requesrs for neccassery operations  from Web Server
type IncomeRegistration interface {
	RegisterEvent(*Event)
	GiveContent(int) *ContentPage
	GiveComments(*Event) *AllEvents
	GiveTitles() *AllTitles
	CreateAcc(*Account) map[string]interface{} 
	LoginAcc(email, password string) map[string]interface{} 
}

// DataAccessLayer is an interface for DAL usage from Application
type DataAccessLayer interface {
	Create(*Event) error
	Read1(int) (*ContentPage, error) //content
	Read2(*Event) (*AllEvents, error) //comments
	Read3() (*AllTitles, error) //titles
	CreateAccDB(*Account) (map[string]interface{} , error)
	LoginAccDB(email, password string) (map[string]interface{} , error)
}

// Application is responsible for all logics and communicates with other layers
type Application struct {
	DB   DataAccessLayer
	errc chan<- error
}

// RegisterEvent sends Event to DAL for saving/registration
func (app *Application) RegisterEvent(currentData *Event) {
	err := app.DB.Create(currentData)

	if err != nil {
		app.errc <- err
		return
	}

	log.Print("New event added to MS SQL server...")
}

func (app *Application) GiveContent(currentID int) *ContentPage {
	allEv := NewContentPage()
	allEv, err := app.DB.Read1(currentID)

	if err != nil {
		app.errc <- err
		return nil
	}

	log.Print("Events readed from MS SQL server...")
	return allEv
}

// RegisterEvent sends Event to DAL for saving/registration
func (app *Application) GiveComments(currentData *Event) *AllEvents {
	allEv := GetEvents()
	allEv, err := app.DB.Read2(currentData)

	if err != nil {
		app.errc <- err
		return nil
	}

	//log.Print("Events readed from MS SQL server...")
	return allEv
}

func (app *Application) GiveTitles() *AllTitles {
	allEv := GetTitles()
	allEv, err := app.DB.Read3()

	if err != nil {
		app.errc <- err
		return nil
	}

	//log.Print("Events readed from MS SQL server...")
	return allEv
}

func (app *Application)  CreateAcc(acc *Account) map[string]interface{} {
	flag, err := app.DB.CreateAccDB(acc)

	if err != nil {
		app.errc <- err
		return map[string]interface{} {"status" : false, "account" : nil}
	}

	log.Print("Created new user...")
	return flag
}

func (app *Application)  LoginAcc(email, password string) map[string]interface{} {
	flag, err := app.DB.LoginAccDB(email, password)

	if err != nil {
		app.errc <- err
		return map[string]interface{} {"status" : false, "account" : nil}
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
