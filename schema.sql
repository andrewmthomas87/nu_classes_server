DROP DATABASE IF EXISTS northwestern;
CREATE DATABASE northwestern;
USE northwestern;

CREATE TABLE terms
(
    id         INT,
    name       VARCHAR(100),
    start_date VARCHAR(30),
    end_date   VARCHAR(30),
    PRIMARY KEY (id)
);

CREATE TABLE schools
(
    symbol VARCHAR(30),
    name   VARCHAR(100),
    PRIMARY KEY (symbol)
);

CREATE TABLE subjects
(
    symbol VARCHAR(30),
    name   VARCHAR(100),
    PRIMARY KEY (symbol)
);

CREATE TABLE subject_availabilities
(
    term    int,
    school  VARCHAR(30),
    subject VARCHAR(30),
    UNIQUE (term, school, subject)
);

CREATE TABLE instructors
(
    id    INT,
    name  VARCHAR(250),
    phone VARCHAR(30),
    PRIMARY KEY (id)
);

CREATE TABLE instructor_subjects
(
    instructor INT,
    subject    VARCHAR(30),
    UNIQUE (instructor, subject)
);

CREATE TABLE buildings
(
    id   INT,
    name VARCHAR(250),
    lat  DOUBLE PRECISION,
    lon  DOUBLE PRECISION,
    PRIMARY KEY (id)
);

CREATE TABLE rooms
(
    id          INT,
    building_id INT,
    name        VARCHAR(250),
    PRIMARY KEY (id)
);

CREATE TABLE courses
(
    id           INT,
    title        VARCHAR(250),
    term         INT,
    school       VARCHAR(30),
    instructor   INT,
    subject      VARCHAR(30),
    catalog_num  VARCHAR(30),
    section      VARCHAR(30),
    room         INT,
    meeting_days VARCHAR(30),
    start_time   VARCHAR(30),
    end_time     VARCHAR(30),
    start_date   VARCHAR(30),
    end_date     VARCHAR(30),
    seats        INT,
    overview     VARCHAR(5000),
    topic        VARCHAR(2500),
    attributes   VARCHAR(2500),
    requirements VARCHAR(2500),
    component    VARCHAR(30),
    class_num    INT,
    course_id    INT,
    PRIMARY KEY (id)
);

CREATE TABLE course_descriptions
(
    course      INT,
    name        VARCHAR(1000),
    description VARCHAR(5000),
    PRIMARY KEY (course)
);

CREATE TABLE course_components
(
    course       INT,
    component    VARCHAR(30),
    meeting_days VARCHAR(30),
    start_time   VARCHAR(30),
    end_time     VARCHAR(30),
    section      VARCHAR(30),
    room         VARCHAR(250),
    UNIQUE (course, component, section)
);
