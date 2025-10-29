-- name: BulkInsertOptionValues :many
INSERT INTO
  option_values (name, no, option_id)
SELECT
  unnest(@names::text[]),
  unnest(@nos::int[]),
  unnest(@option_ids::bigint[])
RETURNING
  id,
  name,
  option_id;

-- Get all option values for a list of options
-- name: GetOptionValuesByOptionIDs :many
SELECT
  id,
  name,
  no,
  option_id
FROM
  option_values
WHERE
  option_id = ANY ($1::bigint[]);

-- name: BulkUpdateOptionValues :exec
UPDATE option_values AS ov
SET
  name = data.name
FROM (
  SELECT  UNNEST(@ids::bigint[]) AS id,
          UNNEST(@names::text[]) AS name
) AS data
WHERE ov.id = data.id
  AND (
    ov.name IS DISTINCT FROM data.name
  );

-- name: DeleteOptionValuesNotInIDs :exec
DELETE FROM option_values
WHERE option_id IN (SELECT UNNEST(@option_ids::bigint[])) 
  AND id NOT IN (SELECT UNNEST(@value_ids::bigint[]));
