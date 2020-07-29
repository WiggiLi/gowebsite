package api

import (
	"gowebsite/app"
	"github.com/valyala/fasthttp" 
	jwt "github.com/dgrijalva/jwt-go"
	"os"
	//"fmt"
	"encoding/json"
)

func JwtAuthentication(next fasthttp.RequestHandler) fasthttp.RequestHandler {
    return func(ctx *fasthttp.RequestCtx) {

		notAuth := []string{"/new", "/login", "/static/index.html"} // list of endpoints that do not require authorization  
		requestPath := string(ctx.Path()) //current path of request
		tokenHeader :=  string(ctx.Request.Header.Cookie("token"))
		//fmt.Println("res "+tokenHeader)

		

		response := make(map[string] interface{})
		
		
		if tokenHeader == "" { 
			for _, value := range notAuth {

				if value == requestPath {
					next(ctx) 
					return
				}
			}

			response["status"] = false
			response["message"] = "Missing auth token"
			ctx.Response.SetStatusCode(fasthttp.StatusForbidden) //  403 http-code Unauthorized
			ctx.SetContentType("application/json")
			json.NewEncoder(ctx).Encode(response)
			return
		} else {
			/*for _, value := range notAuth {

				if value == requestPath {
					ctx.Redirect("card/inner.html", 302)
					return
				}
			}
			*/
		}

		tk := &app.Token{}            

		token, err := jwt.ParseWithClaims(tokenHeader, tk, func(token *jwt.Token) (interface{}, error) {
						return []byte(os.Getenv("token_password")), nil
					})

		if err != nil { 
			response["status"] = false
			response["message"] = "Malformed authentication token"
			ctx.Response.SetStatusCode(fasthttp.StatusForbidden) //403 http-code
			ctx.SetContentType("application/json")
			json.NewEncoder(ctx).Encode(response)
			return
		}

		if !token.Valid { 
			response["status"] = false
			response["message"] = "Token is not valid"
			ctx.Response.SetStatusCode(fasthttp.StatusForbidden)
			ctx.SetContentType("application/json")
			json.NewEncoder(ctx).Encode(response)
			return
		}

		//fmt.Printf("User %d", tk.UserId) 
		next(ctx) 
	};
}
