-- Create a new user called "arl_user" and remove the DBA/Super User privs
-- NOTE: Change the password ('test1') in the quotes below!!!
CREATE USER arl_user with nosuperuser nocreatedb nocreaterole noreplication login inherit encrypted password 'test1';

-- Prevent the "public" role connecting to the 'arl' database
REVOKE connect ON database arl FROM public;
-- Allow the "arl_user" to connect to the database
GRANT connect on database arl to arl_user;

-- Change the table owner to "arl_user"
ALTER TABLE public.alert OWNER TO arl_user;
ALTER TABLE public.current_autoruns OWNER TO arl_user;
ALTER TABLE public.export OWNER TO arl_user;
ALTER TABLE public.instance OWNER TO arl_user;
ALTER TABLE public.previous_autoruns OWNER TO arl_user;
ALTER TABLE public.classification OWNER TO arl_user;

-- Allow the "arl_user" to perform the various actions on the tables
GRANT ALL ON TABLE public.alert TO arl_user;
GRANT ALL ON TABLE public.current_autoruns TO arl_user;
GRANT ALL ON TABLE public.export TO arl_user;
GRANT ALL ON TABLE public.instance TO arl_user;
GRANT ALL ON TABLE public.previous_autoruns TO arl_user;
GRANT ALL ON TABLE public.classification TO arl_user;
