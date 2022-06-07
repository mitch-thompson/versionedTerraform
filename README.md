# Versioned Terraform
A wrapper for terraform to detect the expected version of terraform, 
download, and execute that version

## Requirements
- go compiler (only tested on go1.17)

## Install
`make build install` for installation to local user<br>
`make build` will create an executable file for you to place where you'd like

## Commands
```
All arguments are passed through to terraform
```

## sample usage
`versionedTerraform version` will display the terraform version executed in a folder

## Known Issues
