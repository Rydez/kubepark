# Guest 👥

The guest is a Kubernetes job that simulates a visitor to your amusement park. Each guest enters the park with a set amount of money and explores attractions based on their preferences and available funds.

## 🎯 Overview

Guests are the lifeblood of your park! They:

- Enter the park and explore available attractions
- Make decisions based on their remaining money
- Report their experiences through logs and metrics
- Leave when they run out of money or the park closes

## 🔧 Configuration

### Required Arguments

- `--park-url`: URL of the kubepark service (default: `http://kubepark:80`)

## 📊 Metrics

Each guest exposes Prometheus metrics at `/metrics` on port 9000:

- `money_spent`: Total amount spent on attractions
- `attractions_visited`: Number of attractions experienced
- `time_spent`: Duration of the park visit (in seconds)

## 🔄 Behavior

1. **Park Entry**

   - Makes a POST request to `/enter` on the kubepark service
   - Logs successful entry or failure

2. **Attraction Exploration**

   - Fetches available attractions from kubepark
   - Randomly selects an attraction to visit
   - Sends current money balance to the attraction
   - Attraction validates if guest can afford the fee
   - Updates money and metrics after successful visit
   - Takes a random break between attractions

3. **Visit Completion**
   Guest leaves the park when:
   - Park closes
   - Guest runs out of money

## 📝 Logging

Guests provide detailed logs about:

- Park entry attempts
- Attraction visits and costs
- Insufficient funds
- Visit completion

## 🔍 Monitoring Tips

- Track `money_spent` to understand guest spending patterns
- Monitor `attractions_visited` to identify popular attractions
- Use `time_spent` to optimize park engagement
- Watch for failed entry attempts or insufficient funds
