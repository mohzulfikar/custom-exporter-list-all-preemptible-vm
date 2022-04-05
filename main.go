package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	compute "cloud.google.com/go/compute/apiv1"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/expfmt"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
	"google.golang.org/protobuf/proto"
)

func explicit(jsonPath, projectID string) {
	var isPreemptibleTemp = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "nodes_preemptibles",
			Help: "Preemptible instance",
		},
		[]string{
			// name of each metrics
			"node_name",
			"node_cluster",
			"node_preemptibility",
		},
	)
	reg.MustRegister(isPreemptibleTemp)
	for {
		isPreemptibleTemp.Reset()
		fmt.Printf("WAITING.......")
		time.Sleep(5 * time.Second)
		ctx := context.Background()
		client, err := compute.NewInstancesRESTClient(ctx, option.WithCredentialsFile(jsonPath))
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()
		fmt.Println("Instances:")
		req := &computepb.AggregatedListInstancesRequest{
			Project: projectID,
			Filter:  proto.String(``),
		}
		it := client.AggregatedList(ctx, req)
		for {
			instancePair, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			instances := instancePair.Value.Instances
			if len(instances) > 0 {
				fmt.Printf("\n++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n")
				fmt.Printf("In zone: %s\n", instancePair.Key)
				for _, instance := range instances {
					fmt.Printf("\n================================================================================================\n")
					// fmt.Printf("%+v\n", instance)
					isPreemptible.With(prometheus.Labels{"node_name": instance.GetName(), "node_cluster": instance.GetName(), "node_preemptibility": fmt.Sprintf("%v", *instance.GetScheduling().Preemptible)}).Set(0)
					fmt.Printf("- Name: %s\n- Type: %s\n- Preemptible: %v\n", instance.GetName(), instance.GetMachineType(), *instance.GetScheduling().Preemptible)
				}
			}
		}
		isPreemptible = isPreemptibleTemp
		reg.Unregister(isPreemptibleTemp)
		hasil, err := isPreemptible.GetMetricWith(prometheus.Labels{})
		fmt.Printf("%#v\n", hasil)
		fmt.Printf("%#v\n", isPreemptibleTemp)

		gatherers := prometheus.Gatherers{
			reg,
		}
		gathering, err := gatherers.Gather()
		if err != nil {
			fmt.Println(err)
		}

		out := &bytes.Buffer{}
		for _, mf := range gathering {
			if _, err := expfmt.MetricFamilyToText(out, mf); err != nil {
				panic(err)
			}
		}
		fmt.Print(out.String())
		time.Sleep(10 * time.Second)
	}
}

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
