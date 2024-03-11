CREATE TABLE IF NOT EXISTS funcionario (
    id serial PRIMARY KEY, 
    nome varchar(50) NOT NULL, 
    email varchar(50) NOT NULL UNIQUE, 
    telefone varchar(11) NOT NULL,
    cargo varchar(50) NOT NULL,
	idade integer NOT NULL,   
    cpf varchar(11) NOT NULL UNIQUE,
    senha varchar(100) NOT NULL,
    criadoem timestamp DEFAULT current_timestamp NOT NULL,
    updateem timestamp DEFAULT current_timestamp NOT NULL
);