UPDATE public.passwords
SET password = $1,
nonce = $2,
use_location = $3,
created_on = $4
WHERE user_id = $5
    AND password_id = $6;