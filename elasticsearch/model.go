package elasticsearch

import (
	"github.com/elastic/go-elasticsearch/v8"
)

// elasticsearch
type ElasticSearch struct {
	Scheme string                `validate:"required,oneof=http https" mapstructure:"Scheme" label:"协议"`
	Host   string                `validate:"required" mapstructure:"Host" label:"主机名"`
	Port   int32                 `validate:"required" mapstructure:"Port" label:"端口"`
	Index  string                `validate:"required" mapstructure:"Index" label:"索引"`
	Client *elasticsearch.Client `label:"客户端"`
}
