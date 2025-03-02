package aggregator

import (
	"fmt"

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

type SavingsAggregator struct {
	Collectors []collector.SavingsCollector
}

type SavingsAggregationResult struct {
	Total                 float32
	CollectorsSavingsInfo []collector.SavingsCollectionInfo
}

func (x SavingsAggregator) Aggregate() SavingsAggregationResult {
	collectorsSavingsInfo := []collector.SavingsCollectionInfo{}

	for _, c := range collectors {
		info, err := c.Collect()
		if err != nil {
			fmt.Println(err)
			continue
		}

		collectorsSavingsInfo = append(collectorsSavingsInfo, info)
	}

	var total float32 = 0.0

	for _, c := range collectorsSavingsInfo {
		total += c.Total
	}

	return SavingsAggregationResult{total, collectorsSavingsInfo}
}
