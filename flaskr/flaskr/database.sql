CREATE TABLE Users (
       id serial PRIMARY KEY,
       Username text UNIQUE NOT NULL,
       Password text NOT NULL,
       Mail_Adress text UNIQUE NOT NULL,
       Subscription_Date timestamptz,
       Birthdate timestamptz NOT NULL,
       Profile_Picture_Path text DEFAULT "",
       Rights smallint DEFAULT 1
);

CREATE TABLE Posts (
       id serial PRIMARY KEY,
       Content text NOT NULL,
       User_id bigint references Users(id),
       Post_date timestamptz,
       Likes_number int,
       Picture_path text
);

CREATE TABLE Posts_Likes (
       id serial PRIMARY KEY,
       User_id bigint references Users(id),
       Post_id bigint references Posts(id)
);

CREATE TABLE Posts_Pictures(
       id serial PRIMARY KEY,
       Post_id bigint references Posts(id),
       Picture_path text
);

CREATE TABLE Comments (
       id serial PRIMARY KEY,
       Content text NOT NULL,
       User_id bigint references Users(id),
       Post_date timestamptz,
       Likes_number int,
       Picture_path text,
       Post_id bigint references Posts(id)
);

CREATE TABLE Comments_Likes (
       id serial PRIMARY KEY,
       User_id bigint references Users(id),
       Comment_id bigint references Comments(id)
);

CREATE TABLE Comments_Pictures(
       id serial PRIMARY KEY,
       Comment_id bigint references Comments(id),
       Picture_path text
);

CREATE TABLE Followers(
       id serial PRIMARY KEY,
       Following_user_id bigint references Users(id),
       Followed_user_id bigint references Users(id)
);
