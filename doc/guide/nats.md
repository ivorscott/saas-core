# Debugging NATS

```bash
kubectl port-forward statefulset.apps/nats 4222 # make nats server accessible to nats cli
nats account info # view account info
nats stream info my_stream # view stream info
nats consumer info # view consumer info
nats stream view my_stream # view messages
nats stream rmm # remove message
nats stream purge # remove all messages
```
https://docs.nats.io/nats-concepts/jetstream/js_walkthrough
