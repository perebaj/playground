from loguru import logger
from google.cloud.logging.handlers import StructuredLogHandler
import sys

handler = StructuredLogHandler()
   

logger.add(handler, level="INFO", serialize=True, enqueue=False)
context_logger = logger.bind(ip="192.168.0.1", user="someone")


if __name__ == "__main__":
    for _ in range(10):
        context_logger.info("Contextualize your logger easily")
        context_logger.error("ERROR")


    