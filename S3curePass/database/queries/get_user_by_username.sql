SELECT user_id, firstname, lastname, username, email, password, registered_on
FROM public.users
WHERE username = $1;