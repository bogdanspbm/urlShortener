package utils

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

var promRegisteredLinkCount = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "put_requests",
		Help: "Количество зарегистрированных ссылок",
	},
)

var promReceivedLinkCount = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "get_requests",
		Help: "Количество запросов ссылок",
	},
)

func RegPrometheus() {
	prometheus.MustRegister(promReceivedLinkCount)
	prometheus.MustRegister(promRegisteredLinkCount)
}

func PrometheusPush() {

	if err := push.New("http://localhost:9091/", "db_backup").
		Collector(promRegisteredLinkCount).
		Grouping("urls", "create").
		Push(); err != nil {
		fmt.Println("Could not push completion time to Pushgateway:", err)
	}
}
