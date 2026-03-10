--- Truncate tables
TRUNCATE TABLE role_permissions, users, roles, permissions RESTART IDENTITY CASCADE;


------- Init data
INSERT INTO roles (name, description)
VALUES
('admin', 'System administrator'),
('instructor', 'Course instructor'),
('student', 'Platform student');

INSERT INTO permissions (code, description)
VALUES
('user_read', 'Read users'),
('user_create', 'Create users'),
('user_update', 'Update users'),
('user_delete', 'Delete users'),

('blog_read', 'Read blogs'),
('blog_create', 'Create blogs'),
('blog_update', 'Update blogs'),
('blog_delete', 'Delete blogs'),

('course_read', 'Read courses'),
('course_create', 'Create courses'),
('course_update', 'Update courses'),
('course_delete', 'Delete courses');

-- Admin Role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'admin';

-- Instructor Role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'instructor'
AND p.code IN (
    'blog_read',
    'blog_create',
    'blog_update',
    'blog_delete',
    'course_read',
    'course_create',
    'course_update',
    'course_delete'
);

-- Student Role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'student'
AND p.code IN (
    'blog_read',
    'course_read'
);