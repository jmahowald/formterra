
{{ range $var := .GetAllVars}}
variable "{{$var.VarName}}" { {{ if eq $var.Type "list" }} type="list" {{end}} }
{{ end }}



{{ range $module := .Modules }}

module "{{$module.Name}}" {
  source = "{{$module.URI}}"
  {{ range $key, $mapping := $module.GetVariables }}
  {{$mapping.VarName}} = "{{printf "${%s}" $mapping.VarPath}}"{{ end }}
}

{{ range $remote := .RemoteVariables }}
data "terraform_remote_state" "{{$remote.RemoteSourceName}}" {
  backend="s3"
  config { {{ range $key,$value := $remote.Config }}
    {{$key}} = "{{$value}}"{{end}} }
} {{end}}

{{end}}  

