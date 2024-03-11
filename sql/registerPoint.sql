CREATE TABLE IF NOT EXISTS registro_ponto (
    id serial PRIMARY KEY,
    codigo_funcionario integer NOT NULL,
    criado_em timestamp with time zone NOT NULL DEFAULT current_timestamp,
    tipo_registro varchar(15) NOT NULL,  
    hora_entrada time,
    hora_saida time,
    hora_entrada_almoco time,  
    hora_retorno_almoco time,
    horas_trabalhadas interval,
    CONSTRAINT fk_codigo_funcionario FOREIGN KEY (codigo_funcionario) REFERENCES funcionario(id)
);

-- Criar a função para calcular as horas trabalhadas descontando o intervalo de almoço
CREATE OR REPLACE FUNCTION calcular_horas_trabalhadas(reg_ponto registro_ponto)
RETURNS interval AS $$
DECLARE
    horas_trabalhadas interval;
BEGIN
    IF reg_ponto.hora_entrada IS NOT NULL AND reg_ponto.hora_saida IS NOT NULL THEN
        -- Calcula as horas trabalhadas descontando o intervalo de almoço
        horas_trabalhadas := reg_ponto.hora_saida - reg_ponto.hora_entrada;
        
        -- Desconta o tempo do intervalo de almoço (se aplicável)
        IF reg_ponto.hora_entrada_almoco IS NOT NULL AND reg_ponto.hora_retorno_almoco IS NOT NULL THEN
            horas_trabalhadas := horas_trabalhadas - (reg_ponto.hora_retorno_almoco - reg_ponto.hora_entrada_almoco);
        END IF;
    ELSE
        horas_trabalhadas := NULL;
    END IF;

    RETURN horas_trabalhadas;
END;
$$ LANGUAGE plpgsql;

-- Criar a função para calcular as horas trabalhadas e atualizar o campo 'horas_trabalhadas'
CREATE OR REPLACE FUNCTION atualizar_horas_trabalhadas()
RETURNS TRIGGER AS $$
BEGIN
    NEW.horas_trabalhadas := calcular_horas_trabalhadas(NEW);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


-- Criar a função para atualizar os campos de hora e calcular as horas trabalhadas
CREATE OR REPLACE FUNCTION atualizar_horario_ponto()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.tipo_registro = 'entrada' THEN
        NEW.hora_entrada := CURRENT_TIMESTAMP AT TIME ZONE 'America/Sao_Paulo';

    ELSIF NEW.tipo_registro = 'saida' THEN
        NEW.hora_saida := CURRENT_TIMESTAMP AT TIME ZONE 'America/Sao_Paulo';

    ELSIF NEW.tipo_registro = 'entrada_almoço' THEN
        NEW.hora_entrada_almoco := CURRENT_TIME AT TIME ZONE 'America/Sao_Paulo';

    ELSIF NEW.tipo_registro = 'retorno_almoço' THEN
        NEW.hora_retorno_almoco := CURRENT_TIME AT TIME ZONE 'America/Sao_Paulo';

    END IF;

    NEW.horas_trabalhadas := calcular_horas_trabalhadas(NEW);

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


-- Criar a trigger para atualizar os campos de hora antes de inserir o registro
CREATE TRIGGER atualizar_horario_trigger
BEFORE INSERT OR UPDATE ON registro_ponto
FOR EACH ROW EXECUTE FUNCTION atualizar_horario_ponto();

-- Criar a trigger para calcular as horas trabalhadas antes de inserir o registro
CREATE TRIGGER atualizar_horas_trabalhadas_trigger
BEFORE INSERT OR UPDATE ON registro_ponto
FOR EACH ROW EXECUTE FUNCTION atualizar_horas_trabalhadas();