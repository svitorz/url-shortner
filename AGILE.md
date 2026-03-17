# Atividade do Laboratório

**Aluno:** Vitor Fábio de C. Souza  
**RA:** VP303156X  

---

## 1. Visão do Produto

Muitas das vezes, ao acessar a web, os usuários se deparam com links (URL’s) longas e com estrutura grande, o que pode causar estranhamento em usuários menos habituados a utilizar sistemas.

A ideia do projeto é criar um **Encurtador de URL**, com capacidade de se tornar modular e escalável, priorizando a arquitetura e infraestrutura do software, para que o sistema tenha alta disponibilidade e seja capaz de armazenar milhares de *slugs* que representam URL’s.

---

## 2. Product Backlog

### 1. Criação e configuração do repositório

O projeto ainda não foi criado, portanto, deve ser criado um repositório Git, configurando `.gitignore` e outras configurações. Definir o módulo de programação em GoLang, preparado para receber os pacotes, e criar a infraestrutura com Docker Compose.

**Critérios de aceite:**
- Criação e configuração do projeto Git
- Criação da configuração do `docker-compose` e `Dockerfile`, capaz de compilar os arquivos da pasta `cmd`
- Subir containers do Postgres e Nginx

---

### 2. Criação do domínio de Link e User

Nessa etapa, devem ser criadas duas *models* para o projeto, que representam o núcleo inicial do sistema.

**Critérios de aceite:**
- Criação da Model **Link**, com os campos:
  - ID
  - Slug
  - TargetURL
  - IsActive
  - Timestamps
- Criação da Model **User**, com os campos:
  - Username
  - Password

---

### 3. Criação da camada service e repository

Agora, é necessário definir a camada de trabalho com o banco de dados. Nessa fase, se faz necessário instalar o pacote **GoORM**, configuração de conexão com Postgres e criação das migrations.

**Critérios de aceite:**
- A aplicação se conecta com o banco de dados
- As tabelas **links** e **users** existem no banco de dados

---

### 4. Expor API

Próximo à fase final dessa sprint, agora é necessário criar a API para o projeto. Nessa fase, devem ser criadas as rotas (*resources*) dos links e a exposição dessa API via porta **8000**.

**Critérios de aceite:**
- Todas as rotas criadas para Link (GET, POST, PUT, DELETE)
- A API precisa ser acessada via `localhost:8000` sendo compilada corretamente
- CRUD completo do Link

---

### 5. Criar autenticação

Na fase final do MVP, é necessário criar as rotas de **cadastro**, **login** e a implementação da autenticação via middleware.

**Critérios de aceite:**
- Criação das rotas **register** e **login** funcionando
- Criação de middleware capaz de permitir ou bloquear requisições através do usuário autenticados