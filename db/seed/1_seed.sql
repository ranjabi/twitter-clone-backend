-- +goose Up
-- +goose StatementBegin
--===== USERS =====
INSERT INTO users (full_name, username, email, password, profile_image) 
VALUES ('John Doe', 'johndoe', 'johndoe@example.com', '$2a$14$ZqZ1FmMgZNYvO.Q2rSht3.fGTX4IBq6VJMBoJ7bRXMAaEQk3pAP9i', 'https://twitter-clone-tzjvdg.s3.ap-southeast-1.amazonaws.com/blue-1.png');
INSERT INTO users (full_name, username, email, password, profile_image) 
VALUES ('Jane Smith', 'janesmith', 'janesmith@example.com', '$2a$14$ZqZ1FmMgZNYvO.Q2rSht3.fGTX4IBq6VJMBoJ7bRXMAaEQk3pAP9i', 'https://twitter-clone-tzjvdg.s3.ap-southeast-1.amazonaws.com/purple-2.png');
INSERT INTO users (full_name, username, email, password, profile_image) 
VALUES ('Mike Johnson', 'mikejohnson', 'mikejohnson@example.com', '$2a$14$ZqZ1FmMgZNYvO.Q2rSht3.fGTX4IBq6VJMBoJ7bRXMAaEQk3pAP9i', 'https://twitter-clone-tzjvdg.s3.ap-southeast-1.amazonaws.com/yellow-2.png');
INSERT INTO users (full_name, username, email, password, profile_image) 
VALUES ('Emily Davis', 'emilydavis', 'emilydavis@example.com', '$2a$14$ZqZ1FmMgZNYvO.Q2rSht3.fGTX4IBq6VJMBoJ7bRXMAaEQk3pAP9i', 'https://twitter-clone-tzjvdg.s3.ap-southeast-1.amazonaws.com/purple-3.png');
INSERT INTO users (full_name, username, email, password, profile_image) 
VALUES ('Chris Brown', 'chrisbrown', 'chrisbrown@example.com', '$2a$14$ZqZ1FmMgZNYvO.Q2rSht3.fGTX4IBq6VJMBoJ7bRXMAaEQk3pAP9i', 'https://twitter-clone-tzjvdg.s3.ap-southeast-1.amazonaws.com/yellow-3.png');
INSERT INTO users (full_name, username, email, password, profile_image) 
VALUES ('Sarah Wilson', 'sarahwilson', 'sarahwilson@example.com', '$2a$14$ZqZ1FmMgZNYvO.Q2rSht3.fGTX4IBq6VJMBoJ7bRXMAaEQk3pAP9i', 'https://twitter-clone-tzjvdg.s3.ap-southeast-1.amazonaws.com/blue-4.png');
INSERT INTO users (full_name, username, email, password, profile_image) 
VALUES ('David Lee', 'davidlee', 'davidlee@example.com', '$2a$14$ZqZ1FmMgZNYvO.Q2rSht3.fGTX4IBq6VJMBoJ7bRXMAaEQk3pAP9i', 'https://twitter-clone-tzjvdg.s3.ap-southeast-1.amazonaws.com/blue-2.png');
INSERT INTO users (full_name, username, email, password, profile_image) 
VALUES ('Olivia Martinez', 'oliviamartinez', 'oliviamartinez@example.com', '$2a$14$ZqZ1FmMgZNYvO.Q2rSht3.fGTX4IBq6VJMBoJ7bRXMAaEQk3pAP9i', 'https://twitter-clone-tzjvdg.s3.ap-southeast-1.amazonaws.com/yellow-1.png');
INSERT INTO users (full_name, username, email, password, profile_image) 
VALUES ('Daniel Clark', 'danielclark', 'danielclark@example.com', '$2a$14$ZqZ1FmMgZNYvO.Q2rSht3.fGTX4IBq6VJMBoJ7bRXMAaEQk3pAP9i', 'https://twitter-clone-tzjvdg.s3.ap-southeast-1.amazonaws.com/blue-3.png');
INSERT INTO users (full_name, username, email, password, profile_image) 
VALUES ('Sophia Anderson', 'sophiaanderson', 'sophiaanderson@example.com', '$2a$14$ZqZ1FmMgZNYvO.Q2rSht3.fGTX4IBq6VJMBoJ7bRXMAaEQk3pAP9i', 'https://twitter-clone-tzjvdg.s3.ap-southeast-1.amazonaws.com/purple-1.png');



--===== 3 TWEETS BY EACH USER =====
-- John Doe
INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Just finished a great workout at the gym. Feeling pumped and ready to tackle the day!', 
    (SELECT id FROM users WHERE email = 'johndoe@example.com' LIMIT 1), NOW()
);

INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Weekend road trip with friends! Can’t wait to hit the road and explore new places.', 
    (SELECT id FROM users WHERE email = 'johndoe@example.com' LIMIT 1), NOW() - INTERVAL '1 hour'
);

INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Finally got my new gaming setup! Who’s ready to game this weekend? 🎮', 
    (SELECT id FROM users WHERE email = 'johndoe@example.com' LIMIT 1), NOW() - INTERVAL '2 hour'
);

-- Jane Smith
INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Just had the best lunch ever! Trying out new recipes has been so fun lately!', 
    (SELECT id FROM users WHERE email = 'janesmith@example.com' LIMIT 1), NOW()
);

INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Can’t believe how fast this year has flown by. Time to start planning for the holidays!', 
    (SELECT id FROM users WHERE email = 'janesmith@example.com' LIMIT 1), NOW() - INTERVAL '1 hour'
);

INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Learning how to code is challenging, but I’m loving every minute of it! 💻🚀', 
    (SELECT id FROM users WHERE email = 'janesmith@example.com' LIMIT 1), NOW() - INTERVAL '2 hour'
);

-- Mike Johnson
INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Just watched the latest episode of my favorite show. The plot twists are insane!', 
    (SELECT id FROM users WHERE email = 'mikejohnson@example.com' LIMIT 1), NOW()
);

INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Trying out a new restaurant tonight. The reviews are amazing, hope the food is too!', 
    (SELECT id FROM users WHERE email = 'mikejohnson@example.com' LIMIT 1), NOW() - INTERVAL '1 hour'
);

INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Spent the afternoon at the beach. Nothing beats the sound of the ocean waves.', 
    (SELECT id FROM users WHERE email = 'mikejohnson@example.com' LIMIT 1), NOW() - INTERVAL '2 hour'
);

-- Emily Davis
INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Can’t believe I finally finished that project. Feels great to check it off my list!', 
    (SELECT id FROM users WHERE email = 'emilydavis@example.com' LIMIT 1), NOW()
);

INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Coffee and coding go hand in hand. Productive afternoon ahead!', 
    (SELECT id FROM users WHERE email = 'emilydavis@example.com' LIMIT 1), NOW() - INTERVAL '1 hour'
);

INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Loving this new playlist I found. Perfect background music for getting work done!', 
    (SELECT id FROM users WHERE email = 'emilydavis@example.com' LIMIT 1), NOW() - INTERVAL '2 hour'
);

-- Chris Brown
INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Just wrapped up my latest vlog. Editing took forever, but I’m excited to share it!', 
    (SELECT id FROM users WHERE email = 'chrisbrown@example.com' LIMIT 1), NOW()
);

INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Taking a break from social media for the weekend. Time to unplug and recharge.', 
    (SELECT id FROM users WHERE email = 'chrisbrown@example.com' LIMIT 1), NOW() - INTERVAL '1 hour'
);

INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Started reading a new book today. The first few chapters are already so good!', 
    (SELECT id FROM users WHERE email = 'chrisbrown@example.com' LIMIT 1), NOW() - INTERVAL '2 hour'
);

-- Sarah Wilson
INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Had the most amazing dinner with friends tonight. Great food, even better company!', 
    (SELECT id FROM users WHERE email = 'sarahwilson@example.com' LIMIT 1), NOW()
);

INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Just signed up for a yoga class. Time to relax and focus on my well-being.', 
    (SELECT id FROM users WHERE email = 'sarahwilson@example.com' LIMIT 1), NOW() - INTERVAL '1 hour'
);

INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Woke up early to catch the sunrise this morning. Totally worth it!', 
    (SELECT id FROM users WHERE email = 'sarahwilson@example.com' LIMIT 1), NOW() - INTERVAL '2 hour'
);

-- David Lee
INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Feeling inspired after attending a tech conference. So many cool ideas to explore!', 
    (SELECT id FROM users WHERE email = 'davidlee@example.com' LIMIT 1), NOW()
);

INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Learning a new programming language today. It’s tough, but I’m up for the challenge!', 
    (SELECT id FROM users WHERE email = 'davidlee@example.com' LIMIT 1), NOW() - INTERVAL '1 hour'
);

INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Just fixed a bug that’s been driving me crazy for hours. Success!', 
    (SELECT id FROM users WHERE email = 'davidlee@example.com' LIMIT 1), NOW() - INTERVAL '2 hour'
);

-- Olivia Martinez
INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Just booked my next travel destination! Can’t wait to share the adventure with you all.', 
    (SELECT id FROM users WHERE email = 'oliviamartinez@example.com' LIMIT 1), NOW()
);

INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Trying out a new hobby this weekend. Time to pick up the paintbrush and get creative!', 
    (SELECT id FROM users WHERE email = 'oliviamartinez@example.com' LIMIT 1), NOW() - INTERVAL '1 hour'
);

INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Had a super productive day today. Checking off all the items on my to-do list!', 
    (SELECT id FROM users WHERE email = 'oliviamartinez@example.com' LIMIT 1), NOW() - INTERVAL '2 hour'
);

-- Daniel Clark
INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Excited to announce my new side project. Stay tuned for updates!', 
    (SELECT id FROM users WHERE email = 'danielclark@example.com' LIMIT 1), NOW()
);

INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Spent the afternoon hiking. The views from the top were absolutely breathtaking!', 
    (SELECT id FROM users WHERE email = 'danielclark@example.com' LIMIT 1), NOW() - INTERVAL '1 hour'
);

INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Grabbing dinner with some old friends tonight. Looking forward to catching up!', 
    (SELECT id FROM users WHERE email = 'danielclark@example.com' LIMIT 1), NOW() - INTERVAL '2 hour'
);

-- Sophia Anderson
INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Starting a new project today. Excited to see where this journey takes me!', 
    (SELECT id FROM users WHERE email = 'sophiaanderson@example.com' LIMIT 1), NOW()
);

INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Binge-watching my favorite series tonight. Perfect way to unwind after a long day.', 
    (SELECT id FROM users WHERE email = 'sophiaanderson@example.com' LIMIT 1), NOW() - INTERVAL '1 hour'
);

INSERT INTO tweets (content, user_id, created_at)  
VALUES (
    'Cooking up a storm in the kitchen tonight. Who else loves experimenting with new recipes?', 
    (SELECT id FROM users WHERE email = 'sophiaanderson@example.com' LIMIT 1), NOW() - INTERVAL '2 hour'
);



--===== FOLLOWS =====
-- John Doe follows other 9 users
INSERT INTO follows (follower_id, following_id) 
VALUES (
    (SELECT id FROM users WHERE email = 'johndoe@example.com' LIMIT 1), 
    (SELECT id FROM users WHERE email = 'janesmith@example.com' LIMIT 1)
);

INSERT INTO follows (follower_id, following_id) 
VALUES (
    (SELECT id FROM users WHERE email = 'johndoe@example.com' LIMIT 1), 
    (SELECT id FROM users WHERE email = 'mikejohnson@example.com' LIMIT 1)
);

INSERT INTO follows (follower_id, following_id) 
VALUES (
    (SELECT id FROM users WHERE email = 'johndoe@example.com' LIMIT 1), 
    (SELECT id FROM users WHERE email = 'emilydavis@example.com' LIMIT 1)
);

INSERT INTO follows (follower_id, following_id) 
VALUES (
    (SELECT id FROM users WHERE email = 'johndoe@example.com' LIMIT 1), 
    (SELECT id FROM users WHERE email = 'chrisbrown@example.com' LIMIT 1)
);

INSERT INTO follows (follower_id, following_id) 
VALUES (
    (SELECT id FROM users WHERE email = 'johndoe@example.com' LIMIT 1), 
    (SELECT id FROM users WHERE email = 'sarahwilson@example.com' LIMIT 1)
);

INSERT INTO follows (follower_id, following_id) 
VALUES (
    (SELECT id FROM users WHERE email = 'johndoe@example.com' LIMIT 1), 
    (SELECT id FROM users WHERE email = 'davidlee@example.com' LIMIT 1)
);

INSERT INTO follows (follower_id, following_id) 
VALUES (
    (SELECT id FROM users WHERE email = 'johndoe@example.com' LIMIT 1), 
    (SELECT id FROM users WHERE email = 'oliviamartinez@example.com' LIMIT 1)
);

INSERT INTO follows (follower_id, following_id) 
VALUES (
    (SELECT id FROM users WHERE email = 'johndoe@example.com' LIMIT 1), 
    (SELECT id FROM users WHERE email = 'danielclark@example.com' LIMIT 1)
);

INSERT INTO follows (follower_id, following_id) 
VALUES (
    (SELECT id FROM users WHERE email = 'johndoe@example.com' LIMIT 1), 
    (SELECT id FROM users WHERE email = 'sophiaanderson@example.com' LIMIT 1)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM follows;

DELETE FROM tweets;

DELETE FROM users;
-- +goose StatementEnd