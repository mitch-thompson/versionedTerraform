package versionedTerraform

import (
	"testing"
	"testing/fstest"
)

const (
	firstFile = `
resource "aws_mq_broker" "sample" {
 depends_on = [aws_security_group.mq]
 broker_name = var.name
 engine_type = "ActiveMQ"
 engine_version = var.mqEngineVersion
 host_instance_type = var.hostInstanceType
 security_groups = [aws_security_groups.mq.id]
 apply_immediately = "true"
 deployment_mode = "ACTIVE_STANDBY_MULTI_AZ"
 auto_minor_version_upgrade = "true"
 subnet_ids = ["10.0.0.0/24", "10.0.1.0/24"]
}
`
	secondFile = `
terraform {
 required_version = "~> 0.12.4"
}
`
)

func TestFileHandler(t *testing.T) {
	want := NewVersion("0.12.31", testVersionList())

	fs := fstest.MapFS{
		"main.tf":     {Data: []byte(firstFile)},
		"versions.tf": {Data: []byte(secondFile)},
	}

	version, err := GetVersionFromFile(fs, testVersionList(), true)

	if err != nil {
		t.Fatal(err)
	}

	got := *version

	if got.Version != want.Version {
		t.Errorf("Expected %v, got %v", want.Version, got.Version)
	}
}

func TestEmptyTerraformVersion(t *testing.T) {
	want := NewVersion("1.1.11", testVersionList())

	fs := fstest.MapFS{"main.tf": {Data: []byte(firstFile)}}

	version, err := GetVersionFromFile(fs, testVersionList(), true)

	if err != nil {
		t.Fatal(err)
	}

	got := *version

	if got.Version != want.Version {
		t.Errorf("Expected %v, got %v", want.Version, got.Version)
	}
}
