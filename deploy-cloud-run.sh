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

pushd ./cloud-run

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


gcloud run deploy apigee-key-expiration \
 --quiet \
 --execution-environment=gen2 \
 --service-account="apigee-key-expiration@${PROJECT_ID}.iam.gserviceaccount.com" \
 --port 8080 \
 --timeout 3600 \
 --region="${REGION}" \
 --source=.
