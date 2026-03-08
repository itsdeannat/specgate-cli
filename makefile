APP = specgate
GO  = go

.PHONY: build run clean

build:
	$(GO) build -o $(APP) .

run:
	$(GO) run .

clean:
	rm -f $(APP)