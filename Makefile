    # Go parameters
    GOCMD=go
		CMDDIR=./cmd
    OUTDIR=./out
    BINARY_NAME=btc-cli
    GOBUILD=$(GOCMD) build
    GOCLEAN=$(GOCMD) clean
    GOTEST=$(GOCMD) test
    GOGET=$(GOCMD) get
    GORUN=$(GOCMD) run
    
    all: test 
    build:
			$(GOBUILD) -o $(OUTDIR)/$(BINARY_NAME) -v $(CMDDIR)
    build-linux: deps
			GOOS=linux $(GOBUILD) -o $(OUTDIR)/$(BINARY_NAME) -v $(CMDDIR)
    test: 
			$(GOTEST) -v ./...
    clean: 
			$(GOCLEAN)
			rm -rf $(OUTDIR)
    run: build	
			$(OUTDIR)/$(BINARY_NAME)