.PHONY: up down full fmt plan init clean 

init:
	terraform init

fmt:
	terraform fmt 

plan: init
	terraform plan 

up: fmt 
	terraform apply -auto-approve 

down: 
	terraform destroy -auto-approve

full: down up

clean: 
	rm -rf .terraform terraform.tfstate* 