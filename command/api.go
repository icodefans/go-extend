package command

import (
	"github.com/icodefans/go-extend/service"
)

func ApiServerStart(serverKey string) {
    service.Success("api server start", serverKey)
}
