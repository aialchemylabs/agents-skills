"""TimesFM 2.5 Torch; optional imports only after external scope validation."""

import math


class Adapter:
    def __init__(self, checkpoint, request):
        import torch
        import timesfm

        torch.set_num_threads(request["budgets"]["threads"])
        torch.set_num_interop_threads(1)
        self.request = request
        self.model = timesfm.TimesFM_2p5_200M_torch(torch_compile=False)
        # Upstream chooses CUDA in its constructor; override before checkpoint loading.
        self.model.model.device = torch.device("cpu")
        self.model.model.device_count = 1
        self.model.model.to(device="cpu", dtype=torch.float32)
        self.model.load_checkpoint(str(checkpoint), torch_compile=False)
        assert (
            self.model.model.device.type == "cpu" and self.model.torch_compile is False
        )
        self.config = timesfm.ForecastConfig

    def predict(self, values, horizon, covariates):
        import numpy as np

        r = self.request
        if any(past or future for past, future in covariates):
            raise ValueError("TimesFM adapter excludes covariates/XReg")
        self.model.compile(
            self.config(
                max_context=math.ceil(r["context"] / 32) * 32,
                max_horizon=math.ceil(horizon / 128) * 128,
                per_core_batch_size=r["budgets"]["batch_size"],
                normalize_inputs=True,
                use_continuous_quantile_head=True,
                force_flip_invariance=True,
                infer_is_positive=False,
                fix_quantile_crossing=False,
                return_backcast=False,
            )
        )
        points, channels = self.model.forecast(
            horizon=horizon, inputs=[np.array(v, dtype=np.float32) for v in values]
        )
        if not np.allclose(points, channels[:, :, 5], rtol=0, atol=0):
            raise ValueError("loaded TimesFM primary no longer equals q50 channel5")
        return [
            {
                "point": p.tolist(),
                "statistic": "median",
                "quantiles": {str(q / 10): ch[:, q].tolist() for q in range(1, 10)},
                "runtime_semantics": "primary equals channel5/q50; channel0 not advertised as a mean",
            }
            for p, ch in zip(points, channels)
        ]
