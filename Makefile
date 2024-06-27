APP_NAME=echo
ACR_NAME=echoapp

.PHONY: all build acr_build update_aca full_deploy run add_secrets clean help
 
build:
	go build -o echo

acr_build: increment_version
	@VERSION=$$(cat VERSION) && \
	az acr build --image $(APP_NAME):$${VERSION} --registry $(ACR_NAME) .

restart_aca:
	@export REVISION=$$(az containerapp revision list --name $(APP_NAME) --resource-group $(APP_NAME) --query "[].name" -o tsv) && \
	az containerapp revision restart -n $(APP_NAME) -g $(APP_NAME) --revision $${REVISION}

full_deploy: acr_build update_aca

update_aca:
	@VERSION=$$(cat VERSION) && \
	az containerapp update -n $(APP_NAME) -g $(APP_NAME) --image $(ACR_NAME).azurecr.io/$(APP_NAME):$${VERSION}

run:
	go run .
 
clean:
	go clean
	rm -f $(APP_NAME)

increment_version:
	@VERSION=$$(cat VERSION) && \
	NEW_VERSION=$$(echo $${VERSION} | awk -F. '{printf "%d.%d.%d", $$1, $$2, $$3+1}') && \
	echo $${NEW_VERSION} > VERSION && \
	echo "Updated version to $${NEW_VERSION}"
