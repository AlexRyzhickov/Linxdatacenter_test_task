.PHONY: docker-build
docker-build:
	@docker build -t linxdatacenter_test_task .

.PHONY: run-with-json
run-with-json:
	@docker run --rm -v `pwd`/data:/bin/data linxdatacenter_test_task --filename ./data/db.json

.PHONY: run-with-csv
run-with-csv:
	@docker run --rm -v `pwd`/data:/bin/data linxdatacenter_test_task --filename ./data/db.csv

.PHONY: run-with-help
run-with-help:
	@docker run --rm -v `pwd`/data:/bin/data linxdatacenter_test_task --help

