--
-- PostgreSQL database dump
--

-- Dumped from database version 9.5.5
-- Dumped by pg_dump version 9.5.5

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

SET default_with_oids = false;

--
-- Name: alert; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE alert (
    id bigint NOT NULL,
    instance integer,
    domain text,
    host text,
    "timestamp" timestamp without time zone,
    autorun_id bigint,
    location text,
    item_name text,
    enabled boolean,
    profile text,
    launch_string text,
    description text,
    company text,
    signer text,
    version_number text,
    file_path text,
    file_name text,
    file_directory text,
    "time" timestamp without time zone,
    sha256 text,
    md5 text,
    linked text,
    text text,
    verified smallint
);


--
-- Name: alert_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE alert_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: alert_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE alert_id_seq OWNED BY alert.id;


--
-- Name: classification; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE classification (
    id bigint NOT NULL,
    alert_id bigint,
    user_name text,
    "timestamp" timestamp without time zone
);


--
-- Name: classification_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE classification_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: classification_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE classification_id_seq OWNED BY classification.id;


--
-- Name: current_autoruns; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE current_autoruns (
    id bigint NOT NULL,
    instance integer,
    location text,
    item_name text,
    enabled boolean,
    profile text,
    launch_string text,
    description text,
    company text,
    signer text,
    version_number text,
    file_path text,
    file_name text,
    file_directory text,
    "time" timestamp without time zone,
    sha256 text,
    md5 text,
    verified smallint
);


--
-- Name: current_autoruns_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE current_autoruns_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: current_autoruns_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE current_autoruns_id_seq OWNED BY current_autoruns.id;


--
-- Name: export; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE export (
    id bigint NOT NULL,
    data_type integer,
    file_name text,
    updated timestamp without time zone
);


--
-- Name: export_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE export_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: export_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE export_id_seq OWNED BY export.id;


--
-- Name: instance; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE instance (
    id bigint NOT NULL,
    domain text,
    host text,
    "timestamp" timestamp without time zone
);


--
-- Name: instance_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE instance_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: instance_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE instance_id_seq OWNED BY instance.id;


--
-- Name: previous_autoruns; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE previous_autoruns (
    id bigint NOT NULL,
    instance integer,
    location text,
    item_name text,
    enabled boolean,
    profile text,
    launch_string text,
    description text,
    company text,
    signer text,
    version_number text,
    file_path text,
    file_name text,
    file_directory text,
    "time" timestamp without time zone,
    sha256 text,
    md5 text,
    verified smallint
);


--
-- Name: previous_autoruns_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE previous_autoruns_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: previous_autoruns_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE previous_autoruns_id_seq OWNED BY previous_autoruns.id;


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY alert ALTER COLUMN id SET DEFAULT nextval('alert_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY classification ALTER COLUMN id SET DEFAULT nextval('classification_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY current_autoruns ALTER COLUMN id SET DEFAULT nextval('current_autoruns_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY export ALTER COLUMN id SET DEFAULT nextval('export_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY instance ALTER COLUMN id SET DEFAULT nextval('instance_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY previous_autoruns ALTER COLUMN id SET DEFAULT nextval('previous_autoruns_id_seq'::regclass);


--
-- Name: alert_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY alert
    ADD CONSTRAINT alert_pk PRIMARY KEY (id);


--
-- Name: classification_alert_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY classification
    ADD CONSTRAINT classification_alert_id_key UNIQUE (alert_id);


--
-- Name: classification_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY classification
    ADD CONSTRAINT classification_pkey PRIMARY KEY (id);


--
-- Name: current_autoruns_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY current_autoruns
    ADD CONSTRAINT current_autoruns_pk PRIMARY KEY (id);


--
-- Name: instance_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY instance
    ADD CONSTRAINT instance_pk PRIMARY KEY (id);


--
-- Name: previous_autoruns_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY previous_autoruns
    ADD CONSTRAINT previous_autoruns_pk PRIMARY KEY (id);


--
-- Name: summary_file_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY export
    ADD CONSTRAINT summary_file_name_key UNIQUE (file_name);


--
-- Name: summary_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY export
    ADD CONSTRAINT summary_pkey PRIMARY KEY (id);


--
-- Name: alert_verified_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX alert_verified_idx ON alert USING btree (verified);


--
-- Name: classification_alert_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX classification_alert_id_idx ON classification USING btree (alert_id);


--
-- Name: current_autoruns_file_path_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX current_autoruns_file_path_idx ON current_autoruns USING btree (file_path);


--
-- Name: current_autoruns_instance_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX current_autoruns_instance_idx ON current_autoruns USING btree (instance);


--
-- Name: current_autoruns_sha256_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX current_autoruns_sha256_idx ON current_autoruns USING btree (sha256);


--
-- Name: export_data_type_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX export_data_type_idx ON export USING btree (data_type);


--
-- Name: instance_domain_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX instance_domain_idx ON instance USING btree (domain);


--
-- Name: instance_host_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX instance_host_idx ON instance USING btree (host);


--
-- Name: instance_timestamp_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX instance_timestamp_idx ON instance USING btree ("timestamp");


--
-- Name: previous_autoruns_file_path_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX previous_autoruns_file_path_idx ON previous_autoruns USING btree (file_path);


--
-- Name: previous_autoruns_instance_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX previous_autoruns_instance_idx ON previous_autoruns USING btree (instance);


--
-- Name: previous_autoruns_sha256_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX previous_autoruns_sha256_idx ON previous_autoruns USING btree (sha256);


--
-- PostgreSQL database dump complete
--

