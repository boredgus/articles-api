-- +goose Up
--
-- PostgreSQL database schemadump
--

-- Dumped from database version 16.2 (Debian 16.2-1.pgdg110+2)
-- Dumped by pg_dump version 16.2 (Debian 16.2-1.pgdg110+2)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: articlesdb; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA articlesdb;


ALTER SCHEMA articlesdb OWNER TO postgres;

--
-- Name: on_update_current_timestamp_article(); Type: FUNCTION; Schema: articlesdb; Owner: postgres
--

-- +goose StatementBegin
CREATE FUNCTION articlesdb.on_update_current_timestamp_article() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
   NEW.updated_at = now();
   RETURN NEW;
END;
$$;
-- +goose StatementEnd


ALTER FUNCTION articlesdb.on_update_current_timestamp_article() OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: article; Type: TABLE; Schema: articlesdb; Owner: postgres
--

-- +goose StatementBegin
CREATE TABLE articlesdb.article (
    id serial,
    o_id uuid NOT NULL,
    user_id integer,
    theme varchar(200) NOT NULL,
    text varchar(500) DEFAULT ''::varchar,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone,
    status integer DEFAULT 0 NOT NULL
);
-- +goose StatementEnd


ALTER TABLE articlesdb.article OWNER TO postgres;

-- --
-- -- Name: article_id_seq; Type: SEQUENCE; Schema: articlesdb; Owner: postgres
-- --

-- CREATE SEQUENCE articlesdb.article_id_seq
--     AS integer
--     START WITH 1
--     INCREMENT BY 1
--     NO MINVALUE
--     NO MAXVALUE
--     CACHE 1;


-- ALTER SEQUENCE articlesdb.article_id_seq OWNER TO postgres;

--
-- Name: article_id_seq; Type: SEQUENCE OWNED BY; Schema: articlesdb; Owner: postgres
--

ALTER SEQUENCE articlesdb.article_id_seq OWNED BY articlesdb.article.id;


--
-- Name: article_tag; Type: TABLE; Schema: articlesdb; Owner: postgres
--

CREATE TABLE articlesdb.article_tag (
    article_id integer,
    tag_id integer,
    UNIQUE(article_id, tag_id)
);


ALTER TABLE articlesdb.article_tag OWNER TO postgres;

--
-- Name: tag; Type: TABLE; Schema: articlesdb; Owner: postgres
--

CREATE TABLE articlesdb.tag (
    id serial,
    label varchar(100) NOT NULL
);


ALTER TABLE articlesdb.tag OWNER TO postgres;

-- --
-- -- Name: tag_id_seq; Type: SEQUENCE; Schema: articlesdb; Owner: postgres
-- --

-- CREATE SEQUENCE articlesdb.tag_id_seq
--     AS integer
--     START WITH 1
--     INCREMENT BY 1
--     NO MINVALUE
--     NO MAXVALUE
--     CACHE 1;


-- ALTER SEQUENCE articlesdb.tag_id_seq OWNER TO postgres;

--
-- Name: tag_id_seq; Type: SEQUENCE OWNED BY; Schema: articlesdb; Owner: postgres
--

ALTER SEQUENCE articlesdb.tag_id_seq OWNED BY articlesdb.tag.id;


--
-- Name: user; Type: TABLE; Schema: articlesdb; Owner: postgres
--

-- +goose StatementBegin
CREATE TABLE articlesdb."user" (
    id serial,
    o_id uuid NOT NULL,
    username varchar(50) NOT NULL,
    pswd varchar(60) NOT NULL,
    role integer NOT NULL
);
-- +goose StatementEnd


ALTER TABLE articlesdb."user" OWNER TO postgres;

-- --
-- -- Name: user_id_seq; Type: SEQUENCE; Schema: articlesdb; Owner: postgres
-- --

-- CREATE SEQUENCE articlesdb.user_id_seq
--     AS integer
--     START WITH 1
--     INCREMENT BY 1
--     NO MINVALUE
--     NO MAXVALUE
--     CACHE 1;


-- ALTER SEQUENCE articlesdb.user_id_seq OWNER TO postgres;

--
-- Name: user_id_seq; Type: SEQUENCE OWNED BY; Schema: articlesdb; Owner: postgres
--

ALTER SEQUENCE articlesdb.user_id_seq OWNED BY articlesdb."user".id;


--
-- Name: user_role; Type: TABLE; Schema: articlesdb; Owner: postgres
--

CREATE TABLE articlesdb.user_role (
    id serial,
    label varchar(30) NOT NULL
);


ALTER TABLE articlesdb.user_role OWNER TO postgres;

--
-- Name: article id; Type: DEFAULT; Schema: articlesdb; Owner: postgres
--

ALTER TABLE ONLY articlesdb.article ALTER COLUMN id SET DEFAULT nextval('articlesdb.article_id_seq'::regclass);


--
-- Name: tag id; Type: DEFAULT; Schema: articlesdb; Owner: postgres
--

ALTER TABLE ONLY articlesdb.tag ALTER COLUMN id SET DEFAULT nextval('articlesdb.tag_id_seq'::regclass);


--
-- Name: user id; Type: DEFAULT; Schema: articlesdb; Owner: postgres
--

ALTER TABLE ONLY articlesdb."user" ALTER COLUMN id SET DEFAULT nextval('articlesdb.user_id_seq'::regclass);


--
-- Name: article idx_16473_primary; Type: CONSTRAINT; Schema: articlesdb; Owner: postgres
--

ALTER TABLE ONLY articlesdb.article
    ADD CONSTRAINT idx_16473_primary PRIMARY KEY (id);


--
-- Name: tag idx_16492_primary; Type: CONSTRAINT; Schema: articlesdb; Owner: postgres
--

ALTER TABLE ONLY articlesdb.tag
    ADD CONSTRAINT idx_16492_primary PRIMARY KEY (id);


--
-- Name: user idx_16497_primary; Type: CONSTRAINT; Schema: articlesdb; Owner: postgres
--

ALTER TABLE ONLY articlesdb."user"
    ADD CONSTRAINT idx_16497_primary PRIMARY KEY (id);


--
-- Name: user_role idx_16501_primary; Type: CONSTRAINT; Schema: articlesdb; Owner: postgres
--

ALTER TABLE ONLY articlesdb.user_role
    ADD CONSTRAINT idx_16501_primary PRIMARY KEY (id);


--
-- Name: idx_16473_o_id; Type: INDEX; Schema: articlesdb; Owner: postgres
--

CREATE UNIQUE INDEX idx_16473_o_id ON articlesdb.article USING btree (o_id);


--
-- Name: idx_16473_user_id; Type: INDEX; Schema: articlesdb; Owner: postgres
--

CREATE INDEX idx_16473_user_id ON articlesdb.article USING btree (user_id);


--
-- Name: idx_16482_article_id; Type: INDEX; Schema: articlesdb; Owner: postgres
--

CREATE INDEX idx_16482_article_id ON articlesdb.article_tag USING btree (article_id);


--
-- Name: idx_16482_tag_id; Type: INDEX; Schema: articlesdb; Owner: postgres
--

CREATE INDEX idx_16482_tag_id ON articlesdb.article_tag USING btree (tag_id);


--
-- Name: idx_16492_label; Type: INDEX; Schema: articlesdb; Owner: postgres
--

CREATE UNIQUE INDEX idx_16492_label ON articlesdb.tag USING btree (label);


--
-- Name: idx_16497_o_id; Type: INDEX; Schema: articlesdb; Owner: postgres
--

CREATE UNIQUE INDEX idx_16497_o_id ON articlesdb."user" USING btree (o_id);


--
-- Name: idx_16497_user_role_fk; Type: INDEX; Schema: articlesdb; Owner: postgres
--

CREATE INDEX idx_16497_user_role_fk ON articlesdb."user" USING btree (role);


--
-- Name: idx_16497_username; Type: INDEX; Schema: articlesdb; Owner: postgres
--

CREATE UNIQUE INDEX idx_16497_username ON articlesdb."user" USING btree (username);


--
-- Name: idx_16501_label; Type: INDEX; Schema: articlesdb; Owner: postgres
--

CREATE UNIQUE INDEX idx_16501_label ON articlesdb.user_role USING btree (label);


--
-- Name: article on_update_current_timestamp; Type: TRIGGER; Schema: articlesdb; Owner: postgres
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON articlesdb.article FOR EACH ROW EXECUTE FUNCTION articlesdb.on_update_current_timestamp_article();


--
-- Name: article article_ibfk_1; Type: FK CONSTRAINT; Schema: articlesdb; Owner: postgres
--

ALTER TABLE ONLY articlesdb.article
    ADD CONSTRAINT article_ibfk_1 FOREIGN KEY (user_id) REFERENCES articlesdb."user"(id) ON DELETE CASCADE;


--
-- Name: article_tag article_tag_ibfk_1; Type: FK CONSTRAINT; Schema: articlesdb; Owner: postgres
--

ALTER TABLE ONLY articlesdb.article_tag
    ADD CONSTRAINT article_tag_ibfk_1 FOREIGN KEY (article_id) REFERENCES articlesdb.article(id) ON DELETE CASCADE;


--
-- Name: article_tag article_tag_ibfk_2; Type: FK CONSTRAINT; Schema: articlesdb; Owner: postgres
--

ALTER TABLE ONLY articlesdb.article_tag
    ADD CONSTRAINT article_tag_ibfk_2 FOREIGN KEY (tag_id) REFERENCES articlesdb.tag(id) ON DELETE CASCADE;


--
-- Name: user user_role_fk; Type: FK CONSTRAINT; Schema: articlesdb; Owner: postgres
--

ALTER TABLE ONLY articlesdb."user"
    ADD CONSTRAINT user_role_fk FOREIGN KEY (role) REFERENCES articlesdb.user_role(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA articlesdb CASCADE;
-- +goose StatementEnd
