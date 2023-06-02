#--------create FG
cd fabric-gateway/
echo "=============Initialize FG Resource==========="
terraform init
echo "=============Create FG ==========="
terraform apply -auto-approve -var="fg_name=terra_fg-test"
echo "=============GET FG ==========="
terraform show
result=$(terraform show | grep state)
echo "######### FG $result ##########"

#--------update FG name
echo "=============Update FG Name==========="
terraform apply -refresh -auto-approve -var="fg_name=terra_fg-test-update"
terraform show
result=$(terraform show | grep "name   ")
echo "######### FG  updated $result ######### "

result=$(terraform show | grep fg_result)
fgid=$(echo $result | awk -F "=" '{print $2}' | xargs)

#--------create connection
echo "=============Create FG2port Connection ==========="
cd ../fg2port
terraform init
terraform apply -auto-approve -var="fg_uuid=$fgid"
echo "=============GET Connection ==========="
terraform show
result=$(terraform show | grep equinix_status)
echo "######### Connection $result #########"

result=$(terraform show | grep connection_result)
con_id=$(echo $result | awk -F "=" '{print $2}' | xargs)

#--------config RP direct IP
echo "=============Config Direct RP on FG2port Connection ==========="
cd ../routing-protocol-direct
terraform init
terraform apply -auto-approve -var="connection_uuid=$con_id"
echo "=============GET Connection Direct RP ==========="
terraform show

#--------Update IP address
echo "=============Update Direct RP on FG2port Connection ==========="
terraform apply -refresh -auto-approve -var="connection_uuid=$con_id" -var="equinix_ipv4_ip="190.1.1.1/30"" -var="equinix_ipv6_ip="190::1:1/126""
terraform show
result=$(terraform show | grep equinix_iface_ip)
echo "######### RP $result #########"

#--------config BGP
echo "=============Config BGP on FG2port Connection ==========="
cd ../routing-protocol-bgp
terraform init
terraform apply -auto-approve -var="connection_uuid=$con_id"

echo "=============GET Connection BGP ==========="
terraform show
result=$(terraform show | grep customer_peer_ip)
echo "######### BGP $result #########"


#--------Check connection is PROVISIONED
echo "=============GET Connection Status ==========="
cd ../fg2port
terraform apply -refresh-only  -auto-approve
terraform show
result=$(terraform show | grep equinix_status)
echo "######### Connection $result #########"

#--------Update connection BW
echo "=============Update Connection BW ==========="
terraform apply -refresh -auto-approve -var="bandwidth=50" -var="fg_uuid=$fgid"
terraform show
result=$(terraform show | grep " bandwidth  =")
echo "######### Connection $result #########"

#--------delete connection
echo "=============Delete Connection ==========="
terraform destroy
terraform show

#--------delete FG
echo "=============Delete FG ==========="
cd ../fabric-gateway
terraform destroy -auto-approve
terraform show

echo "DONE"

