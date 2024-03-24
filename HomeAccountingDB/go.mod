module HomeAccountingDB

go 1.22

toolchain go1.22.1

replace TimeSeriesData => ../TimeSeriesData

require (
	TimeSeriesData v0.0.0-00010101000000-000000000000
	github.com/sergz72/expreval v0.0.0-20240324155213-cdc165c776de
)
