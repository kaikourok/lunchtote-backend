.PHONY: dev
dev:
	air

.PHONY: deps
deps:
	go mod download

.PHONY: cli
cli:
	cd cmd/cli && go build -o ../../build/cli

.PHONY: dev-deps
dev-deps:
# VSCode Go拡張用
	go install github.com/cweill/gotests/gotests@latest
	go install github.com/fatih/gomodifytags@latest
	go install github.com/josharian/impl@latest
	go install github.com/haya14busa/goplay/cmd/goplay@latest
	go install github.com/go-delve/delve/cmd/dlv@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install golang.org/x/tools/gopls@latest
# ホットリロード
	go install github.com/cosmtrek/air@latest

.PHONY: test
test:
	go test -v ./...

.PHONY: check
check:
	go vet ./...
	staticcheck ./...

.PHONY: migrate-latest
migrate-latest:
	make cli && ./build/cli migrate-latest

.PHONY: migrate-drop
migrate-drop:
	make cli && ./build/cli migrate-drop

# make migrate-create name=...
.PHONY: migrate-create
migrate-create:
	migrate create -ext sql -dir infrastructure/database/migrations -seq ${name}

.PHONY: init
init:
	make cli && ./build/cli migrate-drop && ./build/cli migrate-latest && ./build/cli init