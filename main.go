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
	port := envVar(local)

	ctx := context.Background()

	srv, err := drive.NewService(ctx, option.WithCredentialsFile(credFile))
	if err != nil {
		log.Fatalf("Unable to parse credantials file: %v", err)
	}

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
		author := c.DefaultQuery("author", "SOTA")
		println(author)
		file,err:=c.GetRawData()
		//file, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			println("2")
			response.NewResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		f, err := os.Create("data")
		if err != nil {
			println("1")
			response.NewResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		defer f.Close()
		f.Write(file)
		//err = c.SaveUploadedFile(photo, "")
		//if err != nil {
		//	response.NewResponse(c, http.StatusInternalServerError, err.Error())
		//	return
		//}
		//err = service.SendPhoto(srv, photo.Filename)
		//if err != nil {
		//	println("3")
		//	response.NewResponse(c, http.StatusInternalServerError, err.Error())
		//	return
		//}
		//defer func(name string) {
		//	err := arch.MyDelete(name)
		//	if err != nil {
		//		response.NewResponse(c, http.StatusInternalServerError, err.Error())
		//		return
		//	}
		//}(photo.Filename)
		//c.String(http.StatusOK, fmt.Sprintf("'%s by %s' uploaded!", photo.Filename, author))

	})

	if err := server.Run(":" + port); err != nil {
		println(err.Error())
	}
	println("done")
}

func envVar(local bool) string {
	if local {
		err := godotenv.Load(".env")
		if err != nil {
			println(err.Error())
			return ""
		}
	}
	return os.Getenv("DRIVEAPI_PORT")
}
