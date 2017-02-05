variable "vpc_id" {}

//By default allow no traffic
resource "aws_network_acl" "main" {
    vpc_id = "${var.vpc_id}"
    ingress {
        protocol = "all"
        rule_no = 1 
        action = "deny"
        from_port=0
        to_port=0
    }
    tags {
        Name = "main"
    }
}


{{ range $rule := .Rules}}
resource "aws_network_acl_rule" "in_{{$rule.RuleNumber}}" {
    network_acl_id = "${aws_network_acl.main.id}"
    egress = "false"
    from_port = "{{$rule.Port}}"
    to_port = "{{$rule.Port}}"
    rule_no = "{{$rule.RuleNumber}}"
    cidr_block = "{{$rule.Cidr}}"
    protcol = "{{$rule.Protocol}}"
}
{{end}}