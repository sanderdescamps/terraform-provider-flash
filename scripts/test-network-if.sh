export PURE_API_VERSION="1.19"
export PURE_USERNAME=pureuser
export PURE_PASSWORD=pureuser
export PURE_TARGET=flasharray1.testdrive.local
export PURE_URL="https://${PURE_TARGET}"

export PURE_API_TOKEN=$(curl -k -XPOST "${PURE_URL}/api/${PURE_API_VERSION}/auth/apitoken" \
    -H 'Content-Type: application/json' \
    -d "{\"password\": \"$PURE_PASSWORD\",\"username\": \"$PURE_USERNAME\"}" | jq -r '.api_token' )

curl -k -XPOST "${PURE_URL}/api/${PURE_API_VERSION}/auth/session" \
    -H 'Content-Type: application/json' \
    -c cookies.txt \
    -d "{\"api_token\": \"$PURE_API_TOKEN\"}"


curl -k -XPUT "${PURE_URL}/api/${PURE_API_VERSION}/network/vir1" \
    -H 'Content-Type: application/json' \
    -d "{\"address\": \"192.168.6.3\",\"gateway\": \"192.168.6.1\",\"netmask\": \"255.255.255.0\",\"mtu\": 1500}" \
    -b cookies.txt -v | jq -r '.'

curl -k -XPUT "${PURE_URL}/api/${PURE_API_VERSION}/network/vir1?address=192.168.6.4&gateway=192.168.6.3&netmask=255.255.255.0&mtu=1500" \
    -H 'Content-Type: application/json' \
    -b cookies.txt -v | jq -r '.'
