
CREATE TABLE tokens
(
  id serial NOT NULL,
  user_id integer NOT NULL,
  token text NOT NULL,
  CONSTRAINT id_primary PRIMARY KEY (id)
)
WITH (
  OIDS=FALSE
);

CREATE INDEX user_id_idx
  ON tokens
  USING btree
  (user_id);

CREATE TABLE users
(
  id serial NOT NULL,
  email character varying NOT NULL,
  name character varying NOT NULL,
  salt bytea NOT NULL,
  hash bytea NOT NULL,
  CONSTRAINT primary_id PRIMARY KEY (id),
  CONSTRAINT unique_email UNIQUE (email)
)
WITH (
  OIDS=FALSE
);
