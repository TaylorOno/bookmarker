.PHONY: sonar
## sonar: *REQUIRES LOCAL SONARQUBE and SONAR-SCANNER run uploads to local sonar instance
sonar: test
	sonar-scanner -Dsonar.projectKey=bookmarker -Dsonar.exclusions=**/*_test.go,**/test_data/*,**/mocks/**,**/main.go -Dsonar.host.url=http://localhost:9000 -Dsonar.source=. -Dsonar.go.coverage.reportPaths=**/coverage.out


.PHONY: run-local
## run-local: *REQUIRES DOCKER COMPOSE runs a local instance of the app and dependencies
run-local:
	export ENVIRONMENT=local
	docker-compose up -d

.PHONY: test
## test: runs all tests
test:
	go test ./... -race -coverprofile=./coverage.out

.PHONY: vet
## vet: runs go vet
vet:
	go vet ./...

.PHONY: fmt
## fmt: runs go fmt
fmt:
	go fmt ./...

.PHONY: pre-release
## pre-release: runs all tests and go tools
pre-release: fmt vet test

.PHONY: help
## help: Prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ":" | sed -e 's/^/ /'