go test ./... -coverpkg ./... -covermode atomic -coverprofile c.out
go tool cover -html=c.out
rm c.out
