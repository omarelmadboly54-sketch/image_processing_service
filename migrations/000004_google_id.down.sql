ALTER Table users
DROP COLUMN IF EXISTS google_id;

ALTER Table users
ALTER COLUMN password SET NOT NULL;