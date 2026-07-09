# Handoff — Executive Protocol R15

## Completed

### Repository Sync

- **Fetched**: `origin` (MDMAtk/TormentNexus) + `origin-backup` (HyperNexusSoft/HyperNexus)
- **Upstream**: 0 ahead/0 behind — fully synced
- **Branches**: Pruned 297 stale `task/*` linked worktrees (all at commit 25a3a95ff, zero unique work)
- **Submodules**: None configured
- **Pushed**: Both remotes updated to 1d37ff261

### Version Bump

- **v1.0.0-alpha.251 → v1.0.0-alpha.252**
- Updated: `VERSION`, `package.json`, `CHANGELOG.md`

### New Files Committed

| File | Purpose |
|------|---------|
| `deploy/tenant-provision.sh` | Provision new org with isolated Docker/data/Nginx |
| `deploy/tenant-deprovision.sh` | Deprovision org, archive data, stop containers |
| `deploy/stripe-webhook-provisioner.sh` | Auto-provision from Stripe webhook |
| `deploy/provision-cert.sh` | Get individual subdomain cert via HTTP-01 |

### Production Deployment (Hetzner 5.161.250.43)

| Service | Port | Status |
|---------|------|--------|
| Go Sidecar (PM2) | 8090 | ✅ v1.0.0-alpha.237 |
| Nginx SSL | 443 | ✅ Wildcard `*.hypernexus.site` |
| Landing Page | — | ✅ <https://hypernexus.site> |
| TN API | — | ✅ <https://hypernexus.site/api/runtime/status> |
| PM2 | — | ✅ `tn-primary`, auto-restart via systemd |

### Wildcard SSL

- **Cert**: `/etc/letsencrypt/live/hypernexus.site-0001/`
- **SANs**: `*.hypernexus.site`, `hypernexus.site`
- **Expires**: 2026-10-07
- **Renewal**: WARNING — manual cert requires DNS TXT record update. Cron at weekly check runs dry-run. For full automation, install certbot-dns-dreamhost plugin or switch to per-tenant HTTP-01 certs.

### Vanilla MCP Question

**Why "1/2 servers" vs "/mcp shows both?**

- `/mcp` shows what's **configured** (both in `mcp.json`)
- The count shows what's **connected** — `browser-use` MCP fails because Python module missing
- `tormentnexus` (direct binary) connects fine
- Fix: `pip install browser-use` or remove from `mcp.json`

## Pending

1. **Add DNS A record** for `demo.hypernexus.site` pointing to `5.161.250.43`
2. **Wire Stripe billing** → auto-provision using `deploy/stripe-webhook-provisioner.sh`
3. **Add API keys** to Hetzner server for provider support
4. **Build & deploy Next.js dashboard** on the server
5. **Install pi-intercom** for subagents: `pi install npm:pi-intercom`
