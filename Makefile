E2E_TEST_PATH=./tests

ENV_LOCAL_TEST=\
	JWT_SECRET=secret
#   POSTGRES_HOST=
#   POSTGRES_DB=
#   POSTGRES_USER=
#   POSTGRES_PASSWORD=
#   JWT_SECRET=

test-e2e:
	$(ENV_LOCAL_TEST) go test ${E2E_TEST_PATH} -v

db.seed.up:
	goose -dir ./db/seed -no-versioning up

db.seed.reset:
	goose -dir ./db/seed -no-versioning reset

db.migrate.up:
	goose -dir ./db/migrations up

db.migrate.reset:
	goose -dir ./db/migrations reset

db.status:
	goose -dir ./db/migrations status