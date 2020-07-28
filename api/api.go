package api

import (
	"gowebsite/app"
	"encoding/json"
	//"fmt"
	"log"
	"time"
	"flag"
	"strconv"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

// WebServer accepts POST requests with payload of XML docs of Receipts
// Then it parses them with XPath and pushes data to Application
type WebServer struct {
	application app.IncomeRegistration
}

// ParseJSON JSON data and pushes it to Application
func (server *WebServer) createEvent(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "POST")

	//log.Print("API got new event. Parsing started...")

	newEvent := app.NewEvent()

	json.Unmarshal(ctx.PostBody(), &newEvent)

	server.application.RegisterEvent(newEvent)
	time.Sleep(2 * time.Second)

	//fmt.Printf("Inserted successfully.\n")
	ctx.Response.SetStatusCode(fasthttp.StatusCreated)

	json.NewEncoder(ctx).Encode(newEvent)
}

func (server *WebServer) getComms(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "POST")

	payLoad := string(ctx.PostBody())
	log.Print(payLoad)

	newEvent := app.NewEvent()

	json.Unmarshal(ctx.PostBody(), &newEvent)

	events := app.GetEvents()
	events = server.application.GiveComments(newEvent)

	json.NewEncoder(ctx).Encode(events)
	ctx.Response.SetStatusCode(fasthttp.StatusOK)
}

func (server *WebServer) getContentByID(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET")

	contentID := ctx.UserValue("id")

	i1, err := strconv.Atoi(contentID.(string))
	if err != nil {
		//fmt.Printf("page ERROR. Type: %T\n", contentID)
	}

	events := app.NewContentPage()
	events = server.application.GiveContent(i1)

	json.NewEncoder(ctx).Encode(events)
	ctx.Response.SetStatusCode(fasthttp.StatusOK)
}

func (server *WebServer) getTitles(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET")

	events := app.GetTitles()
	events = server.application.GiveTitles()

	json.NewEncoder(ctx).Encode(events)
	ctx.Response.SetStatusCode(fasthttp.StatusOK)
}

func (server *WebServer) createAccount(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "POST")

	account := app.NewAccount()
	account.Email = string(ctx.FormValue("email_register")) // Data from the form
    account.Password = string(ctx.FormValue("password_register")) // Data from the form

	ans := server.application.CreateAcc(account)

	if ans["status"]==true {
		acc := app.NewAccount()
		acc = ans["account"].(*app.Account)
		//fmt.Println("ans status true " + acc.Email)
		SetCookie("name", acc.Email, ctx)
		SetCookie("token", acc.Token, ctx)
		ctx.Redirect("static/inner.html", 302)
	}  else {
		ctx.Redirect("static/index.html", 403)
	}
}


func (server *WebServer) authenticate(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "POST")
	
	account := app.NewAccount()

	account.Email = string(ctx.FormValue("email_login")) // Data from the form
    account.Password = string(ctx.FormValue("password_login")) // Data from the form

	ans := server.application.LoginAcc(account.Email, account.Password)

	if ans["status"]==true {
		acc := app.NewAccount()
		acc = ans["account"].(*app.Account)
		//fmt.Println("ans status true " + acc.Email)
		SetCookie("name", acc.Email, ctx)
		SetCookie("token", acc.Token, ctx)
		ctx.Redirect("static/inner.html", 302)
	} else {
		ctx.Redirect("static/index.html", 403)
	}
}

func (server *WebServer) logout(ctx *fasthttp.RequestCtx) {
	ClearCookie("name", ctx)
	ClearCookie("token", ctx)
	ctx.Redirect("static/index.html", 302)
}

func SetCookie(name, value string, ctx *fasthttp.RequestCtx) {
	var c fasthttp.Cookie
	c.SetKey(name)
	c.SetValue(value)
	c.SetPath("/")
	c.SetMaxAge(3600000)
	ctx.Response.Header.SetCookie(&c)
}

func ClearCookie(name string, ctx *fasthttp.RequestCtx) {  
	var c fasthttp.Cookie
	c.SetKey(name)
	c.SetValue("")
	c.SetPath("/")
	c.SetMaxAge(-3600000)
    ctx.Response.Header.SetCookie(&c)
}

// Start initializes Web Server, starts application and begins serving
func (server *WebServer) Start(errc chan<- error) {
	router := fasthttprouter.New()
	flag.Parse()
	hub := newHub()
	go hub.run()
		
	port := ":4200"
	router.POST("/comm", server.createEvent) // create comment
	router.POST("/comms", server.getComms) //get all comments
	router.GET("/content/:id", server.getContentByID) // get content for page
	router.GET("/titles", server.getTitles) // get titles of all articles
	router.GET("/socket", func(ctx *fasthttp.RequestCtx) {
							serveWs(ctx, hub)
						  })///imlement websocket for comments
	
	router.POST("/new", server.createAccount)
	router.POST("/login", server.authenticate)
	router.POST("/logout", server.logout)
	
	router.GET("/", func(ctx *fasthttp.RequestCtx) {
		log.Println("go to index")
		ctx.Redirect("static/index.html", 200)
	  })

	router.ServeFiles("/static/*filepath", "./static")

	log.Print("Server is starting on port ", port)
	errc <- fasthttp.ListenAndServe(port, JwtAuthentication(router.Handler))
}

// NewWebServer constructs Web Server
func NewWebServer(application app.IncomeRegistration) *WebServer {
	res := &WebServer{}
	res.application = application

	return res
}
