# Restroom 🚻

```
namespace: amenities
command: restroom
```

## 🎯 Overview

The restroom is an essential amenity that keeps your guests comfortable and happy. This clean, well-maintained facility:

- Provides necessary comfort for park visitors
- Helps maintain guest satisfaction
- Contributes to longer park stays
- Essential for family-friendly park operations

## 💰 Pricing

- Default fee: $2 per use
- Duration: 2 seconds per use

## 🔧 Configuration

The restroom can be configured with the following arguments:

- `--closed`: Temporarily close the facility (default: false)
- `--fee`: Set a custom usage fee (default: $2)
- `--park-url`: Specify the kubepark service URL (default: http://kubepark:80)

## 📊 Metrics

The restroom exposes Prometheus metrics at `/metrics` on port 9000:

- `revenue`: Total money earned from usage
- `fee`: Current usage fee
- `is_closed`: Facility status (0=open, 1=closed)
- `attempts`: Guest interaction attempts with labels:
  - `success`: true/false
  - `reason`: Detailed explanation of the outcome
