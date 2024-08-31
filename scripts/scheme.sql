pragma foreign_keys = ON;

create table Applications
(
    id       integer primary key autoincrement,
    name     text not null,
    xkey     text not null,
    callback text not null
);

create table Messages
(
    id         integer primary key autoincrement,
    isSilent   integer not null,                 -- bool, сообщение без уведомления
    isForward  integer not null,                 -- bool, сообщение создано как "пересланное"
    reply      integer null references Messages, -- сообщение создано как "ответ" на ID другого сообщения
    username   text    not null,                 -- имя автора сообщения
    text       text    not null,                 -- текст сообщения
    createDate integer not null                  -- время создания сообщения в формате unix
);

create table MessageMap
(
    msgId   integer not null references Messages,
    inAppId text    not null
);

create table Attachments
(
    msgId      integer not null references Messages,
    fileId     integer not null references Files,
    hasSpoiler integer not null, -- bool
    type       integer not null
);

create table Files
(
    id       integer primary key autoincrement,
    fileName text not null
);