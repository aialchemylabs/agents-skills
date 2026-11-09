# Harness adapters

Use only capabilities actually exposed by the current host. Discovery metadata shows potential access, not proof of successful inference. Report actual execution failures. The Python helper implements LiteLLM only; it does not emulate harness tools.

## T3 / Athanor

Discover tool names from the host's available tool catalog. If T3 tools load late, try the named capability once before declaring it absent. Before the first `delegate_task`, call `orchestrator_capabilities` once per thread. Use its provider/model/options information, and any discovery it directs, to build the roster. Refresh after an error or a user-named provider change. Pass supported settings in `target.options`; do not guess an options schema.

Use native delegation for a same-provider child only when the exposed native tool supports that model, the settings need nothing beyond the parent's, and it finishes in this turn. Use `delegate_task` for another provider, user-requested settings, isolated checkout, background work, or an explicit request for that tool. If a required provider/model is unavailable, report it; do not switch to native tools, a CLI, Workflow, or another model unasked.

Each independent assessment, subsequent review, or re-check is a **new** `delegate_task` with the bounded brief, evidence, findings and objections it needs. Do not reuse a conversation via `t3_thread_send` for fresh reviews. Follow the running task's returned `next`/completion instructions; do not poll `task_status`, sleep, or create watchers. Cancel superseded or unwanted work with `task_cancel`, never by interrupting a child thread. Do not start replacement tasks when told to use fewer delegates. Top-level thread launch tools are not subagent tools.

Keep reviewers read-only. For code reviews, provide the exact revision/diff and relevant tests; give no implementation authority to the reviewer. No checkout setup commands in the brief. The implementation owner fixes accepted findings and supplies the resulting evidence for a fresh, bounded re-check.

## Other native hosts

Inspect the exposed delegation tool's schema and supported model list. A model selector is usable only for models the tool actually offers; role names or repeated same-model children do not establish diversity. Honor concurrency, completion/wait, cancellation, inherited-context, read-only and permission rules of that host. Do not copy T3 tool names into hosts that do not expose them.

Pass a common evidence packet with minimal inherited conclusions. If isolated context is unavailable, state the independence limitation. Wait for actual responses using the host's supported completion mechanism. Never present a launched-but-unfinished child as a completed reviewer.

## No usable harness

If an already configured LiteLLM gateway is accessible and within authority, use the helper. Discover accessible aliases with its `models` command before proposing a roster; no manifest or paid inference is needed. If neither backend offers two distinct models, stop with the access limitation. Do not install a gateway, create keys, or claim a multi-model council as a workaround.
