# API de Cotação

Aplicação desenvolvida como desáfio proposto no curso Go Expert.

A aplicação possui um client e um server, o client chama o server no endpoint /cotacao, o server busca na api economia.awesomeapi.com.br a cotação do dolar, salva no sqlite e retorna a cotação para o client, que recebe a cotação e salva em um arquivo de texto chamado cotacao.txt

# Pré requisitos

- Golang 1.19+

# Como executar

- Primeiro executar o arquivo dentro do diretório server, com o comando `go run server.go`
- Com o servidor rodando, executar o arquivo dentro do diretório client com o comando `go run client.go`
- Após executar os passos a cima, dentro do diretório `client` deverá ser gerado um aquivo chamado `cotacao.txt` e na base `desafio.db` dentro do diretório `server` deve ser criado o registro com a contação do dolar