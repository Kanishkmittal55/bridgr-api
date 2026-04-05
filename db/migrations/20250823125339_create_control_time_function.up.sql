-- Create a reusable function to manage created_at and updated_at timestamps.
CREATE OR REPLACE FUNCTION hskip_users.tr_control_time()
RETURNS TRIGGER AS $$
BEGIN
    -- For INSERT operations
    IF TG_OP = 'INSERT' THEN
        IF NEW.created_at IS NULL THEN
            NEW.created_at := NOW();
        END IF;
        IF NEW.updated_at IS NULL THEN
            NEW.updated_at := NEW.created_at;
        END IF;

    -- For UPDATE operations
    ELSIF TG_OP = 'UPDATE' THEN
        -- Prevent created_at from ever being changed
        NEW.created_at = OLD.created_at;
        -- Automatically set updated_at to the current time
        NEW.updated_at = NOW();
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION hskip_users.tr_control_time() IS 'A generic trigger function to manage created_at and updated_at columns.';