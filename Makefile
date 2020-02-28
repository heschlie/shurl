.PHONY: start-postgres
start-postgres: stop-postgres
	docker run --rm --name test-postgres \
	    -e POSTGRES_USER=shurl \
	    -e POSTGRES_PASSWORD=micron \
	    -e POSTGRES_DB=shurl \
	    -p 5432:5432 \
	    -v $(pwd)/sql/schema:/docker-entrypoint-initdb.d \
	    -d postgres

.PHONY: stop-postgres
stop-postgres:
	-docker stop test-postgres
	-docker rm -f test-postgres