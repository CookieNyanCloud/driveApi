package delivery

import (
	"fmt"
	"github.com/CookieNyanCloud/driveApi/internal/service"
	"github.com/CookieNyanCloud/driveApi/pkg/arch"
	"github.com/CookieNyanCloud/driveApi/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type input struct {
	Name string `json:"name"`
}

func (h *Handler) getPhoto(c *gin.Context) {
	println("start")
	var inp input
	if err := c.ShouldBindJSON(&inp); err != nil {
		response.NewResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	println(len(inp.Name))
	if len(inp.Name) <7 || strings.Contains(inp.Name,"."){
		response.NewResponse(c, http.StatusOK, "слишком широкая выборка")
		return
	}
	names, err := service.GetPhoto(h.drive, inp.Name)
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
}

func (h *Handler) sendPhoto(c *gin.Context) {
	dirType := c.PostForm("dirType")
	//dirType := c.Query("dirType")
	file, err := c.FormFile("file")
	if err != nil {
		response.NewResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err = c.SaveUploadedFile(file, file.Filename); err != nil {
		response.NewResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	//author := c.PostForm("author")
	err = service.SendPhoto(h.drive, file.Filename, dirType, h.conf.DrivePeople, h.conf.DriveZag)
	if err = arch.MyDelete(file.Filename); err != nil {
		response.NewResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, fmt.Sprintf("file uploaded!"))
}
