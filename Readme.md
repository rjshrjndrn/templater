# Templater

![License](https://img.shields.io/badge/license-MIT-green.svg) ![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg) ![Version](https://img.shields.io/badge/version-1.0.0-blue.svg)

**Templater** is a cli tool that helps you to quickly and easily create files from templates, using go template.

## How to use

1. Inline variables

```bash
# File using helm templating
cat <<'EOF'>file
---
{{- $name := "John" }}
class:
{{- range $i := until 11 }}
  {{$i}}: {{$name}}
{{- end }}
EOF

templater -i file

# This will be the output

File: out/file
---
class:
  0: John
  1: John
  2: John
  3: John
  4: John
  5: John
  6: John
  7: John
  8: John
  9: John
  10: John

```

2. Passing variables from external yaml, using helm style templating.
```bash
# Value file
cat <<'EOF'>values.yaml
school: Govt. Public School
EOF
# File using helm templating
cat <<'EOF'>file
---
{{- $name := "John" }}
class:
  school: {{ .Values.school }}
{{- range $i := until 11 }}
  {{$i}}: {{$name}}
{{- end }}
EOF

templater -i file -f values.yaml -o out/

# This will be the output

File: out/file
---
class:
  school: Govt. Public School
  0: John
  1: John
  2: John
  3: John
  4: John
  5: John
  6: John
  7: John
  8: John
  9: John
  10: John
```

## Getting Started

### Prerequisites

- Go 1.19 or later

### Installation

1. Binary

Download the latest binary from [Release Page.](https://github.com/rjshrjndrn/templater/releases)

2. Using go get

```bash
go get github.com/rjshrjndrn/templater
```
