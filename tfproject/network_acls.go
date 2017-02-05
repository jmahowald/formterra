package tfproject

import (
	"github.com/fsouza/go-dockerclient/external/github.com/Sirupsen/logrus"
)

// NetworkACLRequest basic request to open the given ports to the given cidrs
type NetworkACLAllowRequest struct {
	Cidrs         []string
	Ports         []string
	TCP           bool
	UDP           bool
	StartingIndex int
}

type rule struct {
	Cidr       string
	Port       string
	Protocol   string
	RuleNumber int
}

// Simple wrapper
type networkACLRules struct {
	Rules []rule
}

func (n NetworkACLAllowRequest) makeRules() networkACLRules {
	i := n.StartingIndex
	protocols := make([]string, 0, 2)
	if n.TCP {
		protocols = append(protocols, "tcp")
	}
	if n.UDP {
		protocols = append(protocols, "udp")
	}
	if len(protocols) == 0 {
		logrus.Fatalf("You must specify at least one of tcp or udp")
	}
	rules := make([]rule, 0, len(n.Cidrs)*len(n.Ports)*len(protocols))
	for _, cidr := range n.Cidrs {
		for _, port := range n.Ports {
			for _, protocol := range protocols {
				r := rule{cidr, port, protocol, i}
				rules = append(rules, r)
				i++
			}
		}
	}
	return networkACLRules{rules}
}

// GetData returns itself for template contexts
func (n NetworkACLAllowRequest) getData() interface{} {
	return n.makeRules()
}

// Create Creates a terraform layer to create an s3 bucket
func (n NetworkACLAllowRequest) Create() (TerraformLayer, bool) {
	layer := TerraformLayer{Name: "networkacls"}
	path, _ := layer.getDir()
	processAssetTemplates(path, []string{"acls", "common"}, n)
	return layer, true
}
