"""Tests for the pure search logic in CelebrityDatabase.

These tests avoid the heavy ML model and the celebrity datasets by constructing
tiny synthetic embedding tables, which is enough to validate the cosine-similarity
ranking that matching depends on, plus the on-disk format invariants the
downstream OpenCV calls rely on.
"""

from pathlib import Path

import numpy as np

from argos_vision.database import REPRESENTATIVE_SIZE, CelebrityDatabase, _to_uniform_uint8


def _normalize(v: np.ndarray) -> np.ndarray:
    return v / np.linalg.norm(v)


def test_search_returns_closest_by_cosine_similarity() -> None:
    embeddings = np.eye(3)
    db = CelebrityDatabase(embeddings, ["X", "Y", "Z"], images=np.empty(3, dtype=object))

    query = _normalize(np.array([0.9, 0.1, 0.0]))
    matches = db.search(query, top_k=1)

    assert matches[0][0] == "X"
    assert matches[0][1] > 0.9


def test_search_returns_top_k_in_descending_order() -> None:
    embeddings = np.eye(3)
    db = CelebrityDatabase(embeddings, ["A", "B", "C"], images=np.empty(3, dtype=object))

    query = _normalize(np.array([0.6, 0.5, 0.3]))
    matches = db.search(query, top_k=3)

    assert [m[0] for m in matches] == ["A", "B", "C"]
    similarities = [m[1] for m in matches]
    assert similarities == sorted(similarities, reverse=True)


def test_search_caps_top_k_at_database_size() -> None:
    embeddings = np.eye(2)
    db = CelebrityDatabase(embeddings, ["A", "B"], images=np.empty(2, dtype=object))

    matches = db.search(np.array([1.0, 0.0]), top_k=10)
    assert len(matches) == 2


def test_to_uniform_uint8_returns_uint8_at_target_size() -> None:
    img = np.random.randint(0, 255, size=(120, 80, 3), dtype=np.uint8)

    out = _to_uniform_uint8(img)

    assert out.dtype == np.uint8
    assert out.shape == (*REPRESENTATIVE_SIZE, 3)


def test_cached_images_roundtrip_as_uint8_not_object(tmp_path: Path) -> None:
    # Regression: prior versions stored images with dtype=object, which numpy
    # preserves on load even when shapes are uniform — that breaks cv2.resize
    # downstream with "src data type = object is not supported".
    embeddings = np.eye(3).astype(np.float32)
    names = np.array(["A", "B", "C"])
    raw = [np.full((100, 100, 3), v, dtype=np.uint8) for v in (10, 20, 30)]
    images = np.stack([_to_uniform_uint8(img) for img in raw])
    cache = tmp_path / "celebs.npz"
    np.savez_compressed(cache, embeddings=embeddings, names=names, images=images)

    data = np.load(cache, allow_pickle=True)
    db = CelebrityDatabase(
        embeddings=data["embeddings"], names=list(data["names"]), images=data["images"]
    )

    assert db.images.dtype == np.uint8
    assert db.images[0].dtype == np.uint8
    assert db.images[0].shape == (*REPRESENTATIVE_SIZE, 3)


def test_search_returns_index_pointing_to_correct_name() -> None:
    embeddings = np.array(
        [
            _normalize(np.array([1.0, 0.0])),
            _normalize(np.array([0.0, 1.0])),
        ]
    )
    db = CelebrityDatabase(embeddings, ["First", "Second"], images=np.empty(2, dtype=object))

    name, _, idx = db.search(_normalize(np.array([0.1, 1.0])), top_k=1)[0]
    assert name == "Second"
    assert db.names[idx] == "Second"
