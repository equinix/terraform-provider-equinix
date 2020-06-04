module terraform-provider-equinix

go 1.14

require (
	ecx-go-client/v3 v3.0.0
	github.com/hashicorp/go-getter v1.4.2-0.20200106182914-9813cbd4eb02 // indirect
	github.com/hashicorp/hcl/v2 v2.3.0 // indirect
	github.com/hashicorp/terraform-config-inspect v0.0.0-20191212124732-c6ae6269b9d7 // indirect
	github.com/hashicorp/terraform-plugin-sdk v1.10.0
	github.com/stretchr/testify v1.6.0
	oauth2-go v1.0.0
)

replace oauth2-go v1.0.0 => ../oauth2-go

replace ecx-go-client/v3 v3.0.0 => ../ecx-go-client
