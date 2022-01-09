.PHONY: build
build:
	sam build

.PHONY: deploy
deploy:
	sam deploy

# prepare resrouces for CI/CD
.PHONY: cicd
cicd:
	aws cloudformation deploy \
		--region ap-northeast-1 \
		--stack-name "holidays-jp-cicd" \
		--template-file "cicd.yaml" \
		--capabilities CAPABILITY_NAMED_IAM
