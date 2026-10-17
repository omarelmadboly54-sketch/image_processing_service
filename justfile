set shell := ["powershell.exe", "-c"]

server:
    go run main.go

DB_SOURCE := env('DB_SOURCE', 'postgresql://postgres:omar@localhost/image_processing_service?sslmode=disable')
MIGRATIONS_DIR := env('MIGRATIONS_DIR', './migrations')

create-migration name:
    migrate create -ext sql -dir {{MIGRATIONS_DIR}} -seq {{name}}

migrate-up:
    migrate -path {{MIGRATIONS_DIR}} -database "{{DB_SOURCE}}" up


migrate-down:
     migrate -path {{MIGRATIONS_DIR}} -database "{{DB_SOURCE}}" down 1
