# kubepark 🎢

### An amusement park tycoon game that runs on Kubernetes!

You are an entrepreneur seeking to build an amazing amusement park by deploying various attractions, monitoring guest behavior, and maximizing profits.

- Deploy and manage various attractions (rides, amenities, concessions)
- Monitor guest satisfaction and park performance
- Handle maintenance and repairs
- Track revenue and expenses

## 🎮 How it works

**Commands:** Actions are performed using the commands defined in the Taskfile. Use these commands to start the game and create your park. The commands will deploy k8s components which function as the different parts of your amusement park.

**Monitoring:** Grafana is used to monitor how your park, rides, and guests are doing. The grafana dashboard will tell you what time it is in your park, how much money you have to spend, which rides are broken down, and much more.

## 🚀 Getting Started

1. **Setup the infrastructure:**

   ```bash
   task setup
   ```

   This will create a local k8s cluster using Kind (Kubernetes in Docker) and start running your monitoring stack (Grafana, Prometheus, Loki).

2. **Start the game:**

   ```bash
   task build
   task deploy -- park
   ```

   This will build the single image needed, and then deploy an emtpy park ready for you to build.

3. **Monitor the park:**

   ```bash
   task open-grafana
   ```

   Grafana is the window into your park. The Grafana dashboard will give you a plethora of information to understand and manage your park.

4. **Create an attraction:**

   ```bash
   task deploy -- carousel
   task deploy -- restroom
   ...
   ```

   Now you're ready to start building. Spend your money wisely.

## 🔒 Remember

This is a game meant to be played through Kubernetes orchestration. Avoid direct HTTP requests or data manipulation to let the system work as designed.

## Developers

For more advanced docs that will help develop this game, but are irrelevant to playing it, see [DEVELOP.md](./DEVELOP.md).
