-- +goose Up
-- +goose StatementBegin
INSERT INTO users (username, email, password) 
VALUES ('username', 'email@email.com', 'password');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users
WHERE email='email@email.com';
-- +goose StatementEnd
