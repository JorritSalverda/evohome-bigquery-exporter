
## Installation

To install this application using Helm run the following commands: 

```bash
helm repo add jorritsalverda https://helm.jorritsalverda.com
kubectl create namespace evohome-bigquery-exporter 

helm upgrade \
  evohome-bigquery-exporter-test \
  jorritsalverda/evohome-bigquery-exporter \
  --install \
  --namespace evohome-bigquery-exporter-test \
  --set config.bqProjectID=your-project-id \
  --set config.bqDataset=your-dataset \
  --set config.bqTable=your-table \
  --set config.outdoorZoneName=outside \
  --set secret.evohomeUsername=yourusername \
  --set secret.evohomePassword=yourpassword \
  --set secret.gcpServiceAccountKeyfile='{abc: blabla}' \
  --wait
```