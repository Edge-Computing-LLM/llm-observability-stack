{{- define "llm-observability-stack.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "llm-observability-stack.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s" (include "llm-observability-stack.name" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{- define "llm-observability-stack.namespace" -}}
{{- default .Release.Namespace .Values.namespace.name -}}
{{- end -}}

{{- define "llm-observability-stack.langsmithSecretName" -}}
{{- if .Values.langsmith.existingSecret -}}
{{- .Values.langsmith.existingSecret -}}
{{- else -}}
langsmith-secrets
{{- end -}}
{{- end -}}

{{- define "llm-observability-stack.webuiSecretName" -}}
{{- if .Values.openWebUI.existingSecret -}}
{{- .Values.openWebUI.existingSecret -}}
{{- else -}}
open-webui-secrets
{{- end -}}
{{- end -}}
