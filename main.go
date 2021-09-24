package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/CookieNyanCloud/driveApi/arch"
	"github.com/CookieNyanCloud/driveApi/response"
	"github.com/CookieNyanCloud/driveApi/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type input struct {
	Name string `json:"name"`
}

const (
	credFile = "driveapisearch.json"
)
func main() {

	var local bool
	flag.BoolVar(&local, "local", false, "хост")
	flag.Parse()
	port, drivePeople, driveZag := envVar(local)
	ctx := context.Background()

	srv, err := drive.NewService(ctx, option.WithCredentialsFile(credFile))
	if err != nil {
		log.Fatalf("Unable to parse credantials file: %v", err)
	}

	server := gin.Default()
	server.MaxMultipartMemory = 8 << 25
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
			defer func() {
				err := arch.MyDelete(names[0])
				if err != nil {
					response.NewResponse(c, http.StatusInternalServerError, err.Error())
					return
				}
			}()
			return
		} else {
			output := "done.zip"
			if err := arch.ZipFiles(output, names); err != nil {
				response.NewResponse(c, http.StatusInternalServerError, err.Error())
				return
			}
			c.File(output)
			defer func() {
				err := arch.MyDelete(output)
				if err != nil {
					response.NewResponse(c, http.StatusInternalServerError, err.Error())
					return
				}
			}()
			defer func() {
				err := arch.AllDelete(names)
				if err != nil {
					response.NewResponse(c, http.StatusInternalServerError, err.Error())
					return
				}
			}()
			return
		}
	})
	server.POST("/sendphoto", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			response.NewResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		if err = c.SaveUploadedFile(file, file.Filename); err != nil {
			response.NewResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		dirType:= c.PostForm("dirType")
		author:= c.PostForm("author")
		err = service.SendPhoto(srv, file.Filename, author,dirType,drivePeople, driveZag)

		c.String(http.StatusOK, fmt.Sprintf("file uploaded!"))
	})

	if err := server.Run(":" + port); err != nil {
		println(err.Error())
	}
	println("done")
}

func envVar(local bool) (string, string, string) {
	if local {
		err := godotenv.Load(".env")
		if err != nil {
			println(err.Error())
			return "", "", ""
		}
	}
	return os.Getenv("DRIVEAPI_PORT"), os.Getenv("DRIVE_PEOPLE"), os.Getenv("DRIVE_ZAG")
}
