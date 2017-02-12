{{- range $var := .GetAllVars}}
variable "{{$var.VarName}}" { {{ if eq $var.Type "list" -}} 
  type="list" {{end -}}
{{- if ne $var.DefaultValue "" }} default = "{{$var.DefaultValue}}" {{end -}}
}{{ end }}

{{ range $module := .Modules -}}
module "{{$module.Name}}" {
  source = "{{$module.LocalLocation}}"
  {{- range $key, $mapping := $module.GetVariables }}
  {{ $mapping.VarName}} = "{{printf "${%s}" $mapping.VarPath}}"{{ end }}
}

{{ range $output:= .Outputs }}
output "{{$output}}" {
  value = "${module.{{$module.Name}}.{{$output}}}"
} {{end}}


{{end}}


{{ range $remote := .GetAllRemotes }}
data "terraform_remote_state" "{{$remote.RemoteSourceName}}" {
  backend="s3"
  config { {{ range $key,$value := $remote.Config }}
    {{$key}} = "{{$value}}"{{end}} }
} {{end}}
