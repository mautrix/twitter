INSERT INTO "user" (bridge_id, mxid, management_room, access_token)
SELECT '', mxid, notice_room, ''
FROM user_old;

INSERT INTO user_login (bridge_id, user_mxid, id, remote_name, space_room, metadata, remote_profile)
SELECT
    '', -- bridge_id
    mxid, -- user_mxid
    CAST(twid AS TEXT), -- id
    '', -- remote_name
    notice_room,
    -- only: postgres
    jsonb_build_object
    -- only: sqlite (line commented)
--  json_object
    (
        'cookies', 'auth_token=' || auth_token || '; ct0=' || csrf_token || ';'
    ), -- metadata
    '{}' -- remote_profile
FROM user_old WHERE twid IS NOT NULL;

INSERT INTO ghost (
    bridge_id, id, name, avatar_id, avatar_hash, avatar_mxc,
    name_set, avatar_set, contact_info_set, is_bot, identifiers, metadata
)
SELECT
    '', -- bridge_id
    CAST(twid AS TEXT), -- id
    COALESCE(name, ''), -- name
    COALESCE(photo_url, ''), -- avatar_id
    '', -- avatar_hash
    COALESCE(photo_mxc, ''), -- avatar_mxc
    CASE WHEN name <> '' THEN 1 ELSE 0 END, -- name_set
    CASE WHEN photo_mxc <> '' THEN 1 ELSE 0 END, -- avatar_set
    contact_info_set, -- contact_info_set
    false, -- is_bot
    '[]', -- identifiers
    '{}' -- metadata
FROM puppet_old;

INSERT INTO portal (
    bridge_id, id, receiver, mxid, other_user_id,
    name, topic, avatar_id, avatar_hash, avatar_mxc, name_set, avatar_set, topic_set, name_is_custom,
    in_space, room_type, metadata
)
SELECT
    '', -- bridge_id
    CAST(twid AS TEXT), -- id
    CASE WHEN receiver<>0 THEN CAST(receiver AS TEXT) ELSE '' END, -- receiver
    mxid, -- mxid
    CASE WHEN conv_type='ONE_TO_ONE' THEN CAST(other_user AS TEXT) END, -- other_user_id
    '', -- name
    '', -- topic
    '', -- avatar_id
    '', -- avatar_hash
    '', -- avatar_mxc
    1, -- name_set
    1, --avatar_set
    false, -- topic_set
    false, -- name_is_custom
    false, -- in_space
    CASE WHEN conv_type='GROUP_DM' THEN 'group_dm' ELSE 'dm' END, -- room_type
    '{}' -- metadata
FROM portal_old;

INSERT INTO ghost (bridge_id, id, name, avatar_id, avatar_hash, avatar_mxc, name_set, avatar_set, contact_info_set, is_bot, identifiers, metadata)
VALUES ('', '', '', '', '', '', false, false, false, false, '[]', '{}')
ON CONFLICT (bridge_id, id) DO NOTHING;

INSERT INTO message (
    bridge_id, id, part_id, mxid, room_id, room_receiver,
    sender_id, sender_mxid, timestamp, edit_count, metadata
)
SELECT
    '', -- bridge_id
    CAST(twid as TEXT), -- id
    '', -- part_id
    mxid, -- mxid
    (SELECT twid FROM portal_old WHERE portal_old.mxid=message_old.mx_room), -- room_id
    CASE WHEN receiver<>0 THEN CAST(receiver AS TEXT) ELSE '' END, -- room_receiver
    '', -- sender_id
    '', -- sender_mxid
    ((twid>>22)+1288834974657)*1000000, -- timestamp
    0, -- edit_count
    '{}' -- metadata
FROM message_old;

INSERT INTO reaction (
    bridge_id, message_id, message_part_id, sender_id, sender_mxid,
    emoji_id, room_id, room_receiver, timestamp,  mxid, emoji, metadata
)
SELECT
    '', -- bridge_id
    tw_msgid, -- message_id
    '', -- message_part_id
    CAST(tw_sender AS TEXT), -- sender_id
    '', -- sender_mxid
    '', -- emoji_id
    (SELECT twid FROM portal_old WHERE portal_old.mxid=reaction_old.mx_room), -- room_id
    CASE WHEN tw_receiver<>0 THEN CAST(tw_receiver AS TEXT) ELSE '' END, -- room_receiver
    ((COALESCE(tw_reaction_id, tw_msgid)>>22)+1288834974657)*1000000, -- timestamp
    mxid,
    reaction, -- emoji
    '{}' -- metadata
FROM reaction_old;

-- Python -> Go mx_ table migration
ALTER TABLE mx_room_state DROP COLUMN is_encrypted;
ALTER TABLE mx_room_state RENAME COLUMN has_full_member_list TO members_fetched;
UPDATE mx_room_state SET members_fetched=false WHERE members_fetched IS NULL;

-- only: postgres until "end only"
ALTER TABLE mx_room_state ALTER COLUMN power_levels TYPE jsonb USING power_levels::jsonb;
ALTER TABLE mx_room_state ALTER COLUMN encryption TYPE jsonb USING encryption::jsonb;
ALTER TABLE mx_room_state ALTER COLUMN members_fetched SET DEFAULT false;
ALTER TABLE mx_room_state ALTER COLUMN members_fetched SET NOT NULL;
-- end only postgres

ALTER TABLE mx_user_profile ADD COLUMN name_skeleton bytea;
CREATE INDEX mx_user_profile_membership_idx ON mx_user_profile (room_id, membership);
CREATE INDEX mx_user_profile_name_skeleton_idx ON mx_user_profile (room_id, name_skeleton);

UPDATE mx_user_profile SET displayname='' WHERE displayname IS NULL;
UPDATE mx_user_profile SET avatar_url='' WHERE avatar_url IS NULL;

CREATE TABLE mx_registrations (
    user_id TEXT PRIMARY KEY
);

UPDATE mx_version SET version=7;

DROP TABLE reaction_old;
DROP TABLE message_old;
DROP TABLE portal_old;
DROP TABLE puppet_old;
DROP TABLE user_old;
