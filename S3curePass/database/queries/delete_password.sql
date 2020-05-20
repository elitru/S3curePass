DELETE FROM public.passwords
WHERE password_id = $1
    AND user_id = $2;