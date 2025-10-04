


CREATE TABLE IF NOT EXISTS VIDEO(
    id VARCHAR(30) PRIMARY KEY,
    titulo VARCHAR(60), 
    descricao TEXT,
    canal VARCHAR(100),
    data_publicacao VARCHAR(30),
    quant_view INTEGER,
    quant_like INTEGER
);

CREATE TABLE IF NOT EXISTS COMENTARIO(
    id VARCHAR(30) PRIMARY KEY,
    video_id VARCHAR(30),
    autor VARCHAR(30),
    texto TEXT,
    data_publicacao VARCHAR(30),
    reply VARCHAR(30),

    FOREIGN KEY (reply) REFERENCES COMENTARIO(id),
    FOREIGN KEY (video_id) REFERENCES METADADO(id)
);

