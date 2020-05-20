-- database -> "securepass"
-- postgres needs extension -> "CREATE EXTENSION IF NOT EXISTS "uuid-ossp";"
CREATE TABLE public.users (
    user_id uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    firstname character varying(50) NOT NULL,
    lastname character varying(50) NOT NULL,
    username character varying(50) NOT NULL UNIQUE,
    password character varying(200) NOT NULL,
    email character varying(100) NOT NULL UNIQUE,
    registered_on timestamp NOT NULL DEFAULT NOW()
);

CREATE TABLE public.passwords (
    password_id uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    password character varying(300) NOT NULL,
    nonce character varying(100) NOT NULL,
    use_location character varying(200) NOT NULL,
    created_on timestamp NOT NULL DEFAULT NOW(),
    user_id uuid NOT NULL,
    CONSTRAINT user_id FOREIGN KEY (user_id)
        REFERENCES public.users (user_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
);

ALTER TABLE public.users
    OWNER to postgres;

ALTER TABLE public.passwords
    OWNER to postgres;