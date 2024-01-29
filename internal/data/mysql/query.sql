/* name: GetCurrentTemplate :one */
SELECT * FROM last_updated_template
ORDER BY id DESC
LIMIT 1;

/* name: ListAllLastUpdatedTemplate :many */
SELECT * FROM consent_template
WHERE id > ?
ORDER BY id DESC;

/* name: CreateCurrentTemplate :execresult */
INSERT INTO last_updated_template (
    consent_template_id, template_name, version
) VALUES (
    ?, ?, ?
);
