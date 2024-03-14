.PHONY: openapi

check-env:
ifndef GOPATH
	$(error GOPATH is undefined)
endif

openapi: check-env
	${GOPATH}/bin/oapi-codegen  -old-config-style -package="api" pkg/api/openapi.yaml > pkg/api/zz_openapi_generated.go