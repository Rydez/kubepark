# kubepark

An amusement park simulation game played on Kubernetes. Deploy rides, concession stands, restrooms, and more for the guests in your park. Monitor your deployments and guests in a visualization software like Grafana. Use visualizations to determine how to improved your park for your guests to make it more profitable.

The kubepark image is full of binaries for attractions to make your park a fun and interesting place. These attractions are Kubernetes services which guests can interact with by calling their HTTP endpoints. The first time you deploy an attraction you will need to pay the build price, and if the attraction breaks down and the pod crashes it will cost a maintenance fee to repair it when the pod is restarted. Each attraction can have a customized entrance fee which can be configured by deployment args.

The kubepark image also contains the binary responsible for running the park simulation. This simulation keeps track of what day it is, how much money you have, what attractions you own, and so on. It's also responsible for creating guests that will enter, explore, and interact with your park. These guests are Kubernetes jobs which will make HTTP requests to the attractions. It's important to monitor the logs and metrics for the guests because they will let you know what they like and don't like.

This game is meant to be played as an orchestrator. You should not modify the data or make your own HTTP requests to the pods. Rather you should be creating and configuring Kubernetes objects and using the logs and metrics from those objects as your interface. Every binary has logs in the default location for docker, and has a prometheus metrics enpoint, `/metrics`, exposed at port `9000`.

### Launching the kubepark binary

Every directory has it's own `README.md`, so head over to the [kubepark](./kubepark) directory so you can launch the simulator in your cluster and start playing. Then checkout the other directories to explore the attractions.

### Launching attractions

Most directory contain an attraction, which is something in your park that your guests will interact with, like rides, amenities, decor, and so on. Just like the kubepark binary, every attraction has logs in the default location for docker, and has a prometheus metrics enpoint, `/metrics`, exposed at port `9000`. All attractions will need to be deployed using a certain command with a service in a certain namespace so that guests can find it. The command and namespace for a specific attraction can be found in the README in the attraction's directory. Also, every attraction has the following args and metrics:

#### Args

- `--closed`, defaults to `false`, but can be set to `true` to prevent guests from using this attraction.
- `--fee`, defaults to `0` dollars, this is the price for guests to use this attraction.

#### Metrics

Prometheus metrics for each pod are available at the `/metrics` endpoint. Also, you can find the logs in the default location for a docker container.

- `revenue` is a counter which tracks how much money this attraction has made.
- `fee` is a gauge set to the price to enter your park.
- `is_closed` is a gauge set to 1 or 0, 0 means that the attraction is closed, 1 means it's open.
- `attempts` is a counter which tracks when a guest tries to use this attraction. This metric has the following labels:
  - `success` which can be `true` or `false`
  - `reason` which when `success=true` describes why a guest used this attraction, and when `success=false` describes why a guest couldn't use the attraction.
