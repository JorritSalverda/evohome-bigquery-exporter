
## Installation

To install this application using Helm run the following commands: 

```bash
helm repo add jorritsalverda https://helm.jorritsalverda.com
kubectl create namespace evohome-bigquery-exporter 

helm upgrade \
  evohome-bigquery-exporter \
  jorritsalverda/evohome-bigquery-exporter \
  --install \
  --namespace evohome-bigquery-exporter \
  --set config.bqProjectID=your-project-id \
  --set config.bqDataset=your-dataset \
  --set config.bqTable=your-table \
  --set config.outdoorZoneName=outside \
  --set secret.evohomeUsername=yourusername \
  --set secret.evohomePassword=yourpassword \
  --set secret.gcpServiceAccountKeyfile='{abc: blabla}' \
  --wait
```

If you later on want to upgrade without specifying all values again you can use

```bash
helm upgrade \
  evohome-bigquery-exporter \
  jorritsalverda/evohome-bigquery-exporter \
  --install \
  --namespace evohome-bigquery-exporter \
  --reuse-values \
  --set cronjob.schedule='*/1 * * * *' \
  --wait
```