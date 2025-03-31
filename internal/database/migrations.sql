-- Assumption here is cases will never be deleted

--- ENUMS ------------

CREATE TYPE user_role_enum AS ENUM ('admin', 'investigator', 'officer', 'auditor');
CREATE TYPE case_levels_enum AS ENUM ('low', 'medium', 'high', 'critical');
CREATE TYPE person_type_enum AS ENUM ('victim', 'suspect', 'witness');
CREATE TYPE gender_enum AS ENUM ('male', 'female');
CREATE TYPE evidence_type_enum AS ENUM ('text', 'image');
CREATE TYPE audit_action_enum AS ENUM ('added', 'updated', 'soft_deleted', 'hard_deleted');
CREATE TYPE case_status_enum AS ENUM ('pending', 'ongoing', 'closed');

---- Sequences -----------------
CREATE SEQUENCE staff_id_seq START 100;
CREATE SEQUENCE case_number_seq START 10000;


----- Functions -------------------------------
CREATE OR REPLACE FUNCTION generate_unique_staff_id() 
RETURNS TEXT AS $$
DECLARE
    new_id TEXT;
    seq_num BIGINT;
BEGIN
    LOOP
        seq_num := nextval('staff_id_seq');
        new_id := 'A' || seq_num::TEXT;
        IF NOT EXISTS (SELECT 1 FROM users WHERE id = new_id) THEN
            RETURN new_id;
        END IF;
    END LOOP;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION generate_unique_case_number()
RETURNS TEXT AS $$
DECLARE
    new_case_number TEXT;
BEGIN
    LOOP
        new_case_number := 'C' || nextval('case_number_seq')::TEXT;
        IF NOT EXISTS (SELECT 1 FROM cases WHERE case_number = new_case_number) THEN
            RETURN new_case_number;
        END IF;
    END LOOP;
END;
$$ LANGUAGE plpgsql;


-- TABLES ---------------------------

CREATE TABLE users (
    id TEXT PRIMARY KEY DEFAULT generate_unique_staff_id(),
    name VARCHAR(100) NOT NULL,
    password TEXT NOT NULL,
    role user_role_enum NOT NULL,
    clearance_level case_levels_enum NOT NULL,
    deleted BOOLEAN DEFAULT FALSE
);


CREATE TABLE cases (
    case_number TEXT PRIMARY KEY DEFAULT generate_unique_case_number(),
    case_name VARCHAR(100) NOT NULL,
    description VARCHAR(100) NOT NULL, -- 100 characters; truncate in Go if needed
    area VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    created_by TEXT REFERENCES users(id) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    case_type TEXT NOT NULL DEFAULT 'criminal',
    status case_status_enum NOT NULL DEFAULT 'pending',
    level case_levels_enum NOT NULL
);

CREATE TABLE case_assignees (
    case_number TEXT NOT NULL REFERENCES cases(case_number),
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (case_number, user_id)
);

CREATE TABLE persons (
    id SERIAL PRIMARY KEY,
    case_number TEXT NOT NULL REFERENCES cases(case_number),
    type person_type_enum NOT NULL,
    name VARCHAR(100) NOT NULL,
    age INT NOT NULL CHECK (age > 0),
    gender gender_enum NOT NULL,
    role VARCHAR(100) NOT NULL
);

CREATE TABLE evidence (
    id SERIAL PRIMARY KEY,
    case_number TEXT NOT NULL REFERENCES cases(case_number),
    officer_id TEXT NOT NULL REFERENCES users(id),
    type evidence_type_enum NOT NULL,
    content TEXT NOT NULL,  -- If image then URL, else text content
    size TEXT,
    remarks TEXT,
    deleted BOOLEAN DEFAULT FALSE
);

CREATE TABLE audit_logs (
    id SERIAL PRIMARY KEY,
    action audit_action_enum NOT NULL,
    evidence_id INT,
    user_id TEXT NOT NULL REFERENCES users(id),
    timestamp TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE reports (
    report_id SERIAL PRIMARY KEY,
    email TEXT NOT NULL,
    civil_id TEXT NOT NULL,
    name VARCHAR(100) NOT NULL,
    role VARCHAR(100) NOT NULL DEFAULT 'Citizen',
    case_number TEXT REFERENCES cases(case_number),
    description TEXT NOT NULL,
    area TEXT NOT NULL,
    city TEXT NOT NULL
);

--- Populate the Tables ------------------------

INSERT INTO users (id, name, password, role, clearance_level) VALUES ('A001', 'Admin User', '$2a$10$J14ZMfR26KYUczdLamBmpOLaEMp8nibjGITnY.AZpkmz6k6h7hh/q', 'admin', 'critical'),
('A101', 'Detective Jane Smith', '$2a$10$J14ZMfR26KYUczdLamBmpOLaEMp8nibjGITnY.AZpkmz6k6h7hh/q', 'investigator', 'high'),
('A102', 'Officer Mike Johnson', '$2a$10$J14ZMfR26KYUczdLamBmpOLaEMp8nibjGITnY.AZpkmz6k6h7hh/q', 'officer', 'medium'),
('A104', 'Auditor Alice Green', '$2a$10$J14ZMfR26KYUczdLamBmpOLaEMp8nibjGITnY.AZpkmz6k6h7hh/q', 'auditor', 'medium');


INSERT INTO cases (case_number, case_name, description, area, city, created_by, created_at, case_type, status, level) VALUES 
('C12345', 'Theft Investigation', 'Investigation of a reported theft at a local store.', 'Downtown', 'New York', 'A001', '2025-03-10T14:30:00Z', 'criminal', 'ongoing', 'high');

INSERT INTO case_assignees (case_number, user_id) VALUES ('C12345', 'A101'), ('C12345', 'A102'), ('C12345', 'A104');

-- Insert persons associated with the case
INSERT INTO persons (case_number, type, name, age, gender, role) VALUES ('C12345', 'suspect', 'Michael Brown', 32, 'male', 'Primary Suspect'), 
('C12345', 'victim', 'Sarah Parker', 28, 'female', 'Store Owner');

INSERT INTO reports (email, civil_id, name, role, case_number, description, area, city)
VALUES ('bob.wilson@gmail.com', 'A12356879', 'Citizen Bob Wilson', 'Citizen', 'C12345', 'Saw there was a theft at XYZ store', 'Downtown', 'New York'); -- Should be 12 however chose to serialise it instead
