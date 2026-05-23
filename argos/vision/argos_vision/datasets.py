"""Celebrity image loaders.

Each loader yields ``(name, list_of_RGB_uint8_images)`` per person. The
``CelebrityDatabase`` consumes any loader to build its embedding index, so
adding a new source is a matter of writing a new generator with this contract.
"""

from __future__ import annotations

import random
from collections.abc import Iterator
from pathlib import Path

import numpy as np
from PIL import Image

PersonImages = tuple[str, list[np.ndarray]]


def _pil_to_uint8(img: Image.Image) -> np.ndarray:
    return np.asarray(img.convert("RGB"), dtype=np.uint8)


def pins_loader(max_per_person: int = 25, seed: int = 0) -> Iterator[PersonImages]:
    """Pins Face Recognition (Kaggle): 105 celebrities, high-quality Pinterest crops."""
    import kagglehub

    root = Path(kagglehub.dataset_download("hereisburak/pins-face-recognition"))
    # The dataset ships images under nested PINS/PINS/pins_<name>/*.jpg directories.
    candidates = [p for p in root.rglob("pins_*") if p.is_dir()]
    if not candidates:
        raise RuntimeError(f"Could not locate PINS person directories under {root}")

    rng = random.Random(seed)
    for person_dir in sorted(candidates):
        name = person_dir.name.removeprefix("pins_").strip()
        image_paths = sorted(person_dir.glob("*.jpg")) + sorted(person_dir.glob("*.jpeg"))
        if len(image_paths) > max_per_person:
            image_paths = rng.sample(image_paths, max_per_person)
        images = [_pil_to_uint8(Image.open(p)) for p in image_paths]
        if images:
            yield name, images


def hf1000_loader(max_per_person: int = 25, seed: int = 0) -> Iterator[PersonImages]:
    """Tony Assi's celebrity-1000 on HuggingFace: 1000 contemporary celebrities.

    Requires the optional ``hf`` extra: ``uv sync --extra hf``.
    """
    try:
        from datasets import load_dataset
    except ImportError as exc:
        raise ImportError(
            "The 'hf' extra is required for the hf1000 dataset. Install with: uv sync --extra hf"
        ) from exc

    ds = load_dataset("tonyassi/celebrity-1000", split="train")
    label_names = ds.features["label"].names

    by_label: dict[int, list[np.ndarray]] = {}
    for example in ds:
        by_label.setdefault(example["label"], []).append(_pil_to_uint8(example["image"]))

    rng = random.Random(seed)
    for label, images in by_label.items():
        if len(images) > max_per_person:
            images = rng.sample(images, max_per_person)
        yield label_names[label], images


# Registry consumed by the CLI. Each entry is (callable, default_cache_name).
LOADERS: dict[str, tuple[callable, str]] = {
    "pins": (pins_loader, "pins"),
    "hf1000": (hf1000_loader, "hf1000"),
}
