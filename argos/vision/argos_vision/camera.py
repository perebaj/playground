"""Webcam capture with a live preview."""

from __future__ import annotations

import cv2
import numpy as np

WINDOW = "Argos — webcam (SPACE to capture, ESC to quit)"


def capture_from_webcam(camera_index: int = 0) -> np.ndarray | None:
    """Open the webcam, show a preview, and return one captured RGB frame.

    Returns the frame on SPACE, or None if the user hits ESC or closes the window.
    """
    cap = cv2.VideoCapture(camera_index)
    if not cap.isOpened():
        raise RuntimeError(f"Could not open camera index {camera_index}")

    captured: np.ndarray | None = None
    try:
        while True:
            ok, frame_bgr = cap.read()
            if not ok:
                continue
            cv2.imshow(WINDOW, frame_bgr)
            key = cv2.waitKey(1) & 0xFF
            if key == 27:  # ESC
                break
            if key == 32:  # SPACE
                captured = cv2.cvtColor(frame_bgr, cv2.COLOR_BGR2RGB)
                break
    finally:
        cap.release()
        cv2.destroyAllWindows()

    return captured
