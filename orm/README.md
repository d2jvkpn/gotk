### Title
---

#### P01
https://github.com/go-gorm/gorm/issues/4135 keep version of gorm.io/driver/postgres in v1.4.5, don't
upgrade to 1.5.0


#### P02: migration
```bash
go install -tags 'postgres mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

go install github.com/mikefarah/yq/v4@latest
```

```bash
migrate create -ext sql -dir migrations -seq alter_key_invocations

echo "export DATABASE_URL=postgresql://hello:world@localhost:5433/simple_bank?sslmode=disable" > .env

. .env

migrate -path migrations -database "$DATABASE_URL" -verbose up

migrate -path migrations -database "$DATABASE_URL" -verbose down -all
```

