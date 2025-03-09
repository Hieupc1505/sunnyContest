ALTER TABLE sf_user_contest DROP CONSTRAINT IF EXISTS sf_user_contest_user_id_fk;
ALTER TABLE sf_user_contest DROP CONSTRAINT IF EXISTS sf_user_contest_contest_id_fk;

-- Lệnh xóa bảng nếu tồn tại
DROP TABLE IF EXISTS sf_user_contest;