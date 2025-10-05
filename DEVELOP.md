## Advanced Commands

Run `task --list` to see all available commands, or just `task` for help:

### Game Commands

- `task deploy -- park` - Start the park (begins the game!)
- `task deploy -- carousel` - Deploy carousel attraction
- `task remove -- carousel` - Remove a carousel instance

### Monitoring

- `task logs` - View park logs
- `task restart -- monitoring` - Restart monitoring stack
- `task restart -- alloy` - Restart Alloy log collection
- `task upgrade-alloy` - Upgrade Alloy with latest configuration

### Cleanup

- `task clean` - Clean up everything

## Cleanup

The cleanup process properly removes the Helm deployment:

```bash
task clean
```

This will uninstall the Alloy Helm release along with all other resources.
