-- DDL: tabela notes
CREATE TABLE IF NOT EXISTS notes (
    id          TEXT         PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    content     JSONB        NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ  NULL
);

CREATE INDEX IF NOT EXISTS idx_notes_created_at ON notes (created_at);
CREATE INDEX IF NOT EXISTS idx_notes_name ON notes (LOWER(name));

-- INSERT de exemplo 1
INSERT INTO notes (
    id, name, content, created_at, updated_at, deleted_at
) VALUES (
    '11111111-1111-1111-1111-111111111111',
    'Home - Meu Site',
    $${
      "url": "https://meusite.com/",
      "contentType": "text/html",
      "html": "<!DOCTYPE html><html><head><title>Meu Site</title></head><body><h1>Bem-vindo ao Meu Site</h1><p>Conteúdo qualquer...</p></body></html>"
    }$$::jsonb,
    NOW(),
    NOW(),
    NULL
);

-- INSERT de exemplo 2
INSERT INTO notes (
    id, name, content, created_at, updated_at, deleted_at
) VALUES (
    '22222222-2222-2222-2222-222222222222',
    'Artigo - Observabilidade em Go',
    $${
      "url": "https://blog.meusite.com/observabilidade-em-go",
      "contentType": "text/html",
      "html": "<html><head><title>Observabilidade em Go</title></head><body><h1>Observabilidade em Go</h1><p>Logs, métricas e traces são pilares importantes...</p><pre><code>func main() { /* ... */ }</code></pre></body></html>",
      "metadata": {
        "author": "Renato",
        "tags": ["go", "observabilidade", "logs"],
        "language": "pt-BR"
      }
    }$$::jsonb,
    NOW(),
    NOW(),
    NULL
);

-- INSERT de exemplo 3
INSERT INTO notes (
    id, name, content, created_at, updated_at, deleted_at
) VALUES (
    '33333333-3333-3333-3333-333333333333',
    'Dump HTML - Página de Login',
    $${
      "url": "https://app.meusite.com/login",
      "contentType": "text/html",
      "blob": "<html><head><title>Login</title></head><body><form><input type='email' name='email'/><input type='password' name='password'/></form></body></html>"
    }$$::jsonb,
    NOW(),
    NOW(),
    NULL
);
