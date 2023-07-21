package storage

type DataSourceHealth string

const (
	OK       DataSourceHealth = "OK"
	DEGRADED DataSourceHealth = "DEGRADED"
	DOWN     DataSourceHealth = "DOWN"
)
