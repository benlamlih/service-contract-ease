-- Enable the uuid-ossp extension
CREATE
EXTENSION IF NOT EXISTS "uuid-ossp";

-- Insert sample students (assuming classrooms table already has data)
-- You can adjust classroom_id manually to match actual existing classroom UUIDs

INSERT INTO students (id, name, classroom_id)
VALUES (uuid_generate_v4(), 'Alice Anderson', NULL),
       (uuid_generate_v4(), 'Bob Brown', NULL),
       (uuid_generate_v4(), 'Charlie Chaplin', NULL);
