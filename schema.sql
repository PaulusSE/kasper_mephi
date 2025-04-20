--
-- PostgreSQL database dump
--

-- Dumped from database version 17.4 (Debian 17.4-1.pgdg120+2)
-- Dumped by pg_dump version 17.4 (Debian 17.4-1.pgdg120+2)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: public; Type: SCHEMA; Schema: -; Owner: postgres
--

-- *not* creating schema, since initdb creates it


ALTER SCHEMA public OWNER TO postgres;

--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


--
-- Name: approval_status; Type: TYPE; Schema: public; Owner: root
--

CREATE TYPE public.approval_status AS ENUM (
    'todo',
    'approved',
    'on review',
    'in progress',
    'empty',
    'failed'
);


ALTER TYPE public.approval_status OWNER TO root;

--
-- Name: classroom_load_type; Type: TYPE; Schema: public; Owner: root
--

CREATE TYPE public.classroom_load_type AS ENUM (
    'practice',
    'lectures',
    'laboratory',
    'exam'
);


ALTER TYPE public.classroom_load_type OWNER TO root;

--
-- Name: conference_status; Type: TYPE; Schema: public; Owner: root
--

CREATE TYPE public.conference_status AS ENUM (
    'registered',
    'performed'
);


ALTER TYPE public.conference_status OWNER TO root;

--
-- Name: individual_students_load_type; Type: TYPE; Schema: public; Owner: root
--

CREATE TYPE public.individual_students_load_type AS ENUM (
    'project practice',
    'bachelor',
    'masters'
);


ALTER TYPE public.individual_students_load_type OWNER TO root;

--
-- Name: patent_type; Type: TYPE; Schema: public; Owner: root
--

CREATE TYPE public.patent_type AS ENUM (
    'software',
    'database'
);


ALTER TYPE public.patent_type OWNER TO root;

--
-- Name: progress_type; Type: TYPE; Schema: public; Owner: root
--

CREATE TYPE public.progress_type AS ENUM (
    'intro',
    'ch. 1',
    'ch. 2',
    'ch. 3',
    'ch. 4',
    'ch. 5',
    'ch. 6',
    'end',
    'literature',
    'abstract'
);


ALTER TYPE public.progress_type OWNER TO root;

--
-- Name: publication_status; Type: TYPE; Schema: public; Owner: root
--

CREATE TYPE public.publication_status AS ENUM (
    'to print',
    'published',
    'in progress'
);


ALTER TYPE public.publication_status OWNER TO root;

--
-- Name: request_status; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.request_status AS ENUM (
    'pending',
    'approved',
    'rejected'
);


ALTER TYPE public.request_status OWNER TO postgres;

--
-- Name: student_status; Type: TYPE; Schema: public; Owner: root
--

CREATE TYPE public.student_status AS ENUM (
    'academic',
    'graduated',
    'studying',
    'expelled'
);


ALTER TYPE public.student_status OWNER TO root;

--
-- Name: user_type; Type: TYPE; Schema: public; Owner: root
--

CREATE TYPE public.user_type AS ENUM (
    'admin',
    'student',
    'supervisor'
);


ALTER TYPE public.user_type OWNER TO root;

--
-- Name: show_sample_data(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.show_sample_data() RETURNS TABLE(table_name text, data jsonb)
    LANGUAGE plpgsql
    AS $$
DECLARE 
    tbl text;
BEGIN
    FOR tbl IN
        SELECT t.table_name  -- Явное указание алиаса таблицы
        FROM information_schema.tables t  -- Добавлен алиас "t"
        WHERE t.table_schema = 'public'
    LOOP
        RETURN QUERY EXECUTE 
            format('SELECT %L::text, jsonb_agg(tab) FROM (SELECT * FROM %I LIMIT 2) tab', tbl, tbl); -- Изменен алиас внутренней таблицы
    END LOOP;
END;
$$;


ALTER FUNCTION public.show_sample_data() OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: additional_load; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.additional_load (
    load_id uuid NOT NULL,
    t_load_id uuid NOT NULL,
    name character varying(512) NOT NULL,
    volume character varying(256),
    comment character varying(1024)
);


ALTER TABLE public.additional_load OWNER TO root;

--
-- Name: authorization_token; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.authorization_token (
    token_id integer NOT NULL,
    user_id uuid,
    is_active boolean DEFAULT false NOT NULL,
    token_number text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now()
);


ALTER TABLE public.authorization_token OWNER TO root;

--
-- Name: authorization_token_token_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.authorization_token_token_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.authorization_token_token_id_seq OWNER TO root;

--
-- Name: authorization_token_token_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.authorization_token_token_id_seq OWNED BY public.authorization_token.token_id;


--
-- Name: classroom_load; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.classroom_load (
    load_id uuid NOT NULL,
    t_load_id uuid NOT NULL,
    hours integer NOT NULL,
    load_type public.classroom_load_type NOT NULL,
    main_teacher character varying(256) NOT NULL,
    group_name character varying(50) NOT NULL,
    subject_name character varying(256) NOT NULL
);


ALTER TABLE public.classroom_load OWNER TO root;

--
-- Name: conferences; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.conferences (
    conference_id uuid NOT NULL,
    works_id uuid NOT NULL,
    status public.conference_status NOT NULL,
    scopus boolean DEFAULT false NOT NULL,
    rinc boolean DEFAULT false NOT NULL,
    wac boolean DEFAULT false NOT NULL,
    wos boolean DEFAULT false NOT NULL,
    conference_name character varying(512) NOT NULL,
    report_name character varying(512) NOT NULL,
    location character varying(512) NOT NULL,
    reported_at timestamp with time zone NOT NULL
);


ALTER TABLE public.conferences OWNER TO root;

--
-- Name: dissertation_commentary; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.dissertation_commentary (
    commentary_id uuid NOT NULL,
    student_id uuid NOT NULL,
    semester integer NOT NULL,
    commentary text
);


ALTER TABLE public.dissertation_commentary OWNER TO root;

--
-- Name: dissertation_plans; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.dissertation_plans (
    plan_id uuid NOT NULL,
    student_id uuid NOT NULL,
    semester integer NOT NULL,
    plan_text text
);


ALTER TABLE public.dissertation_plans OWNER TO root;

--
-- Name: dissertation_titles; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.dissertation_titles (
    title_id uuid NOT NULL,
    student_id uuid NOT NULL,
    title character varying(512) NOT NULL,
    created_at timestamp with time zone NOT NULL,
    status public.approval_status DEFAULT 'empty'::public.approval_status NOT NULL,
    accepted_at timestamp with time zone,
    semester integer NOT NULL,
    research_object character varying(512) NOT NULL,
    research_subject character varying(512) NOT NULL
);


ALTER TABLE public.dissertation_titles OWNER TO root;

--
-- Name: dissertations; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.dissertations (
    dissertation_id uuid NOT NULL,
    student_id uuid NOT NULL,
    status public.approval_status DEFAULT 'empty'::public.approval_status NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone DEFAULT now(),
    semester integer NOT NULL,
    file_name text
);


ALTER TABLE public.dissertations OWNER TO root;

--
-- Name: exam_type; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.exam_type (
    type_id integer NOT NULL,
    exam_name text NOT NULL
);


ALTER TABLE public.exam_type OWNER TO root;

--
-- Name: exams; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.exams (
    exam_id uuid NOT NULL,
    student_id uuid NOT NULL,
    exam_type integer NOT NULL,
    semester integer NOT NULL,
    mark integer NOT NULL,
    set_at timestamp with time zone
);


ALTER TABLE public.exams OWNER TO root;

--
-- Name: feedback; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.feedback (
    feedback_id uuid NOT NULL,
    student_id uuid NOT NULL,
    dissertation_id uuid NOT NULL,
    feedback text,
    semester integer NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.feedback OWNER TO root;

--
-- Name: goose_db_version; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.goose_db_version (
    id integer NOT NULL,
    version_id bigint NOT NULL,
    is_applied boolean NOT NULL,
    tstamp timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.goose_db_version OWNER TO postgres;

--
-- Name: goose_db_version_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.goose_db_version ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.goose_db_version_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: groups; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.groups (
    group_id integer NOT NULL,
    group_name character varying(15) NOT NULL,
    archived boolean DEFAULT false NOT NULL
);


ALTER TABLE public.groups OWNER TO root;

--
-- Name: groups_group_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.groups_group_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.groups_group_id_seq OWNER TO root;

--
-- Name: groups_group_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.groups_group_id_seq OWNED BY public.groups.group_id;


--
-- Name: individual_students_load; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.individual_students_load (
    load_id uuid NOT NULL,
    t_load_id uuid NOT NULL,
    load_type public.individual_students_load_type NOT NULL,
    students_amount integer NOT NULL,
    comment character varying(1024)
);


ALTER TABLE public.individual_students_load OWNER TO root;

--
-- Name: marks; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.marks (
    student_id uuid NOT NULL,
    mark integer NOT NULL,
    semester integer NOT NULL
);


ALTER TABLE public.marks OWNER TO root;

--
-- Name: patents; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.patents (
    patent_id uuid NOT NULL,
    works_id uuid NOT NULL,
    name character varying(256) NOT NULL,
    registration_date timestamp with time zone NOT NULL,
    type public.patent_type NOT NULL,
    add_info text
);


ALTER TABLE public.patents OWNER TO root;

--
-- Name: progressiveness; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.progressiveness (
    progress_id uuid NOT NULL,
    student_id uuid NOT NULL,
    semester integer NOT NULL,
    progressiveness integer DEFAULT 0 NOT NULL
);


ALTER TABLE public.progressiveness OWNER TO root;

--
-- Name: publications; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.publications (
    publication_id uuid NOT NULL,
    works_id uuid NOT NULL,
    name character varying(256) NOT NULL,
    scopus boolean DEFAULT false NOT NULL,
    rinc boolean DEFAULT false NOT NULL,
    wac boolean DEFAULT false NOT NULL,
    wos boolean DEFAULT false NOT NULL,
    impact double precision NOT NULL,
    status public.publication_status NOT NULL,
    output_data text,
    co_authors character varying(256),
    volume integer
);


ALTER TABLE public.publications OWNER TO root;

--
-- Name: registration_requests; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.registration_requests (
    request_id uuid DEFAULT gen_random_uuid() NOT NULL,
    email character varying(128) NOT NULL,
    password_hash text NOT NULL,
    user_type public.user_type NOT NULL,
    payload jsonb NOT NULL,
    status public.request_status DEFAULT 'pending'::public.request_status NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    processed_by uuid,
    processed_at timestamp with time zone
);


ALTER TABLE public.registration_requests OWNER TO postgres;

--
-- Name: research_projects; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.research_projects (
    project_id uuid NOT NULL,
    works_id uuid NOT NULL,
    project_name character varying(512) NOT NULL,
    start_at timestamp with time zone NOT NULL,
    end_at timestamp with time zone NOT NULL,
    add_info character varying(1024),
    grantee character varying(1024)
);


ALTER TABLE public.research_projects OWNER TO root;

--
-- Name: scientific_works_status; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.scientific_works_status (
    works_id uuid NOT NULL,
    student_id uuid NOT NULL,
    semester integer NOT NULL,
    status public.approval_status DEFAULT 'empty'::public.approval_status NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    accepted_at timestamp with time zone
);


ALTER TABLE public.scientific_works_status OWNER TO root;

--
-- Name: semester_count; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.semester_count (
    count_id uuid NOT NULL,
    amount integer NOT NULL
);


ALTER TABLE public.semester_count OWNER TO root;

--
-- Name: semester_progress; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.semester_progress (
    progress_id uuid NOT NULL,
    student_id uuid NOT NULL,
    progress_type public.progress_type NOT NULL,
    first boolean DEFAULT false NOT NULL,
    second boolean DEFAULT false NOT NULL,
    third boolean DEFAULT false NOT NULL,
    forth boolean DEFAULT false NOT NULL,
    fifth boolean DEFAULT false NOT NULL,
    sixth boolean DEFAULT false NOT NULL,
    seventh boolean DEFAULT false NOT NULL,
    eighth boolean DEFAULT false NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    status public.approval_status DEFAULT 'empty'::public.approval_status NOT NULL,
    accepted_at timestamp with time zone
);


ALTER TABLE public.semester_progress OWNER TO root;

--
-- Name: specializations; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.specializations (
    spec_id integer NOT NULL,
    title text NOT NULL,
    archived boolean DEFAULT false NOT NULL
);


ALTER TABLE public.specializations OWNER TO root;

--
-- Name: specializations_spec_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.specializations_spec_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.specializations_spec_id_seq OWNER TO root;

--
-- Name: specializations_spec_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.specializations_spec_id_seq OWNED BY public.specializations.spec_id;


--
-- Name: students; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.students (
    student_id uuid NOT NULL,
    user_id uuid NOT NULL,
    full_name character varying(256) NOT NULL,
    spec_id integer NOT NULL,
    actual_semester integer DEFAULT 1 NOT NULL,
    years integer NOT NULL,
    start_date timestamp with time zone NOT NULL,
    studying_status public.student_status DEFAULT 'studying'::public.student_status NOT NULL,
    group_id integer NOT NULL,
    status public.approval_status DEFAULT 'empty'::public.approval_status NOT NULL,
    can_edit boolean DEFAULT false NOT NULL,
    phone character varying(20) DEFAULT ''::character varying NOT NULL,
    category character varying(30) NOT NULL,
    end_date timestamp with time zone NOT NULL
);


ALTER TABLE public.students OWNER TO root;

--
-- Name: students_commentary; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.students_commentary (
    commentary_id uuid NOT NULL,
    student_id uuid NOT NULL,
    semester integer NOT NULL,
    commentary text,
    commented_at timestamp with time zone NOT NULL
);


ALTER TABLE public.students_commentary OWNER TO root;

--
-- Name: students_supervisors; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.students_supervisors (
    id uuid NOT NULL,
    student_id uuid NOT NULL,
    supervisor_id uuid NOT NULL,
    start_at timestamp with time zone DEFAULT now() NOT NULL,
    end_at timestamp with time zone
);


ALTER TABLE public.students_supervisors OWNER TO root;

--
-- Name: supervisor_marks; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.supervisor_marks (
    mark_id uuid NOT NULL,
    student_id uuid NOT NULL,
    mark integer NOT NULL,
    semester integer NOT NULL,
    supervisor_id uuid NOT NULL
);


ALTER TABLE public.supervisor_marks OWNER TO root;

--
-- Name: supervisors; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.supervisors (
    supervisor_id uuid NOT NULL,
    user_id uuid NOT NULL,
    full_name character varying(256),
    phone character varying(20) DEFAULT ''::character varying NOT NULL,
    archived boolean DEFAULT false NOT NULL,
    faculty character varying(512),
    department character varying(512),
    degree character varying(512),
    rank character varying(512),
    "position" character varying(512)
);


ALTER TABLE public.supervisors OWNER TO root;

--
-- Name: teaching_load_status; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.teaching_load_status (
    loads_id uuid NOT NULL,
    student_id uuid NOT NULL,
    semester integer NOT NULL,
    status public.approval_status DEFAULT 'empty'::public.approval_status NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    accepted_at timestamp with time zone
);


ALTER TABLE public.teaching_load_status OWNER TO root;

--
-- Name: users; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.users (
    user_id uuid NOT NULL,
    email character varying(128) NOT NULL,
    password text NOT NULL,
    kasper_id uuid NOT NULL,
    user_type public.user_type NOT NULL,
    registered boolean DEFAULT false NOT NULL
);


ALTER TABLE public.users OWNER TO root;

--
-- Name: authorization_token token_id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authorization_token ALTER COLUMN token_id SET DEFAULT nextval('public.authorization_token_token_id_seq'::regclass);


--
-- Name: groups group_id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.groups ALTER COLUMN group_id SET DEFAULT nextval('public.groups_group_id_seq'::regclass);


--
-- Name: specializations spec_id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.specializations ALTER COLUMN spec_id SET DEFAULT nextval('public.specializations_spec_id_seq'::regclass);


--
-- Name: additional_load additional_load_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.additional_load
    ADD CONSTRAINT additional_load_pkey PRIMARY KEY (load_id);


--
-- Name: authorization_token authorization_token_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authorization_token
    ADD CONSTRAINT authorization_token_pkey PRIMARY KEY (token_id);


--
-- Name: authorization_token authorization_token_token_number_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authorization_token
    ADD CONSTRAINT authorization_token_token_number_key UNIQUE (token_number);


--
-- Name: classroom_load classroom_load_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.classroom_load
    ADD CONSTRAINT classroom_load_pkey PRIMARY KEY (load_id);


--
-- Name: conferences conferences_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.conferences
    ADD CONSTRAINT conferences_pkey PRIMARY KEY (conference_id);


--
-- Name: dissertation_commentary dissertation_commentary_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.dissertation_commentary
    ADD CONSTRAINT dissertation_commentary_pkey PRIMARY KEY (commentary_id);


--
-- Name: dissertation_commentary dissertation_commentary_student_id_semester_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.dissertation_commentary
    ADD CONSTRAINT dissertation_commentary_student_id_semester_key UNIQUE (student_id, semester);


--
-- Name: dissertation_plans dissertation_plans_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.dissertation_plans
    ADD CONSTRAINT dissertation_plans_pkey PRIMARY KEY (plan_id);


--
-- Name: dissertation_plans dissertation_plans_student_id_semester_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.dissertation_plans
    ADD CONSTRAINT dissertation_plans_student_id_semester_key UNIQUE (student_id, semester);


--
-- Name: dissertation_titles dissertation_titles_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.dissertation_titles
    ADD CONSTRAINT dissertation_titles_pkey PRIMARY KEY (title_id);


--
-- Name: dissertations dissertations_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.dissertations
    ADD CONSTRAINT dissertations_pkey PRIMARY KEY (dissertation_id);


--
-- Name: dissertations dissertations_student_id_semester_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.dissertations
    ADD CONSTRAINT dissertations_student_id_semester_key UNIQUE (student_id, semester);


--
-- Name: exam_type exam_type_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.exam_type
    ADD CONSTRAINT exam_type_pkey PRIMARY KEY (type_id);


--
-- Name: exams exams_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.exams
    ADD CONSTRAINT exams_pkey PRIMARY KEY (exam_id);


--
-- Name: exams exams_student_id_semester_exam_type_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.exams
    ADD CONSTRAINT exams_student_id_semester_exam_type_key UNIQUE (student_id, semester, exam_type);


--
-- Name: feedback feedback_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.feedback
    ADD CONSTRAINT feedback_pkey PRIMARY KEY (feedback_id);


--
-- Name: feedback feedback_student_id_semester_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.feedback
    ADD CONSTRAINT feedback_student_id_semester_key UNIQUE (student_id, semester);


--
-- Name: goose_db_version goose_db_version_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.goose_db_version
    ADD CONSTRAINT goose_db_version_pkey PRIMARY KEY (id);


--
-- Name: groups groups_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.groups
    ADD CONSTRAINT groups_pkey PRIMARY KEY (group_id);


--
-- Name: individual_students_load individual_students_load_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.individual_students_load
    ADD CONSTRAINT individual_students_load_pkey PRIMARY KEY (load_id);


--
-- Name: marks marks_student_id_semester_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.marks
    ADD CONSTRAINT marks_student_id_semester_key UNIQUE (student_id, semester);


--
-- Name: patents patents_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.patents
    ADD CONSTRAINT patents_pkey PRIMARY KEY (patent_id);


--
-- Name: progressiveness progressiveness_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.progressiveness
    ADD CONSTRAINT progressiveness_pkey PRIMARY KEY (progress_id);


--
-- Name: progressiveness progressiveness_student_id_semester_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.progressiveness
    ADD CONSTRAINT progressiveness_student_id_semester_key UNIQUE (student_id, semester);


--
-- Name: publications publications_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.publications
    ADD CONSTRAINT publications_pkey PRIMARY KEY (publication_id);


--
-- Name: registration_requests registration_requests_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.registration_requests
    ADD CONSTRAINT registration_requests_email_key UNIQUE (email);


--
-- Name: registration_requests registration_requests_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.registration_requests
    ADD CONSTRAINT registration_requests_pkey PRIMARY KEY (request_id);


--
-- Name: research_projects research_projects_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.research_projects
    ADD CONSTRAINT research_projects_pkey PRIMARY KEY (project_id);


--
-- Name: scientific_works_status scientific_works_status_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.scientific_works_status
    ADD CONSTRAINT scientific_works_status_pkey PRIMARY KEY (works_id);


--
-- Name: semester_count semester_count_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.semester_count
    ADD CONSTRAINT semester_count_pkey PRIMARY KEY (count_id);


--
-- Name: semester_progress semester_progress_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.semester_progress
    ADD CONSTRAINT semester_progress_pkey PRIMARY KEY (progress_id);


--
-- Name: semester_progress semester_progress_student_id_progress_type_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.semester_progress
    ADD CONSTRAINT semester_progress_student_id_progress_type_key UNIQUE (student_id, progress_type);


--
-- Name: specializations specializations_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.specializations
    ADD CONSTRAINT specializations_pkey PRIMARY KEY (spec_id);


--
-- Name: students_commentary students_commentary_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.students_commentary
    ADD CONSTRAINT students_commentary_pkey PRIMARY KEY (commentary_id);


--
-- Name: students_commentary students_commentary_student_id_semester_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.students_commentary
    ADD CONSTRAINT students_commentary_student_id_semester_key UNIQUE (student_id, semester);


--
-- Name: students students_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.students
    ADD CONSTRAINT students_pkey PRIMARY KEY (student_id);


--
-- Name: students_supervisors students_supervisors_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.students_supervisors
    ADD CONSTRAINT students_supervisors_pkey PRIMARY KEY (id);


--
-- Name: students students_user_id_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.students
    ADD CONSTRAINT students_user_id_key UNIQUE (user_id);


--
-- Name: supervisor_marks supervisor_marks_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.supervisor_marks
    ADD CONSTRAINT supervisor_marks_pkey PRIMARY KEY (mark_id);


--
-- Name: supervisor_marks supervisor_marks_student_id_semester_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.supervisor_marks
    ADD CONSTRAINT supervisor_marks_student_id_semester_key UNIQUE (student_id, semester);


--
-- Name: supervisors supervisors_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.supervisors
    ADD CONSTRAINT supervisors_pkey PRIMARY KEY (supervisor_id);


--
-- Name: supervisors supervisors_user_id_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.supervisors
    ADD CONSTRAINT supervisors_user_id_key UNIQUE (user_id);


--
-- Name: teaching_load_status teaching_load_status_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.teaching_load_status
    ADD CONSTRAINT teaching_load_status_pkey PRIMARY KEY (loads_id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_kasper_id_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_kasper_id_key UNIQUE (kasper_id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (user_id);


--
-- Name: users_email_idx; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX users_email_idx ON public.users USING btree (email);


--
-- Name: users_email_idx1; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX users_email_idx1 ON public.users USING btree (email);


--
-- Name: users_email_idx2; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX users_email_idx2 ON public.users USING btree (email);


--
-- Name: additional_load additional_load_t_load_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.additional_load
    ADD CONSTRAINT additional_load_t_load_id_fkey FOREIGN KEY (t_load_id) REFERENCES public.teaching_load_status(loads_id) ON DELETE CASCADE;


--
-- Name: authorization_token authorization_token_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.authorization_token
    ADD CONSTRAINT authorization_token_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(user_id) ON DELETE CASCADE;


--
-- Name: classroom_load classroom_load_t_load_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.classroom_load
    ADD CONSTRAINT classroom_load_t_load_id_fkey FOREIGN KEY (t_load_id) REFERENCES public.teaching_load_status(loads_id) ON DELETE CASCADE;


--
-- Name: conferences conferences_works_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.conferences
    ADD CONSTRAINT conferences_works_id_fkey FOREIGN KEY (works_id) REFERENCES public.scientific_works_status(works_id) ON DELETE CASCADE;


--
-- Name: dissertation_commentary dissertation_commentary_student_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.dissertation_commentary
    ADD CONSTRAINT dissertation_commentary_student_id_fkey FOREIGN KEY (student_id) REFERENCES public.students(student_id) ON DELETE CASCADE;


--
-- Name: dissertation_plans dissertation_plans_student_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.dissertation_plans
    ADD CONSTRAINT dissertation_plans_student_id_fkey FOREIGN KEY (student_id) REFERENCES public.students(student_id) ON DELETE CASCADE;


--
-- Name: dissertation_titles dissertation_titles_student_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.dissertation_titles
    ADD CONSTRAINT dissertation_titles_student_id_fkey FOREIGN KEY (student_id) REFERENCES public.students(student_id) ON DELETE CASCADE;


--
-- Name: dissertations dissertations_student_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.dissertations
    ADD CONSTRAINT dissertations_student_id_fkey FOREIGN KEY (student_id) REFERENCES public.students(student_id) ON DELETE CASCADE;


--
-- Name: exams exams_exam_type_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.exams
    ADD CONSTRAINT exams_exam_type_fkey FOREIGN KEY (exam_type) REFERENCES public.exam_type(type_id);


--
-- Name: exams exams_student_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.exams
    ADD CONSTRAINT exams_student_id_fkey FOREIGN KEY (student_id) REFERENCES public.students(student_id) ON DELETE CASCADE;


--
-- Name: feedback feedback_dissertation_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.feedback
    ADD CONSTRAINT feedback_dissertation_id_fkey FOREIGN KEY (dissertation_id) REFERENCES public.dissertations(dissertation_id) ON DELETE CASCADE;


--
-- Name: feedback feedback_student_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.feedback
    ADD CONSTRAINT feedback_student_id_fkey FOREIGN KEY (student_id) REFERENCES public.students(student_id) ON DELETE CASCADE;


--
-- Name: individual_students_load individual_students_load_t_load_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.individual_students_load
    ADD CONSTRAINT individual_students_load_t_load_id_fkey FOREIGN KEY (t_load_id) REFERENCES public.teaching_load_status(loads_id) ON DELETE CASCADE;


--
-- Name: marks marks_student_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.marks
    ADD CONSTRAINT marks_student_id_fkey FOREIGN KEY (student_id) REFERENCES public.students(student_id) ON DELETE CASCADE;


--
-- Name: patents patents_works_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.patents
    ADD CONSTRAINT patents_works_id_fkey FOREIGN KEY (works_id) REFERENCES public.scientific_works_status(works_id) ON DELETE CASCADE;


--
-- Name: progressiveness progressiveness_student_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.progressiveness
    ADD CONSTRAINT progressiveness_student_id_fkey FOREIGN KEY (student_id) REFERENCES public.students(student_id) ON DELETE CASCADE;


--
-- Name: publications publications_works_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.publications
    ADD CONSTRAINT publications_works_id_fkey FOREIGN KEY (works_id) REFERENCES public.scientific_works_status(works_id) ON DELETE CASCADE;


--
-- Name: research_projects research_projects_works_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.research_projects
    ADD CONSTRAINT research_projects_works_id_fkey FOREIGN KEY (works_id) REFERENCES public.scientific_works_status(works_id) ON DELETE CASCADE;


--
-- Name: scientific_works_status scientific_works_status_student_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.scientific_works_status
    ADD CONSTRAINT scientific_works_status_student_id_fkey FOREIGN KEY (student_id) REFERENCES public.students(student_id) ON DELETE CASCADE;


--
-- Name: semester_progress semester_progress_student_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.semester_progress
    ADD CONSTRAINT semester_progress_student_id_fkey FOREIGN KEY (student_id) REFERENCES public.students(student_id) ON DELETE CASCADE;


--
-- Name: students_commentary students_commentary_student_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.students_commentary
    ADD CONSTRAINT students_commentary_student_id_fkey FOREIGN KEY (student_id) REFERENCES public.students(student_id) ON DELETE CASCADE;


--
-- Name: students students_group_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.students
    ADD CONSTRAINT students_group_id_fkey FOREIGN KEY (group_id) REFERENCES public.groups(group_id);


--
-- Name: students students_spec_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.students
    ADD CONSTRAINT students_spec_id_fkey FOREIGN KEY (spec_id) REFERENCES public.specializations(spec_id);


--
-- Name: students_supervisors students_supervisors_student_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.students_supervisors
    ADD CONSTRAINT students_supervisors_student_id_fkey FOREIGN KEY (student_id) REFERENCES public.students(student_id) ON DELETE CASCADE;


--
-- Name: students_supervisors students_supervisors_supervisor_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.students_supervisors
    ADD CONSTRAINT students_supervisors_supervisor_id_fkey FOREIGN KEY (supervisor_id) REFERENCES public.supervisors(supervisor_id) ON DELETE CASCADE;


--
-- Name: students students_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.students
    ADD CONSTRAINT students_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(user_id) ON DELETE CASCADE;


--
-- Name: supervisor_marks supervisor_marks_student_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.supervisor_marks
    ADD CONSTRAINT supervisor_marks_student_id_fkey FOREIGN KEY (student_id) REFERENCES public.students(student_id) ON DELETE CASCADE;


--
-- Name: supervisor_marks supervisor_marks_supervisor_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.supervisor_marks
    ADD CONSTRAINT supervisor_marks_supervisor_id_fkey FOREIGN KEY (supervisor_id) REFERENCES public.supervisors(supervisor_id) ON DELETE CASCADE;


--
-- Name: supervisors supervisors_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.supervisors
    ADD CONSTRAINT supervisors_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(user_id) ON DELETE CASCADE;


--
-- Name: teaching_load_status teaching_load_status_student_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.teaching_load_status
    ADD CONSTRAINT teaching_load_status_student_id_fkey FOREIGN KEY (student_id) REFERENCES public.students(student_id) ON DELETE CASCADE;


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE USAGE ON SCHEMA public FROM PUBLIC;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--

