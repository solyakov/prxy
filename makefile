SHELL := /bin/bash

DATA_DIR := data
SERVER_SAN := DNS:localhost,IP:127.0.0.1 # MUST MATCH the hostname or IP that the client uses to connect
CLIENT_SAN := DNS:example.com # Server validates expiry, CA signature, etc., but NOT the hostname

.PHONY: client server keys clean

client:
	go run cmd/client/main.go --certificate $(DATA_DIR)/client.crt --key $(DATA_DIR)/client.key --ca $(DATA_DIR)/ca.crt

server:
	go run cmd/server/main.go --certificate $(DATA_DIR)/server.crt --key $(DATA_DIR)/server.key --ca $(DATA_DIR)/ca.crt

$(DATA_DIR):
	mkdir -p $(DATA_DIR)

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
	    -extfile <(echo "subjectAltName = $(CLIENT_SAN)") \
		-days 3650

clean:
	rm -rf $(DATA_DIR)
