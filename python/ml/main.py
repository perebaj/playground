import pydantic
import modal


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


app = modal.App("jojo-is-awesome")


@app.function()
@modal.web_endpoint(method="POST", docs=True)
def endpoint(request: Request) -> str:
    print(request)
    return "Hello, World!"
