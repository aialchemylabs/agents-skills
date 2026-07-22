# AI Alchemy Agent Skills

Claude Code plugin marketplace for reusable AI Alchemy Labs workflows.

## Install

```text
claude plugin marketplace add aialchemylabs/agents-skills
claude plugin install ai-gateway@aialchemy-agent-skills
```

Run `/reload-plugins` in Claude Code after installing or updating the plugin.

## AI Gateway

The `ai-gateway` plugin contains LiteLLM operational skills. Its first skill is:

- `ai-gateway:configuring-litellm-entra-sso` — configure and harden direct Microsoft Entra OIDC SSO for a LiteLLM Admin UI while keeping API clients on scoped virtual keys.

The marketplace metadata is in `.claude-plugin/marketplace.json`. Each plugin keeps only `plugin.json` inside its `.claude-plugin/` directory and places its skills under `skills/<skill-name>/SKILL.md`.
