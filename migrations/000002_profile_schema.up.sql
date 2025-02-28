-- Tạo bảng sf_profile
CREATE TABLE sf_profile (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    nickname VARCHAR NOT NULL,
    avatar VARCHAR(500) NOT NULL,
    created_time timestamptz  NOT NULL  DEFAULT now(),
    updated_time timestamptz  NOT NULL  DEFAULT now()
);

ALTER TABLE "sf_profile" ADD FOREIGN KEY ("user_id") REFERENCES "sf_user" ("id") ON DELETE CASCADE;

-- Tạo index trên user_id để tối ưu truy vấn
CREATE INDEX idx_sf_profile_user_id ON sf_profile(user_id);
