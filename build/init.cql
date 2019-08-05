CREATE KEYSPACE IF NOT EXISTS stats WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

CREATE TABLE IF NOT EXISTS counters (
    account_id text,
    action_id text,
    count counter,
    PRIMARY KEY((account_id, action_id)),
);
