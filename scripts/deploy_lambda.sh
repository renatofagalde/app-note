#!/usr/bin/env bash
set -euo pipefail

LAMBDA_NAME="${LAMBDA_NAME:?LAMBDA_NAME não definido}"
AWS_REGION="${AWS_REGION:-us-east-1}"
LAMBDA_ROLE_NAME="${LAMBDA_ROLE_NAME:-lambda-exec-${LAMBDA_NAME}}"
ZIP_PATH="${1:-function.zip}"

echo ">> Região: ${AWS_REGION}"
echo ">> Função: ${LAMBDA_NAME}"
echo ">> Role da Lambda: ${LAMBDA_ROLE_NAME}"

export AWS_REGION

# 1) Criar role se não existir
echo ">> Verificando role ${LAMBDA_ROLE_NAME}..."
if aws iam get-role --role-name "${LAMBDA_ROLE_NAME}" >/dev/null 2>&1; then
  echo "✅ Role já existe."
else
  echo "⚠️ Role não existe. Criando..."

  TRUST_POLICY_FILE="$(mktemp)"
  cat > "${TRUST_POLICY_FILE}" <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF

  aws iam create-role \
    --role-name "${LAMBDA_ROLE_NAME}" \
    --assume-role-policy-document "file://${TRUST_POLICY_FILE}"

  aws iam attach-role-policy \
    --role-name "${LAMBDA_ROLE_NAME}" \
    --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole

  echo "✅ Role criada com sucesso."
fi

# 2) Obter ARN da role
LAMBDA_ROLE_ARN="$(aws iam get-role --role-name "${LAMBDA_ROLE_NAME}" --query 'Role.Arn' --output text)"
echo ">> Role ARN: ${LAMBDA_ROLE_ARN}"

# 3) Criar ou atualizar Lambda
echo ">> Verificando se Lambda ${LAMBDA_NAME} existe..."
if aws lambda get-function --function-name "${LAMBDA_NAME}" >/dev/null 2>&1; then
  echo "✅ Função existe. Atualizando código..."
  aws lambda update-function-code \
    --function-name "${LAMBDA_NAME}" \
    --zip-file "fileb://${ZIP_PATH}"
else
  echo "⚠️ Função não existe. Criando..."
  aws lambda create-function \
    --function-name "${LAMBDA_NAME}" \
    --runtime "provided.al2" \
    --handler "bootstrap" \
    --role "${LAMBDA_ROLE_ARN}" \
    --architectures "x86_64" \
    --zip-file "fileb://${ZIP_PATH}"

  echo ">> Aguardando função sair do estado 'Pending'..."

  # espera até ficar Active (ou até estourar o timeout simples)
  for i in {1..20}; do
    state="$(aws lambda get-function-configuration \
      --function-name "${LAMBDA_NAME}" \
      --query 'State' \
      --output text || echo 'Desconhecido')"

    echo "   Estado atual: ${state}"

    if [ "${state}" = "Active" ]; then
      echo "✅ Função está Active."
      break
    fi

    if [ "${state}" = "Failed" ]; then
      echo "❌ Função entrou em estado Failed. Abortando."
      exit 1
    fi

    sleep 5
  done
fi


# 4) Atualizar env vars
echo ">> Atualizando variáveis de ambiente..."
aws lambda update-function-configuration \
  --function-name "${LAMBDA_NAME}" \
  --environment "Variables={
    APP_PORT=\"${APP_PORT}\",
    POSTGRES_HOST=\"${POSTGRES_HOST}\",
    POSTGRES_PORT=\"${POSTGRES_PORT}\",
    POSTGRES_USER=\"${POSTGRES_USER}\",
    POSTGRES_PASSWORD=\"${POSTGRES_PASSWORD}\",
    POSTGRES_DB=\"${POSTGRES_DB}\",
    POSTGRES_SSLMODE=\"${POSTGRES_SSLMODE}\",
    APP_ENV=\"${APP_ENV:-dev}\"
  }"

echo "✅ Deploy finalizado para ${LAMBDA_NAME}"

echo "✅ Deploy finalizado para ${LAMBDA_NAME}"
