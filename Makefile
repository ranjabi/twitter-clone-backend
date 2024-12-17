E2E_TEST_PATH=./tests
export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING=postgres://postgres:123456@localhost:5432/postgres

ENV_LOCAL_TEST=\
	JWT_SECRET=secret
#   POSTGRES_HOST=
#   POSTGRES_DB=
#   POSTGRES_USER=
#   POSTGRES_PASSWORD=
#   JWT_SECRET=

test:
	$(ENV_LOCAL_TEST) go test ${E2E_TEST_PATH} -v

test.no-cache:
	$(ENV_LOCAL_TEST) go test ${E2E_TEST_PATH} -count=1 -v

db.seed.up:
	goose -dir ./db/seed -no-versioning up

db.seed.reset:
	goose -dir ./db/seed -no-versioning reset

db.seed.status:
	goose -dir ./db/seed status

db.up:
	goose -dir ./db/migrations up

db.reset:
	goose -dir ./db/migrations reset

db.clean:
	make db.mig.reset
	make db.mig.up

db.status:
	goose -dir ./db/migrations status