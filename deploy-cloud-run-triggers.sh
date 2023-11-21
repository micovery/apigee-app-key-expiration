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


export REGION=${REGION:-us-west1}

gcloud services enable \
  --quiet \
  eventarc.googleapis.com \
  artifactregistry.googleapis.com \
  run.googleapis.com \
  cloudbuild.googleapis.com \
  artifactregistry.googleapis.com \
  secretmanager.googleapis.com \
  --project "${PROJECT_ID}"

echo "Creating App creation trigger ..."
gcloud eventarc triggers create apigee-key-expiration-trigger-create-app \
 --location=global \
 --service-account="apigee-key-expiration@${PROJECT_ID}.iam.gserviceaccount.com" \
 --destination-run-service=apigee-key-expiration \
 --destination-run-region="${REGION}" \
 --destination-run-path="/" \
 --event-filters="type=google.cloud.audit.log.v1.written" \
 --event-filters="serviceName=apigee.googleapis.com" \
 --event-filters="methodName=google.cloud.apigee.v1.DeveloperApps.CreateDeveloperApp"
echo ""

echo "Creating App Key creation trigger ..."
gcloud eventarc triggers create apigee-key-expiration-trigger-create-app-key \
 --location=global \
 --service-account="apigee-key-expiration@${PROJECT_ID}.iam.gserviceaccount.com" \
 --destination-run-service=apigee-key-expiration \
 --destination-run-region="${REGION}" \
 --destination-run-path="/" \
 --event-filters="type=google.cloud.audit.log.v1.written" \
 --event-filters="serviceName=apigee.googleapis.com" \
 --event-filters="methodName=google.cloud.apigee.v1.DeveloperAppKeys.CreateDeveloperAppKey"
echo ""

echo "Creating App update trigger ..."
gcloud eventarc triggers create apigee-key-expiration-trigger-update-app \
 --location=global \
 --service-account="apigee-key-expiration@${PROJECT_ID}.iam.gserviceaccount.com" \
 --destination-run-service=apigee-key-expiration \
 --destination-run-region="${REGION}" \
 --destination-run-path="/" \
 --event-filters="type=google.cloud.audit.log.v1.written" \
 --event-filters="serviceName=apigee.googleapis.com" \
 --event-filters="methodName=google.cloud.apigee.v1.DeveloperApps.UpdateDeveloperApp"
echo ""

echo "Creating App Key update trigger ..."
 gcloud eventarc triggers create apigee-key-expiration-trigger-update-app-key \
  --location=global \
  --service-account="apigee-key-expiration@${PROJECT_ID}.iam.gserviceaccount.com" \
  --destination-run-service=apigee-key-expiration \
  --destination-run-region="${REGION}" \
  --destination-run-path="/" \
  --event-filters="type=google.cloud.audit.log.v1.written" \
  --event-filters="serviceName=apigee.googleapis.com" \
  --event-filters="methodName=google.cloud.apigee.v1.DeveloperAppKeys.UpdateDeveloperAppKey"
echo ""
