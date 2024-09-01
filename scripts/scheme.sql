pragma foreign_keys = ON;

create table Applications
(
    id   integer primary key autoincrement,
    name text not null,
    xkey text not null,
    host text not null
);

create table Messages
(
    id    integer primary key autoincrement,
    appId integer not null references Applications,
    reply integer null references Messages -- сообщение создано как "ответ" на ID другого сообщения
);

create table Binds
(
    appId      integer not null references Applications,
    msgId      integer not null references Messages,
    msgLocalId text    not null
);

create table Attachments
(
    id         integer primary key autoincrement,
    localId    text    not null,
    appId      integer not null references Applications,
    msgId      integer not null references Messages,
    url        text    not null,
    hasSpoiler integer not null, -- bool
    type       integer not null
);

create table Files
(
    id       integer primary key autoincrement,
    fileName text not null,
    attachmentId int not null references Attachments
);