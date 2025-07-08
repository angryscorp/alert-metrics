package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/angryscorp/alert-metrics/internal/domain"
	"github.com/angryscorp/alert-metrics/internal/http/router"
)

type PingHandler struct {
	storage domain.MetricStorage
}

func NewPingHandler(storage domain.MetricStorage) PingHandler {
	return PingHandler{
		storage: storage,
	}
}

var _ router.PingHandler = (*PingHandler)(nil)

func (handler PingHandler) Ping(c *gin.Context) {
	if err := handler.storage.Ping(c.Request.Context()); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}
