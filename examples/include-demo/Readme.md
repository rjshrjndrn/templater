# Include Function Example

This demonstrates the `include` function that loads and renders external template files.

## Files

- `helpers.tpl` - Common labels template
- `deployment.yaml.tpl` - Kubernetes Deployment using include
- `service.yaml.tpl` - Kubernetes Service with nested includes
- `values.yaml` - Configuration values
- `nested/common-tags.tpl` - Additional tags
- `nested/all-labels.tpl` - Combines multiple includes

## Usage

```bash
# Basic include
templater -i deployment.yaml.tpl -f values.yaml

# Nested includes (templates including other templates)
templater -i service.yaml.tpl -f values.yaml
```

## How it works

The `include` function loads an external template file and renders it with the provided context:

```yaml
{{ include "helpers.tpl" $ }}
```

- First argument: path to template file (relative or absolute)
- Second argument: context data (use `$` for root context, `.` for current context)
- Returns: rendered template content as string

### Context: `.` vs `$`

- **`$`** = Root context (recommended, works everywhere, same as Helm)
- **`.`** = Current context (changes inside `range`, `with` blocks)

```yaml
# Best practice: use $ to always access root context
{{ include "helpers.tpl" $ }}

# Inside a range, . changes but $ stays as root
{{- range .Values.items }}
  item: {{ . }}                    # current item
  app: {{ $.Values.appName }}      # root context via $
  {{ include "helpers.tpl" $ }}    # pass root context
{{- end }}
```

### Relative paths

Paths are resolved relative to the template's directory:
- `include "helpers.tpl" .` - same directory
- `include "../other.tpl" .` - parent directory  
- `include "nested/file.tpl" .` - subdirectory

### Nested includes

Templates can include other templates. Each resolves paths relative to its own directory:

```yaml
# In nested/all-labels.tpl
{{ include "../helpers.tpl" . }}
{{ include "common-tags.tpl" . }}
```
