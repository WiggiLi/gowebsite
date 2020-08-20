package dal

import (
	"context"
	"database/sql"

	"github.com/WiggiLi/gowebsite/app"
	_ "github.com/lib/pq"

	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func (t *PSQL) CreateAccDB(account *app.Account) (bool, *app.Account, error) {
	//Validate incoming user details...
	if !strings.Contains(account.Email, "@") {
		log.Print("Miss @")
		return false, nil, nil
	}

	if len(account.Password) < 6 {
		log.Print("Password is required and contain 6 simbols")
		return false, nil, nil
	}

	//check if email already exist

	ctx := context.Background()
	var err error

	psql := fmt.Sprintf("SELECT email FROM accounts WHERE email='%s' ORDER BY id LIMIT 1", account.Email)

	stmt, err := t.DataBase.Prepare(psql)
	if err != nil {
		log.Println("Prepare error: " + err.Error())
		return false, nil, nil
	}
	defer stmt.Close()

	var ok string

	err = t.DataBase.QueryRowContext(ctx, psql).Scan(&ok)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal("Error reading rows: " + err.Error())
		return false, nil, nil
	}

	if ok != "" {
		log.Println("Account email already exist.")
		return false, nil, nil
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	psql2 := fmt.Sprintf("INSERT INTO accounts (email, password) VALUES ('%s','%s') RETURNING id;", account.Email, account.Password)
	/*
		stmt2, err := t.DataBase.Prepare(psql2)
		if err != nil {
			log.Println("Prepare error" + err.Error())
			return false, nil, nil
		}
		defer stmt2.Close()
	*/
	lastInsertId := 0
	err = t.DataBase.QueryRow(psql2).Scan(&lastInsertId)
	//log.Printf("Acoount Id:  %d\n", lastInsertId)

	if err != nil {
		log.Println("Failed to create account, connection error." + err.Error())
		return false, nil, nil
	}

	//CreareComment new JWT token for the newly registered account
	tk := &app.Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString

	account.Password = "" //delete password

	log.Println("Account has been created")

	return true, account, nil
}

func (t *PSQL) LoginAccDB(email, password string) (map[string]interface{}, error) {
	ctx := context.Background()
	var err error

	if t.DataBase == nil {
		err = errors.New("DB is null")
		fmt.Println("Null DB")
		return map[string]interface{}{"status": false, "account": nil}, nil
	}

	// Check if database is alive.
	err = t.DataBase.PingContext(ctx)
	if err != nil {
		fmt.Println("Error pinging database: " + err.Error())
	}

	tsql := fmt.Sprintf("SELECT email, password FROM accounts WHERE email='%s' ORDER BY id LIMIT 1", email)

	stmt, err := t.DataBase.Prepare(tsql)
	if err != nil {
		log.Println("Prepare error: " + err.Error())
		return map[string]interface{}{"status": false, "account": nil}, nil
	}
	defer stmt.Close()

	rows, err := t.DataBase.QueryContext(ctx, tsql)
	if err != nil {
		log.Fatal("Error reading rows: " + err.Error())
		return map[string]interface{}{"status": false, "account": nil}, nil
	}

	defer rows.Close()

	account := &app.Account{}

	for rows.Next() {
		err := rows.Scan(&account.Email, &account.Password)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("Email address not found" + err.Error())
				return map[string]interface{}{"status": false, "account": nil}, nil
			}
			log.Fatal("Error reading rows: " + err.Error())
			return map[string]interface{}{"status": false, "account": nil}, nil
		}
		//fmt.Printf("RESULT : %s\n" + account.Email)
	}

	if account.Email == "" {
		fmt.Println("Email address not found2")
		return map[string]interface{}{"status": false, "account": nil}, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		fmt.Println("Invalid login credentials. Please try again")

		return map[string]interface{}{"status": false, "account": nil}, nil
	}

	account.Password = ""

	//CreareComment JWT token
	tk := &app.Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString //Store the token in the response

	log.Println("Logged In")

	return map[string]interface{}{"status": true, "account": account}, nil
}
