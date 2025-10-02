# KubePark Walkthrough

Welcome to KubePark! This walkthrough demonstrates a complete Kubernetes observability setup using Grafana Alloy for log and metrics collection.

## What's New: Helm-based Alloy Deployment

We've migrated from manual YAML deployments to using the official Grafana Alloy Helm chart for better management and easier updates.

### Benefits of Using Helm for Alloy:

- **Official Support**: Uses the officially maintained Grafana Helm chart
- **Easier Updates**: Simple `helm upgrade` commands for configuration changes
- **Better Management**: Helm tracks deployment history and allows rollbacks
- **Standardized Configuration**: Follows Helm best practices and conventions
- **Dependency Management**: Automatic handling of RBAC, ServiceAccount, and other resources

## Quick Start

1. **Setup the infrastructure:**

   ```bash
   make setup
   ```

   This now uses Helm to deploy Alloy instead of manual YAML files.

2. **Build and deploy the game:**

   ```bash
   make build
   make deploy-park
   make deploy-carousel
   make deploy-restroom
   ```

3. **Monitor the park:**
   ```bash
   make status
   make open-grafana
   ```

## Alloy Management Commands

The Makefile now includes specific commands for managing Alloy:

- `make restart-alloy` - Restart the Alloy DaemonSet
- `make upgrade-alloy` - Upgrade Alloy with the latest configuration
- `make status` - Shows Alloy pod status (now uses correct Helm labels)

## Configuration Files

- **`k8s/alloy-values.yaml`** - Helm values file with Alloy configuration

## Migration Notes

The Helm deployment maintains the same functionality as the previous manual deployment:

- DaemonSet deployment for log collection from all nodes
- Same Alloy configuration for Prometheus and Loki integration
- Identical RBAC permissions and volume mounts
- Compatible with existing monitoring stack

## Cleanup

The cleanup process now properly removes the Helm deployment:

```bash
make clean
```

This will uninstall the Alloy Helm release along with all other resources.

---

For detailed game instructions and monitoring setup, see the original README sections below.
