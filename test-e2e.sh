go test -failfast -coverpkg ./graph/resolvers/...,./modules/... -coverprofile ./e2e/cover.out ./e2e
go tool cover -html ./e2e/cover.out -o ./e2e/coverage.html