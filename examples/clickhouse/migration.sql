CREATE TABLE IF NOT EXISTS demo.news
(
    `title`     String,
    `content`   String,
    `posted_at` DATETIME
)
    ENGINE = MergeTree() ORDER BY `posted_at`;