import pydantic
import modal
import pickle


class Request(pydantic.BaseModel):
    purchase_intention_max_zip_code_to_rfb_hamming: int
    company_max_relationship_age_days: int
    purchase_intention_div_amount_by_med_by_supplier_id: float
    snowplow_rnk_ip_treated_3d_by_company_fst_zip_code_2d: list[int]
    purchase_intention_div_amount_over_limit_by_med_by_supplier_id: float
    purchase_intention_max_hour_int: int
    login_cnt_purchase_intentions_all: int
    company_pct_share_capital_by_purchase_intention_fst_zip_code_2d: list[int]
    purchase_intention_pct_amount_by_purchase_intention_fst_zip_code_3d_and_company_fst_cnae_4d_and_supplier_id: list[
        int
    ]


app = modal.App("jojo-is-awesome", image=modal.Image.debian_slim().pip_install_from_requirements("ml/requirements.txt"))


# read a file from Google Cloud Storage and download it as bytes
# def read_file(project_id: str, bucket_name: str, blob_name: str) -> bytes:
#     client = storage.Client(project=project_id)
#     bucket = client.get_bucket(bucket_name)
#     blob = bucket.get_blob(blob_name)
#     return blob.download_as_bytes()


@app.function()
@modal.web_endpoint(method="POST", docs=True)
def endpoint(request: Request) -> str:
    # load the model
    f = open("ml/models/misc.pickle", "rb")
    model = pickle.load(f)
    # make a prediction
    return model.predict([])
