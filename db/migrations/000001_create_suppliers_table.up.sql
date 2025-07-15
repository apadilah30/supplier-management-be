
CREATE TABLE suppliers (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    nick_name   VARCHAR(100),
    status      VARCHAR(50) NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE supplier_addresses (
    id          SERIAL PRIMARY KEY,
    supplier_id INTEGER NOT NULL REFERENCES suppliers(id) ON DELETE CASCADE,
    name        VARCHAR(100),
    address     TEXT NOT NULL,
    is_main     BOOLEAN DEFAULT false
);

CREATE TABLE supplier_contacts (
    id            SERIAL PRIMARY KEY,
    supplier_id   INTEGER NOT NULL REFERENCES suppliers(id) ON DELETE CASCADE,
    name          VARCHAR(255) NOT NULL,
    job_position  VARCHAR(100),
    email         VARCHAR(255),
    phone         VARCHAR(50),
    mobile        VARCHAR(50),
    is_main       BOOLEAN DEFAULT false
);

CREATE TABLE supplier_groups (
    id            SERIAL PRIMARY KEY,
    supplier_id   INTEGER NOT NULL REFERENCES suppliers(id) ON DELETE CASCADE,
    group_name    VARCHAR(255) NOT NULL,
    value         VARCHAR(255) NOT NULL,
    is_active     BOOLEAN DEFAULT true
);