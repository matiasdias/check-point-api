CREATE TABLE IF NOT EXISTS funcionario (
    id serial PRIMARY KEY, 
    is_admin boolean DEFAULT false,
    nome varchar(250) NOT NULL, 
    email varchar(100) NOT NULL UNIQUE, 
    telefone varchar(20) NOT NULL,
    cargo varchar(50) NOT NULL,
    genero varchar(1) NOT NULL,
	idade integer NOT NULL,   
    cpf varchar(15) NOT NULL UNIQUE,
    senha varchar(100) NOT NULL,
    criadoem timestamp DEFAULT now() NOT NULL,
    updateem timestamp DEFAULT now() NOT NULL
);

INSERT INTO funcionario (nome, email, telefone, senha, idade, cpf, cargo, genero, is_admin) 
VALUES
('Matias Dias', 'matias@gmail.com', '(88) 99777-4641', '$2y$10$/tnAheWxR7TFmFg7M7UOUudkTGTuhwHFVigILJ.S2k5YizbXYbR3e', 26, '071.749.583-38', 'Desenvolvedor back-end', 'M', true);