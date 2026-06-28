# Clima por CEP com Observabilidade

Sistema distribuído em Go (Serviço A + Serviço B) que retorna a temperatura de uma
cidade a partir do CEP, com tracing distribuído via OpenTelemetry + OTEL Collector + Zipkin.

## Como executar

```bash
cd Deploy
cp .env.example .env      # preencha API_KEY com sua chave da WeatherAPI
docker compose up --build
```

## Requisição (Serviço A)

`POST http://localhost:8080/weatherByCEP` com o CEP no corpo:

```bash
curl -X POST http://localhost:8080/weatherByCEP \
  -H "Content-Type: application/json" \
  -d '{"cep":"29902555"}'
```

Resposta de sucesso (`200 OK`):

```json
{ "city": "São Paulo", "temp_C": 28.5, "temp_F": 83.3, "temp_K": 301.5 }
```

Erros:

| Status | Mensagem               | Quando ocorre                           |
|--------|------------------------|-----------------------------------------|
| 422    | `invalid zipcode`      | CEP com formato inválido (≠ 8 dígitos)  |
| 404    | `can not find zipcode` | CEP com formato válido, mas inexistente |

## Visualizar os traços (Zipkin)

1. Faça ao menos uma requisição (curl acima).
2. Abra o Zipkin em <http://localhost:9411>.
3. Clique em **Run Query** e abra um trace para ver o fluxo `cep → weather`, incluindo
   os spans de busca de CEP e de temperatura.
