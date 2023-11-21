package main

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
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/micovery/apigee-key-expiration/pkg/key_expiration"
	"google.golang.org/api/apigee/v1"
	"io"
	"net/http"
	"time"
)

func main() {
	e := echo.New()
	e.POST("/", func(c echo.Context) error {
		fmt.Printf("Received new request ...\n")

		var jsonBody []byte
		var err error

		if jsonBody, err = io.ReadAll(c.Request().Body); err != nil {
			panic(fmt.Errorf("could read request body"))
		}

		var method string
		if method, err = key_expiration.DetectMethod(jsonBody); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		fmt.Printf("Detected method: %s\n", method)

		switch method {
		case key_expiration.CreateAppMethod:
			fallthrough
		case key_expiration.CreateAppKeyMethod:
			fallthrough
		case key_expiration.UpdateAppMethod:
			fallthrough
		case key_expiration.UpdateAppKeyMethod:
			var path string
			var app *apigee.GoogleCloudApigeeV1DeveloperApp
			if path, app, err = key_expiration.GetApigeeDeveloperApp(method, jsonBody); err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			fmt.Printf("Found app: %s\n", app.Name)

			var updatedKeys []string
			time.Sleep(5 * time.Second) //sleep for some time to allow UI to complete interactions
			if updatedKeys, err = key_expiration.UpdateAPIKeyExpiration(path, app); err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			fmt.Printf("Updated keys: %v\n", updatedKeys)

		default:
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("method %s not supported", method))
		}

		resBytes, _ := json.Marshal(map[string]string{
			"message": "complete",
		})

		res := string(resBytes)
		return c.String(http.StatusOK, res)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
