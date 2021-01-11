# Common Module for Event-driven Architecture

Designed for Command Query Responsibility Segregation (CQRS) and event sourcing.

[Github Repository](https://github.com/ivorscott/devpie-client-events)

## Install

```
pip install devpie-client-events
```

Example:

```
import events

print(events.Commands.ADD_USER.value)
print(events.Events.USER_ADDED.value)
```
