.PHONY: modules.outdated modules.direct vet

#list outdated modules
modules.outdated:
	go list -u -m -f '{{if .Update}}{{.}}{{end}}' all

#list direct modules
modules.direct:
	go list -u -m -f '{{if not .Indirect}}{{.}}{{end}}' all

vet:
	go vet ./...