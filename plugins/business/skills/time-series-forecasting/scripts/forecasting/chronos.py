"""Chronos-2 local safetensors with explicit grouping and median normalization."""


class Adapter:
    def __init__(self, checkpoint, request):
        import torch
        from chronos import Chronos2Pipeline

        torch.set_num_threads(request["budgets"]["threads"])
        torch.set_num_interop_threads(1)
        self.request = request
        self.model = Chronos2Pipeline.from_pretrained(
            str(checkpoint),
            device_map="cpu",
            torch_dtype=torch.float32,
            local_files_only=True,
            trust_remote_code=False,
        )
        self.model.model.eval()
        if (
            self.model.model.device.type != "cpu"
            or self.model.model.dtype != torch.float32
        ):
            raise ValueError("Chronos effective device/precision mismatch")

    def predict(self, values, horizon, covariates):
        import numpy as np
        import torch

        items = [
            {
                "target": np.array(v, dtype=np.float32),
                "past_covariates": {
                    k: np.array(v, dtype=np.float32) for k, v in past.items()
                },
                "future_covariates": {
                    k: np.array(v, dtype=np.float32) for k, v in future.items()
                },
            }
            for v, (past, future) in zip(values, covariates)
        ]
        with torch.inference_mode():
            qs, misnamed_mean = self.model.predict_quantiles(
                items,
                prediction_length=horizon,
                quantile_levels=[q / 10 for q in range(1, 10)],
                context_length=self.request["context"],
                batch_size=self.request["budgets"]["batch_size"],
                cross_learning=self.request["cross_learning"],
                limit_prediction_length=True,
            )
        result = []
        for quantiles, point in zip(qs, misnamed_mean):
            if not torch.equal(point, quantiles[..., 4]):
                raise ValueError("loaded Chronos point no longer equals q50")
            result.append(
                {
                    "point": point[0].cpu().tolist(),
                    "statistic": "median",
                    "quantiles": {
                        str(q / 10): quantiles[0, :, q - 1].cpu().tolist()
                        for q in range(1, 10)
                    },
                    "runtime_semantics": "upstream mean variable equals median/q50",
                }
            )
        return result
