

# How to run

Run this once for creating the secret with credentials and GCP service account keyfile

```bash
curl https://raw.githubusercontent.com/JorritSalverda/evohome-bigquery-exporter/master/k8s/secret.yaml | EVOHOME_USERNAME='<base64 encoded username>' EVOHOME_PASSWORD='<base64 encoded password>' GCP_SERVICE_ACCOUNT_KEYFILE='<base64 encoded service account keyfile>' envsubst \$EVOHOME_USERNAME,\$EVOHOME_PASSWORD,\$GCP_SERVICE_ACCOUNT_KEYFILE | kubectl apply -f -
```

Make sure to base64 encode your username in the following way to avoid newlines to be included and dollar signs not to be expanded

```bash
echo -n '<username>' | base64 -w0
echo -n '<password>' | base64 -w0
```

The service account keyfile can include newlines, since it's mounted as a file; so encode it using

```bash
cat keyfile.json | base64
```

To set up rbac (role-based access control) permissions run

```bash
curl https://raw.githubusercontent.com/JorritSalverda/evohome-bigquery-exporter/master/k8s/rbac.yaml | kubectl apply -f -
```

In order to configure the application run

```bash
curl https://raw.githubusercontent.com/JorritSalverda/evohome-bigquery-exporter/master/k8s/configmap.yaml | BQ_PROJECT_ID='gcp-project-id' BQ_DATASET='my-dataset' BQ_TABLE='my-table' OUTDOOR_ZONE_NAME='outside' envsubst \$BQ_PROJECT_ID,\$BQ_DATASET,\$BQ_TABLE,\$OUTDOOR_ZONE_NAME | kubectl apply -f -
```

And for deploying a new version or changing the schedule run

```bash
curl https://raw.githubusercontent.com/JorritSalverda/evohome-bigquery-exporter/master/k8s/cronjob.yaml | SCHEDULE='*/1 * * * *' CONTAINER_TAG='0.1.19' envsubst \$SCHEDULE,\$CONTAINER_TAG | kubectl apply -f -
```