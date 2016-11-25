
{{ range $key, $var := .GetAllVars}}
variable "{{$var.VarName}}" { }
{{ end }}



{{ range $key, $module := .Modules }}

module "{{$module.Name}}" {
  source = "{{$module.URI}}"
  {{ range $key, $mapping := $module.GetVariables }}
  {{$mapping.VarName}} = "{{printf "${%s}" $mapping.VarPath}}"{{ end }}
}

{{ end }}


