package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// 客户端实例化
func (es *ElasticSearch) New() (err error) {
	es.Client, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{fmt.Sprintf("%s://%s:%d", es.Scheme, es.Host, es.Port)},
		// Transport: &http.Transport{
		// 	MaxIdleConnsPerHost:   50,               // 每个主机最大空闲连接数
		// 	MaxIdleConns:          100,              // 总最大空闲连接数
		// 	IdleConnTimeout:       90 * time.Second, // 空闲连接超时
		// 	ResponseHeaderTimeout: 10 * time.Second, // 响应头超时
		// 	DialContext: (&net.Dialer{
		// 		Timeout:   5 * time.Second,  // 连接建立超时
		// 		KeepAlive: 30 * time.Second, // 保持连接时间
		// 	}).DialContext,
		// },
		// // 各种超时设置
		// DiscoverNodesInterval: 5 * time.Minute, // 节点发现间隔
		// MaxRetries:            3,               // 最大重试次数
		// DisableRetry:          true,            // 禁用超时重试
	})
	return err
}

// 索引初始化
func (es *ElasticSearch) Migrate(indexName string, s any) error {
	res, err := es.Client.Indices.Exists([]string{indexName})
	if err != nil {
		return fmt.Errorf("es.Client.Indices.Exists Err: %s", err)
	}
	_ = res.Body.Close()
	// 创建索引时可以使用 PutMapping 来定义映射
	if res.StatusCode != 404 {
		// 索引已存在
	} else if createRes, err := es.Client.Indices.Create(indexName); err != nil {
		return fmt.Errorf("es.Client.Indices Err: %s", err)
	} else if createRes.IsError() {
		return fmt.Errorf("createRes.IsError: %s", createRes.String())
	} else {
		_ = createRes.Body.Close()
	}
	// 使用结构体标签动态生成Elasticsearch映射
	var mappingJson []byte
	if mapping, err := GenerateMapping(s); err != nil {
		return fmt.Errorf("GenerateMapping Err: %s", err)
	} else if mappingJson, err = json.Marshal(mapping); err != nil {
		return fmt.Errorf("json.MarshalErr: %s", err)
	}
	// 添加映射
	putMappingReq := esapi.IndicesPutMappingRequest{
		Index: []string{indexName},
		Body:  strings.NewReader(string(mappingJson)),
	}
	putMappingRes, err := putMappingReq.Do(context.Background(), es.Client)
	if err != nil {
		return fmt.Errorf("putMappingRes Err: %s", err)
	}
	defer putMappingRes.Body.Close()
	if putMappingRes.IsError() {
		return fmt.Errorf("putMappingRes.IsError: %s", putMappingRes.String())
	}
	return nil
}
