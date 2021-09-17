package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/CookieNyanCloud/driveApi/api"
	"github.com/CookieNyanCloud/driveApi/arch"
	"github.com/CookieNyanCloud/driveApi/response"
	"github.com/CookieNyanCloud/driveApi/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type input struct {
	Name string `json:"name"`
}

func main() {
	var local bool
	flag.BoolVar(&local, "local", false, "хост")
	flag.Parse()
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

	server := gin.Default()
	server.POST("/getphoto", func(c *gin.Context) {
		println("start")
		var inp input
		if err := c.ShouldBindJSON(&inp); err != nil {
			response.NewResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		names, err := service.GetPhoto(srv, inp.Name)
		fmt.Println(names)
		if err != nil {
			response.NewResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		if len(names) == 0 {
			response.NewResponse(c, http.StatusOK, "нет фото")
			return
		} else if len(names) == 1 {
			c.File(names[0])
			defer myDelete(names[0])
			return
		} else {
			output := "done.zip"
			if err := arch.ZipFiles(output, names); err != nil {
				response.NewResponse(c, http.StatusInternalServerError, err.Error())
				return
			}
			c.File(output)
			defer func() {
				err := myDelete(output)
				if err != nil {
					response.NewResponse(c, http.StatusInternalServerError, err.Error())
					return
				}
			}()
			defer func() {
				err := allDelete(names)
				if err != nil {
					response.NewResponse(c, http.StatusInternalServerError, err.Error())
					return
				}
			}()
			return
		}
	})
	if local {
		err = godotenv.Load(".env")
		if err != nil {
			println(err.Error())
			return
		}
	}
	port := os.Getenv("DRIVEAPI_PORT")
	if err := server.Run(":" + port); err != nil {
		println(err.Error())
		return
	}
		println("done")

}

func myDelete(name string) error {
	return os.Remove(name)
}

func allDelete(names []string) error {
	for _, v := range names {
		fmt.Println(v)
		err := os.Remove(v)
		if err != nil {
			return err
		}
	}
	return nil
}
