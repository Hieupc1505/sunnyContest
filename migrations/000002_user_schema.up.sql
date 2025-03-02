CREATE TABLE IF NOT EXISTS sf_user (
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

-- Tạo bảng sf_profile
CREATE TABLE sf_profile (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    nickname VARCHAR NOT NULL,
    avatar VARCHAR(500) NOT NULL,
    created_time timestamptz  NOT NULL  DEFAULT now(),
    updated_time timestamptz  NOT NULL  DEFAULT now()
);


-- Tạo index trên user_id để tối ưu truy vấn
CREATE INDEX idx_sf_profile_user_id ON sf_profile(user_id);

-- Thêm khóa ngoại user_id tham chiếu đến bảng sf_user
ALTER TABLE sf_profile
    ADD CONSTRAINT sf_profile_user_id_fk FOREIGN KEY (user_id) REFERENCES sf_user (id)  ON DELETE CASCADE;
