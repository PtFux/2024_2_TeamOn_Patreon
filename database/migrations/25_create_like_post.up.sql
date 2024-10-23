CREATE TABLE like_post (
    like_post_id uuid PRIMARY KEY,
    post_id uuid REFERENCES post ON DELETE CASCADE NOT NULL,
    user_id uuid REFERENCES people ON DELETE CASCADE NOT NULL,
    posted_date timestamp NOT NULL DEFAULT now()
);