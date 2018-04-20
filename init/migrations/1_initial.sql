-- +migrate Up
CREATE TABLE "Association"
(
  id UUID PRIMARY KEY NOT NULL,
  organisationid UUID NOT NULL,
  version INT NOT NULL,
  isdeleted BOOLEAN NOT NULL,
  islocked BOOLEAN NOT NULL,
  paginationId SERIAL,
  record JSONB
);
CREATE UNIQUE INDEX ON "Association" (id);
CREATE UNIQUE INDEX Association_organisationid ON "Association" (organisationid);
