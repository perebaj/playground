import structlog
from google.cloud.logging.handlers import StructuredLogHandler
# import loggings
import logging
handler = StructuredLogHandler()

processors= [
        structlog.stdlib.filter_by_level,
        structlog.stdlib.add_logger_name,
        structlog.stdlib.add_log_level,
        structlog.processors.dict_tracebacks,
        structlog.processors.JSONRenderer()
    ]

# s
structlog.configure(wrapper_class=structlog.make_filtering_bound_logger(logging.INFO), processors=[
        structlog.stdlib.add_log_level,
        structlog.processors.dict_tracebacks,
        structlog.processors.JSONRenderer()
    ]
)
context_log = structlog.get_logger("Structured Logger").bind(message_id='1234')
context_log.info("This is a test message")


# logging.basicConfig(level=logging.DEBUG, handlers=[handler])
# logger = logging.getLogger()
# logger.info("This is a test message", extra={"message_id": "1234"})