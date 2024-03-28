CREATE TABLE IF NOT EXISTS registro_ponto (
    id serial PRIMARY KEY,
    codigo_funcionario integer NOT NULL,
    criado_em timestamp DEFAULT now() NOT NULL,
    tipo_registro varchar(15) NOT NULL,  
    hora_entrada time,
    hora_saida time,
    hora_entrada_almoco time,  
    hora_retorno_almoco time,
    horas_trabalhadas interval,
    horas_extras interval,
    CONSTRAINT fk_codigo_funcionario FOREIGN KEY (codigo_funcionario) REFERENCES funcionario(id)
);

-- Criar o tipo de registro para representar as horas trabalhadas e as horas extras
CREATE TYPE horas_registro AS (
    horas_trabalhadas interval,
    horas_extras interval
);

-- Criar a função para calcular as horas trabalhadas descontando o intervalo de almoço
CREATE OR REPLACE FUNCTION calcular_horas_trabalhadas(reg_ponto registro_ponto)
RETURNS horas_registro AS $$
DECLARE
    horas_trabalhadas interval;
    horas_extras interval;
    horas_normais interval;
    horas_resultado horas_registro;
BEGIN
    IF reg_ponto.hora_entrada IS NOT NULL AND reg_ponto.hora_saida IS NOT NULL THEN
        horas_trabalhadas := reg_ponto.hora_saida - reg_ponto.hora_entrada;
        
        IF reg_ponto.hora_entrada_almoco IS NOT NULL AND reg_ponto.hora_retorno_almoco IS NOT NULL THEN
            horas_trabalhadas := horas_trabalhadas - (reg_ponto.hora_retorno_almoco - reg_ponto.hora_entrada_almoco);
        END IF;

        -- Definir horas normais de trabalho (por exemplo, 8 horas e meia)
        horas_normais := '08:30:00';

        IF horas_trabalhadas > horas_normais THEN
            horas_extras := horas_trabalhadas - horas_normais;
        ELSE
            horas_extras := '00:00:00';
        END IF;
    ELSE
        horas_trabalhadas := NULL;
        horas_extras := NULL;
    END IF;
    horas_resultado.horas_trabalhadas := horas_trabalhadas;
    horas_resultado.horas_extras := horas_extras;

    RETURN horas_resultado;
END;
$$ LANGUAGE plpgsql;

-- Criar a função para calcular as horas trabalhadas e atualizar o campo 'horas_trabalhadas'
CREATE OR REPLACE FUNCTION atualizar_horas_trabalhadas_extras()
RETURNS TRIGGER AS $$
DECLARE
    resultado horas_registro;
BEGIN
   resultado := calcular_horas_trabalhadas(NEW);

   -- Atribui as horas trabalhadas e as horas extras aos campos correspondentes
    NEW.horas_trabalhadas := resultado.horas_trabalhadas;
    NEW.horas_extras := resultado.horas_extras;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


-- Criar a função para atualizar os campos de hora e calcular as horas trabalhadas
CREATE OR REPLACE FUNCTION atualizar_horario_ponto()
RETURNS TRIGGER AS $$
DECLARE
    resultado horas_registro;
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
    resultado := calcular_horas_trabalhadas(NEW);

    -- Atribui as horas trabalhadas e as horas extras aos campos correspondentes
    NEW.horas_trabalhadas := resultado.horas_trabalhadas;
    NEW.horas_extras := resultado.horas_extras;

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
FOR EACH ROW EXECUTE FUNCTION atualizar_horas_trabalhadas_extras();