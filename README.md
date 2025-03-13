[![CircleCI](https://circleci.com/gh/mitch-thompson/versionedTerraform.svg?style=svg)](https://circleci.com/gh/mitch-thompson/versionedTerraform)

# Versioned Terraform
A wrapper for terraform to detect the expected version of terraform, 
download, and execute that version

## Requirements
- go

## Install
`make build install` for installation to local user<br>
`make build` will create an executable file for you to place where you'd like

## Commands
```
All arguments are passed through to terraform
```

## Sample usage
`versionedTerraform version` will display the terraform version executed in a folder

## Configuration
A configuration file is created in `~/.versionedTerraform`<br><br>

`StableOnly` boolean values: <b>true</b>/false<br>
This value is used to restrict terraform to release versions only defaults to true
## Known Issues
