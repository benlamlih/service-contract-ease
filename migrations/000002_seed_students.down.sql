-- Remove inserted test students
DELETE
FROM students
WHERE name IN ('Alice Anderson', 'Bob Brown', 'Charlie Chaplin');

DROP
EXTENSION IF EXISTS "uuid-ossp";
