"""Face detection and embedding.

MTCNN (Zhang et al. 2016) for detection + alignment, and InceptionResnetV1
pretrained on VGGFace2 for the 512-d embedding. The embedding output is
L2-normalized, so cosine similarity reduces to a dot product.
"""

from __future__ import annotations

import numpy as np
import torch
from facenet_pytorch import MTCNN, InceptionResnetV1
from PIL import Image

# MTCNN runs MPS-unsafe ops (adaptive_avg_pool2d with non-divisible sizes
# during its image pyramid), so pin it to CPU even when the embedder gets MPS.
MTCNN_DEVICE = "cpu"


def _pick_device() -> str:
    if torch.cuda.is_available():
        return "cuda"
    if torch.backends.mps.is_available():
        return "mps"
    return "cpu"


class FaceProcessor:
    """Detect a face in an image and produce a 512-d embedding."""

    def __init__(self, device: str | None = None, image_size: int = 160) -> None:
        self.device = device or _pick_device()
        self.detector = MTCNN(
            image_size=image_size,
            margin=0,
            keep_all=False,
            post_process=True,
            device=MTCNN_DEVICE,
        )
        self.embedder = InceptionResnetV1(pretrained="vggface2").eval().to(self.device)

    def detect(self, image: np.ndarray | Image.Image) -> torch.Tensor | None:
        if isinstance(image, np.ndarray):
            image = Image.fromarray(image)
        return self.detector(image)

    @torch.no_grad()
    def embed(self, face_tensor: torch.Tensor) -> np.ndarray:
        if face_tensor.dim() == 3:
            face_tensor = face_tensor.unsqueeze(0)
        face_tensor = face_tensor.to(self.device)
        emb = self.embedder(face_tensor).cpu().numpy()[0]
        return emb / np.linalg.norm(emb)

    def detect_and_embed(self, image: np.ndarray | Image.Image) -> np.ndarray | None:
        face = self.detect(image)
        if face is None:
            return None
        return self.embed(face)
