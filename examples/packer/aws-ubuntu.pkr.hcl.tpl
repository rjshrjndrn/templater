packer {
  // Specify required plugins and their versions for this Packer configuration.
  required_plugins {
    amazon = {
      version = ">= 0.0.2" // Minimum plugin version requirement.
      source  = "github.com/hashicorp/amazon" // Plugin source location.
    }
  }
}

{{/* Initialize variables from Helm values for reuse. */}}
{{- $version := .Values.version }}
{{- $regions := .Values.regions }}

{{/* Iterate over each region specified in Helm values. */}}
{{- range $_, $region := $regions }}

// Define an Amazon EBS (Elastic Block Store) source configuration for each region.
source "amazon-ebs" "{{ $region }}" {
  ami_name      = "or-ubuntu-arm64" // AMI name pattern.
  instance_type = "t4g.xlarge" // Instance type to use.
  region        = "{{ $region }}" // AWS region from Helm values.
  source_ami_filter {
    filters = {
      name = "ubuntu/images/*ubuntu-jammy-22.04-arm64-server-*" // Filter for Ubuntu AMI.
      root-device-type    = "ebs" // EBS-backed AMI.
      virtualization-type = "hvm" // Hardware Virtual Machine.
    }
    most_recent = true // Use the most recent AMI matching the filter.
    owners      = ["099720109477"] // Canonical as the AMI owner.
  }
  ssh_username = "ubuntu" // Default SSH username.
  tags = {
    version   = "{{ $version }}" // Inject version from Helm values.
    createdBy = "packer" // Tag to identify the creator.
  }
}

{{- end }}

build {
  name = "or-ee-saas-ami" // Build name for the resulting AMI.
  sources = [
    // Inject source references for each region defined in Helm values.
    {{- range $_, $value := $regions}}
    "source.amazon-ebs.{{$value}}",
    {{- end}}
  ]
}
