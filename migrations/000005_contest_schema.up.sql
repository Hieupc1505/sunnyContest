-- Tạo ký hóa contest_state
CREATE TYPE contest_state AS ENUM ('IDLE', 'RUNNING', 'FINISHED', 'WAITING');

CREATE TABLE sf_contest (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    subject_id BIGINT NOT NULL,
    num_question INT NOT NULL,
    time_exam INT NOT NULL,
    time_start_exam timestamptz,
    state contest_state NOT NULL,
    questions TEXT NOT NULL,
    created_time timestamptz  NOT NULL  DEFAULT now(),
    updated_time timestamptz  NOT NULL  DEFAULT now()
);

CREATE INDEX idx_sf_contest_user_id ON sf_contest(user_id);

ALTER TABLE sf_contest
    ADD CONSTRAINT sf_contest_user_id_fk
        FOREIGN KEY (user_id) REFERENCES sf_user(id) ON DELETE CASCADE;