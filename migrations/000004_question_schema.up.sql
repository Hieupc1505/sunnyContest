-- Tạo ký hóa contest_state
CREATE TYPE level_question AS ENUM ('EASY', 'MEDIUM', 'HARD');

CREATE TABLE sf_question (
     id BIGSERIAL PRIMARY KEY,
     subject_id BIGINT NOT NULL,
     user_id BIGINT NOT NULL,
     level level_question NOT NULL,
     question TEXT NOT NULL,
     question_type VARCHAR NOT NULL,
     question_image VARCHAR,
     answers TEXT NOT NULL,
     answer_type VARCHAR NOT NULL ,
     state INT NOT NULL DEFAULT 1,
     created_time timestamptz  NOT NULL  DEFAULT now(),
     updated_time timestamptz  NOT NULL  DEFAULT now()
);

-- Tạo index để tối ưu truy vấn
CREATE INDEX idx_sf_question_subject_id ON sf_question(subject_id);
CREATE INDEX idx_sf_question_user_id ON sf_question(user_id);

ALTER TABLE sf_question
    ADD CONSTRAINT sf_question_subject_id_fk
        FOREIGN KEY (subject_id) REFERENCES sf_subject(id) ON DELETE CASCADE;

ALTER TABLE sf_question
    ADD CONSTRAINT sf_question_user_id_fk
        FOREIGN KEY (user_id) REFERENCES sf_user(id) ON DELETE CASCADE;


