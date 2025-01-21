module github.com/spearkkk/u

go 1.23.5

replace github.com/spearkkk/u/uuid => ./uuid

replace github.com/spearkkk/u/timestamp => ./timestamp

require (
	github.com/deanishe/awgo v0.29.1
	github.com/itchyny/timefmt-go v0.1.6
	github.com/sosodev/duration v1.3.1
)

require (
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/smarty/assertions v1.15.0 // indirect
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/magefile/mage v1.11.0 // indirect
	github.com/smartystreets/goconvey v1.8.1
	go.deanishe.net/env v0.5.1 // indirect
	go.deanishe.net/fuzzy v1.0.0 // indirect
	golang.org/x/text v0.3.6 // indirect
)
