package application

type Metrics interface {
	RecordRequestTotal(router, protocol string, statusCode int)
	RecordDBQueryDuration(database, table, method string, duration float64)
	RecordRequestDuration(router, protocol string, statusCode int, duration float64)
	RecordTransactionTotal(status string)
	RecordAccountTotal()
}
