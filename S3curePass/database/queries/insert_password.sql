INSERT INTO public.passwords(password, nonce, use_location, created_on, user_id)
VALUES ($1, $2, $3, $4, $5);