from loguru import logger
from google.cloud.logging.handlers import StructuredLogHandler
import sys
from fastapi import HTTPException

ENABLE_JSON_LOG = False

class ApplicationException(HTTPException):
    def __init__(
        self,
        key: str,
        message: str = None,
        details: any = None,
        status_code: int = None,
    ):
        """
        Construtor.
        - key: identificador unico do erro dentro da aplicação.
        - message: Mensagem do erro.
        - details: Detalhes do erro.
        - http_status_code: Código do erro HTTP (valor padrão 500).
        """
        self.key = key if key is not None else "application_with_error"
        self.message = message
        self.details = details
        self.status_code = status_code if status_code is not None else 500

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
        if True:
            raise ApplicationException("application_exception", "Application exception", "Details", 500)

    except:
        logger.opt(exception=True).error(f"Error while handling message")

message_handler()