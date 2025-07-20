/* name: GetBook :one */
SELECT * FROM books
WHERE book_id = ?;

/* name: CreateBook :execresult */
INSERT INTO books (
    author_id,
    isbn,
    book_type,
    title,
    yr,
    available,
    tags
) VALUES (
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?
);
