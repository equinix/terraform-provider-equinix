# Terraform examples for Fabric Cloud Router (FCR)

## Create FCR

```cd FCR/
echo "=============Initialize FCR Resource==========="
terraform init

echo "=============Create FCR ==========="

terraform apply -auto-approve -var="fcr_name=terra_fcr-test"

echo "=============GET FCR ==========="

terraform show

result=$(terraform show | grep state)

echo "######### FCR $result ##########"
```
## Update FCR
```#--------update FCR name
echo "=============Update FCR Name==========="

terraform apply -refresh -auto-approve -var="fcr_name=terra_fcr-test-update"

terraform show

result=$(terraform show | grep "name   ")

echo "######### FCR  updated $result ######### "

result=$(terraform show | grep fcr_result)

fcrid=$(echo $result | awk -F "=" '{print $2}' | xargs)
```

## Create Connection
```
echo "=============Create FCR2port Connection ==========="

cd ../fcr2port

terraform init

terraform apply -auto-approve -var="fcr_uuid=$fcrid"

echo "=============GET Connection ==========="

terraform show

result=$(terraform show | grep equinix_status)

echo "######### Connection $result #########"

result=$(terraform show | grep connection_result)

con_id=$(echo $result | awk -F "=" '{print $2}' | xargs)
```

## Configure Routing Protocol(RP) direct IP
```
echo "=============Config Direct RP on FCR2port Connection ==========="

cd ../routing-protocol-direct

terraform init

terraform apply -auto-approve -var="connection_uuid=$con_id"

echo "=============GET Connection Direct RP ==========="

terraform show
```

## Update IP address
```
echo "=============Update Direct RP on FCR2port Connection ==========="

terraform apply -refresh -auto-approve -var="connection_uuid=$con_id" -var="equinix_ipv4_ip="190.1.1.1/30"" -var="equinix_ipv6_ip="190::1:1/126""

terraform show

result=$(terraform show | grep equinix_iface_ip)

echo "######### RP $result #########"
```
## Configure BGP

```echo "=============Config BGP on FCR2port Connection ==========="

cd ../routing-protocol-bgp

terraform init

terraform apply -auto-approve -var="connection_uuid=$con_id"

echo "=============GET Connection BGP ==========="

terraform show

result=$(terraform show | grep customer_peer_ip)

echo "######### BGP $result #########" 
```

## Check connection is PROVISIONED
```
echo "=============GET Connection Status ==========="

cd ../fcr2port

terraform apply -refresh-only  -auto-approve

terraform show

result=$(terraform show | grep equinix_status)

echo "######### Connection $result #########"
```
## Update connection Bandwidth
```
echo "=============Update Connection BW ==========="

terraform apply -refresh -auto-approve -var="bandwidth=50" -var="fcr_uuid=$fcrid"

terraform show

result=$(terraform show | grep " bandwidth  =")

echo "######### Connection $result #########"
```
## Delete connection
```
echo "=============Delete Connection ==========="

terraform destroy

terraform show
```
## Delete FCR
```
echo "=============Delete FCR ==========="

cd ../FCR

terraform destroy -auto-approve

terraform show

echo "DONE" 
```



