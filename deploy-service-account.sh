#!/bin/bash

# Copyright 2023 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e
export PROJECT_ID="$(gcloud config get project)"

if [ -z "${PROJECT_ID}" ] ; then
  echo "No project detected. Run: gcloud config set project ..."
  exit 1
fi

gcloud iam service-accounts create apigee-key-expiration --project "${PROJECT_ID}"

gcloud projects add-iam-policy-binding "${PROJECT_ID}" \
  --member="serviceAccount:apigee-key-expiration@${PROJECT_ID}.iam.gserviceaccount.com" \
  --role=roles/eventarc.eventReceiver

gcloud projects add-iam-policy-binding "${PROJECT_ID}" \
  --member="serviceAccount:apigee-key-expiration@${PROJECT_ID}.iam.gserviceaccount.com" \
  --role=roles/run.invoker

gcloud projects add-iam-policy-binding "${PROJECT_ID}" \
    --member="serviceAccount:apigee-key-expiration@${PROJECT_ID}.iam.gserviceaccount.com" \
    --role=roles/apigee.readOnlyAdmin


gcloud projects add-iam-policy-binding "${PROJECT_ID}" \
    --member="serviceAccount:apigee-key-expiration@${PROJECT_ID}.iam.gserviceaccount.com" \
    --role=roles/apigee.developerAdmin






