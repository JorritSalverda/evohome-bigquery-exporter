

# How to run

```bash
cat kubernetes.yaml | CONTAINER_TAG=0.1.2 EVOHOME_USERNAME=<base64 encoded username> EVOHOME_PASSWORD=<base64 encoded password> envsubst \$CONTAINER_TAG,\$EVOHOME_USERNAME,\$EVOHOME_PASSWORD | kubectl apply -f -
```