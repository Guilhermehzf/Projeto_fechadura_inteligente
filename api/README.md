
# 🚪 API Trancarduino (Backend em Go)

## 📖 Visão Geral

A **API Trancarduino** é um backend desenvolvido em **Go** para o sistema de fechadura inteligente, integrando o hardware com a interface do usuário através de **MQTT** e **PostgreSQL**. A API fornece autenticação de usuários, controle de estado da fechadura e comunicação em tempo real com o hardware (ESP32).

## 🏗️ Arquitetura

A arquitetura da API é composta por três componentes principais:

- **Frontend (Web/Mobile):** Consome os endpoints REST da API para autenticação e controle.
- **Banco de Dados (PostgreSQL):** Armazena os dados de usuários e o histórico de ações da fechadura.
- **Broker MQTT:** Envia comandos para o ESP32 e recebe atualizações de estado.

```plaintext
[ Frontend ] <-----> [ API (Go) ] <-----> [ PostgreSQL ]
                         ^  |
                         |  v
                   [ Broker MQTT ]
                         ^  |
                         |  v
                     [ ESP32 ]
```

## 🚀 Tecnologias Utilizadas

- **Linguagem:** Go (v1.20+)
- **Banco de Dados:** PostgreSQL
- **Comunicação IoT:** MQTT (`eclipse/paho.mqtt.golang`)
- **Servidor HTTP:** `net/http` (Go)
- **Autenticação:** JSON Web Tokens (JWT) (`golang-jwt/jwt/v5`)
- **Driver do Banco de Dados:** `jackc/pgx`
- **Criptografia:** `golang.org/x/crypto/bcrypt`
- **Configuração:** `joho/godotenv`

## 🔧 Configuração do Ambiente

### 1. Pré-requisitos

- **Go:** Versão 1.20 ou superior.
- **Docker** e **Docker Compose** para o banco de dados PostgreSQL.

### 2. Banco de Dados

1. **Iniciar o banco com Docker:**

    ```bash
    docker-compose up -d
    ```

2. **Criar as Tabelas:** Execute o script SQL para criar as tabelas `users` e `lock_history` no banco `tranca_inteligente`.

3. **Variáveis de Ambiente:**

    Crie um arquivo `.env` com as configurações a seguir.

    **Template para `.env`:**

    ```env
    # Configuração do Broker MQTT
    MQTT_BROKER="mqtts://<SEU_BROKER_HIVEMQ_URL>:8883"
    MQTT_CLIENT_ID="api-go-client-final"
    MQTT_USER="<SEU_USUARIO_MQTT>"
    MQTT_PASS="<SUA_SENHA_MQTT>"
    MQTT_TOPIC_COMMANDS="fechadura/comandos"
    MQTT_TOPIC_STATE="fechadura/estado"

    # Configuração do Servidor HTTP
    HTTP_ADDR="0.0.0.0:8088"

    # Configuração do Banco de Dados
    POSTGRES_ENDPOINT="postgres://admin:admin@localhost:5432/tranca_inteligente?sslmode=disable"

    # Configuração de Autenticação JWT
    JWT_SECRET="<SEU_JWT_SECRET_MUITO_SEGURO>"
    JWT_EXP_MINUTES=60

    # Limite de histórico em memória
    MAX_HISTORY=200
    ```

### 3. Dependências do Go

Baixe as dependências:

```bash
go mod tidy
```

## ⚡ Executando a API

### Modo de Desenvolvimento

Para rodar a API em modo de desenvolvimento com **hot-reload**:

```bash
go run ./cmd/server/main.go
```

A API estará disponível em `http://localhost:8088`.

### Build para Produção

Para compilar um executável otimizado:

```bash
go build -o smartlock ./cmd/server/
```

Execute o binário gerado:

```bash
./smartlock
```

## 🔑 Documentação dos Endpoints

<details>
<summary><strong>Clique para expandir/recolher a documentação da API</strong></summary>

### Autenticação

#### `POST /login`

Autentica um usuário e retorna um token JWT.

  - **Corpo da Requisição (`application/json`)**
    ```json
    {
      "email": "user@exemple.com",
      "password": "password"
    }
    ```

  - **Resposta de Sucesso (`200 OK`)**
    ```json
    {
        "success": true,
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "user": {
            "id": "c2a8f8e2-4b1a-4b8f-8c1e-7d9a1b3c4d5e",
            "email": "user@exemple.com",
            "name": "Nome do Usuário"
        }
    }
    ```

-----

### Rotas Protegidas (Exigem `Authorization: Bearer <token>`)

#### `GET /status`

Retorna o estado atual da fechadura. Suporta **long-polling**.

  - **Parâmetros de Query (Opcionais):**
    - `since` (RFC3339 timestamp): Para long-polling.
    - `timeout` (duração, ex: `25s`): Timeout do long-polling.

  - **Resposta de Sucesso (`200 OK`)**
    ```json
    {
        "isLocked": false,
        "isConnected": true,
        "lastUpdate": "2025-09-27T02:35:51Z"
    }
    ```

  - **Resposta com Timeout (`204 No Content`)**: Caso o long-poll expirar sem atualizações.

#### `POST /toggle`

Envia um comando para alternar o estado da fechadura.

  - **Corpo da Requisição:** Vazio.

  - **Resposta de Sucesso (`200 OK`)**
    ```json
    {
        "success": true,
        "isLocked": true
    }
    ```

</details>

## 📦 Como Contribuir

1. Fork o repositório.
2. Crie uma branch para sua feature (`git checkout -b minha-nova-feature`).
3. Faça as alterações necessárias.
4. Commit as mudanças (`git commit -m 'Adiciona nova funcionalidade'`).
5. Envie a branch para o repositório (`git push origin minha-nova-feature`).
6. Abra um pull request.

---

Feito com 💙 por [Guilherme Henrique Zioli](https://portfolio.ghzds.com.br/)
