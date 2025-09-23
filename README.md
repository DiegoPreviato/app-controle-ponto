# App Controle de Ponto

Este projeto é um sistema de controle de ponto de funcionários, desenvolvido com um backend em Go e um frontend em Flutter.

## Visão Geral

O objetivo deste projeto é fornecer uma solução simples para que os funcionários possam registrar suas horas de entrada e saída e para que os administradores possam gerenciar esses registros.

## Arquitetura

O sistema é dividido em duas partes principais:

### 1. Backend (`app_controle_ponto_backend`)

- **Linguagem:** Go (Golang)
- **Banco de Dados:** SQLite (o arquivo `ponto.db` na raiz do backend)
- **Funcionalidades:**
    - API REST para gerenciar usuários e registros de ponto.
    - Autenticação de usuário.
    - CRUD (Create, Read, Update, Delete) para os registros de ponto.

### 2. Frontend (`app_controle_ponto_frontend`)

- **Framework:** Flutter
- **Linguagem:** Dart
- **Plataforma:** Multi-plataforma (Web, Mobile - Android/iOS, Desktop)
- **Funcionalidades:**
    - Interface de usuário para login.
    - Tela para registrar (bater) o ponto.
    - Visualização do histórico de pontos registrados.

## Estrutura de Diretórios

```
/
├── app_controle_ponto_backend/   # Código-fonte do servidor Go
│   ├── database/
│   ├── handlers/
│   ├── middleware/
│   ├── models/
│   ├── go.mod
│   └── main.go
│
├── app_controle_ponto_frontend/  # Código-fonte do aplicativo Flutter
│   ├── lib/
│   ├── android/
│   ├── ios/
│   ├── web/
│   └── pubspec.yaml
│
└── README.md
```

**Nota:** Existem alguns arquivos Go (`main.go`, `database/`, etc.) na raiz do projeto. Eles parecem ser de uma fase inicial de desenvolvimento e devem ser consolidados ou removidos para manter o projeto organizado dentro de `app_controle_ponto_backend`.

## Como Começar

### Pré-requisitos

- [Go](https://golang.org/dl/) instalado.
- [Flutter SDK](https://flutter.dev/docs/get-started/install) instalado.

### Executando o Backend

1.  Navegue até o diretório do backend:
    ```sh
    cd app_controle_ponto_backend
    ```
2.  Instale as dependências:
    ```sh
    go mod tidy
    ```
3.  Inicie o servidor:
    ```sh
    go run main.go
    ```
    O servidor estará em execução em `http://localhost:8080`.

### Executando o Frontend

1.  Navegue até o diretório do frontend:
    ```sh
    cd app_controle_ponto_frontend
    ```
2.  Instale as dependências:
    ```sh
    flutter pub get
    ```
3.  Execute o aplicativo (escolha a plataforma desejada):
    ```sh
    flutter run
    ```