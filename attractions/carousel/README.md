# Carousel 🎠

```
namespace: rides
command: carousel
```

## 🎯 Overview

The classic carousel (merry-go-round) is a timeless attraction that delights guests of all ages. This gently spinning ride features:

- Beautifully themed horses and other decorative seats
- Soothing music that creates a magical atmosphere
- Smooth up-and-down motion for a gentle thrill
- Perfect for families and those seeking a relaxing experience

## 💰 Pricing

- Default fee: $5 per ride
- Duration: 3 seconds per ride

## 🔧 Configuration

The carousel can be configured with the following arguments:

- `--closed`: Temporarily close the attraction (default: false)
- `--fee`: Set a custom entrance fee (default: $5)
- `--park-url`: Specify the kubepark service URL (default: http://kubepark:80)

## 📊 Metrics

The carousel exposes Prometheus metrics at `/metrics` on port 9000:

- `revenue`: Total money earned from rides
- `fee`: Current entrance fee
- `is_closed`: Attraction status (0=open, 1=closed)
- `attempts`: Guest interaction attempts with labels:
  - `success`: true/false
  - `reason`: Detailed explanation of the outcome
