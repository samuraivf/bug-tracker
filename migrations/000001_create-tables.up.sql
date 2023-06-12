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
    admin INT REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE NOT NULL
);

CREATE TABLE projects_members (
    project_id INT REFERENCES projects(id) ON UPDATE CASCADE ON DELETE CASCADE NOT NULL,
    member_id INT REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE NOT NULL,
    CONSTRAINT id PRIMARY KEY (project_id, member_id)
);

CREATE TYPE priority AS ENUM ('low', 'medium', 'high');
CREATE TYPE progress AS ENUM ('TO DO', 'IN PROGRESS', 'IN REVIEW', 'DONE');

CREATE TABLE tasks (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    task_priority priority,
    assignee INT REFERENCES users(id) ON DELETE CASCADE,
    project_id INT REFERENCES projects(id) ON UPDATE CASCADE ON DELETE CASCADE NOT NULL,
    task_type progress,
    created_at TIMESTAMP NOT NULL,
    perform_to TIMESTAMP
);