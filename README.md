# kubepark 🎢

Welcome to kubepark, an innovative amusement park simulation game that runs on Kubernetes! Build, manage, and optimize your virtual amusement park by deploying various attractions, monitoring guest behavior, and maximizing profits.

## 🎯 Overview

kubepark is a unique simulation game that combines the excitement of amusement park management with the power of Kubernetes orchestration. As a park operator, you'll:

- Deploy and manage various attractions (rides, amenities, concessions)
- Monitor guest satisfaction and park performance
- Handle maintenance and repairs
- Track revenue and expenses

## 🎮 How It Works

The simulation consists of three main components:

1. **kubepark Service**: The central park management system that:

   - Tracks park finances and attraction inventory
   - Manages the park's operating hours
   - Creates guest jobs
   - Handles attraction registration and discovery
   - Exposes metrics about the park

2. **Attractions**: Kubernetes services that represent different park features:

   - Each attraction has its own HTTP endpoint for guest interaction
   - Attractions can be configured with custom entrance fees
   - Maintenance costs apply when attractions break down
   - All attractions expose metrics and logs for monitoring

3. **Guests**: Kubernetes jobs that simulate park visitors:
   - Guests enter the park and explore attractions
   - They make decisions based on available money and attraction fees
   - Each guest reports metrics about their experience
   - Guests leave when they run out of money or the park closes

## 🛠️ Technical Details

### Monitoring and Observability

Every component in kubepark exposes:

- Prometheus metrics at `/metrics` on port `9000` about revenue, usage, and guest satisfaction
- Logs at default Docker container location

## 🚀 Getting Started

1. Head to the [kubepark](./kubepark) directory to launch the park simulator
2. Explore the [attractions](./attractions) directory to deploy rides and amenities
3. Monitor your park's performance using tools like Grafana
4. Optimize your park based on insights from logs and metrics

## 🔒 Remember

This is a simulation game meant to be played through Kubernetes orchestration. Avoid direct HTTP requests or data manipulation - let the system work as designed!
