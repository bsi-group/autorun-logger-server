--
-- PostgreSQL database dump
--

-- Dumped from database version 9.5.6
-- Dumped by pg_dump version 9.5.6

-- Started on 2017-04-14 14:59:14 BST

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- TOC entry 1 (class 3079 OID 12393)
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- TOC entry 2207 (class 0 OID 0)
-- Dependencies: 1
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

SET default_with_oids = false;

--
-- TOC entry 181 (class 1259 OID 21734)
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
-- TOC entry 182 (class 1259 OID 21740)
-- Name: alert_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE alert_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2208 (class 0 OID 0)
-- Dependencies: 182
-- Name: alert_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE alert_id_seq OWNED BY alert.id;


--
-- TOC entry 192 (class 1259 OID 21903)
-- Name: classification; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE classification (
    id bigint NOT NULL,
    alert_id bigint,
    user_name text,
    "timestamp" timestamp without time zone
);


--
-- TOC entry 191 (class 1259 OID 21901)
-- Name: classification_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE classification_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2209 (class 0 OID 0)
-- Dependencies: 191
-- Name: classification_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE classification_id_seq OWNED BY classification.id;


--
-- TOC entry 183 (class 1259 OID 21742)
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
-- TOC entry 184 (class 1259 OID 21748)
-- Name: current_autoruns_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE current_autoruns_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2210 (class 0 OID 0)
-- Dependencies: 184
-- Name: current_autoruns_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE current_autoruns_id_seq OWNED BY current_autoruns.id;


--
-- TOC entry 185 (class 1259 OID 21750)
-- Name: export; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE export (
    id bigint NOT NULL,
    data_type integer,
    file_name text,
    updated timestamp without time zone
);


--
-- TOC entry 186 (class 1259 OID 21756)
-- Name: export_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE export_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2211 (class 0 OID 0)
-- Dependencies: 186
-- Name: export_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE export_id_seq OWNED BY export.id;


--
-- TOC entry 187 (class 1259 OID 21758)
-- Name: instance; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE instance (
    id bigint NOT NULL,
    domain text,
    host text,
    "timestamp" timestamp without time zone
);


--
-- TOC entry 188 (class 1259 OID 21764)
-- Name: instance_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE instance_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2212 (class 0 OID 0)
-- Dependencies: 188
-- Name: instance_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE instance_id_seq OWNED BY instance.id;


--
-- TOC entry 189 (class 1259 OID 21766)
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
-- TOC entry 190 (class 1259 OID 21772)
-- Name: previous_autoruns_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE previous_autoruns_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2213 (class 0 OID 0)
-- Dependencies: 190
-- Name: previous_autoruns_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE previous_autoruns_id_seq OWNED BY previous_autoruns.id;


--
-- TOC entry 2053 (class 2604 OID 21774)
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY alert ALTER COLUMN id SET DEFAULT nextval('alert_id_seq'::regclass);


--
-- TOC entry 2058 (class 2604 OID 21906)
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY classification ALTER COLUMN id SET DEFAULT nextval('classification_id_seq'::regclass);


--
-- TOC entry 2054 (class 2604 OID 21775)
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY current_autoruns ALTER COLUMN id SET DEFAULT nextval('current_autoruns_id_seq'::regclass);


--
-- TOC entry 2055 (class 2604 OID 21776)
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY export ALTER COLUMN id SET DEFAULT nextval('export_id_seq'::regclass);


--
-- TOC entry 2056 (class 2604 OID 21777)
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY instance ALTER COLUMN id SET DEFAULT nextval('instance_id_seq'::regclass);


--
-- TOC entry 2057 (class 2604 OID 21778)
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY previous_autoruns ALTER COLUMN id SET DEFAULT nextval('previous_autoruns_id_seq'::regclass);


--
-- TOC entry 2060 (class 2606 OID 21780)
-- Name: alert_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY alert
    ADD CONSTRAINT alert_pk PRIMARY KEY (id);


--
-- TOC entry 2084 (class 2606 OID 21915)
-- Name: classification_alert_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY classification
    ADD CONSTRAINT classification_alert_id_key UNIQUE (alert_id);


--
-- TOC entry 2086 (class 2606 OID 21911)
-- Name: classification_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY classification
    ADD CONSTRAINT classification_pkey PRIMARY KEY (id);


--
-- TOC entry 2065 (class 2606 OID 21782)
-- Name: current_autoruns_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY current_autoruns
    ADD CONSTRAINT current_autoruns_pk PRIMARY KEY (id);


--
-- TOC entry 2075 (class 2606 OID 21784)
-- Name: instance_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY instance
    ADD CONSTRAINT instance_pk PRIMARY KEY (id);


--
-- TOC entry 2080 (class 2606 OID 21786)
-- Name: previous_autoruns_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY previous_autoruns
    ADD CONSTRAINT previous_autoruns_pk PRIMARY KEY (id);


--
-- TOC entry 2069 (class 2606 OID 21788)
-- Name: summary_file_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY export
    ADD CONSTRAINT summary_file_name_key UNIQUE (file_name);


--
-- TOC entry 2071 (class 2606 OID 21790)
-- Name: summary_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY export
    ADD CONSTRAINT summary_pkey PRIMARY KEY (id);


--
-- TOC entry 2061 (class 1259 OID 21834)
-- Name: alert_verified_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX alert_verified_idx ON alert USING btree (verified);


--
-- TOC entry 2082 (class 1259 OID 21913)
-- Name: classification_alert_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX classification_alert_id_idx ON classification USING btree (alert_id);


--
-- TOC entry 2062 (class 1259 OID 21830)
-- Name: current_autoruns_file_path_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX current_autoruns_file_path_idx ON current_autoruns USING btree (file_path);


--
-- TOC entry 2063 (class 1259 OID 21791)
-- Name: current_autoruns_instance_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX current_autoruns_instance_idx ON current_autoruns USING btree (instance);


--
-- TOC entry 2066 (class 1259 OID 21831)
-- Name: current_autoruns_sha256_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX current_autoruns_sha256_idx ON current_autoruns USING btree (sha256);


--
-- TOC entry 2067 (class 1259 OID 21792)
-- Name: export_data_type_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX export_data_type_idx ON export USING btree (data_type);


--
-- TOC entry 2072 (class 1259 OID 21793)
-- Name: instance_domain_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX instance_domain_idx ON instance USING btree (domain);


--
-- TOC entry 2073 (class 1259 OID 21794)
-- Name: instance_host_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX instance_host_idx ON instance USING btree (host);


--
-- TOC entry 2076 (class 1259 OID 21795)
-- Name: instance_timestamp_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX instance_timestamp_idx ON instance USING btree ("timestamp");


--
-- TOC entry 2077 (class 1259 OID 21833)
-- Name: previous_autoruns_file_path_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX previous_autoruns_file_path_idx ON previous_autoruns USING btree (file_path);


--
-- TOC entry 2078 (class 1259 OID 21796)
-- Name: previous_autoruns_instance_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX previous_autoruns_instance_idx ON previous_autoruns USING btree (instance);


--
-- TOC entry 2081 (class 1259 OID 21832)
-- Name: previous_autoruns_sha256_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX previous_autoruns_sha256_idx ON previous_autoruns USING btree (sha256);


-- Completed on 2017-04-14 14:59:14 BST

--
-- PostgreSQL database dump complete
--

