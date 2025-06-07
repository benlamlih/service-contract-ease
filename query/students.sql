-- name: GetStudentByID :one
SELECT *
FROM students
WHERE id = $1;

-- name: GetAllStudents :many
SELECT *
FROM students;

