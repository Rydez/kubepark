# Guest 👥

The guest is a Kubernetes job that simulates a visitor to your amusement park. Each guest enters the park with a set amount of money and explores attractions based on their preferences and available funds. They enter the park and explore available attractions, making decisions based on their remaining money. Throughout their visit, they report their experiences through logs and metrics. When they run out of money or when the park closes, they leave the park.

## 🔧 Configuration

- `--park-url`: URL of the kubepark service (default: `http://kubepark:80`)

## 📊 Metrics

Each guest exposes Prometheus metrics at `/metrics` on port 9000:

- `money_spent`: Total amount spent on attractions
- `attractions_visited`: Number of attractions experienced
- `time_spent`: Duration of the park visit (in seconds)

## 🪵 Logging

Logs can be found in the default location for a docker container.
