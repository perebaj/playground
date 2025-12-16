Anatomia do Subject

og.*.in.otlp.proto.logs
│  │  │   │    │    │
│  │  │   │    │    └── Tipo de dado (logs, traces, metrics)
│  │  │   │    └─────── Formato (proto = protobuf)
│  │  │   └──────────── Protocolo (otlp = OpenTelemetry)
│  │  └──────────────── Direção (in = entrada, out = saída)
│  └─────────────────── Wildcard (* = qualquer org)
└────────────────────── Namespace (og = JJ's Corporation)

Convenção Hierárquica
A estrutura vai do mais geral para o mais específico:

[namespace].[tenant].[direção].[protocolo].[formato].[tipo]
Nível	Exemplo	Propósito
1	og	Namespace/produto
2	* ou acme	Tenant/organização
3	in	Direção do fluxo
4	otlp	Protocolo usado
5	proto	Formato de serialização
6	logs	Tipo de telemetria
Wildcards

# * = um token qualquer
og.*.in.otlp.proto.traces     # Qualquer org, apenas traces
og.acme.*.otlp.proto.traces   # Org acme, qualquer direção

# > = zero ou mais tokens (só no final)
og.*.in.otlp.>                # Qualquer coisa após otlp (traces, metrics, logs)
og.acme.>                     # Tudo da org acme
Por que essa estrutura?
1. Filtragem eficiente

# Consumir apenas traces de uma org
og.acme.in.otlp.proto.traces

# Consumir tudo de uma org
og.acme.>

# Consumir todos os logs de todas as orgs
og.*.in.otlp.proto.logs
2. Isolamento por tenant

og.acme.in.otlp.proto.traces   # Dados da Acme
og.corp.in.otlp.proto.traces   # Dados da Corp

