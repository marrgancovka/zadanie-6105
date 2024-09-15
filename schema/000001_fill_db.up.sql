CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- CREATE TABLE IF NOT EXISTS employee (
--                                         id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
--     username VARCHAR(50) UNIQUE NOT NULL,
--     first_name VARCHAR(50),
--     last_name VARCHAR(50),
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
--     );
--
--
-- DO $$
-- BEGIN
--     IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'organization_type') THEN
-- CREATE TYPE organization_type AS ENUM (
--     'IE',
--     'LLC',
--     'JSC'
-- );
-- END IF;
-- END
-- $$;
--
--
-- CREATE TABLE IF NOT EXISTS organization (
--                                             id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
--     name VARCHAR(100) NOT NULL,
--     description TEXT,
--     type organization_type,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
--     );
--
-- CREATE TABLE IF NOT EXISTS organization_responsible (
--                                                         id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
--     organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
--     user_id UUID REFERENCES employee(id) ON DELETE CASCADE
--     );

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status_type') THEN
CREATE TYPE status_type AS ENUM (
    'Created',
    'Published',
    'Closed'
);
END IF;
END
$$;


DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'type_service') THEN
CREATE TYPE type_service AS ENUM (
    'Construction',
    'Delivery',
    'Manufacture'
);
END IF;
END
$$;


DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'type_author') THEN
CREATE TYPE type_author AS ENUM (
    'User',
    'Organization'
);
END IF;
END
$$;


CREATE TABLE IF NOT EXISTS tender (
                                      id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    status status_type NOT NULL,
    service_type type_service NOT NULL,
    version INTEGER NOT NULL DEFAULT 1 CHECK (version > 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
    creator_username VARCHAR(50) REFERENCES employee(username) ON DELETE CASCADE
    );

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'type_decision') THEN
CREATE TYPE type_decision AS ENUM (
    'Approved',
    'Rejected'
);
END IF;
END
$$;

CREATE TABLE IF NOT EXISTS bid (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    author_type type_author NOT NULL,
    status status_type NOT NULL,
    version INTEGER NOT NULL DEFAULT 1 CHECK (version > 0),
    decision type_decision,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    tender_id UUID REFERENCES tender(id) ON DELETE CASCADE,
    author_id UUID REFERENCES employee(id) ON DELETE CASCADE
    );
