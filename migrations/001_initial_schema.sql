CREATE TABLE sessions (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);


CREATE TABLE votes (
    id UUID PRIMARY KEY,
    session_id UUID NOT NULL REFERENCES sessions(id),
    product_id UUID NOT NULL,
    score INT NOT NULL CHECK (score >= 1 AND score <= 5),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(session_id, product_id)
);

CREATE INDEX votes_session_id_idx ON votes(session_id);
CREATE INDEX votes_product_id_idx ON votes(product_id); 