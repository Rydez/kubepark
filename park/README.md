# kubepark 🎢

kubepark is the central service that orchestrates your entire amusement park simulation. This powerful management system:

- Manages park operations and finances
- Coordinates guest jobs
- Handles attraction registration and discovery
- Provides comprehensive park metrics
- Controls park operating hours

## 🔧 Configuration

kubepark can be configured with the following arguments:

- `--closed`: Temporarily close the park (default: false)
- `--entry-fee`: Set a custom entry fee (default: $20)
- `--open-time`: Park opening hour (default: 9)
- `--close-time`: Park closing hour (default: 21)
- `--metrics-port`: Port for Prometheus metrics (default: 9000)

## 📊 Metrics

kubepark exposes Prometheus metrics at `/metrics` on port 9000:

- `revenue`: Total money earned from entry fees
- `entry_fee`: Current entry fee
- `is_closed`: Park status (0=open, 1=closed)
- `guests`: Number of guests in the park
- `attractions`: Number of registered attractions
- `attempts`: Guest interaction attempts with labels:
  - `success`: true/false
  - `reason`: Detailed explanation of the outcome

## 📝 Logging

kubepark provides detailed logs about:

- Park status changes
- Guest entry and exit events
- Attraction registration
- Financial transactions
- System operations
