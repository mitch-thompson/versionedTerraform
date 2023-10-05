####################################################
# Build
####################################################
build:
	go build -o versionedTerraform ./cmd
####################################################
# Clean
####################################################
clean:
	rm -f $(shell go env GOPATH)/bin/versionedTerraform
####################################################
# Install
####################################################
install:
	mv versionedTerraform $(shell go env GOPATH)/bin/
####################################################
# help feature
####################################################
help:
	@echo ''
	@echo 'Usage: make [TARGET]'
	@echo 'Targets:'
	@echo '  build    go build -o versionedTerraform ./cmd'
	@echo '  clean    removes installed versionedTerraform file'
	@echo '  install  installs versionedTerraform to bin folder in GOPATH'
	@echo '  all      Nothing to do.'
	@echo ''
