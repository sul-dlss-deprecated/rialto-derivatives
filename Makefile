default: package

package: solr postgres

solr:
	GOOS=linux go build -o solr_derivative cmd/solr/main.go
	zip solr_derivative.zip solr_derivative

postgres:
	GOOS=linux go build -o postgres_derivative cmd/postgres/main.go
	zip postgres_derivative.zip postgres_derivative
