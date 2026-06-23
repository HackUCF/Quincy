-- main scores table
-- all data for all boxes, all teams, and all services
CREATE TABLE IF NOT EXISTS scores (
  service   VARCHAR(16)  NOT NULL,
  box       VARCHAR(16)  NOT NULL,
  team_num  INTEGER      NOT NULL,
  status    BOOLEAN      NOT NULL,
  message   TEXT         NOT NULL,
  timestamp BIGINT       NOT NULL, -- unix microseconds as integer
  id        SERIAL       PRIMARY KEY
);

CREATE INDEX IF NOT EXISTS idx_scores_ts             ON scores(timestamp);
CREATE INDEX IF NOT EXISTS idx_scores_team_status    ON scores(team_num, status);
CREATE INDEX IF NOT EXISTS idx_scores_team_ts_status ON scores(team_num, timestamp, status);

-- recent scores table
-- nearly identical: stores only the most recent of each check
-- has different primary key for faster access
CREATE TABLE IF NOT EXISTS recent_scores (
  service   VARCHAR(16) NOT NULL,
  box       VARCHAR(16) NOT NULL,
  team_num  INTEGER     NOT NULL,
  status    BOOLEAN     NOT NULL,
  message   TEXT        NOT NULL,
  timestamp BIGINT      NOT NULL,
  PRIMARY KEY (team_num, service, box)
);

-- final scores table
-- stores totals for each team/box/service
-- this allows for quick final stat generation
CREATE TABLE IF NOT EXISTS final_scores (
  service   VARCHAR(16) NOT NULL,
  box       VARCHAR(16) NOT NULL,
  team_num  INTEGER     NOT NULL,
  passed    INTEGER     NOT NULL,
  total     INTEGER     NOT NULL,
  PRIMARY KEY (team_num, service, box)
);

-- this table exists to persistently store password change requests.
-- the config is loaded into it on startup.
-- already changed passwords persist on reboot
CREATE TABLE IF NOT EXISTS scoring_users (
  team_num  INTEGER NOT NULL,
  user_list TEXT    NOT NULL,
  username  TEXT    NOT NULL,
  password  TEXT,
  domain    TEXT,
  netbios   TEXT,
  PRIMARY KEY (team_num, user_list, username)
);