run:
	@go run cmd/runner/main.go

mig:
	@migrate migrations -database "postgres://postgres.gnqvmormznwyaikdapqy:XtuC1Lm51KC6NiUX@aws-0-eu-central-1.pooler.supabase.com:6543/postgres" -verbose up