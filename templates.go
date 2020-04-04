package commando

var usageTemplate = `
{{ if .IsRootCommand }}{{ .CliDesc }}{{ else }}{{ .Desc }}{{ end }}

Usage:
   {{ .Executable }} {{ with .Args -}}
   {{ range $k := . }}<{{ $k.Name }}> {{ end }}{{ end }}[flags]{{- if .IsRootCommand }}{{ if .Commands }}
   {{ .Executable }} <command> [flags]{{ end }}{{ end -}}


{{- /* commands */ -}}
{{- if .IsRootCommand -}}
{{- with .Commands  }}

Commands: {{ range $k, $v := . }}
   {{ printf "%-30v" $k }}{{ $v.ShortDesc }}
   {{- end -}}
{{- end -}}
{{- end -}}


{{- /* arguments */ -}}
{{- with .Args }}

Arguments: {{ range $k, $v := . }}
   {{ printf "%-30v" $k }}{{ $v.Desc }} {{ if $v.DefaultValue }}(default: {{ $v.DefaultValue }}){{ end }}
   {{- end -}}
{{- end -}}


{{- /* flags */ -}}
{{- with .Flags }}

Flags: {{ range $k, $v := . }}
   {{ if $v.ShortName -}}-{{ $v.ShortName }}, {{ printf "--%-24v" $k -}}
   {{ else -}}{{- printf "--%-28v" $k -}}
   {{- end -}}
   {{- $v.Desc }} {{ if $v.DefaultValue }}(default: {{ $v.DefaultValue }}){{ end }}
   {{- end -}}
{{- end -}}


{{- /* end */ -}}
{{- "" }}
`

var versionTemplate = `
Version: {{ .Version }}

{{- /* end */ -}}
{{- "" }}
`
