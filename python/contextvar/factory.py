from fastapi import FastAPI


def setup_factory(app: FastAPI) -> None:
    app.include_router(api_router, prefix="/api")
