DO
$do$

BEGIN

CREATE TABLE IF NOT EXISTS public."UserConfig"
(
    "Identifier" character varying COLLATE pg_catalog."default" NOT NULL,
    "Config" jsonb NOT NULL,
    "CreatedOn" date NOT NULL,
    "ModifiedOn" date,
)
WITH (
    OIDS = FALSE
);

END
$do$;
