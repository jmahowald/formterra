variable "vpc_id" {}

//By default allow no traffic
resource "aws_network_acl" "main" {
    vpc_id = "${var.vpc_id}"
    ingress {
        protocol = "all"
        rule_no = {{.StartingIndex}}  
        action = "deny"
        from_port=0
        to_port=0
        cidr_block ="0.0.0.0/0"
    }
    tags {
        Name = "main"
    }
}


{{ range $rule := .GetRules.Rules}}
resource "aws_network_acl_rule" "in_{{$rule.RuleNumber}}" {
    network_acl_id = "${aws_network_acl.main.id}"
    egress = "false"
    from_port = "{{$rule.Port}}"
    to_port = "{{$rule.Port}}"
    rule_number = "{{$rule.RuleNumber}}"
    rule_action = "allow"
    cidr_block = "{{$rule.Cidr}}"
    protocol = "{{$rule.Protocol}}"
}
{{end}}