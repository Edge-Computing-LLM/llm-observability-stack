## Summary

- What changed?
- Why was it needed?

## Validation

- [ ] `helm lint .`
- [ ] `helm template llm-observability-stack . >/tmp/rendered-default.yaml`
- [ ] `helm template llm-observability-stack . -f values.local-k3s.example.yaml >/tmp/rendered-local.yaml`
- [ ] Manual runtime checks (if applicable) documented below

## Security Review

- [ ] No secrets or tokens added to git-tracked files
- [ ] Values-file changes preserve `values.yaml` as non-sensitive defaults
- [ ] Any new externally exposed endpoints are documented

## Notes

Include rollout notes, breaking changes, and follow-up tasks.
