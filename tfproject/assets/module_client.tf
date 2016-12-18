{{ define "module"}}

module "{{.Name}}" {
  source = "{{.LocalLocation}}"
  {{ range $key, $mapping := .GetVariables }}
  {{$mapping.VarName}} = "{{printf "${%s}" $mapping.VarPath}}"{{ end }}
}

{{end}}
