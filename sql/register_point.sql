CREATE TABLE IF NOT EXISTS registro_ponto (
    id serial PRIMARY KEY,
    codigo_funcionario integer NOT NULL,
    criado_em timestamp with time zone NOT NULL DEFAULT current_timestamp,
    hora_entrada timestamp with time zone,
    hora_saida timestamp with time zone,
    tipo_registro varchar(10) NOT NULL,
    CONSTRAINT registro_ponto_pk PRIMARY KEY (id),
    CONSTRAINT fk_codigo_funcionario FOREIGN KEY (codigo_funcionario) REFERENCES funcionario(id)
);

-- Criar o trigger function
CREATE OR REPLACE FUNCTION atualizar_horario_ponto()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.tipo_registro = 'entrada' THEN
        NEW.hora_entrada = CURRENT_TIMESTAMP;
    ELSIF NEW.tipo_registro = 'saida' THEN
        NEW.hora_saida = CURRENT_TIMESTAMP;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Criar o trigger na tabela registro_ponto
CREATE TRIGGER atualizar_horario_trigger
BEFORE INSERT ON registro_ponto
FOR EACH ROW EXECUTE FUNCTION atualizar_horario_ponto();