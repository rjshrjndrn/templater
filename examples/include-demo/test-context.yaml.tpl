# Using . (current context)
{{- range .Values.environments }}
Environment: {{ . }}
# Inside range, . is the current item, so we use $ for root context
{{ include "helpers.tpl" $ }}
---
{{- end }}

# Using . outside any scope (same as $)
{{ include "helpers.tpl" . }}