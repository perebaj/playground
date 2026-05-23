"""CLI entrypoint for argos-vision.

Two subcommands:

* ``build``  — download the chosen dataset and pre-compute the embedding cache.
* ``match``  — capture a face from the webcam (or a file) and find the
  closest celebrity. Requires the cache built by ``build``.
"""

from __future__ import annotations

import argparse
import sys
from pathlib import Path

import cv2
import numpy as np

from .camera import capture_from_webcam
from .database import CACHE_DIR, CelebrityDatabase
from .datasets import LOADERS, hf1000_loader, pins_loader
from .face import FaceProcessor


def _cache_path(dataset: str) -> Path:
    return CACHE_DIR / f"{dataset}_embeddings.npz"


def _build_loader(args: argparse.Namespace):
    if args.dataset == "pins":
        return pins_loader(max_per_person=args.max_per_person)
    if args.dataset == "hf1000":
        return hf1000_loader(max_per_person=args.max_per_person)
    raise ValueError(f"Unknown dataset: {args.dataset}")


def show_result(
    query_image: np.ndarray,
    celeb_name: str,
    celeb_image: np.ndarray,
    similarity: float,
) -> None:
    target_h = 400

    def fit(img: np.ndarray) -> np.ndarray:
        # Defensive: caches written before the uint8 fix may resurface as
        # object dtype; cv2.resize rejects those, so coerce on the way in.
        img = np.ascontiguousarray(img, dtype=np.uint8)
        ratio = target_h / img.shape[0]
        return cv2.resize(img, (int(img.shape[1] * ratio), target_h))

    side_by_side = np.hstack(
        [
            cv2.cvtColor(fit(query_image), cv2.COLOR_RGB2BGR),
            cv2.cvtColor(fit(celeb_image), cv2.COLOR_RGB2BGR),
        ]
    )
    label = f"{celeb_name}  similarity={similarity:.3f}"
    cv2.putText(side_by_side, label, (10, 30), cv2.FONT_HERSHEY_SIMPLEX, 0.8, (0, 255, 0), 2)

    print(f"\nYou look like: {celeb_name} (similarity: {similarity:.3f})")
    print("Press any key in the result window to close.")
    cv2.imshow("Argos — match", side_by_side)
    cv2.waitKey(0)
    cv2.destroyAllWindows()


def cmd_build(args: argparse.Namespace) -> int:
    print("Loading face processor (MTCNN + FaceNet/VGGFace2)...")
    processor = FaceProcessor()

    cache_path = _cache_path(args.dataset)
    if cache_path.exists() and not args.force:
        print(f"Cache already exists at {cache_path}. Use --force to rebuild.")
        return 0

    print(f"Building embedding cache for dataset '{args.dataset}'...")
    CelebrityDatabase.load_or_build(
        processor,
        loader=_build_loader(args),
        cache_path=cache_path,
        force_rebuild=args.force,
    )
    return 0


def cmd_match(args: argparse.Namespace) -> int:
    cache_path = _cache_path(args.dataset)
    if not cache_path.exists():
        print(
            f"No cache found at {cache_path}.\n"
            f"Run `argos-vision build --dataset {args.dataset}` first.",
            file=sys.stderr,
        )
        return 1

    print("Loading face processor (MTCNN + FaceNet/VGGFace2)...")
    processor = FaceProcessor()

    print(f"Loading celebrity database: {args.dataset}")
    data = np.load(cache_path, allow_pickle=True)
    db = CelebrityDatabase(
        embeddings=data["embeddings"],
        names=list(data["names"]),
        images=data["images"],
    )
    print(f"  {len(db.names)} celebrities loaded.")

    if args.image:
        img_bgr = cv2.imread(args.image)
        if img_bgr is None:
            print(f"Failed to read image: {args.image}", file=sys.stderr)
            return 1
        query_image = cv2.cvtColor(img_bgr, cv2.COLOR_BGR2RGB)
    else:
        print("Opening webcam — press SPACE to capture, ESC to quit.")
        query_image = capture_from_webcam(args.camera)  # type: ignore
        if query_image is None:
            print("No image captured. Bye.")
            return 0

    emb = processor.detect_and_embed(query_image)
    if emb is None:
        print("No face detected. Try again with better lighting or framing.", file=sys.stderr)
        return 2

    name, similarity, idx = db.search(emb, top_k=1)[0]
    show_result(query_image, name, db.images[idx], similarity)
    return 0


def cmd_list(args: argparse.Namespace) -> int:
    print(f"{'dataset':<10}  {'cache':<40}  status")
    print(f"{'-' * 10}  {'-' * 40}  ------")
    for name in LOADERS:
        cache = _cache_path(name)
        status = "ready" if cache.exists() else "not built"
        print(f"{name:<10}  {str(cache.name):<40}  {status}")
    return 0


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(
        description="Argos vision: face capture, embedding, and celebrity matching."
    )
    sub = parser.add_subparsers(dest="command", required=True)

    def add_dataset_args(p: argparse.ArgumentParser) -> None:
        p.add_argument(
            "--dataset",
            choices=list(LOADERS.keys()),
            default="pins",
            help="Which celebrity catalog to use (default pins).",
        )

    p_build = sub.add_parser("build", help="Build the embedding cache for a dataset.")
    add_dataset_args(p_build)
    p_build.add_argument(
        "--max-per-person",
        type=int,
        default=25,
        help="Max images per celebrity (default 25). Caps build time.",
    )
    p_build.add_argument("--force", action="store_true", help="Rebuild even if cache exists.")

    p_match = sub.add_parser("match", help="Match a face against the cached celebrity database.")
    add_dataset_args(p_match)
    p_match.add_argument("--camera", type=int, default=0, help="Camera index (default 0).")
    p_match.add_argument(
        "--image", type=str, default=None, help="Path to an image file instead of the webcam."
    )

    sub.add_parser("list", help="List datasets and whether their cache is built.")

    args = parser.parse_args(argv)

    if args.command == "build":
        return cmd_build(args)
    if args.command == "match":
        return cmd_match(args)
    if args.command == "list":
        return cmd_list(args)
    parser.error(f"Unknown command: {args.command}")
    return 2


if __name__ == "__main__":
    raise SystemExit(main())
