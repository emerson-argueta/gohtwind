CREATE TABLE products (
    id bigint auto_increment PRIMARY KEY,
    title VARCHAR (50) UNIQUE NOT NULL,
    price numeric(10,2) NOT NULL,
    sku VARCHAR (50) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

INSERT INTO products (title, price, sku, created_at, updated_at)
VALUES
    (
        'Cheese',
        9.99,
        'x12334',
        NOW(),
        NOW()
    ),
    (
        'Bread',
        1.99,
        'x12335',
        NOW(),
        NOW()
    ),
    (
        'Milk',
        2.99,
        'x56788',
        NOW(),
        NOW()
    ),
    (
        'Eggs',
        1.99,
        'x12336',
        NOW(),
        NOW()
    ),
    (
        'Chicken',
        5.99,
        'x12337',
        NOW(),
        NOW()
    );