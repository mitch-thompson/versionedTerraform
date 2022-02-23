####################################################
# Build
####################################################
build:
	go build -o versionedTerraform ./cmd
####################################################
# Clean
####################################################
clean:
	rm -f ~/.local/bin/versionedTerraform
####################################################
# Install
####################################################
install:
	mv versionedTerraform ~/.local/bin/
####################################################
# help feature
####################################################
help:
	@echo ''
	@echo 'Usage: make [TARGET]'
	@echo 'Targets:'
	@echo '  build    		go build -o versionedTerraform ./cmd'
	@echo '  clean		    removes installed versionedTerraform file'
	@echo '  install        installs versionedTerraform to local user bin folder'
	@echo '  all      		Nothing to do.'
	@echo ''