
test:
	@go test -failfast -short .
dep:
	@echo "retrieving dependency" 
	@go get github.com/eaciit/toolkit
	@go get github.com/eaciit/errorlib
