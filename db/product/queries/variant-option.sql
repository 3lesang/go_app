-- name: BulkInsertVariantOption :exec
INSERT INTO
  variant_options (variant_id, option_id, option_value_id)
SELECT
  unnest(@variant_ids::bigint[]),
  unnest(@option_ids::bigint[]),
  unnest(@option_value_ids::bigint[]);

-- name: GetVariantOptionsByVariantIDs :many
SELECT
  vo.variant_id,
  o.id AS option_id,
  o.name AS option_name,
  ov.id AS value_id,
  ov.name AS value_name
FROM
  variant_options vo
  JOIN options o ON o.id = vo.option_id
  JOIN option_values ov ON ov.id = vo.option_value_id
WHERE
  vo.variant_id = ANY ($1::bigint[])
ORDER BY
  vo.variant_id,
  o.no;
