
## Installation

To install this application using Helm run the following commands: 

```bash
helm repo add jorritsalverda https://helm.jorritsalverda.com
kubectl create namespace evohome-bigquery-exporter 

cat << EOF | helm upgrade evohome-bigquery-exporter jorritsalverda/evohome-bigquery-exporter --install --namespace evohome-bigquery-exporter --values -
config:
  bqProjectID: your-gcp-project-id
  bqDataset: your-dataset
  bqTable: your-table
  outdoorZoneName: outside

secret:
  evohomeUsername: yourusername
  evohomePassword: yourpassword
  gcpServiceAccountKeyfile: {}
EOF
```