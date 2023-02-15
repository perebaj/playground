import logging
from pythonjsonlogger import jsonlogger
from google.cloud.logging.handlers import StructuredLogHandler

handler = StructuredLogHandler()

logger = logging.getLogger()

logHandler = logging.StreamHandler()
formatter = jsonlogger.JsonFormatter()
logHandler.setFormatter(formatter)
logger.addHandler(logHandler)
logger.setLevel(logging.INFO)
logger.info("LASKDJLASJDlpoet", extra={"message_id": "1234"})


context = logging.LoggerAdapter(logging.getLogger(), {'foo': 'bar'})

context.info("This is a debug message")