go test -failfast -coverpkg ./... -coverprofile ./e2e/cover.out ./e2e
go tool cover -html ./e2e/cover.out -o ./e2e/cover.html
rm -rf ./e2e/cover.out