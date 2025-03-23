# kubepark 🎢

Welcome to kubepark, an innovative amusement park simulation game that runs on Kubernetes! Build, manage, and optimize your virtual amusement park by deploying various attractions, monitoring guest behavior, and maximizing profits.

- Deploy and manage various attractions (rides, amenities, concessions)
- Monitor guest satisfaction and park performance
- Handle maintenance and repairs
- Track revenue and expenses

## 🎮 How It Works

The simulation consists of three main components:

**Park**: The central park service is responsible for tracking park finances and attraction inventory, managing the park's operating hours, creating guest jobs, handling attraction registration and discovery, and exposing metrics about the park.

**Attractions**: Kubernetes services that represent different park features. Each attraction has its own HTTP endpoint for guest interaction, and can be configured with custom entrance fees. Maintenance costs apply when attractions break down, and all attractions expose metrics and logs for monitoring.

**Guests**: Kubernetes jobs that simulate park visitors. Guests enter the park and explore attractions, making decisions based on their available money and attraction fees. Each guest reports metrics about their experience, and they leave when they run out of money or when the park closes.

## 🛠️ Technical Details

### Monitoring and Observability

Every component in kubepark exposes:

- Prometheus metrics at `/metrics` on port `9000` about revenue, usage, and guest satisfaction
- Logs at default Docker container location

## 🚀 Getting Started

1. Head to the [park](./park) directory to launch the park simulator
2. Explore the [attractions](./attractions) directory to deploy rides and amenities
3. Explore the [guest](./guest) directory to learn more about how guests will interact with your park
4. Monitor your park using tools like Grafana
5. Optimize your park based on insights from logs and metrics

## 🔒 Remember

This is a simulation game meant to be played through Kubernetes orchestration. Avoid direct HTTP requests or data manipulation - let the system work as designed!
