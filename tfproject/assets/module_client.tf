{{ define "module"}}

module "{{.Name}}" {
  source = "{{.URI}}"
  {{ range $key, $mapping := .GetVariables }}
  {{$mapping.VarName}} = "{{printf "${%s}" $mapping.VarPath}}"{{ end }}
}

{{end}}
