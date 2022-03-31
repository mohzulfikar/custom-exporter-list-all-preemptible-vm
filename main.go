package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var reg = prometheus.NewRegistry()

var isPreemptible = promauto.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "nodes_preemptible",
		Help: "Preemptible instance",
	},
	[]string{
		// name of each label
		"node_name",
		"node_cluster",
		"node_preemptibility",
	},
)

func init() {
	reg.MustRegister(opsProcessed)
}

func main() {
	// go recordMetrics()
	go explicit("./creds.json", "")

	http.Handle("/something", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	http.ListenAndServe(":9191", nil)
}
