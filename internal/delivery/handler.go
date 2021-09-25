package delivery

import (
	"github.com/CookieNyanCloud/driveApi/internal/config"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/drive/v3"
)

type Handler struct {
	drive *drive.Service
	conf  *config.Conf
}

func NewHandler(drive *drive.Service, conf *config.Conf) *Handler {
	return &Handler{
		drive: drive,
		conf:  conf,
	}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.Default()
	router.POST("/getphoto", h.getPhoto)
	router.POST("/sendphoto", h.sendPhoto)
	return router
}
