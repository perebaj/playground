import logging
from google.cloud.logging.handlers import StructuredLogHandler

handler = StructuredLogHandler()
logging.basicConfig(level=logging.DEBUG, handlers=[handler], format='%(name)s - %(levelname)s - %(message)s - %(foo)s')

context = logging.LoggerAdapter(logging.getLogger(), {'foo': 'bar'})
context.setLevel("DEBUG")
context.debug("This is a debug message")


