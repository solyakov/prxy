SHELL := /bin/bash

DATA_DIR := data
SERVER_SAN := DNS:localhost,IP:127.0.0.1 # MUST MATCH the hostname or IP of the server
ETC_DIR := /etc/prxy
BIN_DIR := /usr/local/bin

.PHONY: client server keys clean install-server install-client uninstall-server uninstall-client

client:
	go run ./cmd/client/ --certificate $(DATA_DIR)/client.crt --key $(DATA_DIR)/client.key --ca $(DATA_DIR)/ca.crt

client-certificate:
	openssl pkcs12 -export -in $(DATA_DIR)/client.crt -inkey $(DATA_DIR)/client.key -out $(DATA_DIR)/client.p12 -name "Client Certificate" -legacy

server:
	go run cmd/server/main.go --certificate $(DATA_DIR)/server.crt --key $(DATA_DIR)/server.key --ca $(DATA_DIR)/ca.crt

$(DATA_DIR):
	mkdir -p $@

keys: $(DATA_DIR)
	openssl ecparam -name prime256v1 -genkey -noout -out $(DATA_DIR)/ca.key
	openssl req -new -x509 -nodes -sha256 -key $(DATA_DIR)/ca.key -out $(DATA_DIR)/ca.crt \
	    -subj "/CN=CA" \
	    -addext "basicConstraints = CA:TRUE,pathlen:0" \
	    -addext "keyUsage = critical, digitalSignature, cRLSign, keyCertSign" \
		-days 3650

	openssl ecparam -name prime256v1 -genkey -noout -out $(DATA_DIR)/server.key
	openssl req -new -key $(DATA_DIR)/server.key -out $(DATA_DIR)/server.csr \
	    -subj "/CN=Server"
	openssl x509 -req -in $(DATA_DIR)/server.csr -out $(DATA_DIR)/server.crt \
		-CA $(DATA_DIR)/ca.crt \
		-CAkey $(DATA_DIR)/ca.key \
		-CAcreateserial \
	    -extfile <(echo "subjectAltName = $(SERVER_SAN)") \
		-days 3650

	openssl ecparam -name prime256v1 -genkey -noout -out $(DATA_DIR)/client.key
	openssl req -new -key $(DATA_DIR)/client.key -out $(DATA_DIR)/client.csr \
	    -subj "/CN=Client"
	openssl x509 -req -in $(DATA_DIR)/client.csr -out $(DATA_DIR)/client.crt \
		-CA $(DATA_DIR)/ca.crt \
		-CAkey $(DATA_DIR)/ca.key \
		-CAcreateserial \
		-days 3650

install-server:
	go build -o $(DATA_DIR)/prxy-server cmd/server/main.go
	sudo install -d -m 755 $(DATA_DIR) $(ETC_DIR)
	sudo install -m 644 $(DATA_DIR)/server.crt $(ETC_DIR)/server.crt
	sudo install -m 644 $(DATA_DIR)/server.key $(ETC_DIR)/server.key
	sudo install -m 644 $(DATA_DIR)/ca.crt $(ETC_DIR)/ca.crt
	sudo install -m 755 $(DATA_DIR)/prxy-server $(BIN_DIR)/prxy-server
	sudo install -m 644 systemd/prxy-server.service /etc/systemd/system/prxy-server.service
	sudo systemctl daemon-reload
	sudo systemctl enable --now prxy-server

uninstall-server:
	sudo systemctl disable --now prxy-server
	sudo rm -f $(BIN_DIR)/prxy-server
	sudo rm -f $(ETC_DIR)/server.crt
	sudo rm -f $(ETC_DIR)/server.key
	sudo rm -f $(ETC_DIR)/ca.crt
	-sudo rmdir $(ETC_DIR)
	sudo rm -f /etc/systemd/system/prxy-server.service
	sudo systemctl daemon-reload

install-client:
	go build -o $(DATA_DIR)/prxy-client ./cmd/client/
	sudo install -d -m 755 $(DATA_DIR) $(ETC_DIR)
	sudo install -m 644 $(DATA_DIR)/client.crt $(ETC_DIR)/client.crt
	sudo install -m 644 $(DATA_DIR)/client.key $(ETC_DIR)/client.key
	sudo install -m 644 $(DATA_DIR)/ca.crt $(ETC_DIR)/ca.crt
	sudo install -m 755 $(DATA_DIR)/prxy-client $(BIN_DIR)/prxy-client
	sudo install -m 644 systemd/prxy-client.service /etc/systemd/system/prxy-client.service
	sudo systemctl daemon-reload
	sudo systemctl enable --now prxy-client

uninstall-client:
	sudo systemctl disable --now prxy-client
	sudo rm -f $(BIN_DIR)/prxy-client
	sudo rm -f $(ETC_DIR)/client.crt
	sudo rm -f $(ETC_DIR)/client.key
	sudo rm -f $(ETC_DIR)/ca.crt
	-sudo rmdir $(ETC_DIR)
	sudo rm -f /etc/systemd/system/prxy-client.service
	sudo systemctl daemon-reload

windows-client:
	GOOS=windows GOARCH=amd64 go build -o $(DATA_DIR)/prxy-client.exe ./cmd/client/

arm64-client:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -ldflags '-extldflags "-static"' -o prxy-client-arm64 ./cmd/client


HTTPSProxyToggle:
	cd extensions/HTTPSProxyToggle && zip -r ../../data/HTTPSProxyToggle.zip .

clean:
	rm -rf $(DATA_DIR)
