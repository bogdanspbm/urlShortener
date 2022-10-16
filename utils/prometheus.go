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

	if err := push.New("http://217.25.88.166:9091", "register_link").
		Collector(promRegisteredLinkCount).
		Grouping("urls", "create").
		Push(); err != nil {
		fmt.Println("Could not push completion time to Pushgateway:", err)
	}

	if err := push.New("http://217.25.88.166:9091", "receive_link").
		Collector(promReceivedLinkCount).
		Grouping("urls", "get").
		Push(); err != nil {
		fmt.Println("Could not push completion time to Pushgateway:", err)
	}
}
