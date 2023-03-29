from fastapi import FastAPI
from starlette.middleware.base import BaseHTTPMiddleware, RequestResponseEndpoint
from fastapi import Request
from starlette.responses import Response
import contextvars
from loguru import logger
import asyncio
import time

organization_id_context = contextvars.ContextVar("organization_id_context", default=None)
trace_id_context = contextvars.ContextVar("trace_id_context", default=None)

class ContextVariableMiddleware(BaseHTTPMiddleware):
    async def dispatch(self, request: Request, call_next: RequestResponseEndpoint) -> Response:
        if "organization-id" in request.headers:
            organization_id_context.set(request.headers["organization-id"])
        if "trace-id" in request.headers:
            trace_id_context.set(request.headers["trace-id"])

        return await call_next(request)


app = FastAPI()

app.add_middleware(ContextVariableMiddleware)

@app.get("/")
async def root():
    await asyncio.sleep(2)
    logger.info(f"organization_id_context: {organization_id_context.get()}")
    logger.info(f"trace_id_context: {trace_id_context.get()}")
    return {"message": "Hello World"}


