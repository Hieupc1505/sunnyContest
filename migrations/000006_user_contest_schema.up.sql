CREATE TABLE sf_user_contest (
     id BIGINT PRIMARY KEY NOT NULL GENERATED ALWAYS AS IDENTITY,
     contest_id BIGINT NOT NULL,
     user_id BIGINT NOT NULL,
     questions JSONB NOT NULL,
     exam JSONB,
     result JSONB,
     created_time timestamptz  NOT NULL  DEFAULT now(),
     updated_time timestamptz  NOT NULL  DEFAULT now()
);

-- Tạo chỉ mục cho trường contest_id
CREATE INDEX idx_sf_user_contest_contest_id ON sf_user_contest(contest_id);

ALTER TABLE sf_user_contest
    ADD CONSTRAINT sf_user_contest_user_id_fk
        FOREIGN KEY (user_id) REFERENCES sf_user(id) ON DELETE CASCADE;

ALTER TABLE sf_user_contest
    ADD CONSTRAINT sf_user_contest_contest_id_fk
        FOREIGN KEY (contest_id) REFERENCES sf_contest(id) ON DELETE CASCADE;

