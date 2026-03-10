CREATE TABLE IF NOT EXISTS locus_views (
    user_id UUID,
    locus_id UUID,
    viewed_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id, locus_id) -- This prevents duplicates automatically
);
