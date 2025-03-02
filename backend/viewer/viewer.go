package viewer

import (
	"github.com/streletsa/savings-aggreagator/collector"
)

type SavingsViewer interface {
	View(*collector.SavingsCollectionInfo)
}
