package healthlake

type config struct {
	awsRegion, dataStoreID string
}

func NewConfig(awsRegion, dataStoreID string) config {
	return config{
		awsRegion:   awsRegion,
		dataStoreID: dataStoreID,
	}
}
