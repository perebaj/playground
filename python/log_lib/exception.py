from loguru import logger
from google.cloud.logging.handlers import StructuredLogHandler

logger.remove()

handler = StructuredLogHandler()
logger.add(
        handler,
        level="INFO",
        serialize=True,
        enqueue=False,
        format="{message}",
        colorize=False,
        backtrace=True,
        diagnose=False,
    )


try:
    a = 10
    b = 0
    print(a/b)
except Exception as e:
    # logger.error("Error occurred", exc_info=e)
    logger.opt(exception=True).error("Error occurred")