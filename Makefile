.PHONY: all nuke deploy_with_gominprices2

deploy_target := $(DEPLOY_TO)
deploy_branch := $(DEPLOY_BRANCH)

all:
	go get github.com/mattn/gom
	mkdir -p _vendor/bin
	gom install
	gofmt -w ./main.go
	gom build -o _vendor/bin/goqueue ./main.go

deploy_with_gominprices2:
	ssh -A aviasales@gominprices2.int.avs.io "git clone git@github.com:KosyanMedia/goqueue.git; cd ~/goqueue && git clean -df && git fetch && git checkout ${deploy_branch} && git reset --hard origin/${deploy_branch} && make && echo good"
	ssh -A aviasales@${deploy_target} "mkdir -p ~/goqueue/backups && cd ~/goqueue && cp -rf goqueue backups/goqueue.`date +%s` ; scp gominprices2.int.avs.io:'~/goqueue/_vendor/bin/goqueue' goqueue.NEW && mv goqueue.NEW goqueue"

nuke:
	go clean -i