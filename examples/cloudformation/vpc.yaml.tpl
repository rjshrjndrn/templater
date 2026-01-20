AWSTemplateFormatVersion: '2010-09-09'
Description: VPC with dynamic subnet generation using tpl function

Resources:
  VPC:
    Type: AWS::EC2::VPC
    Properties:
      CidrBlock: {{ .Values.vpc.cidr }}
      EnableDnsHostnames: true
      EnableDnsSupport: true

{{- range $index, $cidr := .Values.vpc.publicSubnets.cidrs }}
{{ tpl $.Values.vpc.publicSubnets.placeholder (dict "index" $index "cidr" $cidr) }}
{{- end }}
