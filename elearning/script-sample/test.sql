SELECT * FROM get_instructor_course_purchase_revenue('fd90b491-1aaa-4a98-be3d-36a6f337891a');
SELECT * FROM get_instructor_course_purchase_revenue_by_day('fd90b491-1aaa-4a98-be3d-36a6f337891a', '2026-03-31');
SELECT * FROM get_instructor_course_purchase_revenue_by_month('fd90b491-1aaa-4a98-be3d-36a6f337891a', 2026, 3);
SELECT * FROM get_instructor_course_purchase_revenue_by_year('fd90b491-1aaa-4a98-be3d-36a6f337891a', 2026);

SELECT * FROM get_admin_course_purchase_revenue();
SELECT * FROM get_admin_course_purchase_revenue_by_day('2026-03-31');
SELECT * FROM get_admin_course_purchase_revenue_by_month(2026, 3);
SELECT * FROM get_admin_course_purchase_revenue_by_year(2026);


SELECT * FROM get_instructor_subscription_revenue('fd90b491-1aaa-4a98-be3d-36a6f337891a');
SELECT * FROM get_instructor_subscription_revenue_by_day('fd90b491-1aaa-4a98-be3d-36a6f337891a', '2026-03-01');
SELECT * FROM get_instructor_subscription_revenue_by_month('fd90b491-1aaa-4a98-be3d-36a6f337891a', 2026, 3);
SELECT * FROM get_instructor_subscription_revenue_by_year('fd90b491-1aaa-4a98-be3d-36a6f337891a', 2026);

SELECT * FROM get_admin_subscription_revenue();
SELECT * FROM get_admin_subscription_revenue_by_day('2026-03-01');
SELECT * FROM get_admin_subscription_revenue_by_month(2026, 3);
SELECT * FROM get_admin_subscription_revenue_by_year(2026);