SELECT password
from public.users
WHERE user_id = $1;