# Prxy

## Architecture Diagram

```

                                   localnet | Internet
                                            |
+---------+                   +--------+    |      +--------+    +--------------+
| Browser |--(HTTP CONNECT)-->| Client |==(mTLS)==>| Server |--->| HTTPS Server |
+---------+                   +--------+    |      +--------+    +--------------+
                                            |

```

## Local Setup

To generate the necessary TLS certificates, run the following command:

```bash
make keys
```

To start the proxy client, run:

```bash
make client
```

Make sure to configure your browser to use the proxy client. The client listens on `localhost:8080` by default.

To start the proxy server, run:

```bash
make server
```

The server listens on `localhost:8081` by default.

### TLS Certificates

The proxy client and server use self-signed certificates for TLS. These certificates are generated using the `make keys` command and are stored in the `data/` directory.

Make sure to update the SAN (Subject Alternative Name) fields in the `makefile` file to match your local setup.
