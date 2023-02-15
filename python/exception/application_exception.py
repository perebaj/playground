from loguru import logger
from google.cloud.logging.handlers import StructuredLogHandler
import sys
from fastapi import HTTPException

ENABLE_JSON_LOG = True


logger.remove()
if ENABLE_JSON_LOG:
    handler = StructuredLogHandler()
    logger.add(handler, level="INFO", serialize=True, enqueue=True, backtrace=True, diagnose=False, format="{message}", colorize=False)
else:
    logger.add(sys.stdout, serialize=False, enqueue=True, level="INFO", format="{time} - {level} - {message} - {extra}", colorize=False,  backtrace=True, diagnose=True,)


def message_handler():
    try:
        context_logger = logger.bind(message_id="1209301283190238", jojo="jojo")
        context_logger.info(f"Message handler")
        context_logger.info(f"Calling controller handler")
        suma = 1/0
    except:
        logger.opt(exception=True).error(f"Error while handling message")

message_handler()