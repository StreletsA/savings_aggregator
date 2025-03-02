package collector

type SavingsCollectorSourceType string

const (
	T_BANK_INVESTMENTS_SOURCE SavingsCollectorSourceType = "T_Investments"
)

type SavingsCollectionInfo struct {
	SourceType SavingsCollectorSourceType
	Total      float32
}

type SavingsCollector interface {
	Collect() (SavingsCollectionInfo, error)
	GetSourceType() SavingsCollectorSourceType
}
