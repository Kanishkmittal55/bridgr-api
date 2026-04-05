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
-- Name: hskip_users; Type: SCHEMA; Schema: -; Owner: bridgr
--

CREATE SCHEMA hskip_users;


ALTER SCHEMA hskip_users OWNER TO bridgr;

--
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


--
-- Name: EXTENSION pgcrypto; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION pgcrypto IS 'cryptographic functions';


--
-- Name: tr_control_time(); Type: FUNCTION; Schema: hskip_users; Owner: bridgr
--

CREATE FUNCTION hskip_users.tr_control_time() RETURNS trigger
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


ALTER FUNCTION hskip_users.tr_control_time() OWNER TO bridgr;

--
-- Name: FUNCTION tr_control_time(); Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON FUNCTION hskip_users.tr_control_time() IS 'A generic trigger function to manage created_at and updated_at columns.';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: bridgr_skill_gap_analyses; Type: TABLE; Schema: hskip_users; Owner: bridgr
--

CREATE TABLE hskip_users.bridgr_skill_gap_analyses (
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
    CONSTRAINT chk_bridgr_skill_gap_analyses_status CHECK (((status)::text = ANY ((ARRAY['pending'::character varying, 'extracting'::character varying, 'graphed'::character varying, 'pathed'::character varying, 'completed'::character varying, 'failed'::character varying])::text[])))
);


ALTER TABLE hskip_users.bridgr_skill_gap_analyses OWNER TO bridgr;

--
-- Name: TABLE bridgr_skill_gap_analyses; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON TABLE hskip_users.bridgr_skill_gap_analyses IS 'Bridgr Skill Gap Navigator: analysis run linking user (and optional founder context) to CV/JD extraction and learning-path outputs.';


--
-- Name: COLUMN bridgr_skill_gap_analyses.status; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON COLUMN hskip_users.bridgr_skill_gap_analyses.status IS 'pending, extracting, graphed, pathed, completed, failed';


--
-- Name: COLUMN bridgr_skill_gap_analyses.cv_asset_uri; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON COLUMN hskip_users.bridgr_skill_gap_analyses.cv_asset_uri IS 'Pointer to stored CV (e.g. S3); optional if inline processing only';


--
-- Name: COLUMN bridgr_skill_gap_analyses.jd_asset_uri; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON COLUMN hskip_users.bridgr_skill_gap_analyses.jd_asset_uri IS 'Pointer to stored job description';


--
-- Name: COLUMN bridgr_skill_gap_analyses.cv_fingerprint; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON COLUMN hskip_users.bridgr_skill_gap_analyses.cv_fingerprint IS 'Hash for dedup / idempotency';


--
-- Name: COLUMN bridgr_skill_gap_analyses.jd_fingerprint; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON COLUMN hskip_users.bridgr_skill_gap_analyses.jd_fingerprint IS 'Hash for dedup / idempotency';


--
-- Name: COLUMN bridgr_skill_gap_analyses.extraction_payload; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON COLUMN hskip_users.bridgr_skill_gap_analyses.extraction_payload IS 'Validated structured LLM output (skills graph extraction)';


--
-- Name: COLUMN bridgr_skill_gap_analyses.gap_summary; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON COLUMN hskip_users.bridgr_skill_gap_analyses.gap_summary IS 'Rollup metrics and narrative gap summary';


--
-- Name: COLUMN bridgr_skill_gap_analyses.mermaid_diagram; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON COLUMN hskip_users.bridgr_skill_gap_analyses.mermaid_diagram IS 'Optional v1 DAG visualization text';


--
-- Name: COLUMN bridgr_skill_gap_analyses.sqs_message_id; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON COLUMN hskip_users.bridgr_skill_gap_analyses.sqs_message_id IS 'AWS SQS MessageId after successful enqueue (debugging / dedupe)';


--
-- Name: bridgr_skill_gap_analyses_id_seq; Type: SEQUENCE; Schema: hskip_users; Owner: bridgr
--

CREATE SEQUENCE hskip_users.bridgr_skill_gap_analyses_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE hskip_users.bridgr_skill_gap_analyses_id_seq OWNER TO bridgr;

--
-- Name: bridgr_skill_gap_analyses_id_seq; Type: SEQUENCE OWNED BY; Schema: hskip_users; Owner: bridgr
--

ALTER SEQUENCE hskip_users.bridgr_skill_gap_analyses_id_seq OWNED BY hskip_users.bridgr_skill_gap_analyses.id;


--
-- Name: bridgr_skill_gap_coverage; Type: TABLE; Schema: hskip_users; Owner: bridgr
--

CREATE TABLE hskip_users.bridgr_skill_gap_coverage (
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
    CONSTRAINT bridgr_skill_gap_coverage_kind_chk CHECK (((coverage_kind)::text = ANY ((ARRAY['role_skill'::character varying, 'summary'::character varying, 'aggregate'::character varying])::text[])))
);


ALTER TABLE hskip_users.bridgr_skill_gap_coverage OWNER TO bridgr;

--
-- Name: TABLE bridgr_skill_gap_coverage; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON TABLE hskip_users.bridgr_skill_gap_coverage IS 'Skill coverage / gap rows or summary metrics for an analysis.';


--
-- Name: COLUMN bridgr_skill_gap_coverage.analysis_uuid; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON COLUMN hskip_users.bridgr_skill_gap_coverage.analysis_uuid IS 'App-enforced FK to bridgr_skill_gap_analyses.uuid';


--
-- Name: COLUMN bridgr_skill_gap_coverage.coverage_kind; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON COLUMN hskip_users.bridgr_skill_gap_coverage.coverage_kind IS 'role_skill: pairing row; summary: one row roll-up; aggregate: optional bucket';


--
-- Name: COLUMN bridgr_skill_gap_coverage.match_status; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON COLUMN hskip_users.bridgr_skill_gap_coverage.match_status IS 'covered, gap, partial, surplus, not_applicable, unknown';


--
-- Name: bridgr_skill_gap_coverage_id_seq; Type: SEQUENCE; Schema: hskip_users; Owner: bridgr
--

CREATE SEQUENCE hskip_users.bridgr_skill_gap_coverage_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE hskip_users.bridgr_skill_gap_coverage_id_seq OWNER TO bridgr;

--
-- Name: bridgr_skill_gap_coverage_id_seq; Type: SEQUENCE OWNED BY; Schema: hskip_users; Owner: bridgr
--

ALTER SEQUENCE hskip_users.bridgr_skill_gap_coverage_id_seq OWNED BY hskip_users.bridgr_skill_gap_coverage.id;


--
-- Name: bridgr_skill_gap_edges; Type: TABLE; Schema: hskip_users; Owner: bridgr
--

CREATE TABLE hskip_users.bridgr_skill_gap_edges (
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
    CONSTRAINT bridgr_skill_gap_edges_no_self CHECK ((from_node_uuid <> to_node_uuid))
);


ALTER TABLE hskip_users.bridgr_skill_gap_edges OWNER TO bridgr;

--
-- Name: TABLE bridgr_skill_gap_edges; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON TABLE hskip_users.bridgr_skill_gap_edges IS 'Directed edges between skill-gap nodes in a graph.';


--
-- Name: COLUMN bridgr_skill_gap_edges.graph_uuid; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON COLUMN hskip_users.bridgr_skill_gap_edges.graph_uuid IS 'App-enforced FK to bridgr_skill_gap_graphs.uuid';


--
-- Name: COLUMN bridgr_skill_gap_edges.from_node_uuid; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON COLUMN hskip_users.bridgr_skill_gap_edges.from_node_uuid IS 'App-enforced FK to bridgr_skill_gap_nodes.uuid';


--
-- Name: COLUMN bridgr_skill_gap_edges.to_node_uuid; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON COLUMN hskip_users.bridgr_skill_gap_edges.to_node_uuid IS 'App-enforced FK to bridgr_skill_gap_nodes.uuid';


--
-- Name: COLUMN bridgr_skill_gap_edges.relation; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON COLUMN hskip_users.bridgr_skill_gap_edges.relation IS 'prerequisite, similarity, required_by, etc.';


--
-- Name: bridgr_skill_gap_edges_id_seq; Type: SEQUENCE; Schema: hskip_users; Owner: bridgr
--

CREATE SEQUENCE hskip_users.bridgr_skill_gap_edges_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE hskip_users.bridgr_skill_gap_edges_id_seq OWNER TO bridgr;

--
-- Name: bridgr_skill_gap_edges_id_seq; Type: SEQUENCE OWNED BY; Schema: hskip_users; Owner: bridgr
--

ALTER SEQUENCE hskip_users.bridgr_skill_gap_edges_id_seq OWNED BY hskip_users.bridgr_skill_gap_edges.id;


--
-- Name: bridgr_skill_gap_graphs; Type: TABLE; Schema: hskip_users; Owner: bridgr
--

CREATE TABLE hskip_users.bridgr_skill_gap_graphs (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    id bigint NOT NULL,
    analysis_uuid uuid NOT NULL,
    kind character varying(32) NOT NULL,
    metadata jsonb,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE hskip_users.bridgr_skill_gap_graphs OWNER TO bridgr;

--
-- Name: TABLE bridgr_skill_gap_graphs; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON TABLE hskip_users.bridgr_skill_gap_graphs IS 'Skill graphs for Bridgr: typically two per analysis (candidate skills vs role requirements).';


--
-- Name: COLUMN bridgr_skill_gap_graphs.analysis_uuid; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON COLUMN hskip_users.bridgr_skill_gap_graphs.analysis_uuid IS 'App-enforced FK to bridgr_skill_gap_analyses.uuid';


--
-- Name: COLUMN bridgr_skill_gap_graphs.kind; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON COLUMN hskip_users.bridgr_skill_gap_graphs.kind IS 'candidate = graph from CV; role_requirement = graph from job description';


--
-- Name: COLUMN bridgr_skill_gap_graphs.metadata; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON COLUMN hskip_users.bridgr_skill_gap_graphs.metadata IS 'Optional graph-level notes, model version, or layout hints';


--
-- Name: bridgr_skill_gap_graphs_id_seq; Type: SEQUENCE; Schema: hskip_users; Owner: bridgr
--

CREATE SEQUENCE hskip_users.bridgr_skill_gap_graphs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE hskip_users.bridgr_skill_gap_graphs_id_seq OWNER TO bridgr;

--
-- Name: bridgr_skill_gap_graphs_id_seq; Type: SEQUENCE OWNED BY; Schema: hskip_users; Owner: bridgr
--

ALTER SEQUENCE hskip_users.bridgr_skill_gap_graphs_id_seq OWNED BY hskip_users.bridgr_skill_gap_graphs.id;


--
-- Name: bridgr_skill_gap_learning_paths; Type: TABLE; Schema: hskip_users; Owner: bridgr
--

CREATE TABLE hskip_users.bridgr_skill_gap_learning_paths (
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


ALTER TABLE hskip_users.bridgr_skill_gap_learning_paths OWNER TO bridgr;

--
-- Name: TABLE bridgr_skill_gap_learning_paths; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON TABLE hskip_users.bridgr_skill_gap_learning_paths IS 'Learning plan / path for closing skill gaps from one analysis.';


--
-- Name: COLUMN bridgr_skill_gap_learning_paths.analysis_uuid; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON COLUMN hskip_users.bridgr_skill_gap_learning_paths.analysis_uuid IS 'App-enforced FK to bridgr_skill_gap_analyses.uuid';


--
-- Name: COLUMN bridgr_skill_gap_learning_paths.path_version; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON COLUMN hskip_users.bridgr_skill_gap_learning_paths.path_version IS 'Increment when regenerating path for same analysis';


--
-- Name: bridgr_skill_gap_learning_paths_id_seq; Type: SEQUENCE; Schema: hskip_users; Owner: bridgr
--

CREATE SEQUENCE hskip_users.bridgr_skill_gap_learning_paths_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE hskip_users.bridgr_skill_gap_learning_paths_id_seq OWNER TO bridgr;

--
-- Name: bridgr_skill_gap_learning_paths_id_seq; Type: SEQUENCE OWNED BY; Schema: hskip_users; Owner: bridgr
--

ALTER SEQUENCE hskip_users.bridgr_skill_gap_learning_paths_id_seq OWNED BY hskip_users.bridgr_skill_gap_learning_paths.id;


--
-- Name: bridgr_skill_gap_nodes; Type: TABLE; Schema: hskip_users; Owner: bridgr
--

CREATE TABLE hskip_users.bridgr_skill_gap_nodes (
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


ALTER TABLE hskip_users.bridgr_skill_gap_nodes OWNER TO bridgr;

--
-- Name: TABLE bridgr_skill_gap_nodes; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON TABLE hskip_users.bridgr_skill_gap_nodes IS 'Skill/concept nodes in a skill-gap graph (candidate or role requirement).';


--
-- Name: COLUMN bridgr_skill_gap_nodes.graph_uuid; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON COLUMN hskip_users.bridgr_skill_gap_nodes.graph_uuid IS 'App-enforced FK to bridgr_skill_gap_graphs.uuid';


--
-- Name: COLUMN bridgr_skill_gap_nodes.node_key; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON COLUMN hskip_users.bridgr_skill_gap_nodes.node_key IS 'Stable key for this node within the graph (e.g. normalized skill id)';


--
-- Name: COLUMN bridgr_skill_gap_nodes.display_name; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON COLUMN hskip_users.bridgr_skill_gap_nodes.display_name IS 'Human-readable label';


--
-- Name: COLUMN bridgr_skill_gap_nodes.source; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON COLUMN hskip_users.bridgr_skill_gap_nodes.source IS 'cv, jd, inferred, merged, etc.';


--
-- Name: COLUMN bridgr_skill_gap_nodes.evidence; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON COLUMN hskip_users.bridgr_skill_gap_nodes.evidence IS 'Spans, quotes, or LLM citations supporting the node';


--
-- Name: bridgr_skill_gap_nodes_id_seq; Type: SEQUENCE; Schema: hskip_users; Owner: bridgr
--

CREATE SEQUENCE hskip_users.bridgr_skill_gap_nodes_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE hskip_users.bridgr_skill_gap_nodes_id_seq OWNER TO bridgr;

--
-- Name: bridgr_skill_gap_nodes_id_seq; Type: SEQUENCE OWNED BY; Schema: hskip_users; Owner: bridgr
--

ALTER SEQUENCE hskip_users.bridgr_skill_gap_nodes_id_seq OWNED BY hskip_users.bridgr_skill_gap_nodes.id;


--
-- Name: bridgr_skill_gap_path_step_deps; Type: TABLE; Schema: hskip_users; Owner: bridgr
--

CREATE TABLE hskip_users.bridgr_skill_gap_path_step_deps (
    uuid uuid DEFAULT gen_random_uuid() NOT NULL,
    id bigint NOT NULL,
    path_uuid uuid NOT NULL,
    step_uuid uuid NOT NULL,
    depends_on_step_uuid uuid NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    CONSTRAINT bridgr_skill_gap_path_step_deps_no_self CHECK ((step_uuid <> depends_on_step_uuid))
);


ALTER TABLE hskip_users.bridgr_skill_gap_path_step_deps OWNER TO bridgr;

--
-- Name: TABLE bridgr_skill_gap_path_step_deps; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON TABLE hskip_users.bridgr_skill_gap_path_step_deps IS 'Prerequisite edges between steps in a learning path (DAG).';


--
-- Name: COLUMN bridgr_skill_gap_path_step_deps.path_uuid; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON COLUMN hskip_users.bridgr_skill_gap_path_step_deps.path_uuid IS 'App-enforced FK to bridgr_skill_gap_learning_paths.uuid';


--
-- Name: COLUMN bridgr_skill_gap_path_step_deps.step_uuid; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON COLUMN hskip_users.bridgr_skill_gap_path_step_deps.step_uuid IS 'App-enforced FK to bridgr_skill_gap_path_steps.uuid — dependent step';


--
-- Name: COLUMN bridgr_skill_gap_path_step_deps.depends_on_step_uuid; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON COLUMN hskip_users.bridgr_skill_gap_path_step_deps.depends_on_step_uuid IS 'App-enforced FK to bridgr_skill_gap_path_steps.uuid — must be done first';


--
-- Name: bridgr_skill_gap_path_step_deps_id_seq; Type: SEQUENCE; Schema: hskip_users; Owner: bridgr
--

CREATE SEQUENCE hskip_users.bridgr_skill_gap_path_step_deps_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE hskip_users.bridgr_skill_gap_path_step_deps_id_seq OWNER TO bridgr;

--
-- Name: bridgr_skill_gap_path_step_deps_id_seq; Type: SEQUENCE OWNED BY; Schema: hskip_users; Owner: bridgr
--

ALTER SEQUENCE hskip_users.bridgr_skill_gap_path_step_deps_id_seq OWNED BY hskip_users.bridgr_skill_gap_path_step_deps.id;


--
-- Name: bridgr_skill_gap_path_steps; Type: TABLE; Schema: hskip_users; Owner: bridgr
--

CREATE TABLE hskip_users.bridgr_skill_gap_path_steps (
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


ALTER TABLE hskip_users.bridgr_skill_gap_path_steps OWNER TO bridgr;

--
-- Name: TABLE bridgr_skill_gap_path_steps; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON TABLE hskip_users.bridgr_skill_gap_path_steps IS 'One step in a generated learning path.';


--
-- Name: COLUMN bridgr_skill_gap_path_steps.path_uuid; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON COLUMN hskip_users.bridgr_skill_gap_path_steps.path_uuid IS 'App-enforced FK to bridgr_skill_gap_learning_paths.uuid';


--
-- Name: COLUMN bridgr_skill_gap_path_steps.step_index; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON COLUMN hskip_users.bridgr_skill_gap_path_steps.step_index IS 'Order within the path (0-based or 1-based per app convention)';


--
-- Name: COLUMN bridgr_skill_gap_path_steps.founder_learning_item_uuid; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON COLUMN hskip_users.bridgr_skill_gap_path_steps.founder_learning_item_uuid IS 'App-enforced FK to founder_learning_items.uuid when linked';


--
-- Name: COLUMN bridgr_skill_gap_path_steps.course_lesson_uuid; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON COLUMN hskip_users.bridgr_skill_gap_path_steps.course_lesson_uuid IS 'App-enforced FK to course lesson when linked';


--
-- Name: COLUMN bridgr_skill_gap_path_steps.linked_node_keys; Type: COMMENT; Schema: hskip_users; Owner: bridgr
--

COMMENT ON COLUMN hskip_users.bridgr_skill_gap_path_steps.linked_node_keys IS 'JSON array of graph node_key values this step addresses';


--
-- Name: bridgr_skill_gap_path_steps_id_seq; Type: SEQUENCE; Schema: hskip_users; Owner: bridgr
--

CREATE SEQUENCE hskip_users.bridgr_skill_gap_path_steps_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE hskip_users.bridgr_skill_gap_path_steps_id_seq OWNER TO bridgr;

--
-- Name: bridgr_skill_gap_path_steps_id_seq; Type: SEQUENCE OWNED BY; Schema: hskip_users; Owner: bridgr
--

ALTER SEQUENCE hskip_users.bridgr_skill_gap_path_steps_id_seq OWNED BY hskip_users.bridgr_skill_gap_path_steps.id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: bridgr
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO bridgr;

--
-- Name: bridgr_skill_gap_analyses id; Type: DEFAULT; Schema: hskip_users; Owner: bridgr
--

ALTER TABLE ONLY hskip_users.bridgr_skill_gap_analyses ALTER COLUMN id SET DEFAULT nextval('hskip_users.bridgr_skill_gap_analyses_id_seq'::regclass);


--
-- Name: bridgr_skill_gap_coverage id; Type: DEFAULT; Schema: hskip_users; Owner: bridgr
--

ALTER TABLE ONLY hskip_users.bridgr_skill_gap_coverage ALTER COLUMN id SET DEFAULT nextval('hskip_users.bridgr_skill_gap_coverage_id_seq'::regclass);


--
-- Name: bridgr_skill_gap_edges id; Type: DEFAULT; Schema: hskip_users; Owner: bridgr
--

ALTER TABLE ONLY hskip_users.bridgr_skill_gap_edges ALTER COLUMN id SET DEFAULT nextval('hskip_users.bridgr_skill_gap_edges_id_seq'::regclass);


--
-- Name: bridgr_skill_gap_graphs id; Type: DEFAULT; Schema: hskip_users; Owner: bridgr
--

ALTER TABLE ONLY hskip_users.bridgr_skill_gap_graphs ALTER COLUMN id SET DEFAULT nextval('hskip_users.bridgr_skill_gap_graphs_id_seq'::regclass);


--
-- Name: bridgr_skill_gap_learning_paths id; Type: DEFAULT; Schema: hskip_users; Owner: bridgr
--

ALTER TABLE ONLY hskip_users.bridgr_skill_gap_learning_paths ALTER COLUMN id SET DEFAULT nextval('hskip_users.bridgr_skill_gap_learning_paths_id_seq'::regclass);


--
-- Name: bridgr_skill_gap_nodes id; Type: DEFAULT; Schema: hskip_users; Owner: bridgr
--

ALTER TABLE ONLY hskip_users.bridgr_skill_gap_nodes ALTER COLUMN id SET DEFAULT nextval('hskip_users.bridgr_skill_gap_nodes_id_seq'::regclass);


--
-- Name: bridgr_skill_gap_path_step_deps id; Type: DEFAULT; Schema: hskip_users; Owner: bridgr
--

ALTER TABLE ONLY hskip_users.bridgr_skill_gap_path_step_deps ALTER COLUMN id SET DEFAULT nextval('hskip_users.bridgr_skill_gap_path_step_deps_id_seq'::regclass);


--
-- Name: bridgr_skill_gap_path_steps id; Type: DEFAULT; Schema: hskip_users; Owner: bridgr
--

ALTER TABLE ONLY hskip_users.bridgr_skill_gap_path_steps ALTER COLUMN id SET DEFAULT nextval('hskip_users.bridgr_skill_gap_path_steps_id_seq'::regclass);


--
-- Name: bridgr_skill_gap_analyses bridgr_skill_gap_analyses_id_key; Type: CONSTRAINT; Schema: hskip_users; Owner: bridgr
--

ALTER TABLE ONLY hskip_users.bridgr_skill_gap_analyses
    ADD CONSTRAINT bridgr_skill_gap_analyses_id_key UNIQUE (id);


--
-- Name: bridgr_skill_gap_analyses bridgr_skill_gap_analyses_pkey; Type: CONSTRAINT; Schema: hskip_users; Owner: bridgr
--

ALTER TABLE ONLY hskip_users.bridgr_skill_gap_analyses
    ADD CONSTRAINT bridgr_skill_gap_analyses_pkey PRIMARY KEY (uuid);


--
-- Name: bridgr_skill_gap_coverage bridgr_skill_gap_coverage_id_key; Type: CONSTRAINT; Schema: hskip_users; Owner: bridgr
--

ALTER TABLE ONLY hskip_users.bridgr_skill_gap_coverage
    ADD CONSTRAINT bridgr_skill_gap_coverage_id_key UNIQUE (id);


--
-- Name: bridgr_skill_gap_coverage bridgr_skill_gap_coverage_pkey; Type: CONSTRAINT; Schema: hskip_users; Owner: bridgr
--

ALTER TABLE ONLY hskip_users.bridgr_skill_gap_coverage
    ADD CONSTRAINT bridgr_skill_gap_coverage_pkey PRIMARY KEY (uuid);


--
-- Name: bridgr_skill_gap_edges bridgr_skill_gap_edges_graph_uuid_from_node_uuid_to_node_uu_key; Type: CONSTRAINT; Schema: hskip_users; Owner: bridgr
--

ALTER TABLE ONLY hskip_users.bridgr_skill_gap_edges
    ADD CONSTRAINT bridgr_skill_gap_edges_graph_uuid_from_node_uuid_to_node_uu_key UNIQUE (graph_uuid, from_node_uuid, to_node_uuid, relation);


--
-- Name: bridgr_skill_gap_edges bridgr_skill_gap_edges_id_key; Type: CONSTRAINT; Schema: hskip_users; Owner: bridgr
--

ALTER TABLE ONLY hskip_users.bridgr_skill_gap_edges
    ADD CONSTRAINT bridgr_skill_gap_edges_id_key UNIQUE (id);


--
-- Name: bridgr_skill_gap_edges bridgr_skill_gap_edges_pkey; Type: CONSTRAINT; Schema: hskip_users; Owner: bridgr
--

ALTER TABLE ONLY hskip_users.bridgr_skill_gap_edges
    ADD CONSTRAINT bridgr_skill_gap_edges_pkey PRIMARY KEY (uuid);


--
-- Name: bridgr_skill_gap_graphs bridgr_skill_gap_graphs_id_key; Type: CONSTRAINT; Schema: hskip_users; Owner: bridgr
--

ALTER TABLE ONLY hskip_users.bridgr_skill_gap_graphs
    ADD CONSTRAINT bridgr_skill_gap_graphs_id_key UNIQUE (id);


--
-- Name: bridgr_skill_gap_graphs bridgr_skill_gap_graphs_pkey; Type: CONSTRAINT; Schema: hskip_users; Owner: bridgr
--

ALTER TABLE ONLY hskip_users.bridgr_skill_gap_graphs
    ADD CONSTRAINT bridgr_skill_gap_graphs_pkey PRIMARY KEY (uuid);


--
-- Name: bridgr_skill_gap_learning_paths bridgr_skill_gap_learning_paths_analysis_uuid_path_version_key; Type: CONSTRAINT; Schema: hskip_users; Owner: bridgr
--

ALTER TABLE ONLY hskip_users.bridgr_skill_gap_learning_paths
    ADD CONSTRAINT bridgr_skill_gap_learning_paths_analysis_uuid_path_version_key UNIQUE (analysis_uuid, path_version);


--
-- Name: bridgr_skill_gap_learning_paths bridgr_skill_gap_learning_paths_id_key; Type: CONSTRAINT; Schema: hskip_users; Owner: bridgr
--

ALTER TABLE ONLY hskip_users.bridgr_skill_gap_learning_paths
    ADD CONSTRAINT bridgr_skill_gap_learning_paths_id_key UNIQUE (id);


--
-- Name: bridgr_skill_gap_learning_paths bridgr_skill_gap_learning_paths_pkey; Type: CONSTRAINT; Schema: hskip_users; Owner: bridgr
--

ALTER TABLE ONLY hskip_users.bridgr_skill_gap_learning_paths
    ADD CONSTRAINT bridgr_skill_gap_learning_paths_pkey PRIMARY KEY (uuid);


--
-- Name: bridgr_skill_gap_nodes bridgr_skill_gap_nodes_graph_uuid_node_key_key; Type: CONSTRAINT; Schema: hskip_users; Owner: bridgr
--

ALTER TABLE ONLY hskip_users.bridgr_skill_gap_nodes
    ADD CONSTRAINT bridgr_skill_gap_nodes_graph_uuid_node_key_key UNIQUE (graph_uuid, node_key);


--
-- Name: bridgr_skill_gap_nodes bridgr_skill_gap_nodes_id_key; Type: CONSTRAINT; Schema: hskip_users; Owner: bridgr
--

ALTER TABLE ONLY hskip_users.bridgr_skill_gap_nodes
    ADD CONSTRAINT bridgr_skill_gap_nodes_id_key UNIQUE (id);


--
-- Name: bridgr_skill_gap_nodes bridgr_skill_gap_nodes_pkey; Type: CONSTRAINT; Schema: hskip_users; Owner: bridgr
--

ALTER TABLE ONLY hskip_users.bridgr_skill_gap_nodes
    ADD CONSTRAINT bridgr_skill_gap_nodes_pkey PRIMARY KEY (uuid);


--
-- Name: bridgr_skill_gap_path_step_deps bridgr_skill_gap_path_step_de_path_uuid_step_uuid_depends_o_key; Type: CONSTRAINT; Schema: hskip_users; Owner: bridgr
--

ALTER TABLE ONLY hskip_users.bridgr_skill_gap_path_step_deps
    ADD CONSTRAINT bridgr_skill_gap_path_step_de_path_uuid_step_uuid_depends_o_key UNIQUE (path_uuid, step_uuid, depends_on_step_uuid);


--
-- Name: bridgr_skill_gap_path_step_deps bridgr_skill_gap_path_step_deps_id_key; Type: CONSTRAINT; Schema: hskip_users; Owner: bridgr
--

ALTER TABLE ONLY hskip_users.bridgr_skill_gap_path_step_deps
    ADD CONSTRAINT bridgr_skill_gap_path_step_deps_id_key UNIQUE (id);


--
-- Name: bridgr_skill_gap_path_step_deps bridgr_skill_gap_path_step_deps_pkey; Type: CONSTRAINT; Schema: hskip_users; Owner: bridgr
--

ALTER TABLE ONLY hskip_users.bridgr_skill_gap_path_step_deps
    ADD CONSTRAINT bridgr_skill_gap_path_step_deps_pkey PRIMARY KEY (uuid);


--
-- Name: bridgr_skill_gap_path_steps bridgr_skill_gap_path_steps_id_key; Type: CONSTRAINT; Schema: hskip_users; Owner: bridgr
--

ALTER TABLE ONLY hskip_users.bridgr_skill_gap_path_steps
    ADD CONSTRAINT bridgr_skill_gap_path_steps_id_key UNIQUE (id);


--
-- Name: bridgr_skill_gap_path_steps bridgr_skill_gap_path_steps_path_uuid_step_index_key; Type: CONSTRAINT; Schema: hskip_users; Owner: bridgr
--

ALTER TABLE ONLY hskip_users.bridgr_skill_gap_path_steps
    ADD CONSTRAINT bridgr_skill_gap_path_steps_path_uuid_step_index_key UNIQUE (path_uuid, step_index);


--
-- Name: bridgr_skill_gap_path_steps bridgr_skill_gap_path_steps_pkey; Type: CONSTRAINT; Schema: hskip_users; Owner: bridgr
--

ALTER TABLE ONLY hskip_users.bridgr_skill_gap_path_steps
    ADD CONSTRAINT bridgr_skill_gap_path_steps_pkey PRIMARY KEY (uuid);


--
-- Name: bridgr_skill_gap_graphs uq_bridgr_skill_gap_graphs_analysis_kind; Type: CONSTRAINT; Schema: hskip_users; Owner: bridgr
--

ALTER TABLE ONLY hskip_users.bridgr_skill_gap_graphs
    ADD CONSTRAINT uq_bridgr_skill_gap_graphs_analysis_kind UNIQUE (analysis_uuid, kind);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: bridgr
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: idx_bridgr_skill_gap_analyses_persona; Type: INDEX; Schema: hskip_users; Owner: bridgr
--

CREATE INDEX idx_bridgr_skill_gap_analyses_persona ON hskip_users.bridgr_skill_gap_analyses USING btree (founder_persona_uuid);


--
-- Name: idx_bridgr_skill_gap_analyses_pursuit; Type: INDEX; Schema: hskip_users; Owner: bridgr
--

CREATE INDEX idx_bridgr_skill_gap_analyses_pursuit ON hskip_users.bridgr_skill_gap_analyses USING btree (pursuit_uuid);


--
-- Name: idx_bridgr_skill_gap_analyses_status; Type: INDEX; Schema: hskip_users; Owner: bridgr
--

CREATE INDEX idx_bridgr_skill_gap_analyses_status ON hskip_users.bridgr_skill_gap_analyses USING btree (status);


--
-- Name: idx_bridgr_skill_gap_analyses_user_created; Type: INDEX; Schema: hskip_users; Owner: bridgr
--

CREATE INDEX idx_bridgr_skill_gap_analyses_user_created ON hskip_users.bridgr_skill_gap_analyses USING btree (user_id, created_at DESC);


--
-- Name: idx_bridgr_skill_gap_coverage_analysis_uuid; Type: INDEX; Schema: hskip_users; Owner: bridgr
--

CREATE INDEX idx_bridgr_skill_gap_coverage_analysis_uuid ON hskip_users.bridgr_skill_gap_coverage USING btree (analysis_uuid);


--
-- Name: idx_bridgr_skill_gap_coverage_kind; Type: INDEX; Schema: hskip_users; Owner: bridgr
--

CREATE INDEX idx_bridgr_skill_gap_coverage_kind ON hskip_users.bridgr_skill_gap_coverage USING btree (analysis_uuid, coverage_kind);


--
-- Name: idx_bridgr_skill_gap_edges_from; Type: INDEX; Schema: hskip_users; Owner: bridgr
--

CREATE INDEX idx_bridgr_skill_gap_edges_from ON hskip_users.bridgr_skill_gap_edges USING btree (from_node_uuid);


--
-- Name: idx_bridgr_skill_gap_edges_graph_uuid; Type: INDEX; Schema: hskip_users; Owner: bridgr
--

CREATE INDEX idx_bridgr_skill_gap_edges_graph_uuid ON hskip_users.bridgr_skill_gap_edges USING btree (graph_uuid);


--
-- Name: idx_bridgr_skill_gap_edges_relation; Type: INDEX; Schema: hskip_users; Owner: bridgr
--

CREATE INDEX idx_bridgr_skill_gap_edges_relation ON hskip_users.bridgr_skill_gap_edges USING btree (relation);


--
-- Name: idx_bridgr_skill_gap_edges_to; Type: INDEX; Schema: hskip_users; Owner: bridgr
--

CREATE INDEX idx_bridgr_skill_gap_edges_to ON hskip_users.bridgr_skill_gap_edges USING btree (to_node_uuid);


--
-- Name: idx_bridgr_skill_gap_graphs_analysis; Type: INDEX; Schema: hskip_users; Owner: bridgr
--

CREATE INDEX idx_bridgr_skill_gap_graphs_analysis ON hskip_users.bridgr_skill_gap_graphs USING btree (analysis_uuid);


--
-- Name: idx_bridgr_skill_gap_learning_paths_analysis_uuid; Type: INDEX; Schema: hskip_users; Owner: bridgr
--

CREATE INDEX idx_bridgr_skill_gap_learning_paths_analysis_uuid ON hskip_users.bridgr_skill_gap_learning_paths USING btree (analysis_uuid);


--
-- Name: idx_bridgr_skill_gap_nodes_graph_uuid; Type: INDEX; Schema: hskip_users; Owner: bridgr
--

CREATE INDEX idx_bridgr_skill_gap_nodes_graph_uuid ON hskip_users.bridgr_skill_gap_nodes USING btree (graph_uuid);


--
-- Name: idx_bridgr_skill_gap_nodes_node_key; Type: INDEX; Schema: hskip_users; Owner: bridgr
--

CREATE INDEX idx_bridgr_skill_gap_nodes_node_key ON hskip_users.bridgr_skill_gap_nodes USING btree (node_key);


--
-- Name: idx_bridgr_skill_gap_path_step_deps_depends; Type: INDEX; Schema: hskip_users; Owner: bridgr
--

CREATE INDEX idx_bridgr_skill_gap_path_step_deps_depends ON hskip_users.bridgr_skill_gap_path_step_deps USING btree (depends_on_step_uuid);


--
-- Name: idx_bridgr_skill_gap_path_step_deps_path; Type: INDEX; Schema: hskip_users; Owner: bridgr
--

CREATE INDEX idx_bridgr_skill_gap_path_step_deps_path ON hskip_users.bridgr_skill_gap_path_step_deps USING btree (path_uuid);


--
-- Name: idx_bridgr_skill_gap_path_step_deps_step; Type: INDEX; Schema: hskip_users; Owner: bridgr
--

CREATE INDEX idx_bridgr_skill_gap_path_step_deps_step ON hskip_users.bridgr_skill_gap_path_step_deps USING btree (step_uuid);


--
-- Name: idx_bridgr_skill_gap_path_steps_path_uuid; Type: INDEX; Schema: hskip_users; Owner: bridgr
--

CREATE INDEX idx_bridgr_skill_gap_path_steps_path_uuid ON hskip_users.bridgr_skill_gap_path_steps USING btree (path_uuid);


--
-- Name: bridgr_skill_gap_analyses tr_bridgr_skill_gap_analyses_control_time; Type: TRIGGER; Schema: hskip_users; Owner: bridgr
--

CREATE TRIGGER tr_bridgr_skill_gap_analyses_control_time BEFORE INSERT OR UPDATE ON hskip_users.bridgr_skill_gap_analyses FOR EACH ROW EXECUTE FUNCTION hskip_users.tr_control_time();


--
-- Name: bridgr_skill_gap_coverage tr_bridgr_skill_gap_coverage_control_time; Type: TRIGGER; Schema: hskip_users; Owner: bridgr
--

CREATE TRIGGER tr_bridgr_skill_gap_coverage_control_time BEFORE INSERT OR UPDATE ON hskip_users.bridgr_skill_gap_coverage FOR EACH ROW EXECUTE FUNCTION hskip_users.tr_control_time();


--
-- Name: bridgr_skill_gap_edges tr_bridgr_skill_gap_edges_control_time; Type: TRIGGER; Schema: hskip_users; Owner: bridgr
--

CREATE TRIGGER tr_bridgr_skill_gap_edges_control_time BEFORE INSERT OR UPDATE ON hskip_users.bridgr_skill_gap_edges FOR EACH ROW EXECUTE FUNCTION hskip_users.tr_control_time();


--
-- Name: bridgr_skill_gap_graphs tr_bridgr_skill_gap_graphs_control_time; Type: TRIGGER; Schema: hskip_users; Owner: bridgr
--

CREATE TRIGGER tr_bridgr_skill_gap_graphs_control_time BEFORE INSERT OR UPDATE ON hskip_users.bridgr_skill_gap_graphs FOR EACH ROW EXECUTE FUNCTION hskip_users.tr_control_time();


--
-- Name: bridgr_skill_gap_learning_paths tr_bridgr_skill_gap_learning_paths_control_time; Type: TRIGGER; Schema: hskip_users; Owner: bridgr
--

CREATE TRIGGER tr_bridgr_skill_gap_learning_paths_control_time BEFORE INSERT OR UPDATE ON hskip_users.bridgr_skill_gap_learning_paths FOR EACH ROW EXECUTE FUNCTION hskip_users.tr_control_time();


--
-- Name: bridgr_skill_gap_nodes tr_bridgr_skill_gap_nodes_control_time; Type: TRIGGER; Schema: hskip_users; Owner: bridgr
--

CREATE TRIGGER tr_bridgr_skill_gap_nodes_control_time BEFORE INSERT OR UPDATE ON hskip_users.bridgr_skill_gap_nodes FOR EACH ROW EXECUTE FUNCTION hskip_users.tr_control_time();


--
-- Name: bridgr_skill_gap_path_step_deps tr_bridgr_skill_gap_path_step_deps_control_time; Type: TRIGGER; Schema: hskip_users; Owner: bridgr
--

CREATE TRIGGER tr_bridgr_skill_gap_path_step_deps_control_time BEFORE INSERT OR UPDATE ON hskip_users.bridgr_skill_gap_path_step_deps FOR EACH ROW EXECUTE FUNCTION hskip_users.tr_control_time();


--
-- Name: bridgr_skill_gap_path_steps tr_bridgr_skill_gap_path_steps_control_time; Type: TRIGGER; Schema: hskip_users; Owner: bridgr
--

CREATE TRIGGER tr_bridgr_skill_gap_path_steps_control_time BEFORE INSERT OR UPDATE ON hskip_users.bridgr_skill_gap_path_steps FOR EACH ROW EXECUTE FUNCTION hskip_users.tr_control_time();


--
-- Name: SCHEMA hskip_users; Type: ACL; Schema: -; Owner: bridgr
--

GRANT USAGE ON SCHEMA hskip_users TO hassle_skip_readonly_role;
GRANT USAGE ON SCHEMA hskip_users TO hassle_skip_app_role;
GRANT ALL ON SCHEMA hskip_users TO hassle_skip_migration_role;


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: pg_database_owner
--

GRANT USAGE ON SCHEMA public TO hassle_skip_readonly_role;
GRANT USAGE ON SCHEMA public TO hassle_skip_app_role;
GRANT USAGE ON SCHEMA public TO hassle_skip_migration_role;


--
-- Name: TABLE schema_migrations; Type: ACL; Schema: public; Owner: bridgr
--

GRANT SELECT ON TABLE public.schema_migrations TO hassle_skip_readonly_role;
GRANT SELECT ON TABLE public.schema_migrations TO hassle_skip_app_role;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.schema_migrations TO hassle_skip_migration_role;


--
-- Name: DEFAULT PRIVILEGES FOR SEQUENCES; Type: DEFAULT ACL; Schema: hskip_users; Owner: hassle_skip_migration
--

ALTER DEFAULT PRIVILEGES FOR ROLE hassle_skip_migration IN SCHEMA hskip_users GRANT SELECT,USAGE ON SEQUENCES TO hassle_skip_readonly_role;
ALTER DEFAULT PRIVILEGES FOR ROLE hassle_skip_migration IN SCHEMA hskip_users GRANT SELECT,USAGE ON SEQUENCES TO hassle_skip_app_role;


--
-- Name: DEFAULT PRIVILEGES FOR TABLES; Type: DEFAULT ACL; Schema: hskip_users; Owner: hassle_skip_migration
--

ALTER DEFAULT PRIVILEGES FOR ROLE hassle_skip_migration IN SCHEMA hskip_users GRANT SELECT ON TABLES TO hassle_skip_readonly_role;
ALTER DEFAULT PRIVILEGES FOR ROLE hassle_skip_migration IN SCHEMA hskip_users GRANT SELECT,INSERT,DELETE,UPDATE ON TABLES TO hassle_skip_app_role;


--
-- Name: DEFAULT PRIVILEGES FOR SEQUENCES; Type: DEFAULT ACL; Schema: public; Owner: hassle_skip_migration
--

ALTER DEFAULT PRIVILEGES FOR ROLE hassle_skip_migration IN SCHEMA public GRANT SELECT,USAGE ON SEQUENCES TO hassle_skip_readonly_role;
ALTER DEFAULT PRIVILEGES FOR ROLE hassle_skip_migration IN SCHEMA public GRANT SELECT,USAGE ON SEQUENCES TO hassle_skip_app_role;


--
-- Name: DEFAULT PRIVILEGES FOR TABLES; Type: DEFAULT ACL; Schema: public; Owner: hassle_skip_migration
--

ALTER DEFAULT PRIVILEGES FOR ROLE hassle_skip_migration IN SCHEMA public GRANT SELECT ON TABLES TO hassle_skip_readonly_role;
ALTER DEFAULT PRIVILEGES FOR ROLE hassle_skip_migration IN SCHEMA public GRANT SELECT ON TABLES TO hassle_skip_app_role;


--
-- PostgreSQL database dump complete
--


