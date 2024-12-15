-- +goose Up
-- +goose StatementBegin
INSERT INTO tweets (content, user_id)  VALUES ('Eum aliquam vel sed mollitia id eaque. Et quia aperiam.', (SELECT id FROM users WHERE email = 'email@email.com' LIMIT 1));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM tweets where user_id=(SELECT id FROM users WHERE email = 'email@email.com' LIMIT 1);
-- +goose StatementEnd
