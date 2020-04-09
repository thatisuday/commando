package commando

var usageTemplate = `
{{ if .IsRootCommand }}{{ .CliDesc }}{{ else }}{{ .Desc }}{{ end }}

Usage:
   {{ .Executable }} {{ with .Args -}}
   {{ range $k, $v := . }}{{ if $v.IsRequired }}<{{ $v.ClpArg.Name }}>{{ else }}[{{ $v.ClpArg.Name }}]{{ end }} {{ end }}{{ end }}{flags}{{- if .IsRootCommand }}{{ if .Commands }}
   {{ .Executable }} <command> {flags}{{ end }}{{ end -}}


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
   {{ printf "%-30v" $v.ClpArg.Name }}{{ $v.Desc }}{{ if $v.ClpArg.DefaultValue }} (default: {{ $v.ClpArg.DefaultValue }}){{ end }}{{ if $v.ClpArg.IsVariadic }} {variadic}{{ end }}
   {{- end -}}
{{- end -}}


{{- /* flags */ -}}
{{- with .Flags }}

Flags: {{ range $k, $v := . }}
   {{ if $v.ClpFlag.ShortName -}}-{{ $v.ClpFlag.ShortName }}, {{ if $v.ClpFlag.IsInverted }}{{ printf "--no-%-24v" $k }}{{ else }}{{ printf "--%-24v" $k }}{{ end -}}
   {{ else }}{{ if $v.ClpFlag.IsInverted }}{{ printf "--no-%-25v" $k }}{{ else }}{{ printf "--%-28v" $k }}{{ end -}}
   {{- end -}}
   {{- $v.Desc }} {{ if $v.ClpFlag.DefaultValue }}(default: {{ $v.ClpFlag.DefaultValue }}){{ end }}
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
