PRAGMA FOREIGN_KEYS = ON;

CREATE TABLE IF NOT EXISTS alias(
    i VARCHAR(64) NOT NULL,
    j VARCHAR(64) NOT NULL,
    PRIMARY KEY(i,j),
    CHECK (i != ""),
    CHECK (j != "")
);

CREATE TABLE IF NOT EXISTS pic(
    pid INTEGER PRIMARY KEY NOT NULL,
    uid INTEGER NOT NULL,
    series INTEGER NOT NULL,
    stat INTEGER NOT NULL,
    views INTEGER NOT NULL,
    bookmark INTEGER NOT NULL,
    likes INTEGER NOT NULL,
    timestamp INTEGER NOT NULL,
    title VARCHAR(64) NOT NULL,
    CHECK (pid > 0),
    CHECK (uid > 0),
    CHECK (views >= 0),
    CHECK (bookmark >= 0),
    CHECK (likes >= 0)
);

CREATE TABLE IF NOT EXISTS pictag(
    pid INTEGER NOT NULL,
    name VARCHAR(64) NOT NULL,
    FOREIGN KEY(pid) REFERENCES pic(pid),
    PRIMARY KEY(pid,name),
    CHECK (pid > 0),
    CHECK (name != "")
);
