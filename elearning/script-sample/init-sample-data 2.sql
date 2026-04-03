--- Truncate tables
TRUNCATE TABLE 
    role_permissions,
    follows,
    instructor_profiles,
    lessons,
    chapters,
    enrollments,
    carts,
    course_purchase_details,
    course_purchases,
    coupons,
    courses,
    categories,
    payments,
    subscriptions,
    plans,
    users,
    roles,
    permissions,
    stripe_events
RESTART IDENTITY CASCADE;
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
('course_delete', 'Delete courses'),

('chapter_read', 'Read chapters'),
('chapter_create', 'Create chapters'),
('chapter_update', 'Update chapters'),
('chapter_delete', 'Delete chapters'),

('lesson_read', 'Read lessons'),
('lesson_create', 'Create lessons'),
('lesson_update', 'Update lessons'),
('lesson_delete', 'Delete lessons'),

('follow_read', 'Read follow relationships'),
('follow_create', 'Create follow relationships'),
('follow_delete', 'Delete follow relationships'),

('coupon_read', 'Read course coupons'),
('coupon_create', 'Create course coupons'),
('coupon_update', 'Update course coupons'),
('coupon_delete', 'Delete course coupons'),

('plan_read', 'Read plans'),
('plan_create', 'Create plans'),
('plan_update', 'Update plans'),
('plan_delete', 'Delete plans');

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
    'course_delete',

    'chapter_read',
    'chapter_create',
    'chapter_update',
    'chapter_delete',

    'lesson_read',
    'lesson_create',
    'lesson_update',
    'lesson_delete',

    'follow_read',
    'follow_create',
    'follow_delete',
    
    'coupon_read',
    'coupon_create',
    'coupon_update',
    'coupon_delete'
);

-- Student Role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'student'
AND p.code IN (
    'blog_read',
    'instructor_profile_read',
    'course_read',
    'chapter_read',
    'lesson_read',
    'follow_read',
    'coupon_read'
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
('instructor6@gmail.com','instructor6','$2a$10$B.4r3Ldtso/IQjexBYypgeZAmbZvIYm2vkH.t6C/7yPLhIR2pQhCS','instructor6','https://example.com/avatars/instructor6.jpg',true,(SELECT id FROM roles WHERE name='instructor')),

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
INSERT INTO courses (title, description, image, slug, user_id, duration, category_id, price, average_rate, status, total_student, stripe_product_id, stripe_price_id, stripe_amount, stripe_currency) VALUES
('Go Backend Development Masterclass','Learn scalable backend with Golang','https://example.com/go.jpg','go-backend-development',(SELECT id FROM users WHERE username='instructor1'),600,(SELECT id FROM categories WHERE name='Backend Development'),49.99,4.8,'published',850,'prod_UGD1EECEnSQzWz','price_1THgbRLAb7u1ek8LKb4ZfUVH',4999,'usd'),
('React From Zero to Hero','Complete React course with hooks and advanced patterns','https://example.com/react.jpg','react-zero-to-hero',(SELECT id FROM users WHERE username='instructor2'),720,(SELECT id FROM categories WHERE name='Frontend Development'),39.99,4.7,'published',640,'prod_UGD14O6ffQn2fk','price_1THgbWLAb7u1ek8LGeZ1hbGp',3999,'usd'),
('Docker & Kubernetes Practical Guide','Learn containerization and orchestration','https://example.com/docker.jpg','docker-kubernetes-guide',(SELECT id FROM users WHERE username='instructor3'),540,(SELECT id FROM categories WHERE name='DevOps'),59.99,4.9,'published',1200,'prod_UGD1QuPJwxgJeq','price_1THgbbLAb7u1ek8L3zj9XBvo',5999,'usd'),
('Microservices with Go','Build microservices using Go and gRPC','https://example.com/microservices.jpg','microservices-go',(SELECT id FROM users WHERE username='instructor1'),650,(SELECT id FROM categories WHERE name='Backend Development'),54.99,4.8,'published',430,'prod_UGD1MELHROji9A','price_1THgbiLAb7u1ek8LGqm4E7D8',5499,'usd'),
('Advanced Python for Data Science','Master machine learning with Python and NumPy','https://example.com/python-ds.jpg','python-data-science',(SELECT id FROM users WHERE username='instructor4'),580,(SELECT id FROM categories WHERE name='Data Science'),44.99,4.6,'published',320,'prod_UGD1UtTWxcTU7f','price_1THgbnLAb7u1ek8LQypnhqEN',4499,'usd'),
('Cloud Architecture on AWS','Design and deploy scalable cloud solutions','https://example.com/aws.jpg','aws-cloud-architecture',(SELECT id FROM users WHERE username='instructor5'),500,(SELECT id FROM categories WHERE name='Cloud Computing'),64.99,4.7,'published',520,'prod_UGD17Ze2fuHdW6','price_1THgbtLAb7u1ek8LOAfLqizd',6499,'usd');

------------------- Create sample chapters --------------------
INSERT INTO chapters (title, course_id, "order") VALUES
('Introduction to Go',(SELECT id FROM courses WHERE slug='go-backend-development'),1),
('Go Advanced Concepts',(SELECT id FROM courses WHERE slug='go-backend-development'),2),
('Web Frameworks',(SELECT id FROM courses WHERE slug='go-backend-development'),3),
('Database Integration',(SELECT id FROM courses WHERE slug='go-backend-development'),4),

('React Fundamentals',(SELECT id FROM courses WHERE slug='react-zero-to-hero'),1),
('React Hooks Deep Dive',(SELECT id FROM courses WHERE slug='react-zero-to-hero'),2),
('State Management',(SELECT id FROM courses WHERE slug='react-zero-to-hero'),3),
('Performance Optimization',(SELECT id FROM courses WHERE slug='react-zero-to-hero'),4),

('Docker Basics',(SELECT id FROM courses WHERE slug='docker-kubernetes-guide'),1),
('Docker Advanced',(SELECT id FROM courses WHERE slug='docker-kubernetes-guide'),2),
('Kubernetes Architecture',(SELECT id FROM courses WHERE slug='docker-kubernetes-guide'),3),
('Kubernetes Management',(SELECT id FROM courses WHERE slug='docker-kubernetes-guide'),4);

------------------- Create sample lessons --------------------
INSERT INTO lessons (title, chapter_id, duration, video_url, document_url, is_able_to_preview, "order", "type") VALUES
('Getting Started with Go',(SELECT id FROM chapters WHERE title='Introduction to Go' LIMIT 1),15,'https://example.com/videos/go-intro.mp4',NULL,true,1,CAST('video' AS lesson_type)),
('Variables and Data Types',(SELECT id FROM chapters WHERE title='Introduction to Go' LIMIT 1),20,'https://example.com/videos/go-types.mp4',NULL,true,2,CAST('video' AS lesson_type)),
('Functions and Errors',(SELECT id FROM chapters WHERE title='Introduction to Go' LIMIT 1),22,'https://example.com/videos/go-functions.mp4',NULL,false,3,CAST('video' AS lesson_type)),

('Goroutines Explained',(SELECT id FROM chapters WHERE title='Go Advanced Concepts' LIMIT 1),25,'https://example.com/videos/go-goroutines.mp4',NULL,false,1,CAST('video' AS lesson_type)),
('Channels and Communication',(SELECT id FROM chapters WHERE title='Go Advanced Concepts' LIMIT 1),24,'https://example.com/videos/go-channels.mp4',NULL,false,2,CAST('video' AS lesson_type)),

('JSX Basics',(SELECT id FROM chapters WHERE title='React Fundamentals' LIMIT 1),18,'https://example.com/videos/react-jsx.mp4',NULL,true,1,CAST('video' AS lesson_type)),
('Components and Props',(SELECT id FROM chapters WHERE title='React Fundamentals' LIMIT 1),21,'https://example.com/videos/react-components.mp4',NULL,false,2,CAST('video' AS lesson_type)),

('Hooks Introduction',(SELECT id FROM chapters WHERE title='React Hooks Deep Dive' LIMIT 1),23,'https://example.com/videos/react-hooks.mp4',NULL,true,1,CAST('video' AS lesson_type)),
('Custom Hooks',(SELECT id FROM chapters WHERE title='React Hooks Deep Dive' LIMIT 1),26,'https://example.com/videos/react-custom-hooks.mp4',NULL,false,2,CAST('video' AS lesson_type)),

('What is Docker',(SELECT id FROM chapters WHERE title='Docker Basics' LIMIT 1),17,'https://example.com/videos/docker-intro.mp4',NULL,true,1,CAST('video' AS lesson_type)),
('Docker Images and Registries',(SELECT id FROM chapters WHERE title='Docker Basics' LIMIT 1),19,'https://example.com/videos/docker-images.mp4',NULL,false,2,CAST('video' AS lesson_type));

------------------- Create sample course_events --------------------
INSERT INTO course_events (course_id, name, start_time, end_time, description) VALUES
((SELECT id FROM courses WHERE slug='go-backend-development'),'Live Q&A Session','2026-04-05 14:00:00','2026-04-05 15:00:00','Ask questions about Go backend development'),
((SELECT id FROM courses WHERE slug='react-zero-to-hero'),'Weekly Office Hours','2026-04-10 16:00:00','2026-04-10 17:00:00','Weekly office hours with instructor2'),
((SELECT id FROM courses WHERE slug='docker-kubernetes-guide'),'Kubernetes Workshop','2026-04-15 13:00:00','2026-04-15 16:00:00','Hands-on Kubernetes deployment workshop');

------------------- Create sample plans --------------------
INSERT INTO plans (name, description, billing_cycle, price, stripe_product_id, stripe_price_id, currency, is_active) VALUES
('Premium Monthly','Premium Plan - Subscription hàng tháng',CAST('monthly' AS billing_cycle_enum),9.99,'prod_premium_monthly','price_1TCYIYLAb7u1ek8LlHZQPQUS','usd', true),
('Premium Yearly','Premium Plan - Subscription hàng năm',CAST('yearly' AS billing_cycle_enum),99.99,'prod_premium_yearly','price_1TCYM9LAb7u1ek8LAjkYRCh0','usd', true),
('Pro Monthly','Pro Plan - Full access hàng tháng',CAST('monthly' AS billing_cycle_enum),19.99,'prod_pro_monthly','price_1TCYIYLAb7u1ek8LProMon','usd', true),
('Pro Yearly','Pro Plan - Full access hàng năm',CAST('yearly' AS billing_cycle_enum),199.99,'prod_pro_yearly','price_1TCYIYLAb7u1ek8LProYr','usd', true);

------------------- Create sample coupons --------------------
INSERT INTO coupons (code, discount_type, discount_value, max_redemptions, current_redemptions, is_active, expires_at, stripe_coupon_id, stripe_promotion_code_id) VALUES
('SAVE10','percent',10,100,24,true,'2026-12-31 23:59:59','cou_SAVE10','promo_SAVE10'),
('WELCOME20','percent',20,50,19,true,'2026-06-30 23:59:59','cou_WELCOME20','promo_WELCOME20'),
('FLAT500','amount',500,200,31,true,'2026-08-31 23:59:59','cou_FLAT500','promo_FLAT500');

-- ------------------- Create sample subscriptions --------------------
-- INSERT INTO subscriptions (user_id, plan_id, plan_name, plan_description, plan_price, plan_currency, plan_stripe_price_id, billing_cycle, stripe_subscription_id, status, started_at, current_period_start, current_period_end) VALUES
-- ((SELECT id FROM users WHERE username='student1'),(SELECT id FROM plans WHERE name='Premium Monthly'),(SELECT name FROM plans WHERE name='Premium Monthly'),(SELECT description FROM plans WHERE name='Premium Monthly'),(SELECT price FROM plans WHERE name='Premium Monthly'),(SELECT currency FROM plans WHERE name='Premium Monthly'),(SELECT stripe_price_id FROM plans WHERE name='Premium Monthly'),(SELECT billing_cycle FROM plans WHERE name='Premium Monthly'),'sub_student1_monthly',CAST('active' AS subscription_status_enum),'2026-01-01','2026-01-01','2026-04-01'),
-- ((SELECT id FROM users WHERE username='student2'),(SELECT id FROM plans WHERE name='Pro Yearly'),(SELECT name FROM plans WHERE name='Pro Yearly'),(SELECT description FROM plans WHERE name='Pro Yearly'),(SELECT price FROM plans WHERE name='Pro Yearly'),(SELECT currency FROM plans WHERE name='Pro Yearly'),(SELECT stripe_price_id FROM plans WHERE name='Pro Yearly'),(SELECT billing_cycle FROM plans WHERE name='Pro Yearly'),'sub_student2_yearly',CAST('active' AS subscription_status_enum),'2025-01-01','2025-01-01','2026-01-01'),
-- ((SELECT id FROM users WHERE username='student3'),(SELECT id FROM plans WHERE name='Premium Monthly'),(SELECT name FROM plans WHERE name='Premium Monthly'),(SELECT description FROM plans WHERE name='Premium Monthly'),(SELECT price FROM plans WHERE name='Premium Monthly'),(SELECT currency FROM plans WHERE name='Premium Monthly'),(SELECT stripe_price_id FROM plans WHERE name='Premium Monthly'),(SELECT billing_cycle FROM plans WHERE name='Premium Monthly'),'sub_student3_monthly',CAST('active' AS subscription_status_enum),'2026-02-15','2026-02-15','2026-03-15'),
-- ((SELECT id FROM users WHERE username='student4'),(SELECT id FROM plans WHERE name='Premium Monthly'),(SELECT name FROM plans WHERE name='Premium Monthly'),(SELECT description FROM plans WHERE name='Premium Monthly'),(SELECT price FROM plans WHERE name='Premium Monthly'),(SELECT currency FROM plans WHERE name='Premium Monthly'),(SELECT stripe_price_id FROM plans WHERE name='Premium Monthly'),(SELECT billing_cycle FROM plans WHERE name='Premium Monthly'),'sub_student4_monthly',CAST('active' AS subscription_status_enum),'2026-03-25','2026-03-25','2026-04-25'),
-- ((SELECT id FROM users WHERE username='student5'),(SELECT id FROM plans WHERE name='Pro Monthly'),(SELECT name FROM plans WHERE name='Pro Monthly'),(SELECT description FROM plans WHERE name='Pro Monthly'),(SELECT price FROM plans WHERE name='Pro Monthly'),(SELECT currency FROM plans WHERE name='Pro Monthly'),(SELECT stripe_price_id FROM plans WHERE name='Pro Monthly'),(SELECT billing_cycle FROM plans WHERE name='Pro Monthly'),'sub_student5_monthly',CAST('active' AS subscription_status_enum),'2026-03-01','2026-03-01','2026-04-01'),
-- ((SELECT id FROM users WHERE username='student1'),(SELECT id FROM plans WHERE name='Premium Monthly'),(SELECT name FROM plans WHERE name='Premium Monthly'),(SELECT description FROM plans WHERE name='Premium Monthly'),(SELECT price FROM plans WHERE name='Premium Monthly'),(SELECT currency FROM plans WHERE name='Premium Monthly'),(SELECT stripe_price_id FROM plans WHERE name='Premium Monthly'),(SELECT billing_cycle FROM plans WHERE name='Premium Monthly'),'sub_student1_cycle_202601',CAST('canceled' AS subscription_status_enum),'2026-01-01','2026-01-01','2026-01-31'),
-- ((SELECT id FROM users WHERE username='student1'),(SELECT id FROM plans WHERE name='Premium Monthly'),(SELECT name FROM plans WHERE name='Premium Monthly'),(SELECT description FROM plans WHERE name='Premium Monthly'),(SELECT price FROM plans WHERE name='Premium Monthly'),(SELECT currency FROM plans WHERE name='Premium Monthly'),(SELECT stripe_price_id FROM plans WHERE name='Premium Monthly'),(SELECT billing_cycle FROM plans WHERE name='Premium Monthly'),'sub_student1_cycle_202602',CAST('canceled' AS subscription_status_enum),'2026-02-01','2026-02-01','2026-02-28'),
-- ((SELECT id FROM users WHERE username='student1'),(SELECT id FROM plans WHERE name='Premium Monthly'),(SELECT name FROM plans WHERE name='Premium Monthly'),(SELECT description FROM plans WHERE name='Premium Monthly'),(SELECT price FROM plans WHERE name='Premium Monthly'),(SELECT currency FROM plans WHERE name='Premium Monthly'),(SELECT stripe_price_id FROM plans WHERE name='Premium Monthly'),(SELECT billing_cycle FROM plans WHERE name='Premium Monthly'),'sub_student1_cycle_202603',CAST('active' AS subscription_status_enum),'2026-03-01','2026-03-01','2026-03-31'),
-- ((SELECT id FROM users WHERE username='student2'),(SELECT id FROM plans WHERE name='Premium Monthly'),(SELECT name FROM plans WHERE name='Premium Monthly'),(SELECT description FROM plans WHERE name='Premium Monthly'),(SELECT price FROM plans WHERE name='Premium Monthly'),(SELECT currency FROM plans WHERE name='Premium Monthly'),(SELECT stripe_price_id FROM plans WHERE name='Premium Monthly'),(SELECT billing_cycle FROM plans WHERE name='Premium Monthly'),'sub_student2_cycle_202603',CAST('canceled' AS subscription_status_enum),'2026-03-01','2026-03-01','2026-03-31'),
-- ((SELECT id FROM users WHERE username='student4'),(SELECT id FROM plans WHERE name='Pro Monthly'),(SELECT name FROM plans WHERE name='Pro Monthly'),(SELECT description FROM plans WHERE name='Pro Monthly'),(SELECT price FROM plans WHERE name='Pro Monthly'),(SELECT currency FROM plans WHERE name='Pro Monthly'),(SELECT stripe_price_id FROM plans WHERE name='Pro Monthly'),(SELECT billing_cycle FROM plans WHERE name='Pro Monthly'),'sub_student4_pro_monthly',CAST('active' AS subscription_status_enum),'2026-03-01','2026-03-01','2026-04-01'),
-- ((SELECT id FROM users WHERE username='student3'),(SELECT id FROM plans WHERE name='Pro Monthly'),(SELECT name FROM plans WHERE name='Pro Monthly'),(SELECT description FROM plans WHERE name='Pro Monthly'),(SELECT price FROM plans WHERE name='Pro Monthly'),(SELECT currency FROM plans WHERE name='Pro Monthly'),(SELECT stripe_price_id FROM plans WHERE name='Pro Monthly'),(SELECT billing_cycle FROM plans WHERE name='Pro Monthly'),'sub_student3_pro_cycle_202603',CAST('active' AS subscription_status_enum),'2026-03-16','2026-03-16','2026-04-16');

-- ------------------- Create sample payments --------------------
-- INSERT INTO payments (subscription_id, stripe_invoice_id, stripe_payment_intent, status, amount, currency, stripe_fee, failure_reason, attempt_count, paid_at) VALUES
-- ((SELECT id FROM subscriptions WHERE stripe_subscription_id='sub_student1_monthly'),'in_student1_202601','pi_sub_student1_202601','succeeded',999,'usd',59,'',1,'2026-01-01 10:00:00'),
-- ((SELECT id FROM subscriptions WHERE stripe_subscription_id='sub_student1_monthly'),'in_student1_202602','pi_sub_student1_202602','succeeded',999,'usd',59,'',1,'2026-02-01 10:00:00'),
-- ((SELECT id FROM subscriptions WHERE stripe_subscription_id='sub_student1_monthly'),'in_student1_202603','pi_sub_student1_202603','succeeded',999,'usd',59,'',1,'2026-03-01 10:00:00'),
-- ((SELECT id FROM subscriptions WHERE stripe_subscription_id='sub_student2_yearly'),'in_student2_202501','pi_sub_student2_202501','succeeded',19999,'usd',610,'',1,'2025-01-01 09:00:00'),
-- ((SELECT id FROM subscriptions WHERE stripe_subscription_id='sub_student2_yearly'),'in_student2_202601','pi_sub_student2_202601','succeeded',19999,'usd',610,'',1,'2026-01-01 09:00:00'),
-- ((SELECT id FROM subscriptions WHERE stripe_subscription_id='sub_student3_monthly'),'in_student3_202602','pi_sub_student3_202602','succeeded',999,'usd',59,'',1,'2026-02-15 08:30:00'),
-- ((SELECT id FROM subscriptions WHERE stripe_subscription_id='sub_student3_monthly'),'in_student3_202603','pi_sub_student3_202603','failed',999,'usd',0,'card_declined',2,NULL),
-- ((SELECT id FROM subscriptions WHERE stripe_subscription_id='sub_student4_monthly'),'in_student4_202603','pi_sub_student4_202603','pending',999,'usd',0,'',1,NULL),
-- ((SELECT id FROM subscriptions WHERE stripe_subscription_id='sub_student4_monthly'),'in_student4_202604','pi_sub_student4_202604','succeeded',999,'usd',59,'',1,'2026-04-01 08:00:00'),
-- ((SELECT id FROM subscriptions WHERE stripe_subscription_id='sub_student5_monthly'),'in_student5_202603','pi_sub_student5_202603','succeeded',1999,'usd',88,'',1,'2026-03-01 07:50:00'),
-- ((SELECT id FROM subscriptions WHERE stripe_subscription_id='sub_student5_monthly'),'in_student5_202604','pi_sub_student5_202604','succeeded',1999,'usd',88,'',1,'2026-04-01 07:50:00'),
-- ((SELECT id FROM subscriptions WHERE stripe_subscription_id='sub_student1_cycle_202601'),'in_student1_cycle_202601','pi_student1_cycle_202601','succeeded',999,'usd',59,'',1,'2026-01-01 11:00:00'),
-- ((SELECT id FROM subscriptions WHERE stripe_subscription_id='sub_student1_cycle_202602'),'in_student1_cycle_202602','pi_student1_cycle_202602','succeeded',999,'usd',59,'',1,'2026-02-01 11:00:00'),
-- ((SELECT id FROM subscriptions WHERE stripe_subscription_id='sub_student1_cycle_202603'),'in_student1_cycle_202603','pi_student1_cycle_202603','succeeded',999,'usd',59,'',1,'2026-03-01 11:00:00'),
-- ((SELECT id FROM subscriptions WHERE stripe_subscription_id='sub_student2_cycle_202603'),'in_student2_cycle_202603','pi_student2_cycle_202603','succeeded',999,'usd',59,'',1,'2026-03-01 12:00:00'),
-- ((SELECT id FROM subscriptions WHERE stripe_subscription_id='sub_student4_pro_monthly'),'in_student4_pro_202603','pi_student4_pro_202603','succeeded',1999,'usd',88,'',1,'2026-03-01 09:30:00'),
-- ((SELECT id FROM subscriptions WHERE stripe_subscription_id='sub_student4_pro_monthly'),'in_student4_pro_202604','pi_student4_pro_202604','succeeded',1999,'usd',88,'',1,'2026-04-01 09:30:00'),
-- ((SELECT id FROM subscriptions WHERE stripe_subscription_id='sub_student3_pro_cycle_202603'),'in_student3_pro_202603','pi_student3_pro_202603','succeeded',1999,'usd',88,'',1,'2026-03-16 08:30:00'),
-- ((SELECT id FROM subscriptions WHERE stripe_subscription_id='sub_student3_pro_cycle_202603'),'in_student3_pro_202604','pi_student3_pro_202604','succeeded',1999,'usd',88,'',1,'2026-04-01 08:30:00');

-- ------------------- Create sample course_purchases --------------------
-- INSERT INTO course_purchases (user_id, coupon_id, amount, stripe_fee, purchased_at, status, stripe_checkout_session_id, stripe_payment_intent_id) VALUES
-- ((SELECT id FROM users WHERE username='student1'),NULL,8998,291,'2026-03-31 09:10:00',CAST('paid' AS course_purchase_status_enum),'cs_student1','pi_student1'),
-- ((SELECT id FROM users WHERE username='student2'),(SELECT id FROM coupons WHERE code='SAVE10'),9898,317,'2026-03-31 10:25:00',CAST('paid' AS course_purchase_status_enum),'cs_student2','pi_student2'),
-- ((SELECT id FROM users WHERE username='student3'),NULL,3999,146,'2026-04-01 08:45:00',CAST('paid' AS course_purchase_status_enum),'cs_student3','pi_student3'),
-- ((SELECT id FROM users WHERE username='student4'),(SELECT id FROM coupons WHERE code='WELCOME20'),4799,169,'2026-04-01 11:20:00',CAST('paid' AS course_purchase_status_enum),'cs_student4','pi_student4'),
-- ((SELECT id FROM users WHERE username='instructor1'),NULL,12498,392,'2026-03-31 12:15:00',CAST('paid' AS course_purchase_status_enum),'cs_instructor1_1','pi_instructor1_1'),
-- ((SELECT id FROM users WHERE username='instructor1'),(SELECT id FROM coupons WHERE code='SAVE10'),7648,252,'2026-04-01 09:35:00',CAST('paid' AS course_purchase_status_enum),'cs_instructor1_2','pi_instructor1_2'),
-- ((SELECT id FROM users WHERE username='instructor1'),(SELECT id FROM coupons WHERE code='FLAT500'),4999,175,'2026-04-02 10:40:00',CAST('paid' AS course_purchase_status_enum),'cs_instructor1_3','pi_instructor1_3'),
-- ((SELECT id FROM users WHERE username='instructor1'),NULL,6499,218,'2026-04-02 14:05:00',CAST('paid' AS course_purchase_status_enum),'cs_instructor1_4','pi_instructor1_4'),
-- ((SELECT id FROM users WHERE username='instructor2'),NULL,10498,334,'2026-03-31 15:30:00',CAST('paid' AS course_purchase_status_enum),'cs_instructor2_1','pi_instructor2_1'),
-- ((SELECT id FROM users WHERE username='instructor2'),(SELECT id FROM coupons WHERE code='WELCOME20'),4799,169,'2026-04-01 16:15:00',CAST('paid' AS course_purchase_status_enum),'cs_instructor2_2','pi_instructor2_2'),
-- ((SELECT id FROM users WHERE username='instructor2'),NULL,10998,349,'2026-04-02 09:55:00',CAST('paid' AS course_purchase_status_enum),'cs_instructor2_3','pi_instructor2_3'),
-- ((SELECT id FROM users WHERE username='instructor2'),(SELECT id FROM coupons WHERE code='SAVE10'),3599,134,'2026-04-02 17:25:00',CAST('paid' AS course_purchase_status_enum),'cs_instructor2_4','pi_instructor2_4');

-- ------------------- Create sample course_purchase_details --------------------
-- INSERT INTO course_purchase_details (course_purchase_id, course_id, price) VALUES
-- ((SELECT id FROM course_purchases WHERE stripe_checkout_session_id='cs_student1'),(SELECT id FROM courses WHERE slug='go-backend-development'),4999),
-- ((SELECT id FROM course_purchases WHERE stripe_checkout_session_id='cs_student1'),(SELECT id FROM courses WHERE slug='react-zero-to-hero'),3999),
-- ((SELECT id FROM course_purchases WHERE stripe_checkout_session_id='cs_student2'),(SELECT id FROM courses WHERE slug='go-backend-development'),4999),
-- ((SELECT id FROM course_purchases WHERE stripe_checkout_session_id='cs_student2'),(SELECT id FROM courses WHERE slug='docker-kubernetes-guide'),5999),
-- ((SELECT id FROM course_purchases WHERE stripe_checkout_session_id='cs_student3'),(SELECT id FROM courses WHERE slug='react-zero-to-hero'),3999),
-- ((SELECT id FROM course_purchases WHERE stripe_checkout_session_id='cs_student4'),(SELECT id FROM courses WHERE slug='docker-kubernetes-guide'),5999),

-- ((SELECT id FROM course_purchases WHERE stripe_checkout_session_id='cs_instructor1_1'),(SELECT id FROM courses WHERE slug='docker-kubernetes-guide'),5999),
-- ((SELECT id FROM course_purchases WHERE stripe_checkout_session_id='cs_instructor1_1'),(SELECT id FROM courses WHERE slug='aws-cloud-architecture'),6499),
-- ((SELECT id FROM course_purchases WHERE stripe_checkout_session_id='cs_instructor1_2'),(SELECT id FROM courses WHERE slug='react-zero-to-hero'),3999),
-- ((SELECT id FROM course_purchases WHERE stripe_checkout_session_id='cs_instructor1_2'),(SELECT id FROM courses WHERE slug='python-data-science'),4499),
-- ((SELECT id FROM course_purchases WHERE stripe_checkout_session_id='cs_instructor1_3'),(SELECT id FROM courses WHERE slug='microservices-go'),5499),
-- ((SELECT id FROM course_purchases WHERE stripe_checkout_session_id='cs_instructor1_4'),(SELECT id FROM courses WHERE slug='aws-cloud-architecture'),6499),

-- ((SELECT id FROM course_purchases WHERE stripe_checkout_session_id='cs_instructor2_1'),(SELECT id FROM courses WHERE slug='go-backend-development'),4999),
-- ((SELECT id FROM course_purchases WHERE stripe_checkout_session_id='cs_instructor2_1'),(SELECT id FROM courses WHERE slug='microservices-go'),5499),
-- ((SELECT id FROM course_purchases WHERE stripe_checkout_session_id='cs_instructor2_2'),(SELECT id FROM courses WHERE slug='docker-kubernetes-guide'),5999),
-- ((SELECT id FROM course_purchases WHERE stripe_checkout_session_id='cs_instructor2_3'),(SELECT id FROM courses WHERE slug='python-data-science'),4499),
-- ((SELECT id FROM course_purchases WHERE stripe_checkout_session_id='cs_instructor2_3'),(SELECT id FROM courses WHERE slug='aws-cloud-architecture'),6499),
-- ((SELECT id FROM course_purchases WHERE stripe_checkout_session_id='cs_instructor2_4'),(SELECT id FROM courses WHERE slug='react-zero-to-hero'),3999);

-- ------------------- Create sample enrollments --------------------
-- INSERT INTO enrollments (user_id, course_id, is_completed, canceled_at, type, enrolled_at, created_at, updated_at) VALUES
-- ((SELECT id FROM users WHERE username='student1'),(SELECT id FROM courses WHERE slug='go-backend-development'),false,NULL,'course_purchase','2026-03-31 09:10:00','2026-03-31 09:10:00','2026-03-31 09:10:00'),
-- ((SELECT id FROM users WHERE username='student1'),(SELECT id FROM courses WHERE slug='react-zero-to-hero'),false,NULL,'course_purchase','2026-03-31 09:10:00','2026-03-31 09:10:00','2026-03-31 09:10:00'),
-- ((SELECT id FROM users WHERE username='student2'),(SELECT id FROM courses WHERE slug='go-backend-development'),false,NULL,'course_purchase','2026-03-31 10:25:00','2026-03-31 10:25:00','2026-03-31 10:25:00'),
-- ((SELECT id FROM users WHERE username='student2'),(SELECT id FROM courses WHERE slug='docker-kubernetes-guide'),false,NULL,'course_purchase','2026-03-31 10:25:00','2026-03-31 10:25:00','2026-03-31 10:25:00'),
-- ((SELECT id FROM users WHERE username='student3'),(SELECT id FROM courses WHERE slug='react-zero-to-hero'),false,NULL,'course_purchase','2026-04-01 08:45:00','2026-04-01 08:45:00','2026-04-01 08:45:00'),

-- -- student1 subscriptions: aws-cloud-architecture (instructor5)
-- ((SELECT id FROM users WHERE username='student1'),(SELECT id FROM courses WHERE slug='aws-cloud-architecture'),false,'2026-01-31 23:59:59','subscription','2026-01-01 08:00:00','2026-01-01 08:00:00','2026-01-31 23:59:59'),
-- ((SELECT id FROM users WHERE username='student1'),(SELECT id FROM courses WHERE slug='aws-cloud-architecture'),false,'2026-02-28 23:59:59','subscription','2026-02-01 08:00:00','2026-02-01 08:00:00','2026-02-28 23:59:59'),
-- ((SELECT id FROM users WHERE username='student1'),(SELECT id FROM courses WHERE slug='aws-cloud-architecture'),false,NULL,'subscription','2026-03-01 08:00:00','2026-03-01 08:00:00','2026-03-01 08:00:00'),

-- -- student1 subscriptions: python-data-science (instructor4)
-- ((SELECT id FROM users WHERE username='student1'),(SELECT id FROM courses WHERE slug='python-data-science'),false,'2026-01-31 23:59:59','subscription','2026-01-01 08:05:00','2026-01-01 08:05:00','2026-01-31 23:59:59'),
-- ((SELECT id FROM users WHERE username='student1'),(SELECT id FROM courses WHERE slug='python-data-science'),false,'2026-02-28 23:59:59','subscription','2026-02-01 08:05:00','2026-02-01 08:05:00','2026-02-28 23:59:59'),
-- ((SELECT id FROM users WHERE username='student1'),(SELECT id FROM courses WHERE slug='python-data-science'),false,NULL,'subscription','2026-03-01 08:05:00','2026-03-01 08:05:00','2026-03-01 08:05:00'),

-- -- student2 subscriptions: python-data-science (instructor4) and react-zero-to-hero (instructor2)
-- ((SELECT id FROM users WHERE username='student2'),(SELECT id FROM courses WHERE slug='python-data-science'),false,'2026-03-31 23:59:59','subscription','2026-03-01 09:00:00','2026-03-01 09:00:00','2026-03-31 23:59:59'),
-- ((SELECT id FROM users WHERE username='student2'),(SELECT id FROM courses WHERE slug='react-zero-to-hero'),false,'2026-03-31 23:59:59','subscription','2026-03-01 09:15:00','2026-03-01 09:15:00','2026-03-31 23:59:59'),

-- -- student3 subscriptions: microservices-go (instructor1) and go-backend-development (instructor1)
-- ((SELECT id FROM users WHERE username='student3'),(SELECT id FROM courses WHERE slug='microservices-go'),false,'2026-03-15 23:59:59','subscription','2026-02-15 08:00:00','2026-02-15 08:00:00','2026-03-15 23:59:59'),
-- ((SELECT id FROM users WHERE username='student3'),(SELECT id FROM courses WHERE slug='microservices-go'),false,NULL,'subscription','2026-03-16 08:00:00','2026-03-16 08:00:00','2026-03-16 08:00:00'),
-- ((SELECT id FROM users WHERE username='student3'),(SELECT id FROM courses WHERE slug='go-backend-development'),false,'2026-03-15 23:59:59','subscription','2026-02-15 08:10:00','2026-02-15 08:10:00','2026-03-15 23:59:59'),
-- ((SELECT id FROM users WHERE username='student3'),(SELECT id FROM courses WHERE slug='go-backend-development'),false,NULL,'subscription','2026-03-16 08:10:00','2026-03-16 08:10:00','2026-03-16 08:10:00'),

-- -- student4 premium subscriptions (no active courses in this cycle)
-- ((SELECT id FROM users WHERE username='student4'),(SELECT id FROM courses WHERE slug='docker-kubernetes-guide'),false,NULL,'subscription','2026-03-25 08:00:00','2026-03-25 08:00:00','2026-03-25 08:00:00'),

-- -- student4 pro subscriptions: python-data-science (instructor4) and react-zero-to-hero (instructor2)
-- ((SELECT id FROM users WHERE username='student4'),(SELECT id FROM courses WHERE slug='python-data-science'),false,NULL,'subscription','2026-03-01 09:00:00','2026-03-01 09:00:00','2026-03-01 09:00:00'),
-- ((SELECT id FROM users WHERE username='student4'),(SELECT id FROM courses WHERE slug='react-zero-to-hero'),false,NULL,'subscription','2026-03-01 09:30:00','2026-03-01 09:30:00','2026-03-01 09:30:00'),

-- -- student5 subscriptions: go-backend-development (instructor1)
-- ((SELECT id FROM users WHERE username='student5'),(SELECT id FROM courses WHERE slug='go-backend-development'),false,NULL,'subscription','2026-03-01 07:50:00','2026-03-01 07:50:00','2026-03-01 07:50:00'),

-- -- student3 pro cycle: docker-kubernetes-guide (instructor3) and aws-cloud-architecture (instructor5)
-- ((SELECT id FROM users WHERE username='student3'),(SELECT id FROM courses WHERE slug='docker-kubernetes-guide'),false,NULL,'subscription','2026-03-16 08:00:00','2026-03-16 08:00:00','2026-03-16 08:00:00'),
-- ((SELECT id FROM users WHERE username='student3'),(SELECT id FROM courses WHERE slug='aws-cloud-architecture'),false,NULL,'subscription','2026-03-16 08:30:00','2026-03-16 08:30:00','2026-03-16 08:30:00');

------------------- Create sample blogs --------------------
INSERT INTO blogs (title, slug, content, category_id, image_url, author_id, status, scheduled_at, tags) VALUES
('Getting Started with Go: A Beginner Guide','getting-started-go-guide','<h1>Getting Started with Go</h1><p>Go is a powerful programming language...</p>Go essentials explained',(SELECT id FROM categories WHERE name='Backend Development'),'https://example.com/blog-go.jpg',(SELECT id FROM users WHERE username='instructor1'),'published',NULL,'["go","programming","tutorial"]'),
('React Hooks: Complete Guide','react-hooks-complete-guide','<h1>React Hooks Guide</h1><p>useState, useEffect, and custom hooks...</p>Advanced hook patterns',(SELECT id FROM categories WHERE name='Frontend Development'),'https://example.com/blog-react.jpg',(SELECT id FROM users WHERE username='instructor2'),'published',NULL,'["react","hooks","javascript"]'),
('Docker in Production','docker-in-production','<h1>Docker Production Guide</h1><p>Security, monitoring, and optimization...</p>Production ready setup',(SELECT id FROM categories WHERE name='DevOps'),'https://example.com/blog-docker.jpg',(SELECT id FROM users WHERE username='instructor3'),'published',NULL,'["docker","devops","production"]'),
('Microservices Architecture Patterns','microservices-patterns','<h1>Microservices Patterns</h1><p>CQRS, Event Sourcing, Saga pattern...</p>Real world examples',(SELECT id FROM categories WHERE name='Backend Development'),'https://example.com/blog-micro.jpg',(SELECT id FROM users WHERE username='instructor1'),'published',NULL,'["microservices","architecture","design"]'),
('Data Science with Python','data-science-python','<h1>Data Science Intro</h1><p>NumPy, Pandas, Scikit-learn...</p>Hands-on examples',(SELECT id FROM categories WHERE name='Data Science'),'https://example.com/blog-ds.jpg',(SELECT id FROM users WHERE username='instructor4'),'draft',NULL,'["python","datascience","ml"]'),
('Cloud Architecture Best Practices','cloud-architecture-practices','<h1>Cloud Architecture</h1><p>High availability, disaster recovery...</p>AWS patterns',(SELECT id FROM categories WHERE name='Cloud Computing'),'https://example.com/blog-cloud.jpg',(SELECT id FROM users WHERE username='instructor5'),'published',NULL,'["cloud","aws","architecture"]');

------------------- Create sample follows --------------------
INSERT INTO follows (follower_id, followee_id) VALUES
((SELECT id FROM users WHERE username='student1'),(SELECT id FROM users WHERE username='instructor1')),
((SELECT id FROM users WHERE username='student1'),(SELECT id FROM users WHERE username='instructor2')),
((SELECT id FROM users WHERE username='student2'),(SELECT id FROM users WHERE username='instructor1')),
((SELECT id FROM users WHERE username='student2'),(SELECT id FROM users WHERE username='instructor3')),
((SELECT id FROM users WHERE username='student3'),(SELECT id FROM users WHERE username='instructor2')),
((SELECT id FROM users WHERE username='student4'),(SELECT id FROM users WHERE username='instructor3')),
((SELECT id FROM users WHERE username='student4'),(SELECT id FROM users WHERE username='instructor4')),
((SELECT id FROM users WHERE username='student5'),(SELECT id FROM users WHERE username='instructor1')),
((SELECT id FROM users WHERE username='student5'),(SELECT id FROM users WHERE username='instructor5')),
((SELECT id FROM users WHERE username='instructor1'),(SELECT id FROM users WHERE username='instructor2'));

------------------- Create sample feedbacks --------------------
INSERT INTO feedbacks (user_id, course_id, rate, content) VALUES
((SELECT id FROM users WHERE username='student1'),(SELECT id FROM courses WHERE slug='go-backend-development'),5,'Excellent course! Instructor1 is amazing at explaining complex concepts.'),
((SELECT id FROM users WHERE username='student1'),(SELECT id FROM courses WHERE slug='react-zero-to-hero'),4,'Great content, but could use more projects.'),
((SELECT id FROM users WHERE username='student2'),(SELECT id FROM courses WHERE slug='go-backend-development'),4,'Very detailed and comprehensive. Learned a lot.'),
((SELECT id FROM users WHERE username='student2'),(SELECT id FROM courses WHERE slug='docker-kubernetes-guide'),4,'Best Docker/K8s course I have seen!'),
((SELECT id FROM users WHERE username='student3'),(SELECT id FROM courses WHERE slug='react-zero-to-hero'),4,'Instructor2 teaches React really well.'),
((SELECT id FROM users WHERE username='student4'),(SELECT id FROM courses WHERE slug='docker-kubernetes-guide'),4,'Practical and very useful for production.'),
((SELECT id FROM users WHERE username='student5'),(SELECT id FROM courses WHERE slug='go-backend-development'),4,'Good fundamentals, advanced sections are challenging.');

------------------- Create sample carts --------------------
INSERT INTO carts (user_id, course_id) VALUES
((SELECT id FROM users WHERE username='student1'),(SELECT id FROM courses WHERE slug='microservices-go')),
((SELECT id FROM users WHERE username='student1'),(SELECT id FROM courses WHERE slug='python-data-science')),
((SELECT id FROM users WHERE username='student2'),(SELECT id FROM courses WHERE slug='aws-cloud-architecture')),
((SELECT id FROM users WHERE username='student3'),(SELECT id FROM courses WHERE slug='microservices-go')),
((SELECT id FROM users WHERE username='student4'),(SELECT id FROM courses WHERE slug='python-data-science')),
((SELECT id FROM users WHERE username='student5'),(SELECT id FROM courses WHERE slug='react-zero-to-hero'));

------------------- Create sample stripe_events --------------------
INSERT INTO stripe_events (event_type, event_id) VALUES
('payment_intent.succeeded','evt_1'),
('payment_intent.succeeded','evt_2'),
('payment_intent.succeeded','evt_3'),
('payment_intent.succeeded','evt_4'),
('customer.subscription.created','evt_5'),
('customer.subscription.created','evt_6'),
('invoice.paid','evt_7'),
('invoice.paid','evt_8');





