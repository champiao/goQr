# :skull: README - Excel QR Code Generator (Go)

## :eyes: O que este programa faz
Este programa em Go lê um arquivo Excel (`.xlsx`) e, para cada linha da planilha, cria um QR code contendo todas as colunas e valores dessa linha. Cada QR code é salvo em um PDF separado, centralizado na página, com o texto correspondente abaixo.

## 🛠 Como funciona passo a passo
1. **:mag: Verifica o arquivo Excel**: Confere se o arquivo fornecido existe.
2. **:notebook: Lê o Excel**: Utiliza a biblioteca `excelize` para ler o conteúdo da planilha.
3. **:open_file_folder: Cria a pasta de saída**: Gera a pasta `qrcodes_output` para armazenar os PDFs.
4. **:wrench: Gera o QR Code**: Usa a biblioteca `go-qrcode` para criar um QR code a partir do conteúdo da linha.
5. **:page_facing_up: Gera o PDF**: Insere o QR code no centro da página e imprime o texto abaixo usando a biblioteca `gofpdf`.
6. **:x: Remove arquivos temporários**: Exclui o PNG após inseri-lo no PDF.

## 📋 Pré-requisitos
- Go 1.24.5 ou superior instalado.
- Dependencias:
```bash
go mod tidy
```

## ▶ Como executar
:rocket: No terminal, rode o comando:
```bash
go run cmd/main.go path/to/archive.xlsx
```
Substitua `path/to/archive.xlsx` pelo caminho do seu Excel.

## :heavy_check_mark: Saída
- PDFs gerados na pasta `qrcodes_output`.
- Nomeados como `X.pdf` onde `X` é o número do codigo presente em cada linha da primeira coluna no Excel.

## 💡 Exemplo
Se o Excel tiver:
| Planta| e-mail to  | REP      | Área| Localização |
|-------|------------|----------|
| 0001  | D3@test.com| 4        |
| 0002  | D9@test.com| 3        |

O programa irá gerar:
- `0001.pdf` com QR code e dados de 0001.
- `0002.pdf` com QR code e dados de 0002.
