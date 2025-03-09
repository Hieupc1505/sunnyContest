ALTER TABLE sf_contest DROP CONSTRAINT IF EXISTS sf_contest_user_id_fk;

-- Lệnh xóa bảng nếu tồn tại
DROP TABLE IF EXISTS sf_contest;

DROP TYPE IF EXISTS contest_state;