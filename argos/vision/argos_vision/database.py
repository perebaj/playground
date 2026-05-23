"""Celebrity embedding database with cosine-similarity search."""

from __future__ import annotations

from collections.abc import Iterator
from pathlib import Path

import numpy as np
from PIL import Image

from .datasets import PersonImages
from .face import FaceProcessor

CACHE_DIR = Path(__file__).resolve().parent.parent / "data"

# Representative images are stored at a uniform size so we can keep the cache
# as a single uint8 ndarray. Saving with dtype=object preserves "object" on
# load even when shapes happen to match, which breaks downstream cv2 calls.
REPRESENTATIVE_SIZE = (256, 256)


def _to_uniform_uint8(img: np.ndarray) -> np.ndarray:
    """Resize an RGB image to REPRESENTATIVE_SIZE and return a contiguous uint8 array."""
    pil = Image.fromarray(np.asarray(img, dtype=np.uint8)).resize(REPRESENTATIVE_SIZE)
    return np.ascontiguousarray(np.asarray(pil, dtype=np.uint8))


class CelebrityDatabase:
    """Holds one prototype embedding per celebrity plus a representative image."""

    def __init__(
        self,
        embeddings: np.ndarray,
        names: list[str],
        images: np.ndarray,
    ) -> None:
        self.embeddings = embeddings
        self.names = names
        self.images = images

    def search(self, query_embedding: np.ndarray, top_k: int = 1) -> list[tuple[str, float, int]]:
        """Return up to top_k (name, similarity, index) tuples sorted by similarity."""
        sims = self.embeddings @ query_embedding
        order = np.argsort(-sims)[: min(top_k, len(self.names))]
        return [(self.names[i], float(sims[i]), int(i)) for i in order]

    @classmethod
    def load_or_build(
        cls,
        processor: FaceProcessor,
        loader: Iterator[PersonImages],
        cache_path: Path,
        force_rebuild: bool = False,
    ) -> CelebrityDatabase:
        if cache_path.exists() and not force_rebuild:
            data = np.load(cache_path, allow_pickle=True)
            return cls(
                embeddings=data["embeddings"],
                names=list(data["names"]),
                images=data["images"],
            )

        print(f"Building celebrity database (cache: {cache_path.name})...")
        valid_embeddings: list[np.ndarray] = []
        valid_names: list[str] = []
        valid_images: list[np.ndarray] = []

        for idx, (name, person_images) in enumerate(loader, start=1):
            embs: list[np.ndarray] = []
            representative: np.ndarray | None = None
            for img in person_images:
                emb = processor.detect_and_embed(img)
                if emb is not None:
                    embs.append(emb)
                    if representative is None:
                        representative = img
            if not embs or representative is None:
                continue
            mean_emb = np.mean(embs, axis=0)
            mean_emb /= np.linalg.norm(mean_emb)
            valid_embeddings.append(mean_emb)
            valid_names.append(name)
            valid_images.append(_to_uniform_uint8(representative))
            if idx % 20 == 0:
                print(f"  processed {idx} people  (kept {len(valid_names)})")

        if not valid_embeddings:
            raise RuntimeError("No valid embeddings produced from the loader.")

        embeddings = np.stack(valid_embeddings)
        names = np.array(valid_names)
        images = np.stack(valid_images).astype(np.uint8, copy=False)

        cache_path.parent.mkdir(parents=True, exist_ok=True)
        np.savez_compressed(cache_path, embeddings=embeddings, names=names, images=images)
        print(f"  cached {len(valid_names)} celebrity embeddings to {cache_path}")

        return cls(embeddings=embeddings, names=list(names), images=images)
