-- Enable UUID support in PostgreSQL
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Subject enum
CREATE TYPE subject AS ENUM (
    'math',
    'science',
    'history',
    'language',
    'art'
    );

-- Teachers
CREATE TABLE teachers
(
    id    UUID PRIMARY KEY,
    name  TEXT        NOT NULL,
    email TEXT UNIQUE NOT NULL
);

-- Classrooms
CREATE TABLE classrooms
(
    id         UUID PRIMARY KEY,
    name       TEXT NOT NULL,
    teacher_id UUID REFERENCES teachers (id) ON DELETE CASCADE
);
CREATE INDEX idx_classrooms_teacher ON classrooms (teacher_id);

-- Students
CREATE TABLE students
(
    id           UUID PRIMARY KEY,
    name         TEXT NOT NULL,
    classroom_id UUID REFERENCES classrooms (id) ON DELETE CASCADE
);
CREATE INDEX idx_students_classroom ON students (classroom_id);

-- Assignments
CREATE TABLE assignments
(
    id           UUID PRIMARY KEY,
    title        TEXT    NOT NULL,
    subject      subject NOT NULL,
    classroom_id UUID REFERENCES classrooms (id) ON DELETE CASCADE,
    due_date     DATE
);
CREATE INDEX idx_assignments_classroom ON assignments (classroom_id);

-- Grades
CREATE TABLE grades
(
    id            UUID PRIMARY KEY,
    student_id    UUID REFERENCES students (id) ON DELETE CASCADE,
    assignment_id UUID REFERENCES assignments (id) ON DELETE CASCADE,
    score         INTEGER NOT NULL,
    feedback      TEXT,
    created_at    TIMESTAMP DEFAULT now()
);
CREATE INDEX idx_grades_student_created ON grades (student_id, created_at DESC);
