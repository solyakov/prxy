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

Make sure to configure your browser to use the proxy client. The client listens on `localhost:8080` and forwards the traffic to `localhost:8081` by default.

To start the proxy server, run:

```bash
make server
```

The server listens on `localhost:8081` by default.

### TLS Certificates

The proxy client and server use self-signed certificates for TLS. These certificates are generated using the `make keys` command and are stored in the `data/` directory.

## Installation

1. Update the `PRXY_SERVER` environment variable in the [prxy-client.service](systemd/prxy-client.service) to match the IP address of the server machine.
2. Update the `SERVER_SAN` in the [makefile](makefile) to match the `PRXY_SERVER`.
3. Generate the TLS certificates by running `make keys`.
4. Install the proxy server by running `make install-server` on the server machine.
5. Install the proxy client by running `make install-client` on the client machine.

## Uninstallation

1. Uninstall the proxy client by running `make uninstall-client` on the client machine.
2. Uninstall the proxy server by running `make uninstall-server` on the server machine.

## Configuration

The proxy client and server configurations are stored in the [prxy-client.service](systemd/prxy-client.service) and [prxy-server.service](systemd/prxy-server.service) files, respectively.

## Chrome Extension

For convenience, a Chrome extension is provided in the [extensions](extensions/) directory. To install the extension, follow these steps:

1. Open the Chrome browser and navigate to `chrome://extensions`.
2. Enable the `Developer mode` toggle.
3. Click on the `Load unpacked` button and select the `extensions/HTTPSProxyToggle` directory.

All this extension does is toggle the proxy settings in the browser to use the proxy client. The address of the proxy client is hardcoded to `127.0.0.1:8080`. If the proxy client is running on a different address, update the extension accordingly.