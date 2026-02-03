apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Values.appName }}-{{ .Values.environment }}
  labels:
{{ include "helpers.tpl" $ | indent 4 }}
spec:
  replicas: {{ .Values.replicas }}
  selector:
    matchLabels:
{{ include "helpers.tpl" $ | indent 6 }}
  template:
    metadata:
      labels:
{{ include "helpers.tpl" $ | indent 8 }}
    spec:
      containers:
      - name: {{ .Values.appName }}
        image: {{ .Values.image.repository }}:{{ .Values.image.tag }}
        ports:
        - containerPort: {{ .Values.service.port }}
