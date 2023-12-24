-- +goose Up
CREATE TABLE IF NOT EXISTS article_reaction
(
  `article_id` UUID,
  `rater_id` UUID,
  `reaction` String,
  `votes` int,
  `year` int default toYear(toDateTime(now(),'UTC')),
  `month` int default toMonth(toDateTime(now(),'UTC')),
  `day` int default toDayOfMonth(toDateTime(now(),'UTC')),
  `hour` int default toHour(toDateTime(now(),'UTC')),
  `minute` int default toMinute(toDateTime(now(),'UTC')),
  `second` int default toSecond(toDateTime(now(),'UTC'))
)
ENGINE = SummingMergeTree(votes)
PRIMARY KEY (article_id, rater_id, reaction)
ORDER BY (article_id, rater_id, reaction);

-- +goose Down
DROP TABLE article_reaction;
