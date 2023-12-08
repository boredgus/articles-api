-- +goose Up
GRANT ALL PRIVILEGES ON articlesdb.* TO 'root'@'%';
FLUSH PRIVILEGES;

-- +goose Down
REVOKE ALL PRIVILEGES ON articlesdb.* FROM 'root'@'%';