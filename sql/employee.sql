CREATE TABLE IF NOT EXISTS funcionario (
    id serial PRIMARY KEY, 
    is_admin boolean DEFAULT false,
    nome varchar(250) NOT NULL, 
    email varchar(100) NOT NULL UNIQUE, 
    telefone varchar(11) NOT NULL,
    cargo varchar(50) NOT NULL,
	idade integer NOT NULL,   
    cpf varchar(11) NOT NULL UNIQUE,
    senha varchar(100) NOT NULL,
    criadoem timestamp DEFAULT current_timestamp NOT NULL,
    updateem timestamp DEFAULT current_timestamp NOT NULL
);

INSERT INTO funcionario (nome, email, telefone, senha, idade, cpf, cargo, is_admin) 
VALUES
('Matias Dias', 'matias@gmail.com', '88997774641', '$2y$10$/tnAheWxR7TFmFg7M7UOUudkTGTuhwHFVigILJ.S2k5YizbXYbR3e', 26, '07174958338', 'Desenvolvedor back-end', true);