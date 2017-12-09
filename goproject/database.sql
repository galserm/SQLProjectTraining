CREATE TABLE Users (
       id serial PRIMARY KEY,
       Username text UNIQUE NOT NULL,
       Password text NOT NULL,
       Mail_Adress text UNIQUE NOT NULL,
       Subscription_Date timestamptz DEFAULT CURRENT_DATE,
       Birthdate timestamptz NOT NULL,
       Profile_Picture_Path text DEFAULT "",
       Rights smallint DEFAULT 1
);

CREATE TABLE Posts (
       id serial PRIMARY KEY,
       Content text NOT NULL,
       User_id bigint REFERENCES Users(id) ON DELETE CASCADE,
       Post_date timestamptz DEFAULT CURRENT_DATE,
       Likes_number int DEFAULT 0,
       Picture_path text DEFAULT ""
);

CREATE TABLE Posts_Likes (
       id serial PRIMARY KEY,
       User_id bigint REFERENCES Users(id) ON DELETE CASCADE,
       Post_id bigint REFERENCES Posts(id) ON DELETE CASCADE
);

CREATE TABLE Posts_Pictures(
       id serial PRIMARY KEY,
       Post_id bigint REFERENCES Posts(id) ON DELETE CASCADE,
       Picture_path text DEFAULT ""
);

CREATE TABLE Comments (
       id serial PRIMARY KEY,
       Content text NOT NULL,
       User_id bigint REFERENCES Users(id) ON DELETE CASCADE,
       Comment_date timestamptz DEFAULT CURRENT_DATE,
       Likes_number int DEFAULT 0,
       Picture_path text DEFAULT "",
       Post_id bigint REFERENCES Posts(id) ON DELETE CASCADE
);

CREATE TABLE Comments_Likes (
       id serial PRIMARY KEY,
       User_id bigint REFERENCES Users(id) ON DELETE CASCADE,
       Comment_id bigint REFERENCES Comments(id) ON DELETE CASCADE
);

CREATE TABLE Comments_Pictures(
       id serial PRIMARY KEY,
       Comment_id bigint REFERENCES Comments(id) ON DELETE CASCADE,
       Picture_path text DEFAULT ""
);

CREATE TABLE Followers(
       id serial PRIMARY KEY,
       Following_user_id bigint REFERENCES Users(id) ON DELETE CASCADE, 
       Followed_user_id bigint REFERENCES Users(id) ON DELETE CASCADE
);
