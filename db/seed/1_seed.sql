-- +goose Up
-- +goose StatementBegin
INSERT INTO users (username, email, password) 
VALUES ('username', 'email@email.com', '$2a$14$ZqZ1FmMgZNYvO.Q2rSht3.fGTX4IBq6VJMBoJ7bRXMAaEQk3pAP9i'); -- password
INSERT INTO users (username, email, password) 
VALUES ('username2', 'email2@email.com', '$2a$14$mabfBxkkjs2s6l60tJFo8ucUYGcBtcrH5dBtdmUIC20nArmQNyoyK'); -- password2
INSERT INTO users (username, email, password) 
VALUES ('username3', 'email3@email.com', '$2a$14$xdif3Of1bxQQs3zw8tm/vua5YnoronphqASgIpaaMjIiL1AjufrIW'); -- password3

-- insert 20 tweets
INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin tinciduntlibero nec nulla facilisis.', 
    (SELECT id FROM users WHERE email = 'email@email.com' LIMIT 1), NOW()
);
INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Fusce vehicula quam eget neque venenatis, eget vestibulum metus tristique. Curabitur efficitur lacus.', 
    (SELECT id FROM users WHERE email = 'email@email.com' LIMIT 1), NOW() - INTERVAL '1 hour'
);
INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Suspendisse potenti. Vivamus in velit vitae ligula interdum malesuada. Mauris gravida quam sit amet.', 
    (SELECT id FROM users WHERE email = 'email@email.com' LIMIT 1), NOW() - INTERVAL '2 hour'
);
INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Nulla vel.', 
    (SELECT id FROM users WHERE email = 'email@email.com' LIMIT 1), NOW() - INTERVAL '3 hour'
);
INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Aenean nec leo luctus, ornare arcu sed, ultricies lacus. Nulla at orci eget nunc laoreet varius in.', 
    (SELECT id FROM users WHERE email = 'email@email.com' LIMIT 1), NOW() - INTERVAL '4 hour'
);
INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. In.', 
    (SELECT id FROM users WHERE email = 'email@email.com' LIMIT 1), NOW() - INTERVAL '5 hour'
);
INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Nullam aliquet justo ut tortor scelerisque, vel fringilla elit laoreet. Integer in mi euismod justo.', 
    (SELECT id FROM users WHERE email = 'email@email.com' LIMIT 1), NOW() - INTERVAL '6 hour'
);
INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Morbi venenatis, nisi at condimentum fringilla, dolor justo vulputate augue, non volutpat felis dui.', 
    (SELECT id FROM users WHERE email = 'email@email.com' LIMIT 1), NOW() - INTERVAL '7 hour'
);
INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Etiam fermentum odio sit amet neque lacinia, vel bibendum sapien vestibulum. Suspendisse nec risus.', 
    (SELECT id FROM users WHERE email = 'email@email.com' LIMIT 1), NOW() - INTERVAL '8 hour'
);
INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Donec gravida justo non lectus ultricies, a tincidunt lacus scelerisque. Duis ultricies fermentum.', 
    (SELECT id FROM users WHERE email = 'email@email.com' LIMIT 1), NOW() - INTERVAL '9 hour'
);
INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Mauris nec odio eget justo pharetra iaculis. Phasellus condimentum magna non augue commodo lacinia.', 
    (SELECT id FROM users WHERE email = 'email@email.com' LIMIT 1), NOW() - INTERVAL '10 hour'
);
INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Quisque vel nunc a arcu fermentum fringilla ut eu augue. Sed pharetra sapien sed justo malesuada.', 
    (SELECT id FROM users WHERE email = 'email@email.com' LIMIT 1), NOW() - INTERVAL '11 hour'
);
INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Ut egestas, justo nec pellentesque viverra, nisi elit suscipit nisl, vitae feugiat mi eros ut nulla.', 
    (SELECT id FROM users WHERE email = 'email@email.com' LIMIT 1), NOW() - INTERVAL '12 hour'
);
INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Sed quis lectus vehicula, hendrerit felis vitae, efficitur arcu. Praesent feugiat nulla sed orci ornare.', 
    (SELECT id FROM users WHERE email = 'email@email.com' LIMIT 1), NOW() - INTERVAL '13 hour'
);
INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Aliquam bibendum tortor sed tortor gravida, quis interdum velit tristique. Integer cursus purus sit.', 
    (SELECT id FROM users WHERE email = 'email@email.com' LIMIT 1), NOW() - INTERVAL '14 hour'
);
INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'In convallis, lacus eu vestibulum placerat, felis velit hendrerit libero, at gravida turpis sapien et.', 
    (SELECT id FROM users WHERE email = 'email@email.com' LIMIT 1), NOW() - INTERVAL '15 hour'
);
INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Duis vel ipsum nec nulla eleifend pharetra. Nam scelerisque lorem id elit sollicitudin, nec tempor ex.', 
    (SELECT id FROM users WHERE email = 'email@email.com' LIMIT 1), NOW() - INTERVAL '16 hour'
);
INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Vivamus at velit vitae justo sollicitudin malesuada. Cras vitae urna non ligula dapibus hendrerit.', 
    (SELECT id FROM users WHERE email = 'email@email.com' LIMIT 1), NOW() - INTERVAL '17 hour'
);
INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Aenean feugiat dui ut orci convallis faucibus. Aenean viverra, sapien id posuere sagittis, enim erat.', 
    (SELECT id FROM users WHERE email = 'email@email.com' LIMIT 1), NOW() - INTERVAL '18 hour'
);
INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Praesent eget sapien a nibh interdum tristique. Proin tincidunt magna vitae felis consectetur.', 
    (SELECT id FROM users WHERE email = 'email@email.com' LIMIT 1), NOW() - INTERVAL '19 hour'
);

-- email follow email2
INSERT INTO follows (follower_id, following_id) 
VALUES (
    (SELECT id FROM users WHERE email = 'email@email.com' LIMIT 1), 
    (SELECT id FROM users WHERE email = 'email2@email.com' LIMIT 1)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM follows
WHERE follower_id=(SELECT id FROM users WHERE email = 'email@email.com' LIMIT 1) AND following_id=(SELECT id FROM users WHERE email = 'email2@email.com' LIMIT 1);

DELETE FROM tweets 
WHERE user_id=(SELECT id FROM users WHERE email = 'email@email.com' LIMIT 1);

DELETE FROM users
WHERE email='email@email.com';
DELETE FROM users
WHERE email='email2@email.com';
DELETE FROM users
WHERE email='email3@email.com';
-- +goose StatementEnd