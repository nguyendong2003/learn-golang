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
('user.read', 'Read users'),
('user.create', 'Create users'),
('user.update', 'Update users'),
('user.delete', 'Delete users'),

('blog.read', 'Read blogs'),
('blog.create', 'Create blogs'),
('blog.update', 'Update blogs'),
('blog.delete', 'Delete blogs'),

('course.read', 'Read courses'),
('course.create', 'Create courses'),
('course.update', 'Update courses'),
('course.delete', 'Delete courses');

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
    'blog.read',
    'blog.create',
    'blog.update',
    'blog.delete',
    'course.read',
    'course.create',
    'course.update'
);

-- Student Role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'student'
AND p.code IN (
    'blog.read',
    'course.read'
);