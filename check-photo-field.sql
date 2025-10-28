-- Check if photo_url column exists, if not add it
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'users' AND column_name = 'photo_url'
    ) THEN
        ALTER TABLE users ADD COLUMN photo_url VARCHAR(255);
        RAISE NOTICE 'Added photo_url column to users table';
    ELSE
        RAISE NOTICE 'photo_url column already exists in users table';
    END IF;
END $$;
