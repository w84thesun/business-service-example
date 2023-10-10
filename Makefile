include .env
export

PROJECT_NAME=business-service-example

gen-json:
	cd pkg/api/ && easyjson --all .

gen-proto:
	protoc --gogoslick_out=plugins=grpc:pkg/ -I=. api/grpc/proto/messages/*.proto

download:
	go mod download

tidy:
	go mod tidy

lint:
	golangci-lint run

build-service:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
        go build \
        -installsuffix 'static' \
        -o ./app ./cmd/${PROJECT_NAME}/*

build-docker:
	docker build \
       -t "${CONTAINER_NAME}:${TAG}" \
       --build-arg GOPROXY=${GOPROXY} \
       --build-arg PROJECT_NAME=${PROJECT_NAME} \
       --network host \
       -f ${DOCKER_FILE_PATH} \
       .

run:
	go run cmd/${PROJECT_NAME}

run-docker-deps:
	docker run \
	--network host \
	--env MONGO_INITDB_ROOT_USERNAME=${MONGO_USER} \
	--env MONGO_INITDB_ROOT_PASSWORD=${MONGO_PASS} \
	--env MONGO_INITDB_DATABASE=${MONGO_DB} \
	-d \
	mongo:${MONGO_VERSION}

run-docker-service:
	docker run -d \
	--name ${PROJECT_NAME} \
	--env-file .env \
	--network host \
	"${REPO}:${TAG}"

save-docker-service-logs:
	docker logs -f ${PROJECT_NAME} >> ./functional-test-logs.txt &

stop-docker-service:
	docker rm -f ${PROJECT_NAME}

run-functional-test:
	go test -v  -timeout=30m ${PWD}/tests/functional

migrate-up:
	migrate -verbose -path ./migrations/ -database $(DATABASE) up $(NUMBER)

migrate-down:
	migrate -verbose -path ./migrations/ -database $(DATABASE) down $(NUMBER)

run-local-deps:
	./deployments/local/run_dependencies.sh

migrate-functional-dbs:
	echo "stub"
