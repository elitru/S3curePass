SELECT pw.password_id, pw.password, pw.nonce, pw.use_location, pw.created_on, pw.user_id
FROM public.passwords pw
WHERE pw.user_id = $1
    AND use_location LIKE '%' || $2 || '%';