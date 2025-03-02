package viewer

import (
	"fmt"

	"github.com/streletsa/savings-aggreagator/collector"
)

type ConsoleSavingsViewer struct {
}

func (ConsoleSavingsViewer) View(info *collector.SavingsCollectionInfo) {
	fmt.Printf("%v -> Total: %v\n", info.SourceType, info.Total)
}
