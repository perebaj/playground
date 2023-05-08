from opensearchpy import OpenSearch, OpenSearchException

host = "localhost"
port = 9200
auth = ("admin", "adminadmin")  # For testing only. Don't store credentials in code.

client = OpenSearch(
    hosts=[{"host": host, "port": port}],
    http_auth=auth,
    # client_cert = client_cert_path,
    # client_key = client_key_path,
    use_ssl=False,
)
print("teste")
try:
    index_name = "python-test-index"
    index_body = {"settings": {"index": {"number_of_shards": 4}}}

    response = client.delete(index=index_name, id=id)
    # response = client.indices.create(index_name, body=index_body)
except OpenSearchException as ex:
    # verify if the status code exists
    print(index_name)
    print("ERROR", ex)
    status_code = ex.status_code if hasattr(ex, "status_code") else None
    print("status_code: ", status_code)
    ex.with_traceback()

print(response)
