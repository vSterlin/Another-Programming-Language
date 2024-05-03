go test ./... -coverprofile=build/coverage.out -cover -run ^Test
go tool cover -html=build/coverage.out