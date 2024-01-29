/* name: GetCurrentTemplate :one */
SELECT * FROM last_updated_template
ORDER BY consent_template_id DESC
LIMIT 1;

/* name: ListAllLastUpdatedTemplate :many */
SELECT * FROM consent_template
WHERE id > ?
ORDER BY id;

/* name: CreateCurrentTemplate :execresult */
INSERT INTO last_updated_template (
    consent_template_id, template_name, version
) VALUES (
    ?, ?, ?
);
