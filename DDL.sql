create table osu.BEATMAP
(
    BEATMAP_ID              int(10)                               not null
        primary key,
    BEATMAPSET_ID           int(10)                               not null,
    MODE                    varchar(6)                            null,
    MODE_INT                tinyint(1)                            null,
    STATUS                  varchar(9)                            null,
    RANKED                  tinyint(1)                            null,
    TOTAL_LENGTH            int(10)                               null,
    MAX_COMBO               int(10)                               null,
    DIFFICULTY_RATING       decimal(63, 2)                        null,
    VERSION                 varchar(254)                          null,
    ACCURACY                decimal(63, 2)                        null,
    AR                      decimal(63, 2)                        null,
    CS                      decimal(63, 2)                        null,
    DRAIN                   decimal(63, 2)                        null,
    BPM                     decimal(63, 2)                        null,
    `CONVERT`               tinyint(1)                            null,
    COUNT_CIRCLES           int(10)                               null,
    COUNT_SLIDERS           int(10)                               null,
    COUNT_SPINNERS          int(10)                               null,
    DELETED_AT              datetime                              null,
    HIT_LENGTH              int(10)                               null,
    IS_SCOREABLE            tinyint(1)                            null,
    LAST_UPDATED            datetime                              null,
    PASSCOUNT               int(10)                               null,
    PLAYCOUNT               int(10)                               null,
    CHECKSUM                varchar(32)                           null,
    USER_ID                 int(10)                               null,
    SYSTEM_UPDATE_TIMESTAMP timestamp default current_timestamp() null on update current_timestamp()
)
    collate = utf8mb3_unicode_ci;

create index SEARCH_001
    on osu.BEATMAP (BEATMAPSET_ID);

create index SEARCH_002
    on osu.BEATMAP (MODE_INT);

create index SEARCH_003
    on osu.BEATMAP (CHECKSUM);

create index SEARCH_004
    on osu.BEATMAP (TOTAL_LENGTH);

create index SEARCH_005
    on osu.BEATMAP (MAX_COMBO);

create index SEARCH_006
    on osu.BEATMAP (DIFFICULTY_RATING);

create index SEARCH_007
    on osu.BEATMAP (ACCURACY);

create index SEARCH_008
    on osu.BEATMAP (AR);

create index SEARCH_009
    on osu.BEATMAP (CS);

create index SEARCH_010
    on osu.BEATMAP (DRAIN);

create index SEARCH_011
    on osu.BEATMAP (BPM);

create index SEARCH_012
    on osu.BEATMAP (DELETED_AT);

create table osu.BEATMAPSET
(
    BEATMAPSET_ID                  int(1)                                not null
        primary key,
    ARTIST                         varchar(254)                          null,
    ARTIST_UNICODE                 varchar(254) charset utf8mb4          null,
    CREATOR                        varchar(254) charset utf8mb4          null,
    FAVOURITE_COUNT                int(1)                                null,
    HYPE_CURRENT                   int(1)                                null,
    HYPE_REQUIRED                  int(1)                                null,
    NSFW                           tinyint(1)                            null,
    PLAY_COUNT                     int(1)                                null,
    SOURCE                         varchar(254) charset utf8mb4          null,
    STATUS                         varchar(9) charset utf8mb4            null,
    TITLE                          varchar(254) charset utf8mb4          null,
    TITLE_UNICODE                  varchar(254) charset utf8mb4          null,
    USER_ID                        int(1)                                null,
    VIDEO                          tinyint(1)                            null,
    AVAILABILITY_DOWNLOAD_DISABLED tinyint(1)                            null,
    AVAILABILITY_MORE_INFORMATION  text charset utf8mb4                  null,
    BPM                            decimal(63, 2)                        null,
    CAN_BE_HYPED                   tinyint(1)                            null,
    DISCUSSION_ENABLED             tinyint(1)                            null,
    DISCUSSION_LOCKED              tinyint(1)                            null,
    IS_SCOREABLE                   tinyint(1)                            null,
    LAST_UPDATED                   datetime                              null,
    LEGACY_THREAD_URL              varchar(254) charset utf8mb4          null,
    NOMINATIONS_SUMMARY_CURRENT    int(1)                                null,
    NOMINATIONS_SUMMARY_REQUIRED   int(1)                                null,
    RANKED                         tinyint(1)                            null,
    RANKED_DATE                    datetime                              null,
    DELETED_AT                     datetime                              null,
    STORYBOARD                     tinyint(1)                            null,
    SUBMITTED_DATE                 datetime                              null,
    TAGS                           text charset utf8mb4                  null,
    HAS_FAVOURITED                 tinyint(1)                            null,
    DESCRIPTION                    text charset utf8mb4                  null,
    GENRE_ID                       int(1)                                null,
    GENRE_NAME                     varchar(254) charset utf8mb4          null,
    LANGUAGE_ID                    int(1)                                null,
    LANGUAGE_NAME                  varchar(254) charset utf8mb4          null,
    RATINGS                        varchar(254)                          null,
    SYSTEM_UPDATE_TIMESTAMP        timestamp default current_timestamp() null on update current_timestamp()
)
    collate = utf8mb4_unicode_ci;

create index SEARCH_001
    on osu.BEATMAPSET (RANKED_DATE);

create index SEARCH_002
    on osu.BEATMAPSET (RANKED_DATE desc);

create index SEARCH_003
    on osu.BEATMAPSET (FAVOURITE_COUNT desc);

create index SEARCH_004
    on osu.BEATMAPSET (PLAY_COUNT desc);

create index SEARCH_005
    on osu.BEATMAPSET (LAST_UPDATED);

create index SEARCH_006
    on osu.BEATMAPSET (RANKED);

create index SEARCH_007
    on osu.BEATMAPSET (NSFW);

create index SEARCH_008
    on osu.BEATMAPSET (VIDEO);

create index SEARCH_009
    on osu.BEATMAPSET (STORYBOARD);

create index SEARCH_010
    on osu.BEATMAPSET (DELETED_AT);

create table osu.BEATMAP_PACK
(
    PACK_ID            int                                   not null
        primary key,
    TYPE               varchar(45)                           not null,
    NAME               varchar(1023)                         null,
    CREATOR            varchar(1023)                         null,
    DATE               varchar(45)                           null,
    DOWNLOAD_URL       varchar(2047)                         null,
    SYSTEM_UPDATE_DATE timestamp default current_timestamp() not null on update current_timestamp()
)
    collate = utf8mb3_unicode_ci;

create table osu.BEATMAP_PACK_SETS
(
    PACK_ID            int                                   not null,
    BEATMAPSET_ID      int(1)                                not null,
    SYSTEM_UPDATE_DATE timestamp default current_timestamp() not null on update current_timestamp(),
    primary key (PACK_ID, BEATMAPSET_ID)
)
    collate = utf8mb3_unicode_ci;

create table osu.BLACKLIST
(
    ID                      bigint auto_increment
        primary key,
    IPV4                    char(15)                                                    not null,
    EXPIRED_AT              timestamp(6) default (current_timestamp() + interval 7 day) not null,
    COUNT                   int          default 0                                      not null,
    MEMO                    varchar(254)                                                null,
    SYSTEM_CREATE_TIMESTAMP timestamp(6) default current_timestamp(6)                   not null,
    SYSTEM_UPDATE_TIMESTAMP timestamp(6) default current_timestamp(6)                   not null on update current_timestamp(6)
);

create index BLACKLIST_ID_001
    on osu.BLACKLIST (EXPIRED_AT desc);

create table osu.SEARCH_CACHE_ARTIST
(
    INDEX_KEY     bigint unsigned   not null,
    BEATMAPSET_ID int               not null,
    TMP           tinyint default 1 not null,
    primary key (INDEX_KEY, BEATMAPSET_ID)
)
    collate = utf8mb3_unicode_ci;

create index SEARCH_CACHE_ARTIST_BEATMAPSET_ID
    on osu.SEARCH_CACHE_ARTIST (BEATMAPSET_ID);

create table osu.SEARCH_CACHE_CREATOR
(
    INDEX_KEY     bigint unsigned   not null,
    BEATMAPSET_ID int               not null,
    TMP           tinyint default 1 not null,
    primary key (INDEX_KEY, BEATMAPSET_ID)
)
    collate = utf8mb3_unicode_ci;

create index SEARCH_CACHE_CREATOR_BEATMAPSET_ID
    on osu.SEARCH_CACHE_CREATOR (BEATMAPSET_ID);

create table osu.SEARCH_CACHE_OTHER
(
    INDEX_KEY     bigint unsigned   not null,
    BEATMAPSET_ID int               not null,
    TMP           tinyint default 1 not null,
    primary key (INDEX_KEY, BEATMAPSET_ID)
)
    collate = utf8mb3_unicode_ci;

create index SEARCH_CACHE_OTHER_BEATMAPSET_ID
    on osu.SEARCH_CACHE_OTHER (BEATMAPSET_ID);

create table osu.SEARCH_CACHE_STRING_INDEX
(
    STRING varchar(254) charset utf8mb4         not null,
    ID     bigint unsigned default uuid_short() not null
        primary key,
    TMP    tinyint         default 1            not null,
    constraint STRING_UNIQUE
        unique (STRING)
)
    collate = utf8mb3_unicode_ci;

create table osu.SEARCH_CACHE_TAG
(
    INDEX_KEY     bigint unsigned   not null,
    BEATMAPSET_ID int               not null,
    TMP           tinyint default 1 not null,
    primary key (INDEX_KEY, BEATMAPSET_ID)
)
    collate = utf8mb3_unicode_ci;

create index SEARCH_CACHE_TAG_BEATMAPSET_ID
    on osu.SEARCH_CACHE_TAG (BEATMAPSET_ID);

create table osu.SEARCH_CACHE_TITLE
(
    INDEX_KEY     bigint unsigned   not null,
    BEATMAPSET_ID int               not null,
    TMP           tinyint default 1 not null,
    primary key (INDEX_KEY, BEATMAPSET_ID)
)
    collate = utf8mb3_unicode_ci;

create index SEARCH_CACHE_TITLE_BEATMAPSET_ID
    on osu.SEARCH_CACHE_TITLE (BEATMAPSET_ID);

create table osu.SERVER_CACHE
(
    `KEY` varchar(254)  not null
        primary key,
    VALUE varchar(1024) not null
)
    engine = MEMORY;

