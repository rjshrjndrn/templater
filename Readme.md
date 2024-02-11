# Templater

![License](https://img.shields.io/badge/license-MIT-green.svg) ![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg) ![Version](https://img.shields.io/badge/version-1.0.0-blue.svg)

**Templater** is a cli tool that helps you to quickly and easily create files from templates, using go template.

## How to use

```bash
templater -h

cat <<'EOF'>file
---
{{- $name := "John" }}
class:
{{- range $i := until 11 }}
  {{$i}}: {{$name}}
{{- end }}
EOF

templater -i file
```


## Getting Started

### Prerequisites

- Go 1.19 or later

### Installation

```bash
go get github.com/rjshrjndrn/templater
```
