DO
$do$

BEGIN

IF NOT EXISTS (SELECT FROM pg_extension WHERE extname = 'dblink') THEN
    CREATE EXTENSION dblink;
END IF;

-- CREATE EXTENSION dblink;

IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'zapprconfigdb') THEN
      PERFORM dblink_exec('%s', 'CREATE DATABASE zapprconfigdb');
END IF;

END;
$do$;