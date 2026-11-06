# CLI and JSON contract

Use absolute paths. JSON files and generated outputs are task-local, outside the skill/source repository. `inspect` needs only the `data` member below. It never imports model packages, installs, loads weights or forecasts. `propose` validates a full request without executing predictions. Present the resulting scope plus a human-readable comparison; select only after conversation authorization.

```sh
python3 /absolute/skill/scripts/forecast.py inspect /external/task/request.json
python3 /absolute/skill/scripts/forecast.py propose /external/task/request.json --output /external/task/option-A.json
# After explicit choice, or a verified scoped standing authorization:
python3 /absolute/skill/scripts/forecast.py run /external/task/option-A.json --selected THE_PRINTED_FINGERPRINT
```

A minimal statistical request (replace paths, dates, units and bounds with inspected values):

```json
{
  "data": {
    "path": "/external/task/readings.csv",
    "id": "id", "timestamp": "date", "target": "value",
    "timezone": "UTC", "frequency": "daily", "as_of": "2027-02-28",
    "units": {"meter-a": "kWh"}
  },
  "model": "seasonal-naive",
  "report": {"start": "2027-04-01", "end": "2027-06-30", "calendar_basis": "calendar"},
  "origins": ["2026-08-31", "2026-10-31"],
  "context": 128, "season": 7,
  "baselines": ["mean", "naive", "seasonal-naive", "drift"],
  "primary_loss": "mase", "pooling": "weighted-series-mean",
  "weights": {"meter-a": 1}, "semantics": "flow", "aggregation": "none",
  "covariates": [], "groups": [["meter-a"]], "cross_learning": false,
  "settings": {
    "device": "cpu", "precision": "float32", "compile": false,
    "clipping": "none", "quantile_crossing": "reject",
    "truncation": "last-context", "horizon_extension": "reject"
  },
  "budgets": {
    "seconds": 300, "threads": 4, "batch_size": 8,
    "max_series": 1, "max_origins": 2, "max_horizon": 122,
    "planned_download_bytes": 0, "retries": 0
  },
  "output": "/external/task/report",
  "environment": "/external/environments/forecast-v1",
  "cache": "/external/cache/forecast-models",
  "offline": true, "download_estimate_bytes": 0
}
```

Origins need the entire generated lead window within historical observations. Data is a CSV or JSON row array with mapped columns, optional `available_at` and `complete`. Values are finite numbers or null/blank. `data.max_bytes` defaults to 20 MB. IDs/dates must be unique and all series aligned; `data.units` names every series. `data.calendar`, when supplied, is a strictly increasing complete list of valid history/future timestamps. Without it automatic frequencies are hourly, daily, weekly or month-end monthly. Explicit null rows preserve a regular calendar without pretending missing values were observed. `as_of` must be at or after all supplied observation publication times. Date-only values mean midnight in the specified timezone; use offset timestamps when finer availability matters.

Full requests require exactly the listed top-level keys. Models are `mean`, `naive`, `seasonal-naive`, `drift`, `timesfm`, `chronos`. The model revision and lock are drawn from the installed manifest and bound into the proposal, not supplied as mutable runtime overrides. Foundation choices need a conservatively disclosed `download_estimate_bytes` (at least 1.5 GB, beyond weights/configs), sufficient planned download budget, and a tested host. Both adapter packages share one dependency lock; only the selected checkpoint is downloaded. `offline=false` permits listed setup/downloads only after selection.

`aggregation` is `none`, `direct-total` or `sum-point-rules` for same-unit flows. `groups` partitions series exactly; cross-learning defaults to false in authored options and is supported only by Chronos. A cross-learning group plus its covariates must fit the fixed `batch_size`; it cannot be silently split. `primary_loss` is MAE/RMSE/MASE; unlike units require MASE. Explicit positive weights define the mean across series. No automatic model search.

The only optional quantile correction is `quantile_crossing="clamp-about-median"`: starting at q50, lower quantiles are clamped downward successively and upper quantiles upward. Point medians remain unchanged. It must be selected before evaluation and appears in provenance/status. Default `reject` stops invalid outputs. No clipping, extrapolated horizon or precision/backend switch is supported.

Chronos covariate example:

```json
{"name": "planned_shift", "kind": "known", "values": [
  {"id": "meter-a", "date": "2027-03-01", "value": 2, "available_at": "2027-02-20"}
]}
```

`covariates` is a list of such objects; the name must match a numeric historical column. `kind` is `past` (no future values), `known`, `vintage` or `scenario`. Supply all required generated dates and all backtest dates with origin-visible `available_at`. Future realized targets are prohibited as covariates. Scenario and vintage availability must reflect their actual origin, not an invented historical date. Changing covariates, groups, information-sharing or a scenario is scope drift.

`proposal.scope` includes the complete request, resolved output/environment/cache destinations, data digest, exact generated dates/scoring indices, checkpoint files/digests, dependency-lock digest, script-content digest and artifact estimate. Its canonical SHA-256 must equal `--selected`, and recomputation must agree before setup/imports. Retargeting a destination symlink changes the scope; execution uses the bound resolved destinations. This prevents accidental drift; anyone can compute a hash, so the harness conversation remains the authorization evidence. Internal worker/fetch entrypoints recheck the same scope. They are not an approval service.
