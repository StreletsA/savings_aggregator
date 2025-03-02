package main

import (
	"log"
	"time"

	"github.com/streletsa/savings-aggreagator/aggregator"
	"github.com/streletsa/savings-aggreagator/collector"
	"github.com/streletsa/savings-aggreagator/viewer"
)

var (
	collectors = [...]collector.SavingsCollector{
		collector.TBankInvestmentsSavingsCollector{},
	}
	viewers = [...]viewer.SavingsViewer{
		viewer.ConsoleSavingsViewer{},
		viewer.Esp32SavingsViewer{Config: viewer.LoadConfig("config.yaml")},
	}
)

func main() {
	aggregator := aggregator.SavingsAggregator{Collectors: collectors[:]}

	for {
		aggregationResult := aggregator.Aggregate()

		log.Printf("Calculated total: %.2f", aggregationResult.Total)

		for _, i := range aggregationResult.CollectorsSavingsInfo {
			for _, v := range viewers {
				v.View(&i)
			}
		}

		time.Sleep(5 * time.Second)
	}
}
