package main

import (
	"context"
	"fmt"
	"github.com/CookieNyanCloud/driveApi/api"
	"github.com/CookieNyanCloud/driveApi/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const configsDir = "configs"

type input struct {
	Name string `json:"name"`
}

func main() {
	ctx := context.Background()
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	// If modifying these scopes, delete your previously saved token.json.
	configD, err := google.ConfigFromJSON(b,
		drive.DriveReadonlyScope,
	)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := api.GetClient(configD)
	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))


	server:= gin.Default()
	server.GET("/getphoto", func(c *gin.Context) {
		var inp input
		if err := c.ShouldBindJSON(&inp); err != nil {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		names ,err:=service.GetPhoto(srv,inp.Name)
		if err!= nil {
			newResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		fmt.Println(len(names))
		for i, name:= range names{
			fmt.Println(i,":",name)
		}
		if len(names)==0 {
			newResponse(c, http.StatusOK, "нет фото")
			return
		} else if len(names) == 1 {
			c.File(names[0])
		} else {
			output := "done.zip"
			if err := ZipFiles(output, names); err != nil {
				newResponse(c, http.StatusInternalServerError, err.Error())
				return
			}
			c.File(output)
		}

	})
	err = godotenv.Load(".env")
	if err!=nil {
		println(err.Error())
		return
	}
	port:= os.Getenv("HTTP_PORT")
	if err:=server.Run(":"+port); err!=nil {
		println(err.Error())
		return
	}

}



