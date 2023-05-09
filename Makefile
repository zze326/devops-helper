package:
	GOOS=linux GOARCH=amd64 revel package . -m prod

run-app:
	revel run . -m dev

test-app:
	revel test .

build:
	revel build . ./run

docker-build:
	docker build --platform linux/amd64 -t registry-azj-registry.cn-shanghai.cr.aliyuncs.com/ops/devops-platform/revel-websocket:202303171224 .