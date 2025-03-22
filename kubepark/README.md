# kubepark

This simulation binary must be running in your cluster for the game to work. It's the default command on the kubepark image, and can be configured to run an easy, medium, or hard mode.

### Requirments

You'll need to create a deployment for the kubepark image. No need to set a command, we just want the default. You'll then need to create a service so that attractions and guests can make HTTP requests to your deployment. You may also create a volume so that your game state is saved even if you shut down the kubepark binary.

### Args

`--mode`, defaults to `easy`, but can also be set to `medium` and `hard`.
`--volume`, defaults to empty. Specify the directory where your volume exists so the game state can be saved there.
`--closed`, defaults to `false`, but can be set to `true` to prevent guests from coming into your park until you're ready.
`--entrance-fee`, defaults to `10` dollars, this is the price for guests to enter your park.
`--opens-at`, defaults to `8` in the morning. This determines the hour which your park opens.
`--closes-at`, defaults to `20` at night. This determines the hour which your park closes.

### Metrics

Prometheus metrics for each pod are available at the `/metrics` endpoint. Also, you can find the logs in the default location for a docker container.

`cash` is a gauge which is always set to the current amount of money you have to spend on attractions, repairs, etc.
`time` is a gauge which is set to a unix time stamp representing the current time in your park.
`entrance_fee` is a gauge set to the price to enter your park.
`opens_at` is a gauge set to the hour at which your park opens.
`closes_at` is a gauge set to the hour at which your park closes.
`is_closed` is a gauge set to 1 or 0, 0 means that your park is closed, 1 means it's open.
