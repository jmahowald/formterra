
# If there is an envrc
# simulate direnv by turning variable assignments
# into make assignments.  This is probably incredibly
ifneq ($(wildcard .envrc),)
$(info ************  including .envrc in make ************)
IGNORE := $(shell bash -c "grep '=' .envrc | sed 's/=/:=/' > makeenv")
include makeenv
endif

apply: terraform.tfplan
	terraform apply terraform.tfplan
	rm terraform.tfplan

plan: $(TF_FILES) .terraform
	terraform plan -out terraform.tfplan




.PHONY: plan refresh apply clean
TF_FILES=$(wildcard *.tf)

.terraform: $(TF_FILES)
	terraform get --update


refresh:
	terraform refresh

terraform.tfplan: $(TF_FILES) .terraform
	$(MAKE) plan


clean:
	terraform destroy -force
	rm -rf terraform.tfplan
	# rm -rf terraform.tfstate
