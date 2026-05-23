"""Tests for FaceProcessor wiring.

These tests use monkeypatching to avoid downloading the MTCNN/FaceNet weights
and instantiating PyTorch models. They protect against regressions in the
device-selection wiring — specifically, that MTCNN is pinned to CPU regardless
of the device requested for the embedder, because PyTorch's MPS backend has a
bug in adaptive_avg_pool2d that breaks MTCNN's image-pyramid resampling.
"""

from __future__ import annotations

from typing import Any

import pytest


class _FakeMTCNN:
    last_kwargs: dict[str, Any] = {}

    def __init__(self, **kwargs: Any) -> None:
        type(self).last_kwargs = kwargs


class _FakeEmbedder:
    last_to_device: str | None = None

    def __init__(self, **_: Any) -> None:
        pass

    def eval(self) -> _FakeEmbedder:
        return self

    def to(self, device: str) -> _FakeEmbedder:
        type(self).last_to_device = device
        return self


@pytest.fixture
def fake_models(monkeypatch: pytest.MonkeyPatch) -> None:
    import argos_vision.face as face_mod

    _FakeMTCNN.last_kwargs = {}
    _FakeEmbedder.last_to_device = None
    monkeypatch.setattr(face_mod, "MTCNN", _FakeMTCNN)
    monkeypatch.setattr(face_mod, "InceptionResnetV1", lambda **kwargs: _FakeEmbedder(**kwargs))


def test_mtcnn_is_pinned_to_cpu_even_when_mps_requested(fake_models: None) -> None:
    from argos_vision.face import FaceProcessor

    FaceProcessor(device="mps")

    assert _FakeMTCNN.last_kwargs["device"] == "cpu"


def test_mtcnn_is_pinned_to_cpu_even_when_cuda_requested(fake_models: None) -> None:
    from argos_vision.face import FaceProcessor

    FaceProcessor(device="cuda")

    assert _FakeMTCNN.last_kwargs["device"] == "cpu"


def test_embedder_uses_requested_device(fake_models: None) -> None:
    from argos_vision.face import FaceProcessor

    FaceProcessor(device="mps")

    assert _FakeEmbedder.last_to_device == "mps"


def test_pick_device_returns_known_backend() -> None:
    from argos_vision.face import _pick_device

    assert _pick_device() in {"cpu", "cuda", "mps"}
