--- Truncate tables
TRUNCATE TABLE role_permissions, users, roles, permissions, instructor_profiles, courses, categories RESTART IDENTITY CASCADE;

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

('category_read', 'Read categories'),
('category_create', 'Create categories'),
('category_update', 'Update categories'),
('category_delete', 'Delete categories'),

('instructor_profile_read', 'Read instructor profiles'),
('instructor_profile_create', 'Create instructor profiles'),
('instructor_profile_update', 'Update instructor profiles'),
('instructor_profile_delete', 'Delete instructor profiles'),

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
    
    'instructor_profile_read',
    'instructor_profile_create',
    'instructor_profile_update',
    'instructor_profile_delete',

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
    'instructor_profile_read',
    'course_read'
);

------------------- Create sample users --------------------
INSERT INTO users (email, username, password, name, avatar, is_active, role_id) VALUES

('admin1@gmail.com','admin1','$2a$10$B.4r3Ldtso/IQjexBYypgeZAmbZvIYm2vkH.t6C/7yPLhIR2pQhCS','admin1','https://example.com/avatars/admin1.jpg',true,(SELECT id FROM roles WHERE name='admin')),
('admin2@gmail.com','admin2','$2a$10$B.4r3Ldtso/IQjexBYypgeZAmbZvIYm2vkH.t6C/7yPLhIR2pQhCS','admin2','https://example.com/avatars/admin2.jpg',true,(SELECT id FROM roles WHERE name='admin')),
('admin3@gmail.com','admin3','$2a$10$B.4r3Ldtso/IQjexBYypgeZAmbZvIYm2vkH.t6C/7yPLhIR2pQhCS','admin3','https://example.com/avatars/admin3.jpg',true,(SELECT id FROM roles WHERE name='admin')),
('admin4@gmail.com','admin4','$2a$10$B.4r3Ldtso/IQjexBYypgeZAmbZvIYm2vkH.t6C/7yPLhIR2pQhCS','admin4','https://example.com/avatars/admin4.jpg',true,(SELECT id FROM roles WHERE name='admin')),
('admin5@gmail.com','admin5','$2a$10$B.4r3Ldtso/IQjexBYypgeZAmbZvIYm2vkH.t6C/7yPLhIR2pQhCS','admin5','https://example.com/avatars/admin5.jpg',true,(SELECT id FROM roles WHERE name='admin')),

('instructor1@gmail.com','instructor1','$2a$10$B.4r3Ldtso/IQjexBYypgeZAmbZvIYm2vkH.t6C/7yPLhIR2pQhCS','instructor1','https://example.com/avatars/instructor1.jpg',true,(SELECT id FROM roles WHERE name='instructor')),
('instructor2@gmail.com','instructor2','$2a$10$B.4r3Ldtso/IQjexBYypgeZAmbZvIYm2vkH.t6C/7yPLhIR2pQhCS','instructor2','https://example.com/avatars/instructor2.jpg',true,(SELECT id FROM roles WHERE name='instructor')),
('instructor3@gmail.com','instructor3','$2a$10$B.4r3Ldtso/IQjexBYypgeZAmbZvIYm2vkH.t6C/7yPLhIR2pQhCS','instructor3','https://example.com/avatars/instructor3.jpg',true,(SELECT id FROM roles WHERE name='instructor')),
('instructor4@gmail.com','instructor4','$2a$10$B.4r3Ldtso/IQjexBYypgeZAmbZvIYm2vkH.t6C/7yPLhIR2pQhCS','instructor4','https://example.com/avatars/instructor4.jpg',true,(SELECT id FROM roles WHERE name='instructor')),
('instructor5@gmail.com','instructor5','$2a$10$B.4r3Ldtso/IQjexBYypgeZAmbZvIYm2vkH.t6C/7yPLhIR2pQhCS','instructor5','https://example.com/avatars/instructor5.jpg',true,(SELECT id FROM roles WHERE name='instructor')),

('student1@gmail.com','student1','$2a$10$B.4r3Ldtso/IQjexBYypgeZAmbZvIYm2vkH.t6C/7yPLhIR2pQhCS','student1','https://example.com/avatars/student1.jpg',true,(SELECT id FROM roles WHERE name='student')),
('student2@gmail.com','student2','$2a$10$B.4r3Ldtso/IQjexBYypgeZAmbZvIYm2vkH.t6C/7yPLhIR2pQhCS','student2','https://example.com/avatars/student2.jpg',true,(SELECT id FROM roles WHERE name='student')),
('student3@gmail.com','student3','$2a$10$B.4r3Ldtso/IQjexBYypgeZAmbZvIYm2vkH.t6C/7yPLhIR2pQhCS','student3','https://example.com/avatars/student3.jpg',true,(SELECT id FROM roles WHERE name='student')),
('student4@gmail.com','student4','$2a$10$B.4r3Ldtso/IQjexBYypgeZAmbZvIYm2vkH.t6C/7yPLhIR2pQhCS','student4','https://example.com/avatars/student4.jpg',true,(SELECT id FROM roles WHERE name='student')),
('student5@gmail.com','student5','$2a$10$B.4r3Ldtso/IQjexBYypgeZAmbZvIYm2vkH.t6C/7yPLhIR2pQhCS','student5','https://example.com/avatars/student5.jpg',true,(SELECT id FROM roles WHERE name='student'));


------------------- Create sample instructor_profiles --------------------
INSERT INTO instructor_profiles (user_id, bio, education, rating_avg, total_student, total_course, balance, linkedin_url, youtube_url, instagram_url) VALUES
((SELECT id FROM users WHERE username='instructor1'),'Backend engineer with 10+ years experience in Go','Master of Computer Science - Stanford',4.8,1250,8,5200,'https://linkedin.com/in/instructor1','https://youtube.com/@instructor1','https://instagram.com/instructor1'),
((SELECT id FROM users WHERE username='instructor2'),'Fullstack developer specialized in React and Node','Bachelor of Software Engineering - MIT',4.6,980,5,3100,'https://linkedin.com/in/instructor2','https://youtube.com/@instructor2','https://instagram.com/instructor2'),
((SELECT id FROM users WHERE username='instructor3'),'DevOps engineer teaching Docker and Kubernetes','Bachelor of Information Technology',4.9,2000,10,8800,'https://linkedin.com/in/instructor3','https://youtube.com/@instructor3','https://instagram.com/instructor3'),
((SELECT id FROM users WHERE username='instructor4'),'Cloud architect and microservices specialist','Master of Software Engineering',4.7,1500,7,6400,'https://linkedin.com/in/instructor4','https://youtube.com/@instructor4','https://instagram.com/instructor4'),
((SELECT id FROM users WHERE username='instructor5'),'Senior data engineer teaching data pipelines','Bachelor of Computer Science',4.8,1700,9,7200,'https://linkedin.com/in/instructor5','https://youtube.com/@instructor5','https://instagram.com/instructor5');

------------------- Create sample categories --------------------
INSERT INTO categories (name, description) VALUES
('Backend Development','Courses about backend programming such as Go, Java, Node.js, and system design'),
('Frontend Development','Courses about frontend technologies like HTML, CSS, JavaScript, React, Vue'),
('DevOps','Courses covering Docker, Kubernetes, CI/CD, and cloud infrastructure'),
('Mobile Development','Courses about building mobile apps using Flutter, React Native, Swift, Kotlin'),
('Data Science','Courses related to data analysis, machine learning, and AI'),
('Database','Courses about SQL, PostgreSQL, MySQL, and database design'),
('Cloud Computing','Courses about AWS, GCP, Azure, and distributed systems'),
('Cyber Security','Courses about security fundamentals, penetration testing, and secure coding');

------------------- Create sample courses --------------------
INSERT INTO courses (title, description, image, slug, instructor_id, duration, category_id, price, old_price, average_rate, status, total_student) VALUES
('Go Backend Development Masterclass','Learn scalable backend with Golang','https://example.com/go.jpg','go-backend-development',(SELECT id FROM instructor_profiles WHERE user_id=(SELECT id FROM users WHERE username='instructor1')),600,(SELECT id FROM categories WHERE name='Backend Development'),49.99,79.99,4.8,'published',850),
('React From Zero to Hero','Complete React course with hooks and advanced patterns','https://example.com/react.jpg','react-zero-to-hero',(SELECT id FROM instructor_profiles WHERE user_id=(SELECT id FROM users WHERE username='instructor2')),720,(SELECT id FROM categories WHERE name='Frontend Development'),39.99,69.99,4.7,'published',640),
('Docker & Kubernetes Practical Guide','Learn containerization and orchestration','https://example.com/docker.jpg','docker-kubernetes-guide',(SELECT id FROM instructor_profiles WHERE user_id=(SELECT id FROM users WHERE username='instructor3')),540,(SELECT id FROM categories WHERE name='DevOps'),59.99,89.99,4.9,'published',1200),
('Microservices with Go','Build microservices using Go and gRPC','https://example.com/microservices.jpg','microservices-go',(SELECT id FROM instructor_profiles WHERE user_id=(SELECT id FROM users WHERE username='instructor1')),650,(SELECT id FROM categories WHERE name='Backend Development'),54.99,84.99,4.8,'published',430);