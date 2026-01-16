# URL shortner
#### autor: svitorz
Essa aplicação foi construída com GoLang, VUE (JavaScript), Docker, PostgreSQL, Redis e NGINX.

Para executar, utilize os comandos:
```bash 
cp -r .env.example .env
```

Configure suas variávies.

> [!NOTE]
> Algumas variáveis de "host" já estão definidas, pois referem aos hosts configurados no arquivo docker-compose.yml

Assim que tudo estiver configurado, rode o comando:

```bash
docker compose up -d --build
```

Agora, acesse seu `localhost` na porta em você configurou.

## Rotas da aplicação:

### Usuários
> POST /user - cria um usuário

> PUT /user - atualiza

> DELETE /user - exclui

> POST /login - autentica

### Links:

> POST /links — cria link curto

> GET /links/:slug — detalhes do link

> POST /links/:slug — desativa

> GET /r/:slug — redireciona 

> [!TIP]
> Existe um arquivo `app.http.example` que pode ser utilizado como boilerplate para testes da aplicação.

# Contribuições
Contribuições são bem-vindas! Sinta-se à vontade para abrir issues e enviar pull requests.
# Licença 
Este projeto está licenciado sob a licença MIT. Consulte o arquivo LICENSE para mais detalhes.
