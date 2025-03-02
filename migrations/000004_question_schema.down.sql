DROP TYPE IF EXISTS level_question;

ALTER TABLE sf_question DROP CONSTRAINT IF EXISTS sf_question_subject_id_fk;
ALTER TABLE sf_question DROP CONSTRAINT IF EXISTS sf_question_user_id_fk;

DROP TABLE IF EXISTS sf_question;