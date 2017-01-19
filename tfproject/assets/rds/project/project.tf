
{{ range $var := .GetAllVars}}
variable "{{$var.VarName}}" { {{ if eq $var.Type "" }} {{ else}}
  type="{{$var.Type}}" {{end}} {{ if eq $var.DefaultValue ""}} {{else}}
  default="{{$var.DefaultValue}}" {{ end }}
}
{{ end }}



{{ range $module := .Modules }}

module "{{$module.Name}}" {
  source = "{{$module.LocalLocation}}"
  {{ range $key, $mapping := $module.GetVariables }}
  {{$mapping.VarName}} = "{{printf "${%s}" $mapping.VarPath}}"{{ end }}
}

{{ range $remote := .RemoteVariables }}
data "terraform_remote_state" "{{$remote.RemoteSourceName}}" {
  backend="s3"
  config { {{ range $key,$value := $remote.Config }}
    {{$key}} = "{{$value}}"{{end}} }
} {{end}}



{{ range $output:= .Outputs }}
output "{{$output}}" {
  value = "${module.{{$module.Name}}.{{$output}}}"
} {{end}}


{{end}}
