package main

import (
	_ "github.com/icodefans/go-extend/aliyun"
	_ "github.com/icodefans/go-extend/cache"
	_ "github.com/icodefans/go-extend/cloud"
	_ "github.com/icodefans/go-extend/command"
	_ "github.com/icodefans/go-extend/database"
	_ "github.com/icodefans/go-extend/define"
	_ "github.com/icodefans/go-extend/function"
	_ "github.com/icodefans/go-extend/logger"
	"github.com/icodefans/go-extend/service"
	_ "github.com/icodefans/go-extend/sms"
	_ "github.com/icodefans/go-extend/storage"
)

func main() {
	service.Success("ok", nil, "dev")
}
