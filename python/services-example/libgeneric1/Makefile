.PHONY: poetry-config
poetry-config:
	poetry config repositories.gcp https://us-central1-python.pkg.dev/blup-432616/libtest
	poetry config http-basic.gcp oauth2accesstoken $$(gcloud auth print-access-token)

.PHONY: publish
publish: poetry-config
	poetry publish -r gcp --build
