module github.com/ebenamoafo2/scalable-go-api

go 1.26.3

require (
	github.com/golang-migrate/migrate/v4 v4.19.1
	github.com/joho/godotenv v1.5.1
	github.com/julienschmidt/httprouter v1.3.0
	github.com/lib/pq v1.12.3
	github.com/tomasen/realip v0.0.0-20180522021738-f0c99a92ddce
	github.com/wneessen/go-mail v0.7.3
	golang.org/x/crypto v0.52.0
	golang.org/x/time v0.15.0
)

require (
	github.com/BurntSushi/toml v1.4.1-0.20240526193622-a339e1f7089c // indirect
	golang.org/x/exp/typeparams v0.0.0-20231108232855-2478ac86f678 // indirect
	golang.org/x/mod v0.35.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/text v0.37.0 // indirect
	golang.org/x/tools v0.44.0 // indirect
	honnef.co/go/tools v0.7.0 // indirect
)

tool honnef.co/go/tools/cmd/staticcheck
