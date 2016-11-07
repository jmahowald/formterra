

module "{{.Name}}" {
  source = "{{.SourceURI}}"
  {{ range $key, $mapping := .Mappings }}
  {{$mapping.VarName}} = "{{printf "${%s}" $mapping.VarPath}}"{{ end }}

}
