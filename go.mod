module github.com/rancher/aks-operator

go 1.20

replace (
	github.com/rancher/wrangler => github.com/rancher/wrangler v1.1.1-0.20230831050635-df1bd5aae9df
	k8s.io/client-go => k8s.io/client-go v0.27.5
)

require (
	github.com/Azure/azure-sdk-for-go v55.7.0+incompatible
	github.com/Azure/go-autorest/autorest v0.11.29
	github.com/Azure/go-autorest/autorest/adal v0.9.23
	github.com/rancher/machine v0.15.0-rancher103
	github.com/sirupsen/logrus v1.9.3
)

require (
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest/azure/auth v0.4.2 // indirect
	github.com/Azure/go-autorest/autorest/azure/cli v0.3.1 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/to v0.4.1-0.20210111195520-9fc88b15294e // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.2-0.20210111195520-9fc88b15294e // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/dimchansky/utfbom v1.1.0 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	golang.org/x/crypto v0.13.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
)
