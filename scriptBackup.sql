create table content
(
    contentID int auto_increment
        primary key,
    text      longtext not null,
    file      longblob not null,
    fileName  text     null
)
    row_format = COMPRESSED;

create table user
(
    username    varchar(30)                         null,
    email       varchar(255)                        not null,
    password    binary(16)                          not null,
    create_time timestamp default CURRENT_TIMESTAMP not null,
    userID      int auto_increment
        primary key,
    logo        longblob                            null,
    constraint user_pk
        unique (username),
    constraint user_username_uindex
        unique (username)
)
    row_format = COMPRESSED;

create table connects
(
    userA    int                                 null,
    userB    int                                 null,
    id       int auto_increment
        primary key,
    timeOpen timestamp default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP,
    constraint connects_pk
        unique (userA, userB),
    constraint connects_user_userID_fk
        foreign key (userA) references user (userID)
            on delete cascade,
    constraint connects_user_userID_fk2
        foreign key (userB) references user (userID)
            on delete cascade
)
    row_format = COMPRESSED;

create table message
(
    messageID int auto_increment
        primary key,
    fromID    int                                not null,
    toID      int                                not null,
    date      datetime default CURRENT_TIMESTAMP not null,
    contentID int                                not null,
    constraint message_date_uindex
        unique (date),
    constraint content
        foreign key (contentID) references content (contentID)
            on update cascade on delete cascade,
    constraint `from`
        foreign key (fromID) references user (userid),
    constraint `to`
        foreign key (toID) references user (userid)
)
    row_format = COMPRESSED;

create table favouritemessage
(
    id              int auto_increment
        primary key,
    fromID          int not null,
    messageID       int not null,
    favouriteuserid int not null,
    constraint favouriteMessage_message_messageID_fk
        foreign key (messageID) references message (messageID)
            on update cascade on delete cascade,
    constraint favouriteMessage_user_userID_fk
        foreign key (fromID) references user (userID),
    constraint favouriteMessage_user_userID_fk1
        foreign key (favouriteuserid) references user (userID)
);

create definer = mysql@`%` trigger before_delete_message
    before delete
    on message
    for each row
BEGIN
        delete from content where content.contentID=OLD.contentID;
    END;

create definer = mysql@`%` trigger before_delete_user
    before delete
    on user
    for each row
BEGIN
        delete from message where message.fromID=OLD.userID;
    END;

