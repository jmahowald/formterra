# If there is a terraform.tfvars
# use that as terraform input
# rather than depend on environment variables or user input
TF_VAR_ARGS ?=
ifneq ($(wildcard $(TF_VAR_FILE)),)
TF_VAR_ARGS=--var-file=$(TF_VAR_FILE)
endif

TF_FILES=$(wildcard *.tf)


# If you define a bucket name and bucket prefix
# all will include storing off the remote state 
# protected witha  guard

REMOTE_INIT_FILE = remote_state.txt
REMOTE_STATE_BUCKET ?= 
REMOTE_STATE_PREFIX ?= 
REMOTE_STATE_BUCKET_REGION ?=

# TODO figureo ut how to to do and
ifneq ($(REMOTE_STATE_BUCKET),)

all: $(REMOTE_INIT_FILE) plan

else

all: plan

endif




.PHONY: plan refresh apply clean 

.terraform: $(TF_FILES)
	terraform get --update

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

$(REMOTE_INIT_FILE): 
	@echo "initialise remote statefile"
	terraform remote config -backend=s3 \
	-backend-config="bucket=$(REMOTE_STATE_BUCKET)" \
	-backend-config="key=$(REMOTE_STATE_PREFIX)/terraform.tfstate" \
	-backend-config="region=$(REMOTE_STATE_BUCKET_REGION)"
	echo "remote state initialized with $(REMOTE_STATE_BUCKET) and $(REMOTE_STATE_PREFIX)" > $(REMOTE_INIT_FILE)
