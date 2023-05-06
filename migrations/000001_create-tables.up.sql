CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    username VARCHAR(128) NOT NULL UNIQUE,
    password VARCHAR(128) NOT NULL,
    email TEXT NOT NULL UNIQUE
);

CREATE TABLE projects (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    admin INT REFERENCES users(id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE projects_members (
    project_id INT REFERENCES projects(id) ON UPDATE CASCADE ON DELETE CASCADE NOT NULL,
    member_id INT REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE NOT NULL,
    CONSTRAINT id PRIMARY KEY (project_id, member_id)
);