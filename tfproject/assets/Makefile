# Generated Makefile, created by {{version}}

# If there is a terraform.tfvars
# use that as terraform input
# rather than depend on environment variables or user input
TF_VAR_ARGS ?=
ifneq ($(wildcard $(TF_VAR_FILE)),)
TF_VAR_ARGS=--var-file=$(TF_VAR_FILE)
endif

TF_FILES=$(wildcard *.tf)

.PHONY: plan refresh apply clean 

.terraform: $(TF_FILES)
	terraform get --update
	for i in $$(ls .terraform/modules/*/Makefile); do i=$$(dirname $$i); echo "Trying make in $$i"; make -C $$i; done

plan: $(TF_FILES) .terraform
	terraform plan -out terraform.tfplan $(TF_VAR_ARGS)

refresh:
	terraform refresh $(TF_VAR_ARGS)

terraform.tfplan: $(TF_FILES) .terraform
	$(MAKE) plan

apply: terraform.tfplan
	terraform apply terraform.tfplan

clean:
	terraform destroy -force $(TF_VAR_ARGS)
	rm -rf .terraform
	rm -rf terraform.tfstate
	rm -rf terraform.tfplan


