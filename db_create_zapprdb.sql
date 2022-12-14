DO
$do$

BEGIN

CREATE TABLE IF NOT EXISTS public."DBConfig"
(
    "Identifier" character varying COLLATE pg_catalog."default" NOT NULL,
    "Config" jsonb NOT NULL,
    "CreatedOn" date NOT NULL,
    "ModifiedOn" date,
    CONSTRAINT "DBConfigIdentifier" PRIMARY KEY ("Identifier")
)
WITH (
    OIDS = FALSE
);

END
$do$;
