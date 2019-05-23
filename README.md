

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

For deploying a new version or changing the schedule run

```bash
curl https://raw.githubusercontent.com/JorritSalverda/evohome-bigquery-exporter/master/k8s/cronjob.yaml | SCHEDULE='*/5 * * * *' CONTAINER_TAG='0.1.7' envsubst \$SCHEDULE,\$CONTAINER_TAG | kubectl apply -f -
```