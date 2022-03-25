package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() {
	for {
		// dat, err := explicit()

		if err != nil {
			fmt.Println(err)
			time.Sleep(15 * time.Second)
			continue
		}
		opsProcessed.Inc()
		time.Sleep(2 * time.Second)
	}
}

var reg = prometheus.NewRegistry()

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})
)

func init() {
	reg.MustRegister(opsProcessed)
}

func main() {
	go recordMetrics()

	http.Handle("/something", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	http.ListenAndServe(":9191", nil)
}
