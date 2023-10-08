-- drop table if exists users cascade;
drop table if exists appointments cascade;
drop table if exists appointment_permissions cascade;
drop table if exists schedules cascade;
drop table if exists slots cascade;
drop table if exists reservations cascade;

create table if not exists users (
    id serial primary key,
    email varchar(255) not null unique,
    password varchar(255) not null
);

create table if not exists appointments (
    id serial primary key,
    name varchar(255) not null,
    instructions text not null default '',
    user_id int not null,
    timezone varchar(255) not null,
    token uuid not null default gen_random_uuid(),
    constraint fk_user foreign key(user_id) references users(id) on delete cascade
);

create table if not exists schedules (
    id serial primary key,
    start_date varchar(255) not null,
    end_date varchar(255) not null,
    start_time varchar(255) not null,
    end_time varchar(255) not null,
    timezone varchar(255) not null,
    length int not null,
    spots int not null,
    status varchar(255) not null,
    appointment_id int not null,
    user_id int not null,
    constraint fk_appointment foreign key(appointment_id) references appointments(id) on delete cascade,
    constraint fk_user foreign key(user_id) references users(id) on delete cascade
);

create table if not exists slots (
    id serial primary key,
    start timestamptz not null,
    "end" timestamptz not null,
    spots int not null,
    free int not null,
    appointment_id int not null,
    token uuid not null default gen_random_uuid(),
    constraint fk_appointment foreign key(appointment_id) references appointments(id) on delete cascade
);

create table if not exists reservations (
    id serial primary key,
    slot_id int not null,
    first_name varchar(255) not null,
    last_name varchar(255) not null,
    email varchar(255) not null,
    timezone varchar(255) not null,
    appointment_id int not null,
    token uuid not null default gen_random_uuid(),
    constraint fk_appointment foreign key(appointment_id) references appointments(id) on delete cascade,
    constraint fk_slot foreign key(slot_id) references slots(id) on delete cascade
);

create table if not exists appointment_permissions (
    user_id int not null,
    action varchar(255) not null,
    entity_id int not null,
    constraint fk_user foreign key(user_id) references users(id) on delete cascade,
    constraint fk_appointment foreign key(entity_id) references appointments(id) on delete cascade,
    primary key (user_id, action, entity_id)
);