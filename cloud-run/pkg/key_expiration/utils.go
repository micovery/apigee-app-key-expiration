package key_expiration

// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/api/apigee/v1"
	"google.golang.org/api/googleapi"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func UpdateAPIKeyExpiration(appPath string, app *apigee.GoogleCloudApigeeV1DeveloperApp) ([]string, error) {
	ctx := context.Background()

	var err error

	var updatedKeys []string
	var apigeeService *apigee.Service
	if apigeeService, err = apigee.NewService(ctx); err != nil {
		return nil, fmt.Errorf("could not create Apigee service. error: %s", err.Error())
	}

	for _, credential := range app.Credentials {
		if credential.ExpiresAt > 0 {
			continue
		}

		//Delete the existing API Key
		keyPath := fmt.Sprintf("%s/keys/%s", appPath, credential.ConsumerKey)
		if _, err = apigeeService.Organizations.Developers.Apps.Keys.Delete(keyPath).Do(); err != nil {
			//ignore deletion error
			fmt.Sprintf("could not delete API key %s******. error: %s\n", credential.ConsumerKey[:7], err.Error())
		}

		var expireInSeconds int64
		expireInSecondsEnv := os.Getenv("EXPIRE_IN_SECONDS")
		if expireInSeconds, err = strconv.ParseInt(expireInSecondsEnv, 10, 64); err != nil {
			expireInSeconds = int64(60 * 60 * 24 * 365)
			fmt.Printf("could not parse EXPIRE_IN_SECONDS (%s) as integer. defaulting to: %v\n", expireInSecondsEnv, expireInSeconds)
		}

		//Create a new API Key with same client_id and client_secret, and expiration date
		newAPIKey := &apigee.GoogleCloudApigeeV1DeveloperAppKey{
			ConsumerKey:      credential.ConsumerKey,
			ConsumerSecret:   credential.ConsumerSecret,
			ExpiresInSeconds: expireInSeconds,
			Status:           credential.Status,
		}

		if _, err = apigeeService.Organizations.Developers.Apps.Keys.Create(appPath, newAPIKey).Do(); err != nil {
			fmt.Printf("could not import API key %s******. error: %s\n", credential.ConsumerKey[:7], err.Error())
			continue
		}

		if len(credential.ApiProducts) > 0 {
			//Update the API products on the new API key
			var apiProducts []interface{}
			for _, apiProduct := range credential.ApiProducts {
				apiProducts = append(apiProducts, apiProduct.Apiproduct)
			}

			replaceAPIKey := &apigee.GoogleCloudApigeeV1DeveloperAppKey{
				Scopes:      credential.Scopes,
				ApiProducts: apiProducts,
			}

			if _, err = apigeeService.Organizations.Developers.Apps.Keys.ReplaceDeveloperAppKey(keyPath, replaceAPIKey).Do(); err != nil {
				fmt.Printf("could not replace API key %s******. error: %s\n", credential.ConsumerKey[:7], err.Error())
				continue
			}

			//update the status for each API product
			for _, apiProduct := range credential.ApiProducts {
				productPath := fmt.Sprintf("%s/apiproducts/%s", keyPath, apiProduct.Apiproduct)

				if apiProduct.Status == "approved" || apiProduct.Status == "revoked" {
					action := apiProduct.Status[:len(apiProduct.Status)-1]
					if _, err = apigeeService.Organizations.Developers.Apps.Keys.Apiproducts.UpdateDeveloperAppKeyApiProduct(productPath).Do(googleapi.QueryParameter("action", action)); err != nil {
						fmt.Printf("could not %s API key %s for product %s. error: %s\n", action, credential.ConsumerKey[:7], apiProduct.Apiproduct, err.Error())
					}
				}
			}
		}

		fmt.Printf("updated expiration for key %s to %v\n", credential.ConsumerKey[:7], expireInSeconds)
		updatedKeys = append(updatedKeys, credential.ConsumerKey[:7])
	}

	return updatedKeys, nil
}

func DetectMethod(jsonBody []byte) (string, error) {
	info := EventInfo{}
	json.Unmarshal(jsonBody, &info)

	methodName := info.MethodName()
	if methodName == "" {
		return "", fmt.Errorf("could not detect operation")
	}

	parts := strings.Split(methodName, ".")

	return parts[len(parts)-1], nil
}

func GetApigeeDeveloperApp(method string, jsonBody []byte) (string, *apigee.GoogleCloudApigeeV1DeveloperApp, error) {
	ctx := context.Background()
	var err error
	var path string
	if path = getAppPath(method, jsonBody); path == "" {
		return "", nil, fmt.Errorf("could not determine app path from event")
	}

	var apigeeService *apigee.Service
	if apigeeService, err = apigee.NewService(ctx); err != nil {
		return "", nil, fmt.Errorf("could not create Apigee service")
	}

	var app *apigee.GoogleCloudApigeeV1DeveloperApp
	if app, err = apigeeService.Organizations.Developers.Apps.Get(path).Do(); err != nil {
		return "", nil, fmt.Errorf("could not find app. %s", err.Error())
	}
	return path, app, nil
}

func getAppPath(method string, jsonBody []byte) string {

	appPath := ""

	if method == CreateAppMethod {
		data := CreateAppEventData{}
		json.Unmarshal(jsonBody, &data)
		appPath = fmt.Sprintf("%s/apps/%s", data.ProtoPayload.Request.Parent, data.ProtoPayload.Request.DeveloperApp.Name)
	} else if method == UpdateAppMethod {
		data := UpdateAppEventData{}
		json.Unmarshal(jsonBody, &data)
		appPath = data.ProtoPayload.Request.Name
	} else if method == CreateAppKeyMethod {
		data := CreateAppKeyEventData{}
		json.Unmarshal(jsonBody, &data)
		appPath = data.ProtoPayload.Request.Parent
	} else if method == UpdateAppKeyMethod {
		data := UpdateAppKeyEventData{}
		json.Unmarshal(jsonBody, &data)
		appPath = data.ProtoPayload.Request.Name
		re := regexp.MustCompile("\\/keys\\/.+$")
		appPath = re.ReplaceAllString(appPath, "")
	}

	fmt.Printf("Detected app path: %s\n", appPath)
	return appPath
}
