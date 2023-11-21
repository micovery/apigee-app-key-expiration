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

// ** EVENT INFO ***

type EventInfo struct {
	ProtoPayload struct {
		MethodName string `json:"methodName"`
	} `json:"protoPayload"`
}

func (e *EventInfo) MethodName() string {
	return e.ProtoPayload.MethodName
}

// *** CREATE APP ***

type CreateAppEventData struct {
	ProtoPayload struct {
		Request  CreateAppRequest  `json:"request"`
		Response CreateAppResponse `json:"response"`
	} `json:"protoPayload"`
}

type CreateAppRequest struct {
	Parent       string       `json:"parent"`
	DeveloperApp DeveloperApp `json:"developerApp"`
}

type CreateAppResponse struct {
	AppId       string      `json:"appId"`
	Name        string      `json:"name"`
	DeveloperId string      `json:"developerId"`
	Attributes  []Attribute `json:"attributes"`
	Status      string      `json:"status"`
}

// *** UPDATE APP ***

type UpdateAppEventData struct {
	ProtoPayload struct {
		Request  UpdateAppRequest  `json:"request"`
		Response UpdateAppResponse `json:"response"`
	} `json:"protoPayload"`
}

type UpdateAppRequest struct {
	Name         string       `json:"name"`
	DeveloperApp DeveloperApp `json:"developerApp"`
}

type UpdateAppResponse struct {
	AppId       string      `json:"appId"`
	Name        string      `json:"name"`
	DeveloperId string      `json:"developerId"`
	Attributes  []Attribute `json:"attributes"`
	Status      string      `json:"status"`
}

// *** CREATE APP KEY ***

type CreateAppKeyEventData struct {
	ProtoPayload struct {
		Request  CreateAppKeyRequest  `json:"request"`
		Response CreateAppKeyResponse `json:"response"`
	} `json:"protoPayload"`
}

type CreateAppKeyRequest struct {
	Parent          string          `json:"parent"`
	DeveloperAppKey DeveloperAppKey `json:"developerAppKey"`
}

type CreateAppKeyResponse struct {
	ConsumerKey string `json:"consumerKey"`
	Status      string `json:"status"`
}

// *** UPDATE APP KEY ***

type UpdateAppKeyEventData struct {
	ProtoPayload struct {
		Request  UpdateAppKeyRequest  `json:"request"`
		Response UpdateAppKeyResponse `json:"response"`
	} `json:"protoPayload"`
}

type UpdateAppKeyRequest struct {
	Name            string          `json:"name"`
	DeveloperAppKey DeveloperAppKey `json:"developerAppKey"`
}

type UpdateAppKeyResponse struct {
	ConsumerKey string `json:"consumerKey"`
	Status      string `json:"status"`
}

// *** GENERIC TYPES ***

type DeveloperApp struct {
	Status      string      `json:"status"`
	Name        string      `json:"name"`
	Attributes  []Attribute `json:"attributes"`
	CallbackUrl string      `json:"callbackUrl"`
}

type Attribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type DeveloperAppKey struct {
	ConsumerKey string `json:"consumerKey"`
}
