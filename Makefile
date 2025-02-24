E2E_TEST_PATH=./tests
GOOSE_DRIVER=postgres
GOOSE_DBSTRING=postgres://dev:123456@localhost:5432/twitter_clone

air:
	ENV_NAME=.env.dev.local air

air.test:
	ENV_NAME=.env.dev.local SEED=test air

test.it.all:
	go test ./tests/it -count=1 -v

test.it.file:
	go test ./tests/it/main_test.go ./tests/it/${n} -count=1 -v
# 
test.it.func:
	go test ./tests/it -count=1 -v -run TestRunWithSuite/${n}

test:
	$(ENV_LOCAL_TEST) go test ${E2E_TEST_PATH} -v

# test.partial:
# 	$(ENV_LOCAL_TEST) go test ${E2E_TEST_PATH} -v -run "TestHealthCheck|TestUserRegister|TestUserLogin|TestUserLoginNotExist"

test.no-cache:
	$(ENV_LOCAL_TEST) go test ${E2E_TEST_PATH} -count=1 -v

# format.staging:
# 	git diff --name-only --cached | grep '\.go$$' | xargs -n 1 -I {} sh -c 'go fmt {}; echo "---- {}"'

mig.seed.up:
	GOOSE_DBSTRING=${GOOSE_DBSTRING} GOOSE_DRIVER=${GOOSE_DRIVER} goose -dir ./db/seed -no-versioning up

mig.seed.reset:
	GOOSE_DBSTRING=${GOOSE_DBSTRING} GOOSE_DRIVER=${GOOSE_DRIVER} goose -dir ./db/seed -no-versioning reset

mig.seed.status:
	GOOSE_DBSTRING=${GOOSE_DBSTRING} GOOSE_DRIVER=${GOOSE_DRIVER} goose -dir ./db/seed status

mig.up:
	GOOSE_DBSTRING=${GOOSE_DBSTRING} GOOSE_DRIVER=${GOOSE_DRIVER} goose -dir ./db/migrations up

mig.reset:
	GOOSE_DBSTRING=${GOOSE_DBSTRING} GOOSE_DRIVER=${GOOSE_DRIVER} goose -dir ./db/migrations reset

mig.clean:
	make mig.reset
	make mig.up

mig.init:
	make mig.up
	make mig.seed.up

mig.status:
	goose -dir ./db/migrations status

compose.down.dev:
	docker compose -f docker-compose.dev.yml down --volumes

compose.up.dev:
	make compose.down.dev
	docker compose -f docker-compose.dev.yml --env-file .env.dev.docker up --build