CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(30) NOT NULL,
    last_name VARCHAR(30) NOT NULL,
    dob VARCHAR(10) NOT NULL,
    email VARCHAR(50) NOT NULL UNIQUE,
    phone_number VARCHAR(20) UNIQUE,
    gender VARCHAR(10) CHECK (gender IN('male', 'female')),
    password VARCHAR NOT NULL,
    profile_image_url VARCHAR,
    address VARCHAR,
    type VARCHAR(10) CHECK (type IN('superadmin', 'guest', 'owner')) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users(
    first_name,
    last_name,
    dob,
    email,
    password,
    type
) VALUES(
    'Ibratbek', 
    'Muslimbekov',
    '1999-06-06', 
    'imuslimbekov1421@gmail.com', 
    '$2a$10$tRtne.jx/GQwgL8abzimLer63bn3HiciTBVDfCi5cUlCexFAzrZE6', 
    'superadmin'
);