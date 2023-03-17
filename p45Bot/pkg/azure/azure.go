package azure

import (
	"context"
	"errors"
	"fmt"
	"log"

	"calfinn.io/p45bot/pkg/opts"
	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
	"github.com/microsoftgraph/msgraph-sdk-go/users"
)

func NewAZClient(tenant, clientid, clientsecret string) (*msgraphsdk.GraphServiceClient, error) {
	cred, err := getCred(tenant, clientid, clientsecret)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	client, err := azClientCreds(cred)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return client, nil
}
func getCred(tenant, clientid, clientsecret string) (*azidentity.ClientSecretCredential, error) {
	cred, err := azidentity.NewClientSecretCredential(
		tenant,
		clientid,
		clientsecret,
		nil,
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return cred, nil
}

func azClientCreds(cred *azidentity.ClientSecretCredential) (*msgraphsdk.GraphServiceClient, error) {
	client, err := msgraphsdk.NewGraphServiceClientWithCredentials(cred, []string{"https://graph.microsoft.com/.default"})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return client, nil
}

func printOdataError(err error) {
	switch err.(type) {
	case *odataerrors.ODataError:
		typed := err.(*odataerrors.ODataError)
		fmt.Printf("error:%s", typed.Error())
		if terr := typed.GetError(); terr != nil {
			fmt.Printf("code: %s", *terr.GetCode())
			fmt.Printf("msg: %s", *terr.GetMessage())
		}
	default:
		fmt.Printf("%T > error: %#v", err, err)
	}
}

// Retrieve a user by UPN
func GetUserByUPN(client *msgraphsdk.GraphServiceClient, upn string) (models.UserCollectionResponseable, error) {
	// Create a new Graph API request to retrieve the user by UPN
	//req := client.Users().Request().Filter(fmt.Sprintf("userPrincipalName eq '%s'", upn))
	requestFilter := fmt.Sprintf("(userPrincipalName eq '%s')", upn)
	//fmt.Printf("%+v\n", requestFilter)
	if opts.GetVerbose() {
		log.Printf("Search filter being used is: %s", requestFilter)
	}
	count := true
	requestParameters := &users.UsersRequestBuilderGetQueryParameters{
		Count:  &count,
		Filter: &requestFilter,
	}
	//fmt.Printf("%+v\n", requestParameters)
	//RequestOptions := arg.QueryRequestOptions {
	//    ResultFormat: "objectArray",
	//}
	configuration := &users.UsersRequestBuilderGetRequestConfiguration{
		QueryParameters: requestParameters,
	}
	//fmt.Printf("%+v\n", configuration)
	result, err := client.Users().Get(context.Background(), configuration)
	if err != nil {
		fmt.Printf("Error getting the drive: %v\n", err)
		printOdataError(err)
		log.Println(result)
		log.Println("Error with result from Graph API")
		log.Println(err.Error())
		return nil, err
	}
	return result, err
}

func CheckUserExists(input models.UserCollectionResponseable) (bool, error) {
	res := len(input.GetValue())
	var exists bool
	var err error
	if res == 0 {
		if opts.GetVerbose() {
			log.Println("User search did not return a hit.  User does not exist.")
		}
		exists = false
	} else if res == 1 {
		if opts.GetVerbose() {
		log.Println("User found")
		}
		exists = true
	} else if res > 1 {
		if opts.GetVerbose() {
		log.Println("Multpile users found with same UPN.")
		}
		err = errors.New("Multpile users found with same UPN.")
		return false, err
	}
	return exists, err
}

/*
	fmt.Println("Dumping result to console")
	res := len(result.GetValue())
	fmt.Println("The value of x is", res)
	user := result.GetValue()[0].(models.Userable)
	fmt.Printf("%s\n", *user.GetDisplayName())
	fmt.Println(user)
	//fmt.Println(string(result))
	r := result.GetValue()[0].(models.Userable)
	//fmt.Println(r)
*/

//func azClientToken(token string) (*msgraphsdk.GraphServiceClient, error) {/
//	cred := azcore.AccessToken(token)
//	client, err := msgraphsdk.NewGraphServiceClientWithCredentials(cred, []string{"Directory.ReadAll"})
//	if err != nil {
//		log.Println(err)
//		return nil, err
//	}
//	return client, nil
//}//
/*
func NewGraphClientWithToken(token string) (*graph.Client, error) {
	// Create a new BearerAccess credential using the existing token
	credential := auth.NewBearerAccess(token)

	// Create a new Graph API client using the BearerAccess credential
	client := graph.NewClient(credential)

	// Test the client by fetching the current user's profile
	me, err := client.Me().Request().Get(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get current user: %v", err)
	}
	fmt.Printf("Signed in as: %s\n", *me.DisplayName)

	return client, nil
}


	client, err := msgraphsdk.NewGraphServiceClientWithCredentials(cred, []string{"User.Read"})
	if err != nil {
		log.Fatal()
	}

	result, err := client.Me().Get(nil)
	fmt.Printf(result)
}

//https://graph.microsoft.com/v1.0/users/MeganB@M365x214355.onmicrosoft.com

//token := os.Getenv("TOKEN_ENV_VAR_NAME")
token := "@""
client := &http.Client{}
req, err := http.NewRequest("GET", "https://graph.microsoft.com/v1.0/users/tim..io", nil)
if err != nil {
    // handle error
}
req.Header.Set("Authorization", "Bearer "+token)

resp, err := client.Do(req)
if err != nil {
    // handle error
}
defer resp.Body.Close()

body, err := ioutil.ReadAll(resp.Body)
if err != nil {
    // handle error
}

fmt.Println(string(body))
*/
