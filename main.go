package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/subscription/mgmt/2020-09-01/subscription"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/rancher/machine/drivers/azure/azureutil"
	"github.com/sirupsen/logrus"
)

type Capabilities struct {
	SubscriptionID string
	TenantID       string
	ClientID       string
	ClientSecret   string
	AuthBaseURL    string
	BaseURL        string
	Environment    string
}

var (
	//	ctx  = context.Background()
	cred Capabilities
)

func main() {
	credFile := os.Getenv("AZURE_AUTH_PATH")
	if credFile == "" {
		credFile = "~/aks-credentials.json"
	}
	authInfo, err := readJSON(credFile)
	if err != nil {
		logrus.Fatalf("Failed to read JSON: %+v", err)
	}
	cred.SubscriptionID = (*authInfo)["subscriptionId"].(string)
	cred.TenantID = (*authInfo)["tenantId"].(string)
	cred.ClientID = (*authInfo)["clientId"].(string)
	cred.ClientSecret = (*authInfo)["clientSecret"].(string)
	cred.AuthBaseURL = (*authInfo)["authBaseUrl"].(string)
	cred.BaseURL = (*authInfo)["baseUrl"].(string)
	cred.Environment = (*authInfo)["environment"].(string)

	clientEnvironment := ""
	if cred.Environment != "" {
		clientEnvironment = cred.Environment
	}
	azureEnvironment := GetEnvironment(clientEnvironment)
	logrus.Infof("show azureEnvironment: %v", azureEnvironment)

	cred.BaseURL = azureEnvironment.ResourceManagerEndpoint
	cred.AuthBaseURL = azureEnvironment.ActiveDirectoryEndpoint

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if cred.TenantID == "" {
		cred.TenantID, err = azureutil.FindTenantID(ctx, azureEnvironment, cred.SubscriptionID)
		if err != nil {
			logrus.Errorf("could not find tenant ID for Azure environment %v: %v", azureEnvironment.Name, err)
		}
	}

	logrus.Infof("show credentials: %v", cred)
	client, err := NewSubscriptionServiceClient(&cred)
	if err != nil {
		logrus.Errorf("[AKS] failed to create new subscription client: %v", err)
	}

	subscriptionList, err := client.List(ctx)
	if err != nil {
		logrus.Errorf("[AKS] failed to list subscription details: %v", err)
	}
	logrus.Infof("show subscriptionList: %v", subscriptionList.Values())

	for _, s := range subscriptionList.Values() {
		logrus.Infof("Subscription ID %s Name %v", *s.ID, *s.DisplayName)
	}

	sub, err := client.Get(ctx, cred.SubscriptionID)
	logrus.Infof("show subscription details %v ", *sub.DisplayName)
	if err != nil {
		logrus.Errorf("[AKS] failed to get subscription details: %v", err)
	}

}

func NewSubscriptionServiceClient(cap *Capabilities) (*subscription.SubscriptionsClient, error) {
	authorizer, err := NewAzureClientAuthorizer(cap)
	if err != nil {
		return nil, err
	}

	subscriptionService := subscription.NewSubscriptionsClient()
	subscriptionService.Authorizer = authorizer

	return &subscriptionService, nil
}

func NewAzureClientAuthorizer(cap *Capabilities) (autorest.Authorizer, error) {
	oauthConfig, err := adal.NewOAuthConfig(cap.AuthBaseURL, cap.TenantID)
	if err != nil {
		return nil, err
	}

	spToken, err := adal.NewServicePrincipalToken(*oauthConfig, cap.ClientID, cap.ClientSecret, cap.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("couldn't authenticate to Azure cloud with error: %v", err)
	}

	return autorest.NewBearerAuthorizer(spToken), nil
}

func GetEnvironment(env string) azure.Environment {
	switch env {
	case "AzureGermanCloud":
		return azure.GermanCloud
	case "AzureChinaCloud":
		return azure.ChinaCloud
	case "AzureUSGovernmentCloud":
		return azure.USGovernmentCloud
	default:
		return azure.PublicCloud
	}
}

func readJSON(path string) (*map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("[AKS] failed to read file: %v", err)
	}
	contents := make(map[string]interface{})
	_ = json.Unmarshal(data, &contents)
	return &contents, nil
}
