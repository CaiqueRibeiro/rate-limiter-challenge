# Rate Limiter - Go Expert Challenge

Primeiro desafio do MBA Go Expert - implementação de rate limiter para API que limita o acesso dos usuários por determinado periódo de tempo.

O Rate limiter, por padrão, limita acessos a partir do IP cliente, mas pode conter configuração de tokens de acesso:
1. IP: Ao configurar a variável de ambiente `IP_MAX_REQUESTS`, cada usuário terá uma quantidade igual de requests possíveis no intervalo definido na variável `LIMIT_TIME_WINDOW_MS`.
2. Token de acesso: Ao configurar um ou mais tokens de acesso, cada token poderá possuir determinada quantidade de requests no intervalo definido na variável `LIMIT_TIME_WINDOW_MS`. Esta configuração é feita a partir do módulo CLI contido no projeto.

## Variáveis de Ambiente

- `WEB_SERVER_PORT`: Porta onde será executada a API
- `REDIS_HOST`: Host do Redis, responsável por manter os limites e tokens.
- `REDIS_PORT`: Porta do Redis. Normalmente 6379
- `REDIS_PASSWORD`: Senha do Redis, pode ser executado com `""`
- `REDIS_DB`: Normalmente configurado como 0
- `IP_MAX_REQUESTS`: Máximo de requests que cada IP poderá realizar
- `LIMIT_TIME_WINDOW_MS`: Intervalo em milisegundos para o refresh do limiter (IP e Token)

## Como executar o projeto

1. Crie o arquivo .env na raíz do projeto e popule as variáveis de ambiente. É possível usar os valores de .env.example sem problemas
```
WEB_SERVER_PORT=8080
REDIS_HOST="localhost"
REDIS_PORT=6379
REDIS_PASSWORD=""
REDIS_DB=0
IP_MAX_REQUESTS=10
LIMIT_TIME_WINDOW_MS=1000
```

2. Execute os containers Docker com `docker compose up -d`. Isso inicializará a API e o Redis.

3. Execute chamadas de API em seu gerenciador de preferência para o endpoint **GET** `http://localhost:8080`.
    3.1 Para o limiter de IP, nenhum `body` ou `header` é necessário.
    3.2 Para o limiter por token, é preciso informá-lo no request com o nome `API_KEY`.

## Como cadastrar um token

### Por API
Para registrar um token a partir do endpoint da API, basta enviar uma requisição POST `http://localhost:8080/token` com o `body`:
```
{
    "token": "TOKEN_DESEJADO",
    "max_requests": NUMERO_DE_REQUESTS
}

// Exemplo:
{
	"token": "32bc525d-1535-4025-b070-8f1b51e8f5a1",
	"max_requests": 100
}
```

Caso o token já possua um registro de requisições máximas, uma nova chamada irá sobrescrever esta quantidade.

### A partir do CLI
Para registrar um token a partir da CLI, execute:
```
make save-token token=TOKEN_DESEJADO maxreq=NUMERO_DE_REQUESTS
```

## Como rodar os testes
Para executar os testes, execute:
```
make test
```