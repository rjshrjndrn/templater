packer {
  required_plugins {
    amazon = {
      version = ">= 0.0.2"
      source  = "github.com/hashicorp/amazon"
    }
  }
}

{{$version := .Values.version}}
{{- $regions := .Values.regions }}
{{- range $_, $region := .Values.regions}}
source "amazon-ebs" "{{$region}}" {
  ami_name      = "or-ubuntu-arm64"
  instance_type = "t4g.xlarge"
  region        = "{{$region}}"
  source_ami_filter {
    filters = {
      name = "ubuntu/images/*ubuntu-jammy-22.04-arm64-server-*"
      root-device-type    = "ebs"
      virtualization-type = "hvm"
    }
    most_recent = true
    owners      = ["099720109477"] # Canonical
  }
  ssh_username = "ubuntu"
  tags = {
    version = "{{$version}}"
    createdBy = "packer"
  }
}
{{- end}}

build {
  name = "or-ee-saas-ami"
  sources = [
  {{- range $_, $value := $regions}}
    "source.amazon-ebs.{{$value}}",
  {{- end}}
  ]
}

