routing_map = {
    "408": {"destination_queue": "cred-ext_data-serasa_pf", "description": "Timeout"},
    "200": {"destination_queue": "cred-ext_data-ext_enrich_reprocessor", "description": "Success"},
    "*": {
        "destination_queue": "cred-ext_data-serasa_pf",
        "description": "Qualquer outro caso que não seja sucesso, deve ser enviado para serasa_pf",
    },
}

for key, value in routing_map.items():
    print(key, value)
