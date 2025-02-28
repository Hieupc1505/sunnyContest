CREATE TABLE sf_user (
     id BIGSERIAL PRIMARY KEY,
     username VARCHAR(255) UNIQUE NOT NULL,
     password TEXT NOT NULL,
     role INT NOT NULL DEFAULT 2,  -- 0: admin, 2: user, 4: teacher
     status INT NOT NULL DEFAULT 0, --  0 active, 1 locked, 2 disabled, 3 deleted
     token TEXT,
     token_expired timestamptz,
     created_time timestamptz  NOT NULL  DEFAULT now(),
     updated_time timestamptz  NOT NULL  DEFAULT now()
);