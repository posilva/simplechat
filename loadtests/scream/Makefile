.PHONY: run create_cdeploy run_cdeploy shell_cdeploy clean_cdeploy docker_cp_cdeploy docker_build_cdeploy

run: 
	iex -S mix 

create_cdeploy:
	MIX_ENV=prod mix release  scream --path ../terraform/deploy/appdeploy/scream --force --overwrite --quiet

run_cdeploy:
	../terraform/deploy/appdeploy/scream/bin/scream start 

shell_cdeploy:
	../terraform/deploy/appdeploy/scream/bin/scream start_iex

clean_cdeploy:
	rm -rf ../terraform/deploy/appdeploy/scream/bin/scream

docker_build_cdeploy:
	docker buildx build --load --platform linux/arm64 -t scream --progress plain --pull=false -f code_deploy/Dockerfile .


docker_clean_cdeploy:
	docker kill scream_cp || true && docker rm scream_cp || true

docker_cp_cdeploy: clean_cdeploy
	docker run --name scream_cp -d -it scream start_iex && \
		docker cp scream_cp:/opt/scream ../terraform/deploy/appdeploy && \
		docker stop scream_cp || true 

deploy_cdeploy: docker_clean_cdeploy docker_build_cdeploy docker_cp_cdeploy 
	cd ../terraform/deploy && terraform apply -auto-approve
