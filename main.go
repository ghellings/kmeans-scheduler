package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/muesli/clusters"
	"github.com/muesli/kmeans"
)

func main() {

	sitemap, d := parsecsv("./data.csv")

	// Partition the data points into 16 clusters
	km := kmeans.New()
	clusters, err := km.Partition(d, 3)
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}
	for i, c := range clusters {

		fmt.Printf("Cluster %v\n", i)
		for _, o := range c.Observations {
			statmap := makestatmap(o.Coordinates())
			for _, site := range sitemap[statmap] {
				fmt.Printf("%v\n", site)
			}
		}
	}
}

func parsecsv(file string) (sitemap map[string][]string, d clusters.Observations) {
	sitemap = make(map[string][]string)
	csvfile, err := os.Open(file)
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}
	r := csv.NewReader(csvfile)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		site := record[0]
		record = record[1:]
		stats := []float64{}
		cluster := clusters.Coordinates{}
		for _, val := range record {
			fval, err := strconv.ParseFloat(val, 64)
			if err != nil {
				panic(fmt.Sprintf("Can't convert %v %v", val, err))
			}
			stats = append(stats, fval)
			cluster = append(cluster, fval)
		}
		d = append(d, cluster)
		statmap := makestatmap(stats)
		sitemap[statmap] = append(sitemap[statmap], site)
	}
	return
}

func makestatmap(stats []float64) (statmap string) {
	for _, val := range stats {
		if statmap == "" {
			statmap = fmt.Sprintf("%f", val)
		} else {
			statmap = fmt.Sprintf("%s-%f", statmap, val)
		}
	}
	return
}
