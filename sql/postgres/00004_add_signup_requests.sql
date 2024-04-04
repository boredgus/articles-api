-- +goose Up
-- +goose StatementBegin
ALTER TABLE articlesdb."user"
ALTER COLUMN username TYPE varchar(70);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TYPE articlesdb.user_data
ALTER ATTRIBUTE username TYPE varchar(70);
-- +goose StatementEnd 

-- +goose StatementBegin
DROP FUNCTION articlesdb.GetUserByUsername;
-- +goose StatementEnd 

-- +goose StatementBegin
CREATE FUNCTION articlesdb.GetUserByUsername (
  p_username varchar(70)
)
RETURNS SETOF articlesdb.user_data
LANGUAGE plpgsql
AS $$
BEGIN
  RETURN QUERY
    select u.o_id, u.username, u.pswd, r.label
    from articlesdb."user" as u
    inner join articlesdb."user_role" as r
    on u.role=r.id
    where u.username=p_username;
END $$
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE articlesdb."signup_requests" (
  email varchar(70) UNIQUE,
  pswd varchar(60),
  passcode varchar(60),
  attempted_at timestamp with time zone
);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE articlesdb."signup_requests" OWNER TO postgres;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE PROCEDURE articlesdb.InsertSignupRequest (
  p_username varchar(70),
  p_password varchar(60),
  p_passcode varchar(60)
)
LANGUAGE plpgsql
AS $$
DECLARE
  existed_email varchar(70);
BEGIN
	select req.email into existed_email 
	from articlesdb."signup_requests" as req 
	where email = p_username;

	if existed_email is null then
		insert into articlesdb."signup_requests" 
		values (p_username, p_password, p_passcode, now());
	else
    update articlesdb."signup_requests"
    set pswd = p_password, passcode = p_passcode, attempted_at = now()
    where email = existed_email;
	end if;
END $$
-- +goose StatementEnd

-- +goose StatementBegin
CREATE FUNCTION articlesdb.GetSignupRequest (
  p_email varchar(70)
)
RETURNS SETOF articlesdb."signup_requests"
LANGUAGE plpgsql
AS $$
BEGIN
  RETURN QUERY
    select *
    from articlesdb."signup_requests"
    where email = p_email;
END $$
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE articlesdb."user"
ALTER COLUMN username TYPE varchar(50);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TYPE articlesdb.user_data
ALTER ATTRIBUTE username TYPE varchar(50);
-- +goose StatementEnd 

-- +goose StatementBegin
DROP FUNCTION articlesdb.GetUserByUsername;
-- +goose StatementEnd 

-- +goose StatementBegin
CREATE FUNCTION articlesdb.GetUserByUsername (
  p_username varchar(50)
)
RETURNS SETOF articlesdb.user_data
LANGUAGE plpgsql
AS $$
BEGIN
  RETURN QUERY
    select u.o_id, u.username, u.pswd, r.label
    from articlesdb."user" as u
    inner join articlesdb."user_role" as r
    on u.role=r.id
    where u.username=p_username;
END $$
-- +goose StatementEnd

-- +goose StatementBegin
DROP PROCEDURE articlesdb.InsertSignupRequest;
-- +goose StatementEnd

-- +goose StatementBegin
DROP FUNCTION articlesdb.GetSignupRequest;
-- +goose StatementEnd 

-- +goose StatementBegin
DROP TABLE articlesdb."signup_requests";
-- +goose StatementEnd
