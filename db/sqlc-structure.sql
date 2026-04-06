--
-- PostgreSQL database dump
--


-- Dumped from database version 16.13 (Debian 16.13-1.pgdg13+1)
-- Dumped by pg_dump version 16.13 (Debian 16.13-1.pgdg13+1)

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
-- Name: bridgr; Type: SCHEMA; Schema: -; Owner: bridgr
--

CREATE SCHEMA bridgr;


ALTER SCHEMA bridgr OWNER TO bridgr;

--
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


--
-- Name: EXTENSION pgcrypto; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION pgcrypto IS 'cryptographic functions';


--
-- Name: tr_control_time(); Type: FUNCTION; Schema: bridgr; Owner: bridgr
--

CREATE FUNCTION bridgr.tr_control_time() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
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
$$;


ALTER FUNCTION bridgr.tr_control_time() OWNER TO bridgr;

--
-- Name: FUNCTION tr_control_time(); Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON FUNCTION bridgr.tr_control_time() IS 'A generic trigger function to manage created_at and updated_at columns.';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: analysis_job_link; Type: TABLE; Schema: bridgr; Owner: bridgr
--

CREATE TABLE bridgr.analysis_job_link (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    id bigint NOT NULL,
    user_id integer NOT NULL,
    analysis_uuid uuid NOT NULL,
    job_candidate_uuid uuid NOT NULL,
    link_kind character varying(32) DEFAULT 'from_job_feed'::character varying NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    CONSTRAINT chk_analysis_job_link_kind CHECK (((link_kind)::text = ANY ((ARRAY['from_job_feed'::character varying, 'manual'::character varying, 'import'::character varying])::text[])))
);


ALTER TABLE bridgr.analysis_job_link OWNER TO bridgr;

--
-- Name: TABLE analysis_job_link; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON TABLE bridgr.analysis_job_link IS 'Bridgr: associates a skill_gap_analysis with a discovered job_candidates row.';


--
-- Name: COLUMN analysis_job_link.user_id; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.analysis_job_link.user_id IS 'App-enforced FK to users.id; denormalized for listing';


--
-- Name: COLUMN analysis_job_link.analysis_uuid; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.analysis_job_link.analysis_uuid IS 'App-enforced FK to skill_gap_analyses.uuid';


--
-- Name: COLUMN analysis_job_link.job_candidate_uuid; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.analysis_job_link.job_candidate_uuid IS 'App-enforced FK to job_candidates.uuid';


--
-- Name: COLUMN analysis_job_link.link_kind; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.analysis_job_link.link_kind IS 'from_job_feed, manual, import';


--
-- Name: analysis_job_link_id_seq; Type: SEQUENCE; Schema: bridgr; Owner: bridgr
--

CREATE SEQUENCE bridgr.analysis_job_link_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE bridgr.analysis_job_link_id_seq OWNER TO bridgr;

--
-- Name: analysis_job_link_id_seq; Type: SEQUENCE OWNED BY; Schema: bridgr; Owner: bridgr
--

ALTER SEQUENCE bridgr.analysis_job_link_id_seq OWNED BY bridgr.analysis_job_link.id;


--
-- Name: feed_items; Type: TABLE; Schema: bridgr; Owner: bridgr
--

CREATE TABLE bridgr.feed_items (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    id bigint NOT NULL,
    user_id integer NOT NULL,
    job_candidate_uuid uuid NOT NULL,
    score_uuid uuid,
    verification_uuid uuid,
    composite_score real DEFAULT 0 NOT NULL,
    gap_severity character varying(32),
    title text,
    company text,
    location text,
    job_url text,
    match_summary text,
    gap_summary text,
    feed_status character varying(30) DEFAULT 'new'::character varying NOT NULL,
    surfaced_at timestamp without time zone DEFAULT now() NOT NULL,
    seen_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE bridgr.feed_items OWNER TO bridgr;

--
-- Name: TABLE feed_items; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON TABLE bridgr.feed_items IS 'Bridgr feed: user-facing rows for high-quality job matches (denormalized for fast reads).';


--
-- Name: COLUMN feed_items.user_id; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.feed_items.user_id IS 'App-enforced FK to users.id';


--
-- Name: COLUMN feed_items.job_candidate_uuid; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.feed_items.job_candidate_uuid IS 'App-enforced FK to job_candidates.uuid';


--
-- Name: COLUMN feed_items.score_uuid; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.feed_items.score_uuid IS 'App-enforced FK to job_scores.uuid';


--
-- Name: COLUMN feed_items.verification_uuid; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.feed_items.verification_uuid IS 'Optional future link to a verification record; app-enforced when set';


--
-- Name: COLUMN feed_items.match_summary; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.feed_items.match_summary IS 'Short human-readable why this job matches';


--
-- Name: COLUMN feed_items.gap_summary; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.feed_items.gap_summary IS 'Short human-readable skill-gap note';


--
-- Name: COLUMN feed_items.feed_status; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.feed_items.feed_status IS 'new, seen, saved, dismissed, applied';


--
-- Name: COLUMN feed_items.surfaced_at; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.feed_items.surfaced_at IS 'When this item entered the feed';


--
-- Name: feed_items_id_seq; Type: SEQUENCE; Schema: bridgr; Owner: bridgr
--

CREATE SEQUENCE bridgr.feed_items_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE bridgr.feed_items_id_seq OWNER TO bridgr;

--
-- Name: feed_items_id_seq; Type: SEQUENCE OWNED BY; Schema: bridgr; Owner: bridgr
--

ALTER SEQUENCE bridgr.feed_items_id_seq OWNED BY bridgr.feed_items.id;


--
-- Name: job_candidates; Type: TABLE; Schema: bridgr; Owner: bridgr
--

CREATE TABLE bridgr.job_candidates (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    id bigint NOT NULL,
    user_id integer NOT NULL,
    discovery_run_uuid uuid,
    source_board character varying(80) DEFAULT ''::character varying NOT NULL,
    source_job_id text,
    job_url text NOT NULL,
    url_hash character varying(128) NOT NULL,
    content_hash character varying(128),
    title text,
    company text,
    location text,
    jd_text text,
    jd_s3_uri text,
    fetched_at timestamp without time zone,
    ingestion_status character varying(30) DEFAULT 'pending'::character varying NOT NULL,
    radar_payload jsonb DEFAULT '{}'::jsonb NOT NULL,
    application_url text,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    CONSTRAINT chk_job_candidates_ingestion_status CHECK (((ingestion_status)::text = ANY ((ARRAY['pending'::character varying, 'fetched'::character varying, 'fetch_failed'::character varying, 'partial'::character varying])::text[])))
);


ALTER TABLE bridgr.job_candidates OWNER TO bridgr;

--
-- Name: TABLE job_candidates; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON TABLE bridgr.job_candidates IS 'Bridgr job discovery: one row per unique posting per user (dedupe by url_hash).';


--
-- Name: COLUMN job_candidates.user_id; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_candidates.user_id IS 'App-enforced FK to users.id';


--
-- Name: COLUMN job_candidates.discovery_run_uuid; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_candidates.discovery_run_uuid IS 'App-enforced FK to job_search_discovery_runs.uuid';


--
-- Name: COLUMN job_candidates.url_hash; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_candidates.url_hash IS 'Stable hash of normalized URL for deduplication';


--
-- Name: COLUMN job_candidates.content_hash; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_candidates.content_hash IS 'Optional hash of JD body for refresh/dedup';


--
-- Name: COLUMN job_candidates.jd_text; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_candidates.jd_text IS 'Inline job description text when small enough';


--
-- Name: COLUMN job_candidates.jd_s3_uri; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_candidates.jd_s3_uri IS 'Pointer to stored JD blob when large';


--
-- Name: COLUMN job_candidates.radar_payload; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_candidates.radar_payload IS 'Raw or normalized Radar response slice';


--
-- Name: COLUMN job_candidates.application_url; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_candidates.application_url IS 'Direct apply URL when distinct from job_url';


--
-- Name: job_candidates_id_seq; Type: SEQUENCE; Schema: bridgr; Owner: bridgr
--

CREATE SEQUENCE bridgr.job_candidates_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE bridgr.job_candidates_id_seq OWNER TO bridgr;

--
-- Name: job_candidates_id_seq; Type: SEQUENCE OWNED BY; Schema: bridgr; Owner: bridgr
--

ALTER SEQUENCE bridgr.job_candidates_id_seq OWNED BY bridgr.job_candidates.id;


--
-- Name: job_enrichments; Type: TABLE; Schema: bridgr; Owner: bridgr
--

CREATE TABLE bridgr.job_enrichments (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    id bigint NOT NULL,
    user_id integer NOT NULL,
    job_candidate_uuid uuid NOT NULL,
    status character varying(30) DEFAULT 'pending'::character varying NOT NULL,
    required_skills jsonb,
    experience_range jsonb,
    salary_range jsonb,
    remote_policy character varying(64),
    visa_sponsorship boolean,
    company_size character varying(64),
    industry character varying(128),
    application_deadline timestamp without time zone,
    structured_jd jsonb,
    llm_model character varying(128),
    prompt_version character varying(64),
    error_detail text,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE bridgr.job_enrichments OWNER TO bridgr;

--
-- Name: TABLE job_enrichments; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON TABLE bridgr.job_enrichments IS 'Bridgr feed: LLM- or rules-based structured fields extracted from a job posting for matching.';


--
-- Name: COLUMN job_enrichments.user_id; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_enrichments.user_id IS 'App-enforced FK to users.id';


--
-- Name: COLUMN job_enrichments.job_candidate_uuid; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_enrichments.job_candidate_uuid IS 'App-enforced FK to job_candidates.uuid';


--
-- Name: COLUMN job_enrichments.required_skills; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_enrichments.required_skills IS 'JSON array of {skill, level, mandatory} or similar';


--
-- Name: COLUMN job_enrichments.experience_range; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_enrichments.experience_range IS 'JSON {min_years, max_years, seniority}';


--
-- Name: COLUMN job_enrichments.salary_range; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_enrichments.salary_range IS 'JSON {min, max, currency, period}';


--
-- Name: COLUMN job_enrichments.structured_jd; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_enrichments.structured_jd IS 'Full normalized job description blob for scoring';


--
-- Name: job_enrichments_id_seq; Type: SEQUENCE; Schema: bridgr; Owner: bridgr
--

CREATE SEQUENCE bridgr.job_enrichments_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE bridgr.job_enrichments_id_seq OWNER TO bridgr;

--
-- Name: job_enrichments_id_seq; Type: SEQUENCE OWNED BY; Schema: bridgr; Owner: bridgr
--

ALTER SEQUENCE bridgr.job_enrichments_id_seq OWNED BY bridgr.job_enrichments.id;


--
-- Name: job_harvest_schedules; Type: TABLE; Schema: bridgr; Owner: bridgr
--

CREATE TABLE bridgr.job_harvest_schedules (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    id bigint NOT NULL,
    user_id integer NOT NULL,
    profile_uuid uuid NOT NULL,
    enabled boolean DEFAULT true NOT NULL,
    cadence_minutes integer DEFAULT 360 NOT NULL,
    boards_rotation jsonb DEFAULT '[]'::jsonb NOT NULL,
    last_run_at timestamp without time zone,
    next_run_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE bridgr.job_harvest_schedules OWNER TO bridgr;

--
-- Name: TABLE job_harvest_schedules; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON TABLE bridgr.job_harvest_schedules IS 'Bridgr feed: when and how to re-run discovery for a job_search_profile.';


--
-- Name: COLUMN job_harvest_schedules.user_id; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_harvest_schedules.user_id IS 'App-enforced FK to users.id';


--
-- Name: COLUMN job_harvest_schedules.profile_uuid; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_harvest_schedules.profile_uuid IS 'App-enforced FK to job_search_profiles.uuid';


--
-- Name: COLUMN job_harvest_schedules.cadence_minutes; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_harvest_schedules.cadence_minutes IS 'Minimum interval between automated harvest runs';


--
-- Name: COLUMN job_harvest_schedules.boards_rotation; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_harvest_schedules.boards_rotation IS 'JSON array of board ids to rotate through on successive runs';


--
-- Name: COLUMN job_harvest_schedules.last_run_at; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_harvest_schedules.last_run_at IS 'When the last scheduled harvest started or completed';


--
-- Name: COLUMN job_harvest_schedules.next_run_at; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_harvest_schedules.next_run_at IS 'When the scheduler should enqueue the next harvest';


--
-- Name: job_harvest_schedules_id_seq; Type: SEQUENCE; Schema: bridgr; Owner: bridgr
--

CREATE SEQUENCE bridgr.job_harvest_schedules_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE bridgr.job_harvest_schedules_id_seq OWNER TO bridgr;

--
-- Name: job_harvest_schedules_id_seq; Type: SEQUENCE OWNED BY; Schema: bridgr; Owner: bridgr
--

ALTER SEQUENCE bridgr.job_harvest_schedules_id_seq OWNED BY bridgr.job_harvest_schedules.id;


--
-- Name: job_notifications; Type: TABLE; Schema: bridgr; Owner: bridgr
--

CREATE TABLE bridgr.job_notifications (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    id bigint NOT NULL,
    user_id integer NOT NULL,
    job_candidate_uuid uuid NOT NULL,
    channel character varying(32) DEFAULT 'in_app'::character varying NOT NULL,
    status character varying(30) DEFAULT 'pending'::character varying NOT NULL,
    payload jsonb DEFAULT '{}'::jsonb NOT NULL,
    sent_at timestamp without time zone,
    seen_at timestamp without time zone,
    error_detail text,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    CONSTRAINT chk_job_notifications_channel CHECK (((channel)::text = ANY ((ARRAY['in_app'::character varying, 'email'::character varying, 'push'::character varying])::text[]))),
    CONSTRAINT chk_job_notifications_status CHECK (((status)::text = ANY ((ARRAY['pending'::character varying, 'sent'::character varying, 'failed'::character varying, 'seen'::character varying, 'skipped'::character varying])::text[])))
);


ALTER TABLE bridgr.job_notifications OWNER TO bridgr;

--
-- Name: TABLE job_notifications; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON TABLE bridgr.job_notifications IS 'Bridgr job discovery: notification rows for new surfaced or relevant jobs.';


--
-- Name: COLUMN job_notifications.user_id; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_notifications.user_id IS 'App-enforced FK to users.id';


--
-- Name: COLUMN job_notifications.job_candidate_uuid; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_notifications.job_candidate_uuid IS 'App-enforced FK to job_candidates.uuid';


--
-- Name: COLUMN job_notifications.channel; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_notifications.channel IS 'in_app, email, push';


--
-- Name: COLUMN job_notifications.status; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_notifications.status IS 'pending, sent, failed, seen, skipped';


--
-- Name: COLUMN job_notifications.payload; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_notifications.payload IS 'Title snippet, deep link hints, template vars';


--
-- Name: job_notifications_id_seq; Type: SEQUENCE; Schema: bridgr; Owner: bridgr
--

CREATE SEQUENCE bridgr.job_notifications_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE bridgr.job_notifications_id_seq OWNER TO bridgr;

--
-- Name: job_notifications_id_seq; Type: SEQUENCE OWNED BY; Schema: bridgr; Owner: bridgr
--

ALTER SEQUENCE bridgr.job_notifications_id_seq OWNED BY bridgr.job_notifications.id;


--
-- Name: job_scores; Type: TABLE; Schema: bridgr; Owner: bridgr
--

CREATE TABLE bridgr.job_scores (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    id bigint NOT NULL,
    user_id integer NOT NULL,
    job_candidate_uuid uuid NOT NULL,
    enrichment_uuid uuid,
    skill_match_score real DEFAULT 0 NOT NULL,
    experience_match_score real DEFAULT 0 NOT NULL,
    location_match_score real DEFAULT 0 NOT NULL,
    recency_score real DEFAULT 0 NOT NULL,
    board_quality_score real DEFAULT 0 NOT NULL,
    composite_score real DEFAULT 0 NOT NULL,
    matched_skills jsonb,
    gap_skills jsonb,
    gap_severity character varying(32),
    scoring_model character varying(128),
    scoring_version character varying(64),
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE bridgr.job_scores OWNER TO bridgr;

--
-- Name: TABLE job_scores; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON TABLE bridgr.job_scores IS 'Bridgr feed: per-user relevance and skill-gap scores for a discovered job.';


--
-- Name: COLUMN job_scores.user_id; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_scores.user_id IS 'App-enforced FK to users.id';


--
-- Name: COLUMN job_scores.job_candidate_uuid; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_scores.job_candidate_uuid IS 'App-enforced FK to job_candidates.uuid';


--
-- Name: COLUMN job_scores.enrichment_uuid; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_scores.enrichment_uuid IS 'Optional app-enforced FK to job_enrichments.uuid used for this score';


--
-- Name: COLUMN job_scores.composite_score; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_scores.composite_score IS 'Weighted blend of dimension scores; primary feed sort key';


--
-- Name: COLUMN job_scores.matched_skills; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_scores.matched_skills IS 'JSON: skills the user satisfies for this role';


--
-- Name: COLUMN job_scores.gap_skills; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_scores.gap_skills IS 'JSON: required skills still missing for the user';


--
-- Name: COLUMN job_scores.gap_severity; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_scores.gap_severity IS 'none, minor, moderate, major';


--
-- Name: job_scores_id_seq; Type: SEQUENCE; Schema: bridgr; Owner: bridgr
--

CREATE SEQUENCE bridgr.job_scores_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE bridgr.job_scores_id_seq OWNER TO bridgr;

--
-- Name: job_scores_id_seq; Type: SEQUENCE OWNED BY; Schema: bridgr; Owner: bridgr
--

ALTER SEQUENCE bridgr.job_scores_id_seq OWNED BY bridgr.job_scores.id;


--
-- Name: job_search_discovery_runs; Type: TABLE; Schema: bridgr; Owner: bridgr
--

CREATE TABLE bridgr.job_search_discovery_runs (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    id bigint NOT NULL,
    user_id integer NOT NULL,
    status character varying(30) DEFAULT 'pending'::character varying NOT NULL,
    request_params jsonb DEFAULT '{}'::jsonb NOT NULL,
    radar_meta jsonb DEFAULT '{}'::jsonb NOT NULL,
    raw_candidate_count integer DEFAULT 0 NOT NULL,
    new_candidate_count integer DEFAULT 0 NOT NULL,
    started_at timestamp without time zone,
    completed_at timestamp without time zone,
    error_code character varying(80),
    error_detail text,
    sqs_message_id character varying(200),
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    CONSTRAINT chk_job_search_discovery_runs_status CHECK (((status)::text = ANY ((ARRAY['pending'::character varying, 'queued'::character varying, 'running'::character varying, 'completed'::character varying, 'failed'::character varying, 'cancelled'::character varying])::text[])))
);


ALTER TABLE bridgr.job_search_discovery_runs OWNER TO bridgr;

--
-- Name: TABLE job_search_discovery_runs; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON TABLE bridgr.job_search_discovery_runs IS 'Bridgr job discovery: audit trail for one Radar-backed discovery execution.';


--
-- Name: COLUMN job_search_discovery_runs.user_id; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_search_discovery_runs.user_id IS 'App-enforced FK to users.id';


--
-- Name: COLUMN job_search_discovery_runs.status; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_search_discovery_runs.status IS 'pending, queued, running, completed, failed, cancelled';


--
-- Name: COLUMN job_search_discovery_runs.request_params; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_search_discovery_runs.request_params IS 'Resolved query/location/boards caps sent to worker/Radar';


--
-- Name: COLUMN job_search_discovery_runs.radar_meta; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_search_discovery_runs.radar_meta IS 'Timings, per-board errors, debug payload';


--
-- Name: COLUMN job_search_discovery_runs.sqs_message_id; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_search_discovery_runs.sqs_message_id IS 'Optional queue id for async job discovery';


--
-- Name: job_search_discovery_runs_id_seq; Type: SEQUENCE; Schema: bridgr; Owner: bridgr
--

CREATE SEQUENCE bridgr.job_search_discovery_runs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE bridgr.job_search_discovery_runs_id_seq OWNER TO bridgr;

--
-- Name: job_search_discovery_runs_id_seq; Type: SEQUENCE OWNED BY; Schema: bridgr; Owner: bridgr
--

ALTER SEQUENCE bridgr.job_search_discovery_runs_id_seq OWNED BY bridgr.job_search_discovery_runs.id;


--
-- Name: job_search_profiles; Type: TABLE; Schema: bridgr; Owner: bridgr
--

CREATE TABLE bridgr.job_search_profiles (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    id bigint NOT NULL,
    user_id integer NOT NULL,
    target_roles jsonb DEFAULT '[]'::jsonb NOT NULL,
    locations jsonb DEFAULT '[]'::jsonb NOT NULL,
    boards_enabled jsonb DEFAULT '[]'::jsonb NOT NULL,
    matching jsonb DEFAULT '{}'::jsonb NOT NULL,
    canonical_cv_analysis_uuid uuid,
    max_surfaced_jobs integer DEFAULT 3 NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    CONSTRAINT chk_job_search_profiles_max_surfaced CHECK (((max_surfaced_jobs > 0) AND (max_surfaced_jobs <= 20)))
);


ALTER TABLE bridgr.job_search_profiles OWNER TO bridgr;

--
-- Name: TABLE job_search_profiles; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON TABLE bridgr.job_search_profiles IS 'Bridgr job discovery: user preferences for Radar search and surfacing.';


--
-- Name: COLUMN job_search_profiles.user_id; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_search_profiles.user_id IS 'App-enforced FK to users.id';


--
-- Name: COLUMN job_search_profiles.target_roles; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_search_profiles.target_roles IS 'JSON array of role strings or objects';


--
-- Name: COLUMN job_search_profiles.locations; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_search_profiles.locations IS 'JSON array of location descriptors (city, remote flags, etc.)';


--
-- Name: COLUMN job_search_profiles.boards_enabled; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_search_profiles.boards_enabled IS 'JSON array of board ids (e.g. linkedin, indeed)';


--
-- Name: COLUMN job_search_profiles.matching; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_search_profiles.matching IS 'Policy blob: strict_no_skill_gaps, tiers, etc.';


--
-- Name: COLUMN job_search_profiles.canonical_cv_analysis_uuid; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.job_search_profiles.canonical_cv_analysis_uuid IS 'Optional app-enforced FK to skill_gap_analyses.uuid for CV fingerprint context';


--
-- Name: job_search_profiles_id_seq; Type: SEQUENCE; Schema: bridgr; Owner: bridgr
--

CREATE SEQUENCE bridgr.job_search_profiles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE bridgr.job_search_profiles_id_seq OWNER TO bridgr;

--
-- Name: job_search_profiles_id_seq; Type: SEQUENCE OWNED BY; Schema: bridgr; Owner: bridgr
--

ALTER SEQUENCE bridgr.job_search_profiles_id_seq OWNED BY bridgr.job_search_profiles.id;


--
-- Name: schema_migrations; Type: TABLE; Schema: bridgr; Owner: bridgr
--

CREATE TABLE bridgr.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE bridgr.schema_migrations OWNER TO bridgr;

--
-- Name: skill_gap_analyses; Type: TABLE; Schema: bridgr; Owner: bridgr
--

CREATE TABLE bridgr.skill_gap_analyses (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    id bigint NOT NULL,
    user_id integer NOT NULL,
    founder_persona_uuid uuid,
    pursuit_uuid uuid,
    title text,
    status character varying(30) DEFAULT 'pending'::character varying NOT NULL,
    cv_asset_uri text,
    jd_asset_uri text,
    cv_fingerprint character varying(128),
    jd_fingerprint character varying(128),
    llm_model character varying(120),
    prompt_version character varying(80),
    extraction_payload jsonb,
    gap_summary jsonb,
    mermaid_diagram text,
    error_code character varying(80),
    error_detail text,
    sqs_message_id character varying(200),
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    CONSTRAINT chk_skill_gap_analyses_status CHECK (((status)::text = ANY ((ARRAY['pending'::character varying, 'extracting'::character varying, 'graphed'::character varying, 'pathed'::character varying, 'completed'::character varying, 'failed'::character varying])::text[])))
);


ALTER TABLE bridgr.skill_gap_analyses OWNER TO bridgr;

--
-- Name: TABLE skill_gap_analyses; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON TABLE bridgr.skill_gap_analyses IS 'Bridgr Skill Gap Navigator: analysis run linking user (and optional founder context) to CV/JD extraction and learning-path outputs.';


--
-- Name: COLUMN skill_gap_analyses.status; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.skill_gap_analyses.status IS 'pending, extracting, graphed, pathed, completed, failed';


--
-- Name: COLUMN skill_gap_analyses.cv_asset_uri; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.skill_gap_analyses.cv_asset_uri IS 'Pointer to stored CV (e.g. S3); optional if inline processing only';


--
-- Name: COLUMN skill_gap_analyses.jd_asset_uri; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.skill_gap_analyses.jd_asset_uri IS 'Pointer to stored job description';


--
-- Name: COLUMN skill_gap_analyses.cv_fingerprint; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.skill_gap_analyses.cv_fingerprint IS 'Hash for dedup / idempotency';


--
-- Name: COLUMN skill_gap_analyses.jd_fingerprint; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.skill_gap_analyses.jd_fingerprint IS 'Hash for dedup / idempotency';


--
-- Name: COLUMN skill_gap_analyses.extraction_payload; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.skill_gap_analyses.extraction_payload IS 'Validated structured LLM output (skills graph extraction)';


--
-- Name: COLUMN skill_gap_analyses.gap_summary; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.skill_gap_analyses.gap_summary IS 'Rollup metrics and narrative gap summary';


--
-- Name: COLUMN skill_gap_analyses.mermaid_diagram; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.skill_gap_analyses.mermaid_diagram IS 'Optional v1 DAG visualization text';


--
-- Name: COLUMN skill_gap_analyses.sqs_message_id; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.skill_gap_analyses.sqs_message_id IS 'AWS SQS MessageId after successful enqueue (debugging / dedupe)';


--
-- Name: skill_gap_analyses_id_seq; Type: SEQUENCE; Schema: bridgr; Owner: bridgr
--

CREATE SEQUENCE bridgr.skill_gap_analyses_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE bridgr.skill_gap_analyses_id_seq OWNER TO bridgr;

--
-- Name: skill_gap_analyses_id_seq; Type: SEQUENCE OWNED BY; Schema: bridgr; Owner: bridgr
--

ALTER SEQUENCE bridgr.skill_gap_analyses_id_seq OWNED BY bridgr.skill_gap_analyses.id;


--
-- Name: skill_gap_coverage; Type: TABLE; Schema: bridgr; Owner: bridgr
--

CREATE TABLE bridgr.skill_gap_coverage (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    id bigint NOT NULL,
    analysis_uuid uuid NOT NULL,
    coverage_kind character varying(32) DEFAULT 'role_skill'::character varying NOT NULL,
    role_skill_key character varying(256),
    candidate_skill_key character varying(256),
    match_status character varying(32) DEFAULT 'unknown'::character varying NOT NULL,
    summary text,
    metrics jsonb DEFAULT '{}'::jsonb,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    CONSTRAINT skill_gap_coverage_kind_chk CHECK (((coverage_kind)::text = ANY ((ARRAY['role_skill'::character varying, 'summary'::character varying, 'aggregate'::character varying])::text[])))
);


ALTER TABLE bridgr.skill_gap_coverage OWNER TO bridgr;

--
-- Name: TABLE skill_gap_coverage; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON TABLE bridgr.skill_gap_coverage IS 'Skill coverage / gap rows or summary metrics for an analysis.';


--
-- Name: COLUMN skill_gap_coverage.analysis_uuid; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.skill_gap_coverage.analysis_uuid IS 'App-enforced FK to skill_gap_analyses.uuid';


--
-- Name: COLUMN skill_gap_coverage.coverage_kind; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.skill_gap_coverage.coverage_kind IS 'role_skill: pairing row; summary: one row roll-up; aggregate: optional bucket';


--
-- Name: COLUMN skill_gap_coverage.match_status; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.skill_gap_coverage.match_status IS 'covered, gap, partial, surplus, not_applicable, unknown';


--
-- Name: skill_gap_coverage_id_seq; Type: SEQUENCE; Schema: bridgr; Owner: bridgr
--

CREATE SEQUENCE bridgr.skill_gap_coverage_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE bridgr.skill_gap_coverage_id_seq OWNER TO bridgr;

--
-- Name: skill_gap_coverage_id_seq; Type: SEQUENCE OWNED BY; Schema: bridgr; Owner: bridgr
--

ALTER SEQUENCE bridgr.skill_gap_coverage_id_seq OWNED BY bridgr.skill_gap_coverage.id;


--
-- Name: skill_gap_edges; Type: TABLE; Schema: bridgr; Owner: bridgr
--

CREATE TABLE bridgr.skill_gap_edges (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    id bigint NOT NULL,
    graph_uuid uuid NOT NULL,
    from_node_uuid uuid NOT NULL,
    to_node_uuid uuid NOT NULL,
    relation character varying(64) DEFAULT 'related'::character varying NOT NULL,
    weight numeric(10,4) DEFAULT 1.0,
    metadata jsonb DEFAULT '{}'::jsonb,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    CONSTRAINT skill_gap_edges_no_self CHECK ((from_node_uuid <> to_node_uuid))
);


ALTER TABLE bridgr.skill_gap_edges OWNER TO bridgr;

--
-- Name: TABLE skill_gap_edges; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON TABLE bridgr.skill_gap_edges IS 'Directed edges between skill-gap nodes in a graph.';


--
-- Name: COLUMN skill_gap_edges.graph_uuid; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.skill_gap_edges.graph_uuid IS 'App-enforced FK to skill_gap_graphs.uuid';


--
-- Name: COLUMN skill_gap_edges.from_node_uuid; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.skill_gap_edges.from_node_uuid IS 'App-enforced FK to skill_gap_nodes.uuid';


--
-- Name: COLUMN skill_gap_edges.to_node_uuid; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.skill_gap_edges.to_node_uuid IS 'App-enforced FK to skill_gap_nodes.uuid';


--
-- Name: COLUMN skill_gap_edges.relation; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.skill_gap_edges.relation IS 'prerequisite, similarity, required_by, etc.';


--
-- Name: skill_gap_edges_id_seq; Type: SEQUENCE; Schema: bridgr; Owner: bridgr
--

CREATE SEQUENCE bridgr.skill_gap_edges_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE bridgr.skill_gap_edges_id_seq OWNER TO bridgr;

--
-- Name: skill_gap_edges_id_seq; Type: SEQUENCE OWNED BY; Schema: bridgr; Owner: bridgr
--

ALTER SEQUENCE bridgr.skill_gap_edges_id_seq OWNED BY bridgr.skill_gap_edges.id;


--
-- Name: skill_gap_graphs; Type: TABLE; Schema: bridgr; Owner: bridgr
--

CREATE TABLE bridgr.skill_gap_graphs (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    id bigint NOT NULL,
    analysis_uuid uuid NOT NULL,
    kind character varying(32) NOT NULL,
    metadata jsonb,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE bridgr.skill_gap_graphs OWNER TO bridgr;

--
-- Name: TABLE skill_gap_graphs; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON TABLE bridgr.skill_gap_graphs IS 'Skill graphs for Bridgr: typically two per analysis (candidate skills vs role requirements).';


--
-- Name: COLUMN skill_gap_graphs.analysis_uuid; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.skill_gap_graphs.analysis_uuid IS 'App-enforced FK to skill_gap_analyses.uuid';


--
-- Name: COLUMN skill_gap_graphs.kind; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.skill_gap_graphs.kind IS 'candidate = graph from CV; role_requirement = graph from job description';


--
-- Name: COLUMN skill_gap_graphs.metadata; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.skill_gap_graphs.metadata IS 'Optional graph-level notes, model version, or layout hints';


--
-- Name: skill_gap_graphs_id_seq; Type: SEQUENCE; Schema: bridgr; Owner: bridgr
--

CREATE SEQUENCE bridgr.skill_gap_graphs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE bridgr.skill_gap_graphs_id_seq OWNER TO bridgr;

--
-- Name: skill_gap_graphs_id_seq; Type: SEQUENCE OWNED BY; Schema: bridgr; Owner: bridgr
--

ALTER SEQUENCE bridgr.skill_gap_graphs_id_seq OWNED BY bridgr.skill_gap_graphs.id;


--
-- Name: skill_gap_learning_paths; Type: TABLE; Schema: bridgr; Owner: bridgr
--

CREATE TABLE bridgr.skill_gap_learning_paths (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    id bigint NOT NULL,
    analysis_uuid uuid NOT NULL,
    path_version integer DEFAULT 1 NOT NULL,
    algorithm character varying(64),
    title text,
    path_metadata jsonb DEFAULT '{}'::jsonb,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE bridgr.skill_gap_learning_paths OWNER TO bridgr;

--
-- Name: TABLE skill_gap_learning_paths; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON TABLE bridgr.skill_gap_learning_paths IS 'Learning plan / path for closing skill gaps from one analysis.';


--
-- Name: COLUMN skill_gap_learning_paths.analysis_uuid; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.skill_gap_learning_paths.analysis_uuid IS 'App-enforced FK to skill_gap_analyses.uuid';


--
-- Name: COLUMN skill_gap_learning_paths.path_version; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.skill_gap_learning_paths.path_version IS 'Increment when regenerating path for same analysis';


--
-- Name: skill_gap_learning_paths_id_seq; Type: SEQUENCE; Schema: bridgr; Owner: bridgr
--

CREATE SEQUENCE bridgr.skill_gap_learning_paths_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE bridgr.skill_gap_learning_paths_id_seq OWNER TO bridgr;

--
-- Name: skill_gap_learning_paths_id_seq; Type: SEQUENCE OWNED BY; Schema: bridgr; Owner: bridgr
--

ALTER SEQUENCE bridgr.skill_gap_learning_paths_id_seq OWNED BY bridgr.skill_gap_learning_paths.id;


--
-- Name: skill_gap_nodes; Type: TABLE; Schema: bridgr; Owner: bridgr
--

CREATE TABLE bridgr.skill_gap_nodes (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    id bigint NOT NULL,
    graph_uuid uuid NOT NULL,
    node_key character varying(256) NOT NULL,
    display_name text NOT NULL,
    description text,
    proficiency_hint text,
    source character varying(32),
    evidence jsonb DEFAULT '{}'::jsonb,
    metadata jsonb DEFAULT '{}'::jsonb,
    position_x integer DEFAULT 0,
    position_y integer DEFAULT 0,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE bridgr.skill_gap_nodes OWNER TO bridgr;

--
-- Name: TABLE skill_gap_nodes; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON TABLE bridgr.skill_gap_nodes IS 'Skill/concept nodes in a skill-gap graph (candidate or role requirement).';


--
-- Name: COLUMN skill_gap_nodes.graph_uuid; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.skill_gap_nodes.graph_uuid IS 'App-enforced FK to skill_gap_graphs.uuid';


--
-- Name: COLUMN skill_gap_nodes.node_key; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.skill_gap_nodes.node_key IS 'Stable key for this node within the graph (e.g. normalized skill id)';


--
-- Name: COLUMN skill_gap_nodes.display_name; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.skill_gap_nodes.display_name IS 'Human-readable label';


--
-- Name: COLUMN skill_gap_nodes.source; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.skill_gap_nodes.source IS 'cv, jd, inferred, merged, etc.';


--
-- Name: COLUMN skill_gap_nodes.evidence; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.skill_gap_nodes.evidence IS 'Spans, quotes, or LLM citations supporting the node';


--
-- Name: skill_gap_nodes_id_seq; Type: SEQUENCE; Schema: bridgr; Owner: bridgr
--

CREATE SEQUENCE bridgr.skill_gap_nodes_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE bridgr.skill_gap_nodes_id_seq OWNER TO bridgr;

--
-- Name: skill_gap_nodes_id_seq; Type: SEQUENCE OWNED BY; Schema: bridgr; Owner: bridgr
--

ALTER SEQUENCE bridgr.skill_gap_nodes_id_seq OWNED BY bridgr.skill_gap_nodes.id;


--
-- Name: skill_gap_path_step_deps; Type: TABLE; Schema: bridgr; Owner: bridgr
--

CREATE TABLE bridgr.skill_gap_path_step_deps (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    id bigint NOT NULL,
    path_uuid uuid NOT NULL,
    step_uuid uuid NOT NULL,
    depends_on_step_uuid uuid NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    CONSTRAINT skill_gap_path_step_deps_no_self CHECK ((step_uuid <> depends_on_step_uuid))
);


ALTER TABLE bridgr.skill_gap_path_step_deps OWNER TO bridgr;

--
-- Name: TABLE skill_gap_path_step_deps; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON TABLE bridgr.skill_gap_path_step_deps IS 'Prerequisite edges between steps in a learning path (DAG).';


--
-- Name: COLUMN skill_gap_path_step_deps.path_uuid; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.skill_gap_path_step_deps.path_uuid IS 'App-enforced FK to skill_gap_learning_paths.uuid';


--
-- Name: COLUMN skill_gap_path_step_deps.step_uuid; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.skill_gap_path_step_deps.step_uuid IS 'App-enforced FK to skill_gap_path_steps.uuid — dependent step';


--
-- Name: COLUMN skill_gap_path_step_deps.depends_on_step_uuid; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.skill_gap_path_step_deps.depends_on_step_uuid IS 'App-enforced FK to skill_gap_path_steps.uuid — must be done first';


--
-- Name: skill_gap_path_step_deps_id_seq; Type: SEQUENCE; Schema: bridgr; Owner: bridgr
--

CREATE SEQUENCE bridgr.skill_gap_path_step_deps_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE bridgr.skill_gap_path_step_deps_id_seq OWNER TO bridgr;

--
-- Name: skill_gap_path_step_deps_id_seq; Type: SEQUENCE OWNED BY; Schema: bridgr; Owner: bridgr
--

ALTER SEQUENCE bridgr.skill_gap_path_step_deps_id_seq OWNED BY bridgr.skill_gap_path_step_deps.id;


--
-- Name: skill_gap_path_steps; Type: TABLE; Schema: bridgr; Owner: bridgr
--

CREATE TABLE bridgr.skill_gap_path_steps (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    id bigint NOT NULL,
    path_uuid uuid NOT NULL,
    step_index integer NOT NULL,
    title text NOT NULL,
    rationale text,
    estimated_hours numeric(10,2),
    resource_uri text,
    resource_kind character varying(64),
    founder_learning_item_uuid uuid,
    course_lesson_uuid uuid,
    linked_node_keys jsonb DEFAULT '[]'::jsonb,
    metadata jsonb DEFAULT '{}'::jsonb,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE bridgr.skill_gap_path_steps OWNER TO bridgr;

--
-- Name: TABLE skill_gap_path_steps; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON TABLE bridgr.skill_gap_path_steps IS 'One step in a generated learning path.';


--
-- Name: COLUMN skill_gap_path_steps.path_uuid; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.skill_gap_path_steps.path_uuid IS 'App-enforced FK to skill_gap_learning_paths.uuid';


--
-- Name: COLUMN skill_gap_path_steps.step_index; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.skill_gap_path_steps.step_index IS 'Order within the path (0-based or 1-based per app convention)';


--
-- Name: COLUMN skill_gap_path_steps.founder_learning_item_uuid; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.skill_gap_path_steps.founder_learning_item_uuid IS 'App-enforced FK to founder_learning_items.uuid when linked';


--
-- Name: COLUMN skill_gap_path_steps.course_lesson_uuid; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.skill_gap_path_steps.course_lesson_uuid IS 'App-enforced FK to course lesson when linked';


--
-- Name: COLUMN skill_gap_path_steps.linked_node_keys; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.skill_gap_path_steps.linked_node_keys IS 'JSON array of graph node_key values this step addresses';


--
-- Name: skill_gap_path_steps_id_seq; Type: SEQUENCE; Schema: bridgr; Owner: bridgr
--

CREATE SEQUENCE bridgr.skill_gap_path_steps_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE bridgr.skill_gap_path_steps_id_seq OWNER TO bridgr;

--
-- Name: skill_gap_path_steps_id_seq; Type: SEQUENCE OWNED BY; Schema: bridgr; Owner: bridgr
--

ALTER SEQUENCE bridgr.skill_gap_path_steps_id_seq OWNED BY bridgr.skill_gap_path_steps.id;


--
-- Name: supported_job_boards; Type: TABLE; Schema: bridgr; Owner: bridgr
--

CREATE TABLE bridgr.supported_job_boards (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    id bigint NOT NULL,
    board_id character varying(64) NOT NULL,
    display_name character varying(128) NOT NULL,
    engine character varying(64) NOT NULL,
    site_type character varying(32) NOT NULL,
    region character varying(32) DEFAULT 'global'::character varying NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    config jsonb DEFAULT '{}'::jsonb NOT NULL,
    sort_order integer DEFAULT 0 NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE bridgr.supported_job_boards OWNER TO bridgr;

--
-- Name: TABLE supported_job_boards; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON TABLE bridgr.supported_job_boards IS 'Bridgr/Radar: canonical list of supported job boards (CRUD via API; seeds from former supported_boards.yaml).';


--
-- Name: COLUMN supported_job_boards.board_id; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.supported_job_boards.board_id IS 'Stable slug (e.g. linkedin, indeed); matches Discovery source_site / JobSpy site key';


--
-- Name: COLUMN supported_job_boards.display_name; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.supported_job_boards.display_name IS 'Human-readable label for UI';


--
-- Name: COLUMN supported_job_boards.engine; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.supported_job_boards.engine IS 'Radar engine: jobspy, smartextract, crawl4ai, workday, indeed, etc.';


--
-- Name: COLUMN supported_job_boards.site_type; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.supported_job_boards.site_type IS 'How the board is crawled: search (query+location URL), static (fixed career pages), etc.';


--
-- Name: COLUMN supported_job_boards.region; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.supported_job_boards.region IS 'Primary region: global, us, uk, eu, …';


--
-- Name: COLUMN supported_job_boards.is_active; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.supported_job_boards.is_active IS 'When false, board is hidden from selection and not used for new crawls';


--
-- Name: COLUMN supported_job_boards.config; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.supported_job_boards.config IS 'Engine-specific options (e.g. fetch_description, proxy hints); extensible JSON';


--
-- Name: COLUMN supported_job_boards.sort_order; Type: COMMENT; Schema: bridgr; Owner: bridgr
--

COMMENT ON COLUMN bridgr.supported_job_boards.sort_order IS 'Lower values list first in admin and user pickers';


--
-- Name: supported_job_boards_id_seq; Type: SEQUENCE; Schema: bridgr; Owner: bridgr
--

CREATE SEQUENCE bridgr.supported_job_boards_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE bridgr.supported_job_boards_id_seq OWNER TO bridgr;

--
-- Name: supported_job_boards_id_seq; Type: SEQUENCE OWNED BY; Schema: bridgr; Owner: bridgr
--

ALTER SEQUENCE bridgr.supported_job_boards_id_seq OWNED BY bridgr.supported_job_boards.id;


--
-- Name: analysis_job_link id; Type: DEFAULT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.analysis_job_link ALTER COLUMN id SET DEFAULT nextval('bridgr.analysis_job_link_id_seq'::regclass);


--
-- Name: feed_items id; Type: DEFAULT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.feed_items ALTER COLUMN id SET DEFAULT nextval('bridgr.feed_items_id_seq'::regclass);


--
-- Name: job_candidates id; Type: DEFAULT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.job_candidates ALTER COLUMN id SET DEFAULT nextval('bridgr.job_candidates_id_seq'::regclass);


--
-- Name: job_enrichments id; Type: DEFAULT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.job_enrichments ALTER COLUMN id SET DEFAULT nextval('bridgr.job_enrichments_id_seq'::regclass);


--
-- Name: job_harvest_schedules id; Type: DEFAULT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.job_harvest_schedules ALTER COLUMN id SET DEFAULT nextval('bridgr.job_harvest_schedules_id_seq'::regclass);


--
-- Name: job_notifications id; Type: DEFAULT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.job_notifications ALTER COLUMN id SET DEFAULT nextval('bridgr.job_notifications_id_seq'::regclass);


--
-- Name: job_scores id; Type: DEFAULT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.job_scores ALTER COLUMN id SET DEFAULT nextval('bridgr.job_scores_id_seq'::regclass);


--
-- Name: job_search_discovery_runs id; Type: DEFAULT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.job_search_discovery_runs ALTER COLUMN id SET DEFAULT nextval('bridgr.job_search_discovery_runs_id_seq'::regclass);


--
-- Name: job_search_profiles id; Type: DEFAULT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.job_search_profiles ALTER COLUMN id SET DEFAULT nextval('bridgr.job_search_profiles_id_seq'::regclass);


--
-- Name: skill_gap_analyses id; Type: DEFAULT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.skill_gap_analyses ALTER COLUMN id SET DEFAULT nextval('bridgr.skill_gap_analyses_id_seq'::regclass);


--
-- Name: skill_gap_coverage id; Type: DEFAULT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.skill_gap_coverage ALTER COLUMN id SET DEFAULT nextval('bridgr.skill_gap_coverage_id_seq'::regclass);


--
-- Name: skill_gap_edges id; Type: DEFAULT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.skill_gap_edges ALTER COLUMN id SET DEFAULT nextval('bridgr.skill_gap_edges_id_seq'::regclass);


--
-- Name: skill_gap_graphs id; Type: DEFAULT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.skill_gap_graphs ALTER COLUMN id SET DEFAULT nextval('bridgr.skill_gap_graphs_id_seq'::regclass);


--
-- Name: skill_gap_learning_paths id; Type: DEFAULT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.skill_gap_learning_paths ALTER COLUMN id SET DEFAULT nextval('bridgr.skill_gap_learning_paths_id_seq'::regclass);


--
-- Name: skill_gap_nodes id; Type: DEFAULT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.skill_gap_nodes ALTER COLUMN id SET DEFAULT nextval('bridgr.skill_gap_nodes_id_seq'::regclass);


--
-- Name: skill_gap_path_step_deps id; Type: DEFAULT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.skill_gap_path_step_deps ALTER COLUMN id SET DEFAULT nextval('bridgr.skill_gap_path_step_deps_id_seq'::regclass);


--
-- Name: skill_gap_path_steps id; Type: DEFAULT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.skill_gap_path_steps ALTER COLUMN id SET DEFAULT nextval('bridgr.skill_gap_path_steps_id_seq'::regclass);


--
-- Name: supported_job_boards id; Type: DEFAULT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.supported_job_boards ALTER COLUMN id SET DEFAULT nextval('bridgr.supported_job_boards_id_seq'::regclass);


--
-- Name: analysis_job_link analysis_job_link_id_key; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.analysis_job_link
    ADD CONSTRAINT analysis_job_link_id_key UNIQUE (id);


--
-- Name: analysis_job_link analysis_job_link_pkey; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.analysis_job_link
    ADD CONSTRAINT analysis_job_link_pkey PRIMARY KEY (uuid);


--
-- Name: feed_items feed_items_id_key; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.feed_items
    ADD CONSTRAINT feed_items_id_key UNIQUE (id);


--
-- Name: feed_items feed_items_pkey; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.feed_items
    ADD CONSTRAINT feed_items_pkey PRIMARY KEY (uuid);


--
-- Name: job_candidates job_candidates_id_key; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.job_candidates
    ADD CONSTRAINT job_candidates_id_key UNIQUE (id);


--
-- Name: job_candidates job_candidates_pkey; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.job_candidates
    ADD CONSTRAINT job_candidates_pkey PRIMARY KEY (uuid);


--
-- Name: job_enrichments job_enrichments_id_key; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.job_enrichments
    ADD CONSTRAINT job_enrichments_id_key UNIQUE (id);


--
-- Name: job_enrichments job_enrichments_pkey; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.job_enrichments
    ADD CONSTRAINT job_enrichments_pkey PRIMARY KEY (uuid);


--
-- Name: job_harvest_schedules job_harvest_schedules_id_key; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.job_harvest_schedules
    ADD CONSTRAINT job_harvest_schedules_id_key UNIQUE (id);


--
-- Name: job_harvest_schedules job_harvest_schedules_pkey; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.job_harvest_schedules
    ADD CONSTRAINT job_harvest_schedules_pkey PRIMARY KEY (uuid);


--
-- Name: job_notifications job_notifications_id_key; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.job_notifications
    ADD CONSTRAINT job_notifications_id_key UNIQUE (id);


--
-- Name: job_notifications job_notifications_pkey; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.job_notifications
    ADD CONSTRAINT job_notifications_pkey PRIMARY KEY (uuid);


--
-- Name: job_scores job_scores_id_key; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.job_scores
    ADD CONSTRAINT job_scores_id_key UNIQUE (id);


--
-- Name: job_scores job_scores_pkey; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.job_scores
    ADD CONSTRAINT job_scores_pkey PRIMARY KEY (uuid);


--
-- Name: job_search_discovery_runs job_search_discovery_runs_id_key; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.job_search_discovery_runs
    ADD CONSTRAINT job_search_discovery_runs_id_key UNIQUE (id);


--
-- Name: job_search_discovery_runs job_search_discovery_runs_pkey; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.job_search_discovery_runs
    ADD CONSTRAINT job_search_discovery_runs_pkey PRIMARY KEY (uuid);


--
-- Name: job_search_profiles job_search_profiles_id_key; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.job_search_profiles
    ADD CONSTRAINT job_search_profiles_id_key UNIQUE (id);


--
-- Name: job_search_profiles job_search_profiles_pkey; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.job_search_profiles
    ADD CONSTRAINT job_search_profiles_pkey PRIMARY KEY (uuid);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: skill_gap_analyses skill_gap_analyses_id_key; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.skill_gap_analyses
    ADD CONSTRAINT skill_gap_analyses_id_key UNIQUE (id);


--
-- Name: skill_gap_analyses skill_gap_analyses_pkey; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.skill_gap_analyses
    ADD CONSTRAINT skill_gap_analyses_pkey PRIMARY KEY (uuid);


--
-- Name: skill_gap_coverage skill_gap_coverage_id_key; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.skill_gap_coverage
    ADD CONSTRAINT skill_gap_coverage_id_key UNIQUE (id);


--
-- Name: skill_gap_coverage skill_gap_coverage_pkey; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.skill_gap_coverage
    ADD CONSTRAINT skill_gap_coverage_pkey PRIMARY KEY (uuid);


--
-- Name: skill_gap_edges skill_gap_edges_graph_uuid_from_node_uuid_to_node_uuid_rela_key; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.skill_gap_edges
    ADD CONSTRAINT skill_gap_edges_graph_uuid_from_node_uuid_to_node_uuid_rela_key UNIQUE (graph_uuid, from_node_uuid, to_node_uuid, relation);


--
-- Name: skill_gap_edges skill_gap_edges_id_key; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.skill_gap_edges
    ADD CONSTRAINT skill_gap_edges_id_key UNIQUE (id);


--
-- Name: skill_gap_edges skill_gap_edges_pkey; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.skill_gap_edges
    ADD CONSTRAINT skill_gap_edges_pkey PRIMARY KEY (uuid);


--
-- Name: skill_gap_graphs skill_gap_graphs_id_key; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.skill_gap_graphs
    ADD CONSTRAINT skill_gap_graphs_id_key UNIQUE (id);


--
-- Name: skill_gap_graphs skill_gap_graphs_pkey; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.skill_gap_graphs
    ADD CONSTRAINT skill_gap_graphs_pkey PRIMARY KEY (uuid);


--
-- Name: skill_gap_learning_paths skill_gap_learning_paths_analysis_uuid_path_version_key; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.skill_gap_learning_paths
    ADD CONSTRAINT skill_gap_learning_paths_analysis_uuid_path_version_key UNIQUE (analysis_uuid, path_version);


--
-- Name: skill_gap_learning_paths skill_gap_learning_paths_id_key; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.skill_gap_learning_paths
    ADD CONSTRAINT skill_gap_learning_paths_id_key UNIQUE (id);


--
-- Name: skill_gap_learning_paths skill_gap_learning_paths_pkey; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.skill_gap_learning_paths
    ADD CONSTRAINT skill_gap_learning_paths_pkey PRIMARY KEY (uuid);


--
-- Name: skill_gap_nodes skill_gap_nodes_graph_uuid_node_key_key; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.skill_gap_nodes
    ADD CONSTRAINT skill_gap_nodes_graph_uuid_node_key_key UNIQUE (graph_uuid, node_key);


--
-- Name: skill_gap_nodes skill_gap_nodes_id_key; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.skill_gap_nodes
    ADD CONSTRAINT skill_gap_nodes_id_key UNIQUE (id);


--
-- Name: skill_gap_nodes skill_gap_nodes_pkey; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.skill_gap_nodes
    ADD CONSTRAINT skill_gap_nodes_pkey PRIMARY KEY (uuid);


--
-- Name: skill_gap_path_step_deps skill_gap_path_step_deps_id_key; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.skill_gap_path_step_deps
    ADD CONSTRAINT skill_gap_path_step_deps_id_key UNIQUE (id);


--
-- Name: skill_gap_path_step_deps skill_gap_path_step_deps_path_uuid_step_uuid_depends_on_ste_key; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.skill_gap_path_step_deps
    ADD CONSTRAINT skill_gap_path_step_deps_path_uuid_step_uuid_depends_on_ste_key UNIQUE (path_uuid, step_uuid, depends_on_step_uuid);


--
-- Name: skill_gap_path_step_deps skill_gap_path_step_deps_pkey; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.skill_gap_path_step_deps
    ADD CONSTRAINT skill_gap_path_step_deps_pkey PRIMARY KEY (uuid);


--
-- Name: skill_gap_path_steps skill_gap_path_steps_id_key; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.skill_gap_path_steps
    ADD CONSTRAINT skill_gap_path_steps_id_key UNIQUE (id);


--
-- Name: skill_gap_path_steps skill_gap_path_steps_path_uuid_step_index_key; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.skill_gap_path_steps
    ADD CONSTRAINT skill_gap_path_steps_path_uuid_step_index_key UNIQUE (path_uuid, step_index);


--
-- Name: skill_gap_path_steps skill_gap_path_steps_pkey; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.skill_gap_path_steps
    ADD CONSTRAINT skill_gap_path_steps_pkey PRIMARY KEY (uuid);


--
-- Name: supported_job_boards supported_job_boards_id_key; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.supported_job_boards
    ADD CONSTRAINT supported_job_boards_id_key UNIQUE (id);


--
-- Name: supported_job_boards supported_job_boards_pkey; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.supported_job_boards
    ADD CONSTRAINT supported_job_boards_pkey PRIMARY KEY (uuid);


--
-- Name: analysis_job_link uq_analysis_job_link_pair; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.analysis_job_link
    ADD CONSTRAINT uq_analysis_job_link_pair UNIQUE (analysis_uuid, job_candidate_uuid);


--
-- Name: feed_items uq_feed_items_candidate_user; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.feed_items
    ADD CONSTRAINT uq_feed_items_candidate_user UNIQUE (job_candidate_uuid, user_id);


--
-- Name: job_enrichments uq_job_enrichments_candidate_user; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.job_enrichments
    ADD CONSTRAINT uq_job_enrichments_candidate_user UNIQUE (job_candidate_uuid, user_id);


--
-- Name: job_harvest_schedules uq_job_harvest_schedules_user_profile; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.job_harvest_schedules
    ADD CONSTRAINT uq_job_harvest_schedules_user_profile UNIQUE (user_id, profile_uuid);


--
-- Name: job_scores uq_job_scores_candidate_user; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.job_scores
    ADD CONSTRAINT uq_job_scores_candidate_user UNIQUE (job_candidate_uuid, user_id);


--
-- Name: job_search_profiles uq_job_search_profiles_user; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.job_search_profiles
    ADD CONSTRAINT uq_job_search_profiles_user UNIQUE (user_id);


--
-- Name: skill_gap_graphs uq_skill_gap_graphs_analysis_kind; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.skill_gap_graphs
    ADD CONSTRAINT uq_skill_gap_graphs_analysis_kind UNIQUE (analysis_uuid, kind);


--
-- Name: supported_job_boards uq_supported_job_boards_board_id; Type: CONSTRAINT; Schema: bridgr; Owner: bridgr
--

ALTER TABLE ONLY bridgr.supported_job_boards
    ADD CONSTRAINT uq_supported_job_boards_board_id UNIQUE (board_id);


--
-- Name: idx_analysis_job_link_analysis; Type: INDEX; Schema: bridgr; Owner: bridgr
--

CREATE INDEX idx_analysis_job_link_analysis ON bridgr.analysis_job_link USING btree (analysis_uuid);


--
-- Name: idx_analysis_job_link_candidate; Type: INDEX; Schema: bridgr; Owner: bridgr
--

CREATE INDEX idx_analysis_job_link_candidate ON bridgr.analysis_job_link USING btree (job_candidate_uuid);


--
-- Name: idx_analysis_job_link_user; Type: INDEX; Schema: bridgr; Owner: bridgr
--

CREATE INDEX idx_analysis_job_link_user ON bridgr.analysis_job_link USING btree (user_id, created_at DESC);


--
-- Name: idx_feed_items_candidate; Type: INDEX; Schema: bridgr; Owner: bridgr
--

CREATE INDEX idx_feed_items_candidate ON bridgr.feed_items USING btree (job_candidate_uuid);


--
-- Name: idx_feed_items_user_score; Type: INDEX; Schema: bridgr; Owner: bridgr
--

CREATE INDEX idx_feed_items_user_score ON bridgr.feed_items USING btree (user_id, composite_score DESC);


--
-- Name: idx_feed_items_user_status_surfaced; Type: INDEX; Schema: bridgr; Owner: bridgr
--

CREATE INDEX idx_feed_items_user_status_surfaced ON bridgr.feed_items USING btree (user_id, feed_status, surfaced_at DESC);


--
-- Name: idx_job_candidates_discovery_run; Type: INDEX; Schema: bridgr; Owner: bridgr
--

CREATE INDEX idx_job_candidates_discovery_run ON bridgr.job_candidates USING btree (discovery_run_uuid);


--
-- Name: idx_job_candidates_user_created; Type: INDEX; Schema: bridgr; Owner: bridgr
--

CREATE INDEX idx_job_candidates_user_created ON bridgr.job_candidates USING btree (user_id, created_at DESC);


--
-- Name: idx_job_enrichments_candidate; Type: INDEX; Schema: bridgr; Owner: bridgr
--

CREATE INDEX idx_job_enrichments_candidate ON bridgr.job_enrichments USING btree (job_candidate_uuid);


--
-- Name: idx_job_enrichments_user_status; Type: INDEX; Schema: bridgr; Owner: bridgr
--

CREATE INDEX idx_job_enrichments_user_status ON bridgr.job_enrichments USING btree (user_id, status, updated_at DESC);


--
-- Name: idx_job_harvest_schedules_next_run; Type: INDEX; Schema: bridgr; Owner: bridgr
--

CREATE INDEX idx_job_harvest_schedules_next_run ON bridgr.job_harvest_schedules USING btree (next_run_at) WHERE (enabled = true);


--
-- Name: idx_job_harvest_schedules_user; Type: INDEX; Schema: bridgr; Owner: bridgr
--

CREATE INDEX idx_job_harvest_schedules_user ON bridgr.job_harvest_schedules USING btree (user_id, enabled);


--
-- Name: idx_job_notifications_candidate; Type: INDEX; Schema: bridgr; Owner: bridgr
--

CREATE INDEX idx_job_notifications_candidate ON bridgr.job_notifications USING btree (job_candidate_uuid);


--
-- Name: idx_job_notifications_user_status; Type: INDEX; Schema: bridgr; Owner: bridgr
--

CREATE INDEX idx_job_notifications_user_status ON bridgr.job_notifications USING btree (user_id, status, created_at DESC);


--
-- Name: idx_job_scores_candidate; Type: INDEX; Schema: bridgr; Owner: bridgr
--

CREATE INDEX idx_job_scores_candidate ON bridgr.job_scores USING btree (job_candidate_uuid);


--
-- Name: idx_job_scores_user_composite; Type: INDEX; Schema: bridgr; Owner: bridgr
--

CREATE INDEX idx_job_scores_user_composite ON bridgr.job_scores USING btree (user_id, composite_score DESC);


--
-- Name: idx_job_search_discovery_runs_status; Type: INDEX; Schema: bridgr; Owner: bridgr
--

CREATE INDEX idx_job_search_discovery_runs_status ON bridgr.job_search_discovery_runs USING btree (user_id, status);


--
-- Name: idx_job_search_discovery_runs_user_created; Type: INDEX; Schema: bridgr; Owner: bridgr
--

CREATE INDEX idx_job_search_discovery_runs_user_created ON bridgr.job_search_discovery_runs USING btree (user_id, created_at DESC);


--
-- Name: idx_skill_gap_analyses_persona; Type: INDEX; Schema: bridgr; Owner: bridgr
--

CREATE INDEX idx_skill_gap_analyses_persona ON bridgr.skill_gap_analyses USING btree (founder_persona_uuid);


--
-- Name: idx_skill_gap_analyses_pursuit; Type: INDEX; Schema: bridgr; Owner: bridgr
--

CREATE INDEX idx_skill_gap_analyses_pursuit ON bridgr.skill_gap_analyses USING btree (pursuit_uuid);


--
-- Name: idx_skill_gap_analyses_status; Type: INDEX; Schema: bridgr; Owner: bridgr
--

CREATE INDEX idx_skill_gap_analyses_status ON bridgr.skill_gap_analyses USING btree (status);


--
-- Name: idx_skill_gap_analyses_user_created; Type: INDEX; Schema: bridgr; Owner: bridgr
--

CREATE INDEX idx_skill_gap_analyses_user_created ON bridgr.skill_gap_analyses USING btree (user_id, created_at DESC);


--
-- Name: idx_skill_gap_coverage_analysis_uuid; Type: INDEX; Schema: bridgr; Owner: bridgr
--

CREATE INDEX idx_skill_gap_coverage_analysis_uuid ON bridgr.skill_gap_coverage USING btree (analysis_uuid);


--
-- Name: idx_skill_gap_coverage_kind; Type: INDEX; Schema: bridgr; Owner: bridgr
--

CREATE INDEX idx_skill_gap_coverage_kind ON bridgr.skill_gap_coverage USING btree (analysis_uuid, coverage_kind);


--
-- Name: idx_skill_gap_edges_from; Type: INDEX; Schema: bridgr; Owner: bridgr
--

CREATE INDEX idx_skill_gap_edges_from ON bridgr.skill_gap_edges USING btree (from_node_uuid);


--
-- Name: idx_skill_gap_edges_graph_uuid; Type: INDEX; Schema: bridgr; Owner: bridgr
--

CREATE INDEX idx_skill_gap_edges_graph_uuid ON bridgr.skill_gap_edges USING btree (graph_uuid);


--
-- Name: idx_skill_gap_edges_relation; Type: INDEX; Schema: bridgr; Owner: bridgr
--

CREATE INDEX idx_skill_gap_edges_relation ON bridgr.skill_gap_edges USING btree (relation);


--
-- Name: idx_skill_gap_edges_to; Type: INDEX; Schema: bridgr; Owner: bridgr
--

CREATE INDEX idx_skill_gap_edges_to ON bridgr.skill_gap_edges USING btree (to_node_uuid);


--
-- Name: idx_skill_gap_graphs_analysis; Type: INDEX; Schema: bridgr; Owner: bridgr
--

CREATE INDEX idx_skill_gap_graphs_analysis ON bridgr.skill_gap_graphs USING btree (analysis_uuid);


--
-- Name: idx_skill_gap_learning_paths_analysis_uuid; Type: INDEX; Schema: bridgr; Owner: bridgr
--

CREATE INDEX idx_skill_gap_learning_paths_analysis_uuid ON bridgr.skill_gap_learning_paths USING btree (analysis_uuid);


--
-- Name: idx_skill_gap_nodes_graph_uuid; Type: INDEX; Schema: bridgr; Owner: bridgr
--

CREATE INDEX idx_skill_gap_nodes_graph_uuid ON bridgr.skill_gap_nodes USING btree (graph_uuid);


--
-- Name: idx_skill_gap_nodes_node_key; Type: INDEX; Schema: bridgr; Owner: bridgr
--

CREATE INDEX idx_skill_gap_nodes_node_key ON bridgr.skill_gap_nodes USING btree (node_key);


--
-- Name: idx_skill_gap_path_step_deps_depends; Type: INDEX; Schema: bridgr; Owner: bridgr
--

CREATE INDEX idx_skill_gap_path_step_deps_depends ON bridgr.skill_gap_path_step_deps USING btree (depends_on_step_uuid);


--
-- Name: idx_skill_gap_path_step_deps_path; Type: INDEX; Schema: bridgr; Owner: bridgr
--

CREATE INDEX idx_skill_gap_path_step_deps_path ON bridgr.skill_gap_path_step_deps USING btree (path_uuid);


--
-- Name: idx_skill_gap_path_step_deps_step; Type: INDEX; Schema: bridgr; Owner: bridgr
--

CREATE INDEX idx_skill_gap_path_step_deps_step ON bridgr.skill_gap_path_step_deps USING btree (step_uuid);


--
-- Name: idx_skill_gap_path_steps_path_uuid; Type: INDEX; Schema: bridgr; Owner: bridgr
--

CREATE INDEX idx_skill_gap_path_steps_path_uuid ON bridgr.skill_gap_path_steps USING btree (path_uuid);


--
-- Name: idx_supported_job_boards_active_sort; Type: INDEX; Schema: bridgr; Owner: bridgr
--

CREATE INDEX idx_supported_job_boards_active_sort ON bridgr.supported_job_boards USING btree (is_active, sort_order, display_name);


--
-- Name: idx_supported_job_boards_engine; Type: INDEX; Schema: bridgr; Owner: bridgr
--

CREATE INDEX idx_supported_job_boards_engine ON bridgr.supported_job_boards USING btree (engine) WHERE (is_active = true);


--
-- Name: uq_job_candidates_user_url_hash; Type: INDEX; Schema: bridgr; Owner: bridgr
--

CREATE UNIQUE INDEX uq_job_candidates_user_url_hash ON bridgr.job_candidates USING btree (user_id, url_hash);


--
-- Name: analysis_job_link tr_analysis_job_link_control_time; Type: TRIGGER; Schema: bridgr; Owner: bridgr
--

CREATE TRIGGER tr_analysis_job_link_control_time BEFORE INSERT OR UPDATE ON bridgr.analysis_job_link FOR EACH ROW EXECUTE FUNCTION bridgr.tr_control_time();


--
-- Name: feed_items tr_feed_items_control_time; Type: TRIGGER; Schema: bridgr; Owner: bridgr
--

CREATE TRIGGER tr_feed_items_control_time BEFORE INSERT OR UPDATE ON bridgr.feed_items FOR EACH ROW EXECUTE FUNCTION bridgr.tr_control_time();


--
-- Name: job_candidates tr_job_candidates_control_time; Type: TRIGGER; Schema: bridgr; Owner: bridgr
--

CREATE TRIGGER tr_job_candidates_control_time BEFORE INSERT OR UPDATE ON bridgr.job_candidates FOR EACH ROW EXECUTE FUNCTION bridgr.tr_control_time();


--
-- Name: job_enrichments tr_job_enrichments_control_time; Type: TRIGGER; Schema: bridgr; Owner: bridgr
--

CREATE TRIGGER tr_job_enrichments_control_time BEFORE INSERT OR UPDATE ON bridgr.job_enrichments FOR EACH ROW EXECUTE FUNCTION bridgr.tr_control_time();


--
-- Name: job_harvest_schedules tr_job_harvest_schedules_control_time; Type: TRIGGER; Schema: bridgr; Owner: bridgr
--

CREATE TRIGGER tr_job_harvest_schedules_control_time BEFORE INSERT OR UPDATE ON bridgr.job_harvest_schedules FOR EACH ROW EXECUTE FUNCTION bridgr.tr_control_time();


--
-- Name: job_notifications tr_job_notifications_control_time; Type: TRIGGER; Schema: bridgr; Owner: bridgr
--

CREATE TRIGGER tr_job_notifications_control_time BEFORE INSERT OR UPDATE ON bridgr.job_notifications FOR EACH ROW EXECUTE FUNCTION bridgr.tr_control_time();


--
-- Name: job_scores tr_job_scores_control_time; Type: TRIGGER; Schema: bridgr; Owner: bridgr
--

CREATE TRIGGER tr_job_scores_control_time BEFORE INSERT OR UPDATE ON bridgr.job_scores FOR EACH ROW EXECUTE FUNCTION bridgr.tr_control_time();


--
-- Name: job_search_discovery_runs tr_job_search_discovery_runs_control_time; Type: TRIGGER; Schema: bridgr; Owner: bridgr
--

CREATE TRIGGER tr_job_search_discovery_runs_control_time BEFORE INSERT OR UPDATE ON bridgr.job_search_discovery_runs FOR EACH ROW EXECUTE FUNCTION bridgr.tr_control_time();


--
-- Name: job_search_profiles tr_job_search_profiles_control_time; Type: TRIGGER; Schema: bridgr; Owner: bridgr
--

CREATE TRIGGER tr_job_search_profiles_control_time BEFORE INSERT OR UPDATE ON bridgr.job_search_profiles FOR EACH ROW EXECUTE FUNCTION bridgr.tr_control_time();


--
-- Name: skill_gap_analyses tr_skill_gap_analyses_control_time; Type: TRIGGER; Schema: bridgr; Owner: bridgr
--

CREATE TRIGGER tr_skill_gap_analyses_control_time BEFORE INSERT OR UPDATE ON bridgr.skill_gap_analyses FOR EACH ROW EXECUTE FUNCTION bridgr.tr_control_time();


--
-- Name: skill_gap_coverage tr_skill_gap_coverage_control_time; Type: TRIGGER; Schema: bridgr; Owner: bridgr
--

CREATE TRIGGER tr_skill_gap_coverage_control_time BEFORE INSERT OR UPDATE ON bridgr.skill_gap_coverage FOR EACH ROW EXECUTE FUNCTION bridgr.tr_control_time();


--
-- Name: skill_gap_edges tr_skill_gap_edges_control_time; Type: TRIGGER; Schema: bridgr; Owner: bridgr
--

CREATE TRIGGER tr_skill_gap_edges_control_time BEFORE INSERT OR UPDATE ON bridgr.skill_gap_edges FOR EACH ROW EXECUTE FUNCTION bridgr.tr_control_time();


--
-- Name: skill_gap_graphs tr_skill_gap_graphs_control_time; Type: TRIGGER; Schema: bridgr; Owner: bridgr
--

CREATE TRIGGER tr_skill_gap_graphs_control_time BEFORE INSERT OR UPDATE ON bridgr.skill_gap_graphs FOR EACH ROW EXECUTE FUNCTION bridgr.tr_control_time();


--
-- Name: skill_gap_learning_paths tr_skill_gap_learning_paths_control_time; Type: TRIGGER; Schema: bridgr; Owner: bridgr
--

CREATE TRIGGER tr_skill_gap_learning_paths_control_time BEFORE INSERT OR UPDATE ON bridgr.skill_gap_learning_paths FOR EACH ROW EXECUTE FUNCTION bridgr.tr_control_time();


--
-- Name: skill_gap_nodes tr_skill_gap_nodes_control_time; Type: TRIGGER; Schema: bridgr; Owner: bridgr
--

CREATE TRIGGER tr_skill_gap_nodes_control_time BEFORE INSERT OR UPDATE ON bridgr.skill_gap_nodes FOR EACH ROW EXECUTE FUNCTION bridgr.tr_control_time();


--
-- Name: skill_gap_path_step_deps tr_skill_gap_path_step_deps_control_time; Type: TRIGGER; Schema: bridgr; Owner: bridgr
--

CREATE TRIGGER tr_skill_gap_path_step_deps_control_time BEFORE INSERT OR UPDATE ON bridgr.skill_gap_path_step_deps FOR EACH ROW EXECUTE FUNCTION bridgr.tr_control_time();


--
-- Name: skill_gap_path_steps tr_skill_gap_path_steps_control_time; Type: TRIGGER; Schema: bridgr; Owner: bridgr
--

CREATE TRIGGER tr_skill_gap_path_steps_control_time BEFORE INSERT OR UPDATE ON bridgr.skill_gap_path_steps FOR EACH ROW EXECUTE FUNCTION bridgr.tr_control_time();


--
-- Name: supported_job_boards tr_supported_job_boards_control_time; Type: TRIGGER; Schema: bridgr; Owner: bridgr
--

CREATE TRIGGER tr_supported_job_boards_control_time BEFORE INSERT OR UPDATE ON bridgr.supported_job_boards FOR EACH ROW EXECUTE FUNCTION bridgr.tr_control_time();


--
-- Name: SCHEMA bridgr; Type: ACL; Schema: -; Owner: bridgr
--

GRANT USAGE ON SCHEMA bridgr TO bridgr_readonly_role;
GRANT USAGE ON SCHEMA bridgr TO bridgr_app_role;
GRANT ALL ON SCHEMA bridgr TO bridgr_migration_role;


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: pg_database_owner
--

GRANT USAGE ON SCHEMA public TO bridgr_readonly_role;
GRANT USAGE ON SCHEMA public TO bridgr_app_role;
GRANT USAGE ON SCHEMA public TO bridgr_migration_role;


--
-- Name: TABLE schema_migrations; Type: ACL; Schema: bridgr; Owner: bridgr
--

GRANT SELECT ON TABLE bridgr.schema_migrations TO bridgr_readonly_role;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE bridgr.schema_migrations TO bridgr_app_role;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE bridgr.schema_migrations TO bridgr_migration_role;


--
-- Name: DEFAULT PRIVILEGES FOR SEQUENCES; Type: DEFAULT ACL; Schema: bridgr; Owner: bridgr_migration
--

ALTER DEFAULT PRIVILEGES FOR ROLE bridgr_migration IN SCHEMA bridgr GRANT SELECT,USAGE ON SEQUENCES TO bridgr_readonly_role;
ALTER DEFAULT PRIVILEGES FOR ROLE bridgr_migration IN SCHEMA bridgr GRANT SELECT,USAGE ON SEQUENCES TO bridgr_app_role;


--
-- Name: DEFAULT PRIVILEGES FOR TABLES; Type: DEFAULT ACL; Schema: bridgr; Owner: bridgr_migration
--

ALTER DEFAULT PRIVILEGES FOR ROLE bridgr_migration IN SCHEMA bridgr GRANT SELECT ON TABLES TO bridgr_readonly_role;
ALTER DEFAULT PRIVILEGES FOR ROLE bridgr_migration IN SCHEMA bridgr GRANT SELECT,INSERT,DELETE,UPDATE ON TABLES TO bridgr_app_role;


--
-- Name: DEFAULT PRIVILEGES FOR SEQUENCES; Type: DEFAULT ACL; Schema: public; Owner: bridgr_migration
--

ALTER DEFAULT PRIVILEGES FOR ROLE bridgr_migration IN SCHEMA public GRANT SELECT,USAGE ON SEQUENCES TO bridgr_readonly_role;
ALTER DEFAULT PRIVILEGES FOR ROLE bridgr_migration IN SCHEMA public GRANT SELECT,USAGE ON SEQUENCES TO bridgr_app_role;


--
-- Name: DEFAULT PRIVILEGES FOR TABLES; Type: DEFAULT ACL; Schema: public; Owner: bridgr_migration
--

ALTER DEFAULT PRIVILEGES FOR ROLE bridgr_migration IN SCHEMA public GRANT SELECT ON TABLES TO bridgr_readonly_role;
ALTER DEFAULT PRIVILEGES FOR ROLE bridgr_migration IN SCHEMA public GRANT SELECT ON TABLES TO bridgr_app_role;


--
-- PostgreSQL database dump complete
--


