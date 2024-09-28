import pickle
import pydantic
import pandas as pd
import os
from fastapi import FastAPI


# LogisticType1Model represents the input data from the ML model LogisticType1Model
class LogisticType1Model(pydantic.BaseModel):
    # Distancy of Hamming between the zip code of the purchase intention and the zip code of the company's registration in the Federal Revenue (RFB).
    purchase_intention_max_zip_code_to_rfb_hamming: float
    # Age of the company in days.
    company_max_relationship_age_days: int
    # Purchase amount in relation to the median of the supplier
    purchase_intention_div_amount_by_med_by_supplier_id: float
    # Ip rank by company zip code
    snowplow_rnk_ip_treated_3d_by_company_fst_zip_code_2d: int
    purchase_intention_div_amount_over_limit_by_med_by_supplier_id: float
    purchase_intention_fst_email_domain_standardized: str
    purchase_intention_max_hour_int: int
    login_cnt_purchase_intentions_all: int
    company_pct_share_capital_by_purchase_intention_fst_zip_code_2d: float
    purchase_intention_pct_amount_by_purchase_intention_fst_zip_code_3d_and_company_fst_cnae_4d_and_supplier_id: float


# Prediction output of the LogisticType1Model
class LogisticType1ModelOutput(pydantic.BaseModel):
    pedicted_probability: list[float]


app = FastAPI()


@app.post("/predict", response_model=LogisticType1ModelOutput, summary="Predict the probability using the LogisticType1Model")
def read_root(body: LogisticType1Model) -> LogisticType1ModelOutput:
    f = open(os.path.join(os.path.dirname(__file__), "logistic_type_1_model.pickle"), "rb")
    model = pickle.load(f)
    f.close()

    input_data = pd.DataFrame([body.model_dump()])

    prediction = model.predict_proba(input_data)

    response = LogisticType1ModelOutput(pedicted_probability=prediction[0])

    response_json = response.model_dump()

    return response_json
