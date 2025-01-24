-- +goose Up
-- +goose StatementBegin
-- ----- USERS -----
INSERT INTO users (full_name, username, email, password, profile_image) 
VALUES ('Test test', 'test', 'test@example.com', '$2a$14$ZqZ1FmMgZNYvO.Q2rSht3.fGTX4IBq6VJMBoJ7bRXMAaEQk3pAP9i', 'https://twitter-clone-tzjvdg.s3.ap-southeast-1.amazonaws.com/purple-1.png');
INSERT INTO users (full_name, username, email, password, profile_image) 
VALUES ('Test test 2', 'test2', 'test2@example.com', '$2a$14$ZqZ1FmMgZNYvO.Q2rSht3.fGTX4IBq6VJMBoJ7bRXMAaEQk3pAP9i', 'https://twitter-clone-tzjvdg.s3.ap-southeast-1.amazonaws.com/purple-1.png');



-- ----- FOLLOWS -----
-- INSERT INTO follows (follower_id, following_id) 
-- VALUES (
--     (SELECT id FROM users WHERE email = 'test@example.com' LIMIT 1), 
--     (SELECT id FROM users WHERE email = 'test2@example.com' LIMIT 1)
-- );

-- ----- TWEETS -----
INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'content', 
    (SELECT id FROM users WHERE email = 'test@example.com' LIMIT 1), NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM follows;

DELETE FROM tweets;

DELETE FROM users;
-- +goose StatementEnd