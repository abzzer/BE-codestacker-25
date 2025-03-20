-- Assumption here is cases will never be deleted

CREATE SEQUENCE staff_id_seq START 100;
CREATE TABLE users (
    id TEXT PRIMARY KEY UNIQUE NOT NULL DEFAULT 'A' || nextval('staff_id_seq'),
    name VARCHAR(100) NOT NULL,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role TEXT CHECK (role IN ('admin', 'investigator', 'officer', 'auditor')) NOT NULL,
    clearance_level TEXT CHECK (clearance_level IN ('low', 'medium', 'high', 'critical')) NOT NULL
);

INSERT INTO users (id, name, username, password, role, clearance_level) VALUES ('A001', 'Admin User', 'admin', crypt('123', gen_salt('bf')), 'admin', 'critical');


CREATE SEQUENCE case_number_seq START 10000;
CREATE TABLE cases (
    id SERIAL PRIMARY KEY,
    case_number VARCHAR(10) UNIQUE NOT NULL DEFAULT 'C' || nextval('case_number_seq'),
    case_name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL CHECK (length(description) <= 100),
    area VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    created_by TEXT NOT NULL REFERENCES users(id),
    created_by_role TEXT NOT NULL,
    created_by_name TEXT NOT NULL ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    case_type TEXT CHECK (case_type IN ('criminal', 'civil')) NOT NULL,
    level TEXT CHECK (level IN ('low', 'medium', 'high', 'critical')) NOT NULL,
    reported_by UUID NOT NULL REFERENCES citizens(id) ON DELETE CASCADE
);

UPDATE cases
SET created_by = 'DELETED'
WHERE created_by IS NULL;


CREATE TABLE case_assignees (
    case_id INT NOT NULL REFERENCES cases(id),
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (case_id, user_id)
);

CREATE TABLE persons (
    id SERIAL PRIMARY KEY,
    case_id INT NOT NULL REFERENCES cases(id),
    type TEXT CHECK (type IN ('victim', 'suspect', 'witness')) NOT NULL,
    name VARCHAR(100) NOT NULL,
    age INT NOT NULL CHECK (age > 0),
    gender TEXT CHECK (gender IN ('male', 'female', 'other')) NOT NULL,
    role TEXT NOT NULL
);

CREATE TABLE evidence (
    id SERIAL PRIMARY KEY,
    case_id INT NOT NULL REFERENCES cases(id) ON DELETE CASCADE,
    officer_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type TEXT CHECK (type IN ('text', 'image')) NOT NULL,
    content TEXT NOT NULL,  -- if image then url else the text file
    size TEXT,
    remarks TEXT,
    deleted BOOLEAN DEFAULT FALSE
);

-- remember that we also want a table for the TOP 10 WORDS too!.
-- Add another table if we every query the API

CREATE TABLE audit_logs (
    id SERIAL PRIMARY KEY,
    action TEXT NOT NULL CHECK (action IN ('added', 'updated', 'soft_deleted', 'hard_deleted')),
    evidence_id INT NOT NULL REFERENCES evidence(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);



CREATE SEQUENCE report_number_seq START 20000;
CREATE TABLE reports (
    id SERIAL PRIMARY KEY,
    report_id VARCHAR(10) UNIQUE NOT NULL DEFAULT 'R' || nextval('report_number_seq'),
    citizen_id UUID NOT NULL REFERENCES citizens(id) ON DELETE CASCADE,
    case_id INT REFERENCES cases(id) ON DELETE SET NULL,
    status TEXT CHECK (status IN ('pending', 'ongoing', 'closed')) NOT NULL DEFAULT 'pending',
    reported_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);



-- CREATE FUNCTION soft_delete_evidence(evidence_id INT, user_id TEXT) RETURNS VOID AS $$
-- BEGIN
--     UPDATE evidence SET deleted = TRUE WHERE id = evidence_id;
--     INSERT INTO audit_logs (action, evidence_id, user_id) VALUES ('soft_deleted', evidence_id, user_id);
-- END;
-- $$ LANGUAGE plpgsql;
