# Demonstrating . vs $ in include

# Top level: . and $ are the same
Top level access: {{ .Values.appName }}

# Inside range: . changes to current item, $ stays as root
{{- range $index, $env := .Values.environments }}
---
environment_{{ $index }}:
  name: {{ $env }}
  # Using $ to access root context
  app: {{ $.Values.appName }}
  # Include with $ passes root context (correct)
  labels:
{{ include "helpers.tpl" $ | indent 4 }}
{{- end }}
