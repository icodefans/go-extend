package logger

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/sirupsen/logrus"
)

// ElasticHook 定义 Elasticsearch 的 Logrus Hook
type ElasticHook struct {
	client    *elasticsearch.Client
	host      string
	index     string
	levels    []logrus.Level
	ctx       context.Context
	ctxCancel context.CancelFunc
}

// NewElasticHook 创建新的 Elasticsearch Hook
func NewElasticHook(client *elasticsearch.Client, host, index string, levels []logrus.Level) *ElasticHook {
	ctx, cancel := context.WithCancel(context.TODO())

	return &ElasticHook{
		client:    client,
		host:      host,
		index:     index,
		levels:    levels,
		ctx:       ctx,
		ctxCancel: cancel,
	}
}

// Fire 实现 Logrus Hook 接口
func (hook *ElasticHook) Fire(entry *logrus.Entry) error {
	level := entry.Level.String()

	msg := struct {
		Level     string        `json:"level"`
		Host      string        `json:"host"`
		Timestamp string        `json:"timestamp"`
		Message   string        `json:"message"`
		Fields    logrus.Fields `json:"fields"`
	}{
		level,
		hook.host,
		entry.Time.UTC().Format(time.RFC3339Nano),
		entry.Message,
		entry.Data,
	}

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = hook.client.Index(
		hook.index,
		strings.NewReader(string(jsonMsg)),
		hook.client.Index.WithContext(hook.ctx),
	)

	return err
}

// Levels 实现 Logrus Hook 接口
func (hook *ElasticHook) Levels() []logrus.Level {
	return hook.levels
}
