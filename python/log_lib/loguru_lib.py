from loguru import logger
from google.cloud.logging.handlers import StructuredLogHandler
import sys
ENABLE_JSON_LOG = False

def fmt(record):
    elapsed = record["elapsed"]
    level = record["level"]
    message = record["message"]

    msg = f"{elapsed} : {level} : "
    msgparts = []
    for key, val in record["extra"].items():
        msgparts.append(key + "=" + val)

    msgparts.append(message)
    msg += " ".join(msgparts)
    return msg + "\n"


if ENABLE_JSON_LOG:
    handler = StructuredLogHandler()
    logger.remove()
    logger.add(handler, level="INFO", serialize=True, enqueue=True)
else:
    logger.remove()
    logger.add(sys.stdout, serialize=False, enqueue=True, level="INFO", format="{time} - {level} - {message} - {extra}")


def message_handler():
    context_logger = logger.bind(message_id="1209301283190238", jojo="jojo")
    context_logger.info(f"Message handler")
    context_logger.info(f"Calling controller handler")

message_handler()