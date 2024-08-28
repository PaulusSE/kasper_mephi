create type user_type as enum ('admin', 'student', 'supervisor');

create type approval_status as enum ('todo', 'approved', 'on review', 'in progress', 'empty', 'failed');

create table users
(
    user_id    uuid primary key,
    email      varchar(128) not null unique,
    password   text         not null,
    kasper_id  uuid unique  not null,
    user_type  user_type    not null,
    registered bool         not null default false
);

create table authorization_token
(
    token_id     serial primary key,
    user_id      uuid references users (user_id) on delete cascade not null,
    is_active    bool                                              not null default false,
    token_number text                                              not null unique,
    created_at   timestamptz                                       not null,
    updated_at   timestamptz                                                default now()
);

create index on users (email);

create table supervisors
(
    supervisor_id uuid primary key,
    user_id       uuid references users (user_id) on delete cascade unique not null,
    full_name     varchar(256),
    phone         varchar(20)                                              not null default '',
    archived      bool                                                     not null default false,
    faculty       varchar(512),
    department    varchar(512),
    degree        varchar(512),
    rank          varchar(512),
    position      varchar(512)
);

create table groups
(
    group_id   serial primary key,
    group_name varchar(15) not null,
    archived   bool        not null default false
);

create table specializations
(
    spec_id  serial primary key,
    title    text not null,
    archived bool not null default false
);

create type student_status as enum ('academic', 'graduated', 'studying', 'expelled');

create table students
(
    student_id      uuid primary key,
    user_id         uuid references users (user_id) on delete cascade not null unique,
    full_name       varchar(256)                                      not null,
--     department      varchar(256)                                      not null,
    spec_id         int references specializations (spec_id)          not null,
    actual_semester int                                               not null default 1,
    years           int                                               not null,
    start_date      timestamptz                                       not null,
    studying_status student_status                                    not null default 'studying',
    group_id        int references groups (group_id)                  not null,
    status          approval_status                                   not null default 'empty',
    can_edit        bool                                              not null default false,
    phone           varchar(20)                                       not null default '',
    category        varchar(30)                                       not null,
    end_date        timestamptz                                       not null
);

create table progressiveness
(
    progress_id     uuid primary key,
    student_id      uuid references students (student_id) on delete cascade not null,
    semester        int                                                     not null,
    progressiveness int                                                     not null default 0,
    unique (student_id, semester)
);

create table marks
(
    student_id uuid references students (student_id) on delete cascade not null,
    mark       int                                                     not null,
    semester   int                                                     not null,
    unique (student_id, semester)
);

create table students_supervisors
(
    id            uuid primary key,
    student_id    uuid references students (student_id) on delete cascade       not null,
    supervisor_id uuid references supervisors (supervisor_id) on delete cascade not null,
    start_at      timestamptz                                                   not null default now(),
    end_at        timestamptz
);

create type progress_type as
    enum ('intro', 'ch. 1', 'ch. 2', 'ch. 3', 'ch. 4', 'ch. 5', 'ch. 6', 'end', 'literature', 'abstract');

create table semester_progress
(
    progress_id   uuid primary key,
    student_id    uuid references students (student_id) on delete cascade not null,
    progress_type progress_type                                           not null,
    first         boolean     default false                               not null,
    second        boolean     default false                               not null,
    third         boolean     default false                               not null,
    forth         boolean     default false                               not null,
    fifth         boolean     default false                               not null,
    sixth         boolean     default false                               not null,
    seventh       boolean     default false                               not null,
    eighth        boolean     default false                               not null,
    updated_at    timestamptz default now()                               not null,
    status        approval_status                                         not null default 'empty',
    accepted_at   timestamptz,
    unique (student_id, progress_type)
);

create table scientific_works_status
(
    works_id    uuid primary key,
    student_id  uuid references students (student_id) on delete cascade not null,
    semester    int                                                     not null,
    status      approval_status default 'empty'                         not null,
    updated_at  timestamptz     default now()                           not null,
    accepted_at timestamptz
);

create type publication_status as enum ('to print', 'published', 'in progress');

create table publications
(
    publication_id uuid primary key,
    works_id       uuid references scientific_works_status (works_id) on delete cascade not null,
    name           varchar(256)                                                         not null,
    scopus         bool                                                                 not null default false,
    rinc           bool                                                                 not null default false,
    wac            bool                                                                 not null default false,
    wos            bool                                                                 not null default false,
    impact         float                                                                not null,
    status         publication_status                                                   not null,
    output_data    text,
    co_authors     varchar(256),
    volume         int
);

create type conference_status as enum ('registered', 'performed');

create table conferences
(
    conference_id   uuid primary key,
    works_id        uuid references scientific_works_status (works_id) on delete cascade not null,
    status          conference_status                                                    not null,
    scopus          bool                                                                 not null default false,
    rinc            bool                                                                 not null default false,
    wac             bool                                                                 not null default false,
    wos             bool                                                                 not null default false,
    conference_name varchar(512)                                                         not null,
    report_name     varchar(512)                                                         not null,
    location        varchar(512)                                                         not null,
    reported_at     timestamptz                                                          not null
);

create table research_projects
(
    project_id   uuid primary key,
    works_id     uuid references scientific_works_status (works_id) on delete cascade not null,
    project_name varchar(512)                                                         not null,
    start_at     timestamptz                                                          not null,
    end_at       timestamptz                                                          not null,
    add_info     varchar(1024),
    grantee      varchar(1024)
);

create table teaching_load_status
(
    loads_id    uuid primary key,
    student_id  uuid references students (student_id) on delete cascade not null,
    semester    int                                                     not null,
    status      approval_status default 'empty'                         not null,
    updated_at  timestamptz     default now()                           not null,
    accepted_at timestamptz
);

create type classroom_load_type as enum ('practice', 'lectures', 'laboratory', 'exam');

create table classroom_load
(
    load_id      uuid primary key,
    t_load_id    uuid references teaching_load_status (loads_id) on delete cascade not null,
    hours        int                                                               not null,
    load_type    classroom_load_type                                               not null,
    main_teacher varchar(256)                                                      not null,
    group_name   varchar(50)                                                       not null,
    subject_name varchar(256)                                                      not null
);

create type individual_students_load_type as enum ('project practice', 'bachelor', 'masters');

create table individual_students_load
(
    load_id         uuid primary key,
    t_load_id       uuid references teaching_load_status (loads_id) on delete cascade not null,
    load_type       individual_students_load_type                                     not null,
    students_amount int                                                               not null,
    comment         varchar(1024)
);

create table additional_load
(
    load_id   uuid primary key,
    t_load_id uuid references teaching_load_status (loads_id) on delete cascade not null,
    name      varchar(512)                                                      not null,
    volume    varchar(256),
    comment   varchar(1024)
);

create table dissertations
(
    dissertation_id uuid primary key,
    student_id      uuid references students (student_id) on delete cascade not null,
    status          approval_status default 'empty'                         not null,
    created_at      timestamptz,
    updated_at      timestamptz     default now(),
    semester        int                                                     not null,
    file_name       text,
    unique (student_id, semester)
);

create table dissertation_titles
(
    title_id         uuid primary key,
    student_id       uuid references students (student_id) on delete cascade not null,
    title            varchar(512)                                            not null,
    created_at       timestamptz                                             not null,
    status           approval_status default 'empty'                         not null,
    accepted_at      timestamptz,
    semester         int                                                     not null,
    research_object  varchar(512)                                            not null,
    research_subject varchar(512)                                            not null
);

create table feedback
(
    feedback_id     uuid primary key,
    student_id      uuid references students (student_id) on delete cascade           not null,
    dissertation_id uuid references dissertations (dissertation_id) on delete cascade not null,
    feedback        text,
    semester        int                                                               not null,
    created_at      timestamptz default now()                                         not null,
    updated_at      timestamptz default now()                                         not null,
    unique (student_id, semester)
);

create table students_commentary
(
    commentary_id uuid primary key,
    student_id    uuid references students (student_id) on delete cascade not null,
    semester      int                                                     not null,
    commentary    text,
    commented_at  timestamptz                                             not null,
    unique (student_id, semester)
);

create table dissertation_commentary
(
    commentary_id uuid primary key,
    student_id    uuid references students (student_id) on delete cascade not null,
    semester      int                                                     not null,
    commentary    text,
    unique (student_id, semester)
);

create table dissertation_plans
(
    plan_id    uuid primary key,
    student_id uuid references students (student_id) on delete cascade not null,
    semester   int                                                     not null,
    plan_text  text,
    unique (student_id, semester)
);

create table exam_type
(
    type_id   int primary key,
    exam_name text not null
);

create table exams
(
    exam_id    uuid primary key,
    student_id uuid references students (student_id) on delete cascade not null,
    exam_type  int references exam_type (type_id)                      not null,
    semester   int                                                     not null,
    mark       int                                                     not null,
    set_at     timestamptz,
    unique (student_id, semester, exam_type)
);

create table supervisor_marks
(
    mark_id       uuid primary key,
    student_id    uuid references students (student_id) on delete cascade       not null,
    mark          int                                                           not null,
    semester      int                                                           not null,
    supervisor_id uuid references supervisors (supervisor_id) on delete cascade not null,
    unique (student_id, semester)
);

create table semester_count
(
    count_id uuid primary key,
    amount   int not null
);

create type patent_type as enum ('software', 'database');

create table patents
(
    patent_id         uuid primary key,
    works_id          uuid references scientific_works_status (works_id) on delete cascade not null,
    name              varchar(256)                                                         not null,
    registration_date timestamptz                                                          not null,
    type              patent_type                                                          not null,
    add_info          text
);

-- CREATE TABLE presentations
-- (
--     presentation_id UUID PRIMARY KEY,
--     presenter_id    UUID REFERENCES users (user_id) ON DELETE CASCADE NOT NULL,
--     created_at      TIMESTAMPTZ DEFAULT now() NOT NULL,
--     file_name       TEXT,
--     comments        TEXT
-- );


insert into users
values ('9bb06a04-3518-4362-bdf1-823591154464', 'chenpasha31@gmail.com',
        '$2a$10$0maw/NL4yvpjAIxDPakEKu8md3ifOVb2E4NHlO6dFPtDFIiSKoJoK', 'cc5f3793-6bd9-46cc-ab4b-8f257d37b9ac',
        'student',
        false);
insert into users
values ('d5aad0a3-e23d-4ae7-a94b-8dc7faf50613', 'chenpavel31@gmail.com',
        '$2a$10$0maw/NL4yvpjAIxDPakEKu8md3ifOVb2E4NHlO6dFPtDFIiSKoJoK', 'e39dca91-186f-456f-96bc-3bdaac3cb597',
        'supervisor',
        false);

alter table supervisors add column rank varchar(512);
alter table supervisors add column position varchar(512);

insert into users
values ('f4453590-9b99-4ea7-b3f0-446e75ff31f8', 'admin@mail.ru',
        '$2a$10$0maw/NL4yvpjAIxDPakEKu8md3ifOVb2E4NHlO6dFPtDFIiSKoJoK', 'ea37d93e-1e61-4a1f-8fef-c65a25604296', 'admin',
        true);


insert into exam_type values (1, 'Английский');
insert into exam_type values (2, 'Философия');
insert into exam_type values (3, 'Специальность');