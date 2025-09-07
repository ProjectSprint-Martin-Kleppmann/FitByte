-- Drop foreign key constraint
ALTER TABLE activities DROP CONSTRAINT IF EXISTS fk_activities_user_id;

-- Drop indexes
DROP INDEX IF EXISTS idx_activities_calories_burned;
DROP INDEX IF EXISTS idx_activities_activity_type;
DROP INDEX IF EXISTS idx_activities_done_at;
DROP INDEX IF EXISTS idx_activities_activity_id;
DROP INDEX IF EXISTS idx_activities_user_id;

-- Drop the activities table
DROP TABLE IF EXISTS activities;
