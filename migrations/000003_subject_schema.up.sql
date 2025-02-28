CREATE TABLE sf_subject (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL ,
    name VARCHAR NOT NULL,
    description TEXT,
    tags VARCHAR,
    state INT NOT NULL DEFAULT 1,
    created_time timestamptz  NOT NULL  DEFAULT now(),
    updated_time timestamptz  NOT NULL  DEFAULT now()
);

-- Tạo index cho user_id để tối ưu truy vấn
CREATE INDEX idx_sf_subject_user_id ON sf_subject(user_id);

-- Thêm khóa ngoại user_id tham chiếu đến bảng sf_user
ALTER TABLE sf_subject
    ADD CONSTRAINT sf_subject_user_id_fk FOREIGN KEY (user_id) REFERENCES sf_user (id)  ON DELETE CASCADE;