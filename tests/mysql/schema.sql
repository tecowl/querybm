-- https://github.com/sqlc-dev/sqlc/blob/main/examples/booktest/mysql/schema.sql
CREATE TABLE authors (
          author_id integer NOT NULL AUTO_INCREMENT PRIMARY KEY,
          name text NOT NULL,
          INDEX idx_authors_name (name(255))
) ENGINE=InnoDB;

CREATE TABLE books (
          book_id integer NOT NULL AUTO_INCREMENT PRIMARY KEY,
          author_id integer NOT NULL,
          isbn varchar(255) NOT NULL DEFAULT '' UNIQUE,
          book_type ENUM('MAGAZINE', 'PAPERBACK', 'HARDCOVER') NOT NULL DEFAULT 'PAPERBACK',
          title text NOT NULL,
          yr integer NOT NULL DEFAULT 2000,
          available datetime NOT NULL DEFAULT NOW(),
          tags text NOT NULL,
          INDEX idx_books_title_yr (title(255), yr),
          CONSTRAINT fk_books_author_id FOREIGN KEY (author_id) REFERENCES authors(author_id)
) ENGINE=InnoDB;
