# park 🎢

The park is the central service that orchestrates your entire amusement park simulation. This service manages park operations and finances, coordinates guest jobs, handles attraction registration and discovery, provides park-level logs and metrics, and controls park operating hours.

## 🚀 Launch

The park requires a deployment. No need to set a container command, we want the default one. You'll also need a service so that other components can make HTTP requests to the park deployment.

## 🔧 Configuration

kubepark can be configured with the following arguments:

- `--closed`: Temporarily close the park (default: false)
- `--entry-fee`: Set a custom entry fee (default: $20)
- `--open-time`: Park opening hour (default: 9)
- `--close-time`: Park closing hour (default: 21)
- `--metrics-port`: Port for Prometheus metrics (default: 9000)

## 📊 Metrics

kubepark exposes Prometheus metrics at `/metrics` on port 9000:

- `park_revenue`: Total money earned from entry fees
- `park_entry_fee`: Current entry fee
- `park_is_closed`: Park status (0=open, 1=closed)
- `park_guests`: Number of guests in the park
- `park_attractions`: Number of registered attractions
- `park_attempts`: Guest interaction attempts with labels:
  - `success`: true/false
  - `reason`: Detailed explanation of the outcome

## 🪵 Logging

Logs can be found in the default location for a docker container.
