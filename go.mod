module github.com/Ruslan-kulm/snmp_monitor

go 1.16

require (
	github.com/gosnmp/gosnmp v1.32.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/prometheus/client_golang v1.11.0
)

//replace github.com/gosnmp/gosnmp v1.32.0 => github.com/Ruslan-kulm/gosnmp v1.32.1-0.20210902105148-268d68a56382
