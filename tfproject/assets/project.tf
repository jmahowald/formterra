
{{ range $key, $var := .GetAllVars}}
variable "{{$var.VarName}}" { }
{{ end }}




{{ range $key, $module := .Modules }}
{{template "module" $module}}
{{ end }}
