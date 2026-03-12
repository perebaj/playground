"""Minimal OTLP HTTP receiver that logs request payload sizes and decompression time.

Listens on four ports, one per compression format:
  - :4318 — gzip
  - :4319 — none (uncompressed)
  - :4320 — snappy
  - :4321 — zstd

Reports wire size, decompressed size, compression ratio, and decompression time.
"""

import gzip
import threading
import time
import zlib
from http.server import HTTPServer, BaseHTTPRequestHandler

snappy_decode = None
zstd_decompress = None

try:
    import cramjam
    snappy_decode = lambda data: bytes(cramjam.snappy.decompress_raw(data))
    zstd_decompress = lambda data: bytes(cramjam.zstd.decompress(data))
except ImportError:
    pass

PORT_LABELS = {
    4318: "gzip",
    4319: "none",
    4320: "snappy",
    4321: "zstd",
}


class PayloadSizerHandler(BaseHTTPRequestHandler):
    def do_POST(self):
        length = int(self.headers.get("Content-Length", 0))
        body = self.rfile.read(length)
        encoding = self.headers.get("Content-Encoding", "none")
        port = self.server.server_address[1]
        label = PORT_LABELS.get(port, str(port))

        # Extract signal type from path (/v1/logs, /v1/traces, /v1/metrics)
        signal = self.path.rsplit("/", 1)[-1] if "/" in self.path else "unknown"

        decompressed_size = length
        decompress_us = 0

        if encoding == "gzip":
            t0 = time.monotonic()
            decompressed_size = len(gzip.decompress(body))
            decompress_us = (time.monotonic() - t0) * 1_000_000
        elif encoding == "zlib":
            t0 = time.monotonic()
            decompressed_size = len(zlib.decompress(body))
            decompress_us = (time.monotonic() - t0) * 1_000_000
        elif encoding == "snappy" and snappy_decode:
            t0 = time.monotonic()
            decompressed_size = len(snappy_decode(body))
            decompress_us = (time.monotonic() - t0) * 1_000_000
        elif encoding == "zstd" and zstd_decompress:
            t0 = time.monotonic()
            decompressed_size = len(zstd_decompress(body))
            decompress_us = (time.monotonic() - t0) * 1_000_000

        ratio = f"{decompressed_size / length:.1f}x" if length > 0 else "N/A"

        decompress_str = f"{decompress_us:.0f}us" if decompress_us > 0 else "N/A"

        print(
            f"[{label:>6}] [{signal:>7}] "
            f"wire: {length:>12,} bytes ({length / 1024 / 1024:>6.2f} MB) | "
            f"decompressed: {decompressed_size:>12,} bytes ({decompressed_size / 1024 / 1024:>6.2f} MB) | "
            f"ratio: {ratio:>6} | "
            f"decompress: {decompress_str}",
            flush=True,
        )

        self.send_response(200)
        self.send_header("Content-Type", "application/json")
        self.end_headers()
        self.wfile.write(b"{}")

    def log_message(self, format, *args):
        pass


def run_server(port):
    server = HTTPServer(("0.0.0.0", port), PayloadSizerHandler)
    print(f"Payload sizer listening on :{port} ({PORT_LABELS.get(port, '?')})", flush=True)
    server.serve_forever()


if __name__ == "__main__":
    for port in [4318, 4319, 4320]:
        threading.Thread(target=run_server, args=(port,), daemon=True).start()
    run_server(4321)
