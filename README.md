## Apigee Update Developer App Expiration

This repo contains a sample Cloud Run application that can be used for auto updating the
expiration date for Application Credentials in Apigee X.

The Cloud Run execution is triggered (through GCP EventArc) whenever a Developer application is created or updated.

In Apigee, once an Application Credential has been created, it's not possible to change the expiration date.
However, it is possible to delete the existing credential, and recreate it with the desired expiration date.
That's what this Cloud Run application does. The steps are as follows:

1. Delete existing Application API Credential
2. Create new Application API Credential with same client_id, client_secret, and status 
3. Update the API Products list on the new Application Credential (to match original state)
4. Update the status for each API Product as approved / revoked / pending (to match original state)


## Deploy the Cloud Run Application


First, let's create a GCP Service account that will be used for the Cloud Run execution. 
This service account has access to the Apigee Administrator API

```shell
export PROJECT_ID="YOUR_GPC_PROJECT"
./deploy-service-acccount.sh
```

Then, let's deploy the Cloud Run application

```shell
export PROJECT_ID="YOUR_GPC_PROJECT"
export REGION="us-west1"
export EXPIRE_IN_SECONDS=31536000
./deploy-cloud-run.sh
```

Finally, let's create the Event-Arc triggers that are used for detecting when apps are created / update

```shell
export PROJECT_ID="YOUR_GPC_PROJECT"
export REGION="us-west1"
./deploy-cloud-run-triggers.sh
```


### Not Google Product Clause

This is not an officially supported Google product.