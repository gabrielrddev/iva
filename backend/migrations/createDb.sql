-- Cria a tabela de usuários
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    phone TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    number_question INTEGER NOT NULL DEFAULT 0
);

-- Exemplo de inserção de usuário
INSERT INTO users (name, email, phone, password, number_question)
VALUES ('Gabriel Rodrigues Dias', 'gabriel@email.com', '123456789', '1234', 1);
