-- +goose Up
-- +goose StatementBegin
INSERT INTO users (username, email, password) 
VALUES ('username', 'email@email.com', 'password');
INSERT INTO users (username, email, password) 
VALUES ('username2', 'email2@email.com', 'password2');

INSERT INTO tweets (content, user_id)  
VALUES ('Eum aliquam vel sed mollitia id eaque. Et quia aperiam.', (SELECT id FROM users WHERE email = 'email@email.com' LIMIT 1));

-- email follow email2
INSERT INTO follows (followers_id, following_id) 
VALUES (
    (SELECT id FROM users WHERE email = 'email@email.com' LIMIT 1), 
    (SELECT id FROM users WHERE email = 'email2@email.com' LIMIT 1)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM follows
WHERE followers_id=(SELECT id FROM users WHERE email = 'email@email.com' LIMIT 1) AND following_id=(SELECT id FROM users WHERE email = 'email2@email.com' LIMIT 1);

DELETE FROM tweets 
WHERE user_id=(SELECT id FROM users WHERE email = 'email@email.com' LIMIT 1);

DELETE FROM users
WHERE email='email@email.com';
-- +goose StatementEnd