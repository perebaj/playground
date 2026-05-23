"""Tests for the pure search logic in CelebrityDatabase.

These tests avoid the heavy ML model and the LFW dataset by constructing tiny
synthetic embedding tables, which is enough to validate the cosine-similarity
ranking that matching depends on.
"""

import numpy as np

from argos_vision.database import CelebrityDatabase


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
