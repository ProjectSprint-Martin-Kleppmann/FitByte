CREATE TABLE IF NOT EXISTS activities (
    id BIGSERIAL PRIMARY KEY,
    activity_id VARCHAR(255) UNIQUE NOT NULL,
    user_id BIGINT NOT NULL,
    activity_type VARCHAR(50) NOT NULL CHECK (activity_type IN ('Walking', 'Yoga', 'Stretching', 'Cycling', 'Swimming', 'Dancing', 'Hiking', 'Running', 'HIIT', 'JumpRope')),
    done_at TIMESTAMP WITH TIME ZONE NOT NULL,
    duration_in_minutes INTEGER NOT NULL CHECK (duration_in_minutes > 0),
    calories_burned INTEGER NOT NULL CHECK (calories_burned > 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_activities_user_id ON activities(user_id);
CREATE INDEX IF NOT EXISTS idx_activities_activity_id ON activities(activity_id);
CREATE INDEX IF NOT EXISTS idx_activities_done_at ON activities(done_at);
CREATE INDEX IF NOT EXISTS idx_activities_activity_type ON activities(activity_type);
CREATE INDEX IF NOT EXISTS idx_activities_calories_burned ON activities(calories_burned);

-- Add foreign key constraint to profiles table (assuming profiles table exists)
ALTER TABLE activities ADD CONSTRAINT fk_activities_user_id 
    FOREIGN KEY (user_id) REFERENCES profiles(id) ON DELETE CASCADE;
