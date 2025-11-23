.DEFAULT_GOAL := main

.PHONY: init
init:
	limactl start --name=default template://ubuntu-lts

.PHONY: exec-limactl
exec-limactl:
	limactl shell default

.PHONY: main
main:
	go run main.go run /bin/bash
