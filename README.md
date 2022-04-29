# alura-challenge-backend



## Descrição

Aplicação para upload e análise de transações financeiras. Suporte arquivos .csv com os campos abaixo (não utilizar cabeçalho):

banco-origem
agencia-origem
conta-origem
banco-destino
agencia-destino
conta-destino
valor
datahora-transacao

### Detalhes das atividades do Projeto

Finalizado:

* [Semana 1](https://trello.com/b/6BVMlCYd/challenge-backend-3-semana-1)

Não Iniciado:

* [Semana 2](https://trello.com/b/nUN64cpL/challenge-backend-3-semana-2)

* [Semanas 3 e 4](https://trello.com/b/Z5fKD7ly/challenge-backend-3-semana-3)

## Stack

Desenvolvido utilizando:

* GO
* Echo
* MongoDB

## Como Desenvolver

Para desenvolver as próximas semanas deste projeto, é necessário possuir [docker](https://www.docker.com/products/docker-desktop/) e [go](https://go.dev/learn/) instalados no computador.

* Iniciando a aplicação.

```bash
docker-compose up -d

go run cmd/main.go
```

* Acessando os recursos:

No browser, acessar `http://localhost/upload` para enviar o arquivo .csv. Os resultados de import ficam no mesmo arquivo.

* Limpando os recursos criados

```bash
#fechar a aplicação, pressionando ctrl+c

docker-compose down -v 
```