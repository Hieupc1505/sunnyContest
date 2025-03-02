-- Xóa ràng buộc khóa ngoại trước
ALTER TABLE sf_profile DROP CONSTRAINT IF EXISTS sf_profile_user_id_fkey;
DROP TABLE IF EXISTS sf_profile;

-- db/migrations/000001_create_users_table.down.sql
DROP TABLE IF EXISTS sf_user;

