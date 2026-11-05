# Repository ownership

This repository owns reusable agent skills, their supporting resources,
plugin metadata, and skill distribution tooling.

Standalone MCP servers and their runtime/deployment configuration belong in
`aialchemylabs/aialchemy-mcp-servers`. Keep executable server code out of this repository.

# Public repository hygiene

Keep private operating plans, deployment-specific agent profiles and prompts,
machine configuration, customer information, credentials, and repository backups
outside this public repository. Use fictional examples in reusable skills.

Before publishing, inspect staged changes and author/committer identities. Use
an intentional public-safe Git email, never an automatically derived machine
address. Keep GitHub secret scanning and push protection enabled.

Deleting a file does not remove it from Git history. For an authorized history
cleanup, retain a private backup outside the checkout, verify the result from a
fresh clone, and push only the intended refs with an explicit expected-old-SHA
force-with-lease. Never mirror-push a development checkout or publish local agent
checkpoint refs. Update downstream catalog pins and installations afterward.
