CREATE TABLE remedies (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name varchar(50) NOT NULL,
    description varchar(100),
    notes varchar(100)
);

CREATE TABLE causes (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    code varchar(50) NOT NULL UNIQUE,
    name varchar(50) NOT NULL,
    type varchar(50) NOT NULL,
    prior float NOT NULL,
    description varchar(100),
    remedy_id bigint,

    CONSTRAINT fk_remedy FOREIGN KEY (remedy_id) REFERENCES remedies(id) ON DELETE SET NULL,
    CONSTRAINT chk_type CHECK (type IN ('poison','disease'))
);

CREATE TABLE cases (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title varchar(50) NOT NULL,
    description varchar(100),
    status varchar(50) NOT NULL,
    
    CONSTRAINT chk_status CHECK (status IN ('open','closed'))
);

CREATE TABLE factors (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    code varchar(50) NOT NULL UNIQUE,
    name varchar(50) NOT NULL,
    description varchar(100)
);

CREATE TABLE factor_profiles (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    cause_id bigint NOT NULL,
    factor_id bigint NOT NULL,
    modifier numeric(10,2) NOT NULL,

    CONSTRAINT uq_factor_cause UNIQUE (factor_id, cause_id),
    CONSTRAINT fk_cause FOREIGN KEY (cause_id) REFERENCES causes(id) ON DELETE CASCADE,
    CONSTRAINT fk_factor FOREIGN KEY (factor_id) REFERENCES factors(id) ON DELETE CASCADE,
    CONSTRAINT chk_modifier CHECK (modifier > 0 AND modifier <= 1)

);


CREATE TABLE categories (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name varchar(50) NOT NULL
);
CREATE TABLE symptoms (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    code varchar(50) NOT NULL UNIQUE,
    name varchar(50) NOT NULL, 
    description varchar(50),
    category_id bigint NOT NULL,

    CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES categories(id)
);

CREATE TABLE symptom_profiles (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    cause_id bigint NOT NULL,
    symptom_id bigint NOT NULL,
    likelihood numeric(10,2) NOT NULL,

    CONSTRAINT uq_symptom_cause UNIQUE (symptom_id, cause_id),
    CONSTRAINT fk_cause FOREIGN KEY (cause_id) REFERENCES causes(id) ON DELETE CASCADE,
    CONSTRAINT fk_symptom FOREIGN KEY (symptom_id) REFERENCES symptoms(id) ON DELETE CASCADE,
    CONSTRAINT chk_likelihood CHECK (likelihood > 0 AND likelihood <= 1)
);

CREATE TABLE observations (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    case_id bigint NOT NULL,
    kind varchar(50),
    ref_id bigint NOT NULL,
    observed_at timestamptz NOT NULL DEFAULT NOW(),
    note varchar(100),

    CONSTRAINT fk_case FOREIGN KEY (case_id) REFERENCES cases(id) ON DELETE CASCADE,
    CONSTRAINT chk_kind CHECK (kind IN ('symptom','factor'))
);