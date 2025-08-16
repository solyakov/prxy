#!/bin/bash

# This example script is used to configure iptables to redirect incoming HTTPS traffic from
# the br-lan interface to the prxy-client running on the router in transparent proxy mode.
#
#                                             localnet | Internet
#                                                      |
# +---------+            +------------------------+    |      +--------+    +--------------+
# | Browser |--(HTTPS)-->| --(redirect)--> Client |==(mTLS)==>| Server |--->| HTTPS Server |
# +---------+            +------------------------+    |      +--------+    +--------------+
#                                                      |

cd $(dirname $0)

SERVER_IP=<PRXY SERVER IP>

export PRXY_SERVER=$SERVER_IP:443
export PRXY_CERTIFICATE=client.crt
export PRXY_KEY=client.key
export PRXY_CA=ca.crt
export PRXY_LISTEN=0.0.0.0:8080
export PRXY_TIMEOUT=60s
export PRXY_BUFFER=32768

function log() {
  echo $(date) $*
}

function configure_iptables() {
  # Create PRXY chain
  iptables -t nat -N PRXY
  
  # Skip localhost and prxy server (to avoid loops)
  iptables -t nat -A PRXY -d 127.0.0.0/8 -j RETURN
  iptables -t nat -A PRXY -d $SERVER_IP -p tcp --dport 443 -j RETURN
  
  # Redirect all other HTTPS to transparent proxy on port 8080
  iptables -t nat -A PRXY -p tcp --dport 443 -j REDIRECT --to-ports 8080

  # Allow tcp/8080 on INPUT from br-lan (otherwise redirected packets get dropped)
  iptables -I INPUT 1 -i br-lan -p tcp --dport 8080 -j ACCEPT
  
  # Apply to br-lan clients (PREROUTING)
  iptables -t nat -I PREROUTING -i br-lan -p tcp --dport 443 -j PRXY
  
  # Apply to router HTTPS connections (OUTPUT) 
  iptables -t nat -I OUTPUT -p tcp --dport 443 -j PRXY
}

function run_forever() {
  while sleep 1s; do
    log starting client
    ./prxy-client-arm64 -T
  done
}

configure_iptables
run_forever
