

migration:
	goose -dir sql/migration postgres "postgres://root:root@localhost:5432/users" up