# KubePark Walkthrough

Welcome to KubePark! This walkthrough demonstrates a complete Kubernetes observability setup using Grafana Alloy for log and metrics collection.

## Quick Start

1. **Setup the infrastructure:**

   ```bash
   task setup
   ```

   This uses Helm to deploy Alloy and sets up the complete monitoring stack.

2. **Build and deploy the game:**

   ```bash
   task build
   task deploy-park
   task deploy-carousel
   task deploy-restroom
   ```

3. **Monitor the park:**
   ```bash
   task status
   task open-grafana
   ```

## Available Commands

Run `task --list` to see all available commands, or just `task` for help:

### Setup & Build

- `task setup` - Create Kind cluster and start monitoring stack
- `task build` - Build and push the game image

### Game Commands

- `task deploy-park` - Start the park (begins the game!)
- `task deploy-carousel` - Deploy carousel attraction
- `task deploy-restroom` - Deploy restroom attraction
- `task list-carousels` - Show all deployed carousel instances
- `task list-restrooms` - Show all deployed restroom instances
- `task remove-restrooms` - Remove all restroom instances

### Monitoring

- `task status` - Show current park status
- `task logs` - View park logs
- `task open-grafana` - Open Grafana dashboard
- `task restart-monitoring` - Restart monitoring stack
- `task restart-alloy` - Restart Alloy log collection
- `task upgrade-alloy` - Upgrade Alloy with latest configuration

### Cleanup

- `task clean` - Clean up everything

## Alloy Management

Task includes specific commands for managing Alloy:

- `task restart-alloy` - Restart the Alloy DaemonSet
- `task upgrade-alloy` - Upgrade Alloy with the latest configuration
- `task status` - Shows Alloy pod status (uses correct Helm labels)

## Configuration Files

- **`k8s/alloy-values.yaml`** - Helm values file with Alloy configuration

## Migration Notes

The Helm deployment maintains the same functionality as the previous manual deployment:

- DaemonSet deployment for log collection from all nodes
- Same Alloy configuration for Prometheus and Loki integration
- Identical RBAC permissions and volume mounts
- Compatible with existing monitoring stack

## Cleanup

The cleanup process properly removes the Helm deployment:

```bash
task clean
```

This will uninstall the Alloy Helm release along with all other resources.
