# $GOROOT/bin/go
NAME = tess
cmd_add:
	cobra-cli add storeBin
run:
	/home/ant/.cache/JetBrains/GoLand2024.2/tmp/GoLand/___go_build_goph_keeper_cli login -l=as -p=asd

migration_create:
	cd ./migrations && \
	$(GOPATH)/bin/goose $(GOOSE_DRIVER) $(GOOSE_DBSTRING) create table sql

migration_drop_all:
	cd /home/ant/go-work/go-musthave-diploma-tpl/migrations && \
	$(GOPATH)/bin/goose $(GOOSE_DRIVER) $(GOOSE_DBSTRING) down-to 0

mock_gen:
	$(GOPATH)/bin/mockgen -destination=mocks/mock_sender.go -package=mocks goph_keeper/cli/remote IRemote
	$(GOPATH)/bin/mockgen -destination=mocks/mock_store.go -package=mocks goph_keeper/cli/storage IStore
#build cli with tags
build_cli:
	$(GOPATH)/bin/go build -ldflags "-X main.Version=v1.0.1 -X main.Date=11-12-25"