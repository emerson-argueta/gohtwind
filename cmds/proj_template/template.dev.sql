CREATE TABLE products
(
    id         bigint auto_increment PRIMARY KEY,
    title      VARCHAR(50) UNIQUE NOT NULL,
    price      numeric(10, 2)     NOT NULL,
    sku        VARCHAR(50) UNIQUE NOT NULL,
    created_at TIMESTAMP          NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP          NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

create table school_years
(
    id            bigint auto_increment PRIMARY KEY,
    active        tinyint(1)  null,
    year          int         null,
    school_start  date        null,
    school_end    date        null,
    reg_start     date        null,
    reg_end       date        null,
    summer_start  date        null,
    summer_end    date        null,
    testing_start date        null,
    testing_end   date        null,
    description   varchar(10) null,
    created_at    TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

create table charters
(
    id            bigint auto_increment PRIMARY KEY,
    name       varchar(255) not null,
    display    varchar(255) not null,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE users
(
    id            bigint auto_increment PRIMARY KEY,
    email      VARCHAR(255) NOT NULL UNIQUE,
    password   VARCHAR(255) NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

create table employees
(
    id            bigint auto_increment PRIMARY KEY,
    user_id        bigint null,
    first_name     varchar(100) null,
    last_name      varchar(100) null,
    middle_name    varchar(100) null,
    gender         varchar(100) null,
    address        varchar(100) null,
    phone          varchar(100) null,
    home_email     varchar(100) null,
    date_of_birth  date         null,
    hire_date      date         null,
    term_date      date         null,
    suspended_date date         null,
    created_at     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    constraint employees_UN
        unique (user_id),
    constraint employees_user_id_foreign
        foreign key (user_id) references users (id)
            on update cascade on delete cascade
);
create index employees_email_index
    on employees (home_email);
create index employees_hire_date_index
    on employees (hire_date);
create index employees_last_name_index
    on employees (last_name);
create index employees_term_date_index
    on employees (term_date);
create index user_id
    on employees (user_id);

create table students
(
    id            bigint auto_increment PRIMARY KEY,
    user_id              bigint       null,
    charter_id           bigint        default 1 not null,
    teacher_of_record_id bigint,
    enroll_status_id     bigint          default 3 not null,
    enroll_start_date    date         null,
    enroll_exit_date     date         null,
    latest_schedule_year int          null,
    first_name           varchar(255) null,
    middle_name          varchar(255) null,
    last_name            varchar(255) null,
    birthdate            date         null,
    gender               varchar(255) null,
    grade_level          varchar(255) null,
    address              varchar(255) null,
    parent_email         varchar(255) null,
    parent_phone         varchar(255) null,
    created_at           TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    constraint STUDENT_FK_CHARTER_ID
        foreign key (charter_id) references charters (id)
            on update cascade on delete set null,
    constraint STUDENT_FK_ENROLL_STATUS_ID
        foreign key (enroll_status_id) references enroll_statuses (id),
    constraint STUDENT_FK_TOR_ID
        foreign key (teacher_of_record_id) references employees (id)
            on update cascade on delete set null,
    constraint students_FK
        foreign key (user_id) references users (id)
            on update cascade on delete set null
);
create index charter_id
    on students (charter_id);
create index enroll_status_id
    on students (enroll_status_id);
create index students_tor_id_IDX
    on students (teacher_of_record_id);
create index students_user_id_IDX
    on students (user_id, id, enroll_status_id);

create table sites
(
    id            bigint auto_increment PRIMARY KEY,
    charter_id             bigint                                  null,
    support_id             bigint                               null,
    employee_id            bigint                                  null,
    active                 tinyint(1)  default 1                null,
    name                   text                                 not null,
    display                varchar(255)                         not null,
    description            varchar(255)                         null,
    subnet                 varchar(18)                          null,
    address                varchar(255)                         null,
    phone                  varchar(15) default '(951) 000-0000' null,
    public_ip_address      varchar(255)                         null,
    is_primary             tinyint(1)  default 0                null,
    access_control_enabled tinyint(1)  default 0                not null,
    created_at             TIMESTAMP                            NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at             TIMESTAMP                            NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    constraint sites_name_unique
        unique (name),
    constraint sites_display_unique
        unique (display),
    constraint sites_charter_id_foreign
        foreign key (charter_id) references charters (id)
            on update cascade on delete cascade,
    constraint sites_support_id_foreign
        foreign key (support_id) references users (id)
            on update cascade on delete set null,
    constraint sites_employee_id_foreign
        foreign key (employee_id) references users (id)
            on update cascade on delete set null
);
create index sites_charter_id_index
    on sites (charter_id);
create index sites_employee_id_index
    on sites (employee_id);
create index sites_support_id_index
    on sites (support_id);

create table course_types
(
    id            bigint auto_increment PRIMARY KEY,
    name       varchar(100) null,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

create table attendance_types
(
    id            bigint auto_increment PRIMARY KEY,
    name       varchar(100) null,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

create table absence_types
(
    id            bigint auto_increment PRIMARY KEY,
    name       varchar(100) null,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

create table tardy_types
(
    id            bigint auto_increment PRIMARY KEY,
    name       varchar(100) null,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

create table enrollment_types
(
    id            bigint auto_increment PRIMARY KEY,
    name       varchar(100) null,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

create table enroll_statuses
(
    id            bigint auto_increment PRIMARY KEY,
    name       varchar(255) null,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
create index enroll_statuses_description_index
    on enroll_statuses (name);

create table courses
(
    id            bigint auto_increment PRIMARY KEY,
    name           varchar(255) null,
    school_year_id bigint       null,
    employee_id    bigint       null,
    site_id        bigint null,
    course_type_id bigint null,
    created_at     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    constraint courses_school_year_id_foreign
        foreign key (school_year_id) references school_years (id)
            on update cascade on delete set null,
    constraint courses_employee_id_foreign
        foreign key (employee_id) references employees (id)
            on update cascade on delete set null,
    constraint courses_site_id_foreign
        foreign key (site_id) references sites (id)
            on update cascade on delete set null,
    constraint courses_course_type_id_foreign
        foreign key (course_type_id) references course_types (id)
            on update cascade on delete set null
);

create table sections
(
    id            bigint auto_increment PRIMARY KEY,
    course_id  bigint null,
    active     tinyint(1)         default 1 null,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    constraint sections_courses_id_fk
        foreign key (course_id) references courses (id)
            on update cascade on delete cascade
);
create index sections_course_id_index
    on sections (course_id);

create table enrollments
(
    id            bigint auto_increment PRIMARY KEY,
    added_by_id        bigint          null,
    student_id         bigint null,
    employee_id        bigint          null,
    course_id          bigint          null,
    section_id         bigint          null,
    enrollment_type_id bigint          not null,
    start_date         date         null,
    end_date           date         null,
    created_at         TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at         TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    constraint enrollments_enrollment_type_id_fk
        foreign key (enrollment_type_id) references enrollment_types (id)
            on update cascade on delete set null,
    constraint enrollments_added_by_id_fk
        foreign key (added_by_id) references users (id)
            on update cascade on delete set null,
    constraint enrollments_course_id_fk
        foreign key (course_id) references courses (id)
            on update cascade on delete set null,
    constraint enrollments_employee_id_fk
        foreign key (employee_id) references employees (id)
            on update cascade on delete set null,
    constraint enrollments_section_id_fk
        foreign key (section_id) references sections (id)
            on update cascade on delete set null,
    constraint enrollments_student_id_fk
        foreign key (student_id) references students (id)
            on update cascade on delete set null
);
create index enrollments_enrollment_type_id_index
    on enrollments (enrollment_type_id);
create index enrollments_section_id_index
    on enrollments (section_id);

create table attendance
(
    id                 bigint auto_increment PRIMARY KEY,
    enrollment_id      bigint    null,
    attendance_type_id bigint    null,
    absence_type_id    bigint    null,
    tardy_type_id      bigint    null,
    attendance_date    date      null,
    created_at         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    constraint attendance_FK
        foreign key (absence_type_id) references absence_types (id)
            on update cascade on delete set null,
    constraint attendance_FK_1
        foreign key (attendance_type_id) references attendance_types (id)
            on update cascade on delete set null,
    constraint attendance_FK_2
        foreign key (tardy_type_id) references tardy_types (id)
            on update cascade on delete set null,
    constraint attendance_FK_3
        foreign key (enrollment_id) references enrollments (id)
            on update cascade on delete cascade
);
create index index_attendance_on_enrollment_id
    on attendance (enrollment_id);

-- seed with fake data
INSERT INTO products (title, price, sku)
VALUES ('Cheese',9.99,'x12334'),
       ('Bread',1.99,'x12335'),
       ('Milk',2.99,'x12336'),
       ('Ice Cream',3.99,'x12337');

INSERT INTO school_years (active, year, school_start, school_end, reg_start, reg_end, summer_start, summer_end, testing_start, testing_end, description)
VALUES (0, 2020, '2020-08-01', '2021-06-30', '2020-07-01', '2020-07-31', '2021-07-01', '2021-07-31', '2021-08-01', '2021-08-31', '2020-2021'),
       (1, 2021, '2021-08-01', '2022-06-30', '2021-07-01', '2021-07-31', '2022-07-01', '2022-07-31', '2022-08-01', '2022-08-31', '2021-2022');

INSERT INTO charters (name, display)
VALUES ('Charter1', 'Test Charter1'),
       ('Charter2', 'Test Charter2'),
       ('Charter3', 'Test Charter3');

INSERT INTO users (email, password)
VALUES ('teststudent1@test.com', 'password1'),
       ('teststudent2@test.com', 'password1'),
       ('teststudent3@test.com', 'password1'),
       ('testemployee1@test.com', 'password1'),
       ('testemployee2@test.com', 'password1'),
       ('testemployee3@test.com', 'password1'),
       ('testparent1@test.com', 'password1'),
       ('testparent2@test.com', 'password1'),
       ('testparent3@test.com', 'password1');

INSERT INTO employees (id, user_id, first_name, last_name, middle_name, gender, address, phone, home_email, date_of_birth, hire_date, term_date, suspended_date)
VALUES (1, 4, 'Test', 'Employee1', 'M', 'M', '123 Test St', '555-555-5555', '777-777-7777', '1980-01-01', '2020-01-01', '2020-12-31', '2020-12-31'),
       (2, 5, 'Test', 'Employee2', 'M', 'M', '123 Test St', '555-555-5555', '777-777-7777', '1980-01-01', '2020-01-01', '2020-12-31', '2020-12-31'),
       (3, 6, 'Test', 'Employee3', 'M', 'M', '123 Test St', '555-555-5555', '777-777-7777', '1980-01-01', '2020-01-01', '2020-12-31', '2020-12-31');

INSERT INTO students (id, user_id, charter_id, teacher_of_record_id, enroll_status_id, enroll_start_date, enroll_exit_date, latest_schedule_year, first_name, middle_name, last_name, birthdate, gender, grade_level, address, parent_email, parent_phone)
VALUES (1, 1, 1, 4, 1, '2020-08-01', '2021-06-30', 2021, 'Test', 'Student1', 'M', '2000-01-01', 'M', '1', '123 Test St', 'parent1@yahoo.com', '555-555-5555'),
       (2, 2, 1, 5, 1, '2020-08-01', '2021-06-30', 2021, 'Test', 'Student2', 'M', '2000-01-01', 'M', '1', '123 Test St', 'parent2@yahoo.com', '555-555-5555'),
       (3, 3, 1, 6, 1, '2020-08-01', '2021-06-30', 2021, 'Test', 'Student3', 'M', '2000-01-01', 'M', '1', '123 Test St', 'parent3@yahoo.com', '555-555-5555');

INSERT INTO sites (id, charter_id, support_id, employee_id, active, name, display, description, subnet, address, phone, public_ip_address, is_primary, access_control_enabled)
VALUES (1, 1, 4, 4, 1, 'Test Site1', 'Test Site1', 'Test Site1', '255', '123 Test St', '555-555-5555', '127.0.0.1', 1, 1),
       (2, 1, 5, 5, 1, 'Test Site2', 'Test Site2', 'Test Site2', '255', '123 Test St', '555-555-5555', '127.0.0.1', 1, 1),
       (3, 1, 6, 6, 1, 'Test Site3', 'Test Site3', 'Test Site3', '255', '123 Test St', '555-555-5555', '127.0.0.1', 1, 1);

INSERT INTO course_types (id, name)
VALUES (1, 'Online'),
       (2, 'Site Based'),
       (3, 'Hybrid');

INSERT INTO attendance_types (id, name)
VALUES (1, 'Present'),
       (2, 'Absent'),
       (3, 'Tardy');

INSERT INTO absence_types (id, name)
VALUES (1, 'Excused'),
       (2, 'Unexcused'),
       (3, 'Suspended');

INSERT INTO tardy_types (id, name)
VALUES (1, 'less than 30 minutes'),
       (2, 'more than 30 minutes');

INSERT INTO enrollment_types (id, name)
VALUES (1, 'student'),
       (2, 'teacher'),
       (3, 'additional teacher');

INSERT INTO enroll_statuses (id, name)
VALUES (1, 'Active'),
       (2, 'Pending'),
       (3, 'Withdrawn'),
       (4, 'Graduated');

INSERT INTO courses (id, name, school_year_id, employee_id, site_id, course_type_id)
VALUES (1, 'Test Course1', 2, 4, 1, 1),
       (2, 'Test Course2', 2, 5, 2, 2),
       (3, 'Test Course3', 2, 6, 3, 3);

INSERT INTO sections (id, course_id, active)
VALUES (1, 1, 1),
       (2, 2, 1),
       (3, 3, 1);

INSERT INTO enrollments (id, added_by_id, student_id, employee_id, course_id, section_id, enrollment_type_id, start_date, end_date)
VALUES (1, 4, 1, 4, 1, 1, 1, '2020-08-01', '2021-06-30'),
       (2, 5, 2, 5, 2, 2, 2, '2020-08-01', '2021-06-30'),
       (3, 6, 3, 6, 3, 3, 3, '2020-08-01', '2021-06-30');
