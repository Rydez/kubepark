# Attractions 🎉

## 🚀 Launch

Each attraction requires a deployment with a container command specifying which attraction to launch. Each attraction directory specifies its command. You'll also need a service so that other components can make HTTP requests to the attaction deployments.

## 🔧 Configuration

All attractions can be configured with the following arguments:

- `--closed`: Temporarily close the attraction (default: false)
- `--fee`: Set a custom entrance fee (default: $5)
- `--park-url`: Specify the kubepark service URL (default: http://kubepark:80)

## 📊 Metrics

All attractions expose Prometheus metrics at `/metrics` on port 9000:

- `revenue`: Total money earned from rides
- `fee`: Current entrance fee
- `is_closed`: Attraction status (0=open, 1=closed)
- `attempts`: Guest interaction attempts with labels:
  - `success`: true/false
  - `reason`: Detailed explanation of the outcome

## 🪵 Logging

Logs can be found in the default location for a docker container.
