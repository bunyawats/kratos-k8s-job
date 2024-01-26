/* name: GetCurrentTemplate :one */
SELECT * FROM current_template
WHERE id = ? LIMIT 1;

/* name: CreateCurrentTemplate :execresult */
INSERT INTO current_template (
    template_name, version
) VALUES (
    ?, ?
);
