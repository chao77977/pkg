#
# Main Trunk Makefile
#

GOBIN = $(shell pwd)/build/bin
COVERAGE_OUT=coverage.out
COVERAGE_HTML=coverage.html

all: alltools

alltools:
	@echo "Running unit tests"
	@build/env.sh go test -v -race -tags kqueue \
		--cover -coverprofile=${COVERAGE_OUT} ./...
	@build/env.sh go tool cover -html=${COVERAGE_OUT} \
		-o ${COVERAGE_HTML} 

clean:
	@echo "Cleaning up all the generated files"
	@find . -name '*.test' | xargs rm -fv
	@find . -name '*~' | xargs rm -fv
	@rm -fv ${COVERAGE_OUT} ${COVERAGE_HTML}
	@rm -rfv build/_workspace
	@rm -rfv build/bin

