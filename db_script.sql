DO
$do$

BEGIN

CREATE TABLE IF NOT EXISTS public."APIKey"
(
    "Identifier" character varying COLLATE pg_catalog."default" NOT NULL,
    "Name" character varying COLLATE pg_catalog."default" NOT NULL,
    "Secret" character varying COLLATE pg_catalog."default" NOT NULL,
    "TenantIdentifier" character varying COLLATE pg_catalog."default" NOT NULL,
    "CreatedOn" date NOT NULL,
    "ModifiedOn" date,
    "Scopes" character varying[] COLLATE pg_catalog."default" NOT NULL,
    "UserIdentifier" character varying COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT "APIKeyIdentifier" PRIMARY KEY ("Identifier")
)
WITH (
    OIDS = FALSE
);

CREATE TABLE IF NOT EXISTS public."ExternalConnector"
(
    "Identifier" character varying COLLATE pg_catalog."default" NOT NULL,
    "Metadata" jsonb NOT NULL,
    "TenantIdentifier" character varying COLLATE pg_catalog."default" NOT NULL,
    "CreatedOn" date NOT NULL,
    "ModifiedOn" date,
    CONSTRAINT "ExternalConnectorIdentifier" PRIMARY KEY ("Identifier")
)
WITH (
    OIDS = FALSE
);

CREATE TABLE IF NOT EXISTS public."ExternalConnectorExecution"
(
    "Identifier" character varying COLLATE pg_catalog."default" NOT NULL,
    "StatusCode" integer NOT NULL,
    "ResponseMessage" character varying COLLATE pg_catalog."default" NOT NULL,
    "ExternalConnectorIdentifier" character varying COLLATE pg_catalog."default" NOT NULL,
    "CreatedOn" date NOT NULL,
    "ModifiedOn" date,
    CONSTRAINT "ExternalConnectorExecutionIdentifier" PRIMARY KEY ("Identifier")
)
WITH (
    OIDS = FALSE
);

CREATE TABLE IF NOT EXISTS public."OAuthProvider"
(
    "Name" character varying COLLATE pg_catalog."default" NOT NULL,
    "Identifier" character varying COLLATE pg_catalog."default" NOT NULL,
    "Metadata" jsonb NOT NULL,
    "TenantIdentifier" character varying COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT "OAuthProviderIdentifier" PRIMARY KEY ("Identifier")
)
WITH (
    OIDS = FALSE
);

CREATE TABLE IF NOT EXISTS public."Role"
(
    "Identifier" character varying COLLATE pg_catalog."default" NOT NULL,
    "Name" character varying COLLATE pg_catalog."default" NOT NULL,
    "CreatedOn" date NOT NULL,
    "ModifiedOn" date,
    "TenantIdentifier" character varying COLLATE pg_catalog."default" NOT NULL,
    "Scopes" character varying[] COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT "MasterRoleIdentifier" PRIMARY KEY ("Identifier")
)
WITH (
    OIDS = FALSE
);

CREATE TABLE IF NOT EXISTS public."SMTPSettings"
(
    "Identifier" character varying COLLATE pg_catalog."default" NOT NULL,
    "Server" character varying COLLATE pg_catalog."default" NOT NULL,
    "Username" character varying COLLATE pg_catalog."default" NOT NULL,
    "Password" character varying COLLATE pg_catalog."default" NOT NULL,
    "TLSPort" integer,
    "NonTLSPort" integer,
    "TenantIdentifier" character varying COLLATE pg_catalog."default" NOT NULL,
    "CreatedOn" date NOT NULL,
    "ModifiedOn" date,
    CONSTRAINT "SMTPSettingsIdentifier" PRIMARY KEY ("Identifier")
)
WITH (
    OIDS = FALSE
);

CREATE TABLE IF NOT EXISTS public."StaticStorage"
(
    "Identifier" character varying COLLATE pg_catalog."default" NOT NULL,
    "URI" character varying COLLATE pg_catalog."default" NOT NULL,
    "TenantIdentifier" character varying COLLATE pg_catalog."default" NOT NULL,
    "CreatedOn" date NOT NULL,
    "ModifiedOn" date,
    "ProviderName" character varying COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT "StaticStorageIdentifier" PRIMARY KEY ("Identifier")
)
WITH (
    OIDS = FALSE
);

CREATE TABLE IF NOT EXISTS public."Tenant"
(
    "Identifier" character varying COLLATE pg_catalog."default" NOT NULL,
    "Name" character varying COLLATE pg_catalog."default" NOT NULL,
    "CreatedOn" date NOT NULL,
    "ModifiedOn" date,
    "AdminEmail" character varying COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT "TenantIdentifier" PRIMARY KEY ("Identifier")
)
WITH (
    OIDS = FALSE
);

CREATE TABLE IF NOT EXISTS public."User"
(
    "Identifier" character varying COLLATE pg_catalog."default" NOT NULL,
    "Email" character varying COLLATE pg_catalog."default" NOT NULL,
    "Password" character varying COLLATE pg_catalog."default",
    "Name" character varying COLLATE pg_catalog."default" NOT NULL,
    "IsExternalOAuthUser" boolean NOT NULL DEFAULT false,
    "Locale" character varying COLLATE pg_catalog."default" NOT NULL,
    "TenantIdentifier" character varying COLLATE pg_catalog."default" NOT NULL,
    "CreatedOn" date NOT NULL,
    "ModifiedOn" date,
    "DeletedAt" date,
    "Metadata" jsonb,
    "ProfilePictureURL" character varying COLLATE pg_catalog."default",
    "OAuthProvider" character varying COLLATE pg_catalog."default",
    CONSTRAINT "UserIdentifier" PRIMARY KEY ("Identifier")
)
WITH (
    OIDS = FALSE
);

CREATE TABLE IF NOT EXISTS public."UserMetadata"
(
    "Identifier" character varying COLLATE pg_catalog."default" NOT NULL,
    "Metadata" jsonb NOT NULL,
    "TenantIdentifier" character varying COLLATE pg_catalog."default" NOT NULL,
    "EntityName" character varying COLLATE pg_catalog."default" NOT NULL,
    "CreatedOn" date NOT NULL,
    "ModifiedOn" date,
    CONSTRAINT "UserMetadataIdentifier" PRIMARY KEY ("Identifier")
)
WITH (
    OIDS = FALSE
);

CREATE TABLE IF NOT EXISTS public."UserRole"
(
    "Identifier" character varying COLLATE pg_catalog."default" NOT NULL,
    "UserIdentifier" character varying COLLATE pg_catalog."default" NOT NULL,
    "CreatedOn" date NOT NULL,
    "ModifiedOn" date,
    "RoleIdentifier" character varying COLLATE pg_catalog."default" NOT NULL,
    "Scopes" character varying[] COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT "RoleIdentifier" PRIMARY KEY ("Identifier")
)
WITH (
    OIDS = FALSE
);

CREATE TABLE IF NOT EXISTS public."Webhook"
(
    "Identifier" character varying COLLATE pg_catalog."default" NOT NULL,
    "PostURL" character varying COLLATE pg_catalog."default" NOT NULL,
    "TriggeringEvent" character varying COLLATE pg_catalog."default" NOT NULL,
    "User" character varying COLLATE pg_catalog."default",
    "Password" character varying COLLATE pg_catalog."default",
    "TenantIdentifier" character varying COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT "WebhookIdentifier" PRIMARY KEY ("Identifier")
)
WITH (
    OIDS = FALSE
);

CREATE TABLE IF NOT EXISTS public."WebhookExecution"
(
    "Identifier" character varying COLLATE pg_catalog."default" NOT NULL,
    "StatusCode" character varying COLLATE pg_catalog."default" NOT NULL,
    "ResponseMessage" character varying COLLATE pg_catalog."default",
    "WebhookIdentifier" character varying COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT "WebhookExecutionIdentifier" PRIMARY KEY ("Identifier")
)
WITH (
    OIDS = FALSE
);

IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'zapprdb') THEN

    ALTER TABLE IF EXISTS public."APIKey"
        ADD CONSTRAINT "APIKeyTenantIdentifier" FOREIGN KEY ("TenantIdentifier")
        REFERENCES public."Tenant" ("Identifier") MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION;


    ALTER TABLE IF EXISTS public."APIKey"
        ADD CONSTRAINT "APIKeyUserIdentifier" FOREIGN KEY ("UserIdentifier")
        REFERENCES public."User" ("Identifier") MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID;


    ALTER TABLE IF EXISTS public."ExternalConnector"
        ADD CONSTRAINT "ExternalConnectorTenantIdentifier" FOREIGN KEY ("TenantIdentifier")
        REFERENCES public."Tenant" ("Identifier") MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION;


    ALTER TABLE IF EXISTS public."ExternalConnectorExecution"
        ADD CONSTRAINT "ExternalConnectorExecutionConnectorIdentifier" FOREIGN KEY ("ExternalConnectorIdentifier")
        REFERENCES public."ExternalConnector" ("Identifier") MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION;


    ALTER TABLE IF EXISTS public."OAuthProvider"
        ADD CONSTRAINT "OAuthProviderTenantIdentifier" FOREIGN KEY ("TenantIdentifier")
        REFERENCES public."Tenant" ("Identifier") MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID;


    ALTER TABLE IF EXISTS public."Role"
        ADD CONSTRAINT "TenantIdentifier" FOREIGN KEY ("TenantIdentifier")
        REFERENCES public."Tenant" ("Identifier") MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID;


    ALTER TABLE IF EXISTS public."SMTPSettings"
        ADD CONSTRAINT "SMTPSettingsTenantIdentifier" FOREIGN KEY ("TenantIdentifier")
        REFERENCES public."Tenant" ("Identifier") MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION;


    ALTER TABLE IF EXISTS public."StaticStorage"
        ADD CONSTRAINT "StaticStorageTenantIdentifier" FOREIGN KEY ("Identifier")
        REFERENCES public."Tenant" ("Identifier") MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION;
    CREATE INDEX IF NOT EXISTS "StaticStorageIdentifierIndex"
        ON public."StaticStorage"("Identifier");


    ALTER TABLE IF EXISTS public."User"
        ADD CONSTRAINT "UserTenantIdentifier" FOREIGN KEY ("TenantIdentifier")
        REFERENCES public."Tenant" ("Identifier") MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION;


    ALTER TABLE IF EXISTS public."UserMetadata"
        ADD CONSTRAINT "UserMetadataTenantIdentifier" FOREIGN KEY ("TenantIdentifier")
        REFERENCES public."Tenant" ("Identifier") MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION;


    ALTER TABLE IF EXISTS public."UserRole"
        ADD CONSTRAINT "MasterRoleIdentifier" FOREIGN KEY ("RoleIdentifier")
        REFERENCES public."Role" ("Identifier") MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID;


    ALTER TABLE IF EXISTS public."UserRole"
        ADD CONSTRAINT "RoleUserIdentifier" FOREIGN KEY ("UserIdentifier")
        REFERENCES public."User" ("Identifier") MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION;


    ALTER TABLE IF EXISTS public."Webhook"
        ADD CONSTRAINT "WebhookTenantIdentifier" FOREIGN KEY ("TenantIdentifier")
        REFERENCES public."Tenant" ("Identifier") MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION;


    ALTER TABLE IF EXISTS public."WebhookExecution"
        ADD CONSTRAINT "WebhookExecutionWebhookIdentifier" FOREIGN KEY ("WebhookIdentifier")
        REFERENCES public."Webhook" ("Identifier") MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION;
END IF;

END
$do$;
