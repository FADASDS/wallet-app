CREATE TABLE IF NOT EXISTS public.transactions_tbl
(
    tx_id uuid NOT NULL DEFAULT gen_random_uuid(),
    from_wallet uuid NOT NULL,
    to_wallet uuid NOT NULL,
    crtn_date timestamp with time zone NOT NULL DEFAULT now(),
    amount numeric(13,2) NOT NULL,
    CONSTRAINT transactions_tbl_pkey PRIMARY KEY (tx_id)
)

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.transactions_tbl
    OWNER to postgres;

CREATE TABLE IF NOT EXISTS public.wallet_tbl
(
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    wallet_owner character varying(50) COLLATE pg_catalog."default",
    balance numeric(13,2) NOT NULL,
    CONSTRAINT wallet_tbl_pkey PRIMARY KEY (id),
    CONSTRAINT wallet_tbl_balance_check CHECK (balance >= 0::numeric)
    )

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.wallet_tbl
    OWNER to postgres;

INSERT INTO public.wallet_tbl (id, wallet_owner, balance) VALUES
(gen_random_uuid(),'Nikolai Petrov',100.00),
(gen_random_uuid(),'Ivan Ivanov',100.00),
(gen_random_uuid(),'Gleb Petrov',100.00),
(gen_random_uuid(),'Alexander Sidorov',100.00),
(gen_random_uuid(),'Evgeniy Gilin',100.00),
(gen_random_uuid(),'Mister Incognito',100.00),
(gen_random_uuid(),'Bob Davis',100.00),
(gen_random_uuid(),'James Cameron',100.00),
(gen_random_uuid(),'David Den',100.00),
(gen_random_uuid(),'Michael Braiton',100.00)