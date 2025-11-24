## 1# Criar o projeto e inicializar o módulo

```shell
mkdir app-notes
cd app-notes
go mod init bootstrap

```


## 2# Adicionar dependências

```shell
go get github.com/gin-gonic/gin@v1.10.0
go get gorm.io/gorm@v1.25.9
go get gorm.io/driver/postgres@v1.5.7
go get github.com/google/uuid@v1.6.0
go get github.com/aws/aws-lambda-go@v1.48.0
go get github.com/awslabs/aws-lambda-go-api-proxy@v0.16.2
```

## 3# Estrutura padrão
```jshelllanguage
mkdir -p cmd/api
mkdir -p cmd/lambda

mkdir -p internal/config
mkdir -p internal/db
mkdir -p internal/http
mkdir -p internal/notes

mkdir -p sql
mkdir -p tests/testdata
```