apiVersion: v1
kind: Service
metadata:
  name: {{ .Values.appName }}-service
  labels:
{{ include "nested/all-labels.tpl" $ | indent 4 }}
spec:
  type: {{ .Values.service.type }}
  ports:
  - port: {{ .Values.service.port }}
    targetPort: {{ .Values.service.port }}
    protocol: TCP
  selector:
{{ include "helpers.tpl" $ | indent 4 }}