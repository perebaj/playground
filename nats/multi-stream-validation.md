
# Guia de Testes Locais: Streams NATS Separados

Guia simples para testar streams NATS separados por tipo de telemetria.

---

## 1. Instalação

```bash
# macOS
brew install nats-server nats-io/nats-tools/nats

# Verificar
nats-server --version
nats --version
```

---

## 2. Iniciar NATS Server

```bash
# Iniciar com JetStream habilitado
nats-server --jetstream

# Verificar (em outro terminal)
nats account info
```

---

## 3. Criar Streams

```bash
# Stream para TRACES
nats stream add og-traces \
  --subjects="og.*.in.otlp.proto.traces" \
  --retention=limits \
  --storage=file \
  --max-age=24h \
  --defaults

# Stream para METRICS
nats stream add og-metrics \
  --subjects="og.*.in.otlp.proto.metrics" \
  --retention=limits \
  --storage=file \
  --max-age=12h \
  --defaults

# Stream para LOGS
nats stream add og-logs \
  --subjects="og.*.in.otlp.proto.logs" \
  --retention=limits \
  --storage=memory \
  --max-age=6h \
  --defaults

# Stream INTERNO (ripeness + insights)
nats stream add og-internal \
  --subjects="og.*.internal.>,og.*.new-insight.>" \
  --retention=workqueue \
  --storage=memory \
  --max-age=1h \
  --defaults
```

Verificar:

```bash
nats stream ls
```

---

## 4. Criar Consumers

```bash
# Consumer de traces (usado pelo gaps)
nats consumer add og-traces gaps-traces \
  --filter="og.*.in.otlp.proto.traces" \
  --ack=explicit \
  --deliver=last \
  --defaults

# Consumer de ripeness (usado pelo gaps)
nats consumer add og-internal gaps-ripeness \
  --filter="og.*.internal.gaps.check-ripeness" \
  --ack=explicit \
  --deliver=all \
  --defaults
```

Verificar:

```bash
nats consumer ls og-traces
nats consumer ls og-internal
```

---

## 5. Testar Isolamento

Publicar mensagens de diferentes tipos:

```bash
nats pub "og.acme.in.otlp.proto.traces" "trace-1"
nats pub "og.acme.in.otlp.proto.metrics" "metric-1"
nats pub "og.acme.in.otlp.proto.logs" "log-1"
nats pub "og.acme.internal.gaps.check-ripeness" "ripeness-1"
```

Verificar que cada stream recebeu apenas seu tipo:

```bash
nats stream report
```

**Resultado esperado:**

```
Stream       Messages
og-traces    1
og-metrics   1
og-logs      1
og-internal  1
```

---

## 6. Testar Consumer

Publicar mais traces:

```bash
nats pub "og.acme.in.otlp.proto.traces" "trace-2"
nats pub "og.acme.in.otlp.proto.traces" "trace-3"
```

Consumir do consumer de traces:

```bash
nats consumer next og-traces gaps-traces --count=5
```

O consumer deve receber apenas as mensagens de traces.

---

## 7. Testar Múltiplas Organizações

```bash
nats pub "og.acme.in.otlp.proto.traces" "trace-acme"
nats pub "og.corp.in.otlp.proto.traces" "trace-corp"
nats pub "og.startup.in.otlp.proto.traces" "trace-startup"
```

Verificar:

```bash
nats stream info og-traces
```

Todas as orgs vão para o mesmo stream `og-traces`.

---

## 8. Ver Mensagens

```bash
# Ver últimas mensagens do stream
nats stream view og-traces

# Ver mensagens de uma org específica
nats stream view og-traces --subject="og.acme.in.otlp.proto.traces"
```

---

## 9. Monitorar

```bash
# Status geral
nats stream report
nats consumer report

# Detalhes de um stream
nats stream info og-traces

# Detalhes de um consumer
nats consumer info og-traces gaps-traces
```

---

## 10. Limpar

```bash
# Limpar mensagens (mantém stream)
nats stream purge og-traces -f
nats stream purge og-metrics -f
nats stream purge og-logs -f
nats stream purge og-internal -f

# Deletar tudo
nats stream delete og-traces -f
nats stream delete og-metrics -f
nats stream delete og-logs -f
nats stream delete og-internal -f
```

Parar o servidor: `Ctrl+C` no terminal do nats-server.

---

## Resumo de Comandos

| Ação | Comando |
|------|---------|
| Iniciar NATS | `nats-server --jetstream` |
| Listar streams | `nats stream ls` |
| Info do stream | `nats stream info og-traces` |
| Publicar | `nats pub "og.acme.in.otlp.proto.traces" "data"` |
| Consumir | `nats consumer next og-traces gaps-traces` |
| Ver mensagens | `nats stream view og-traces` |
| Status | `nats stream report` |
