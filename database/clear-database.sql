TRUNCATE TABLE alert;
ALTER SEQUENCE alert_id_seq RESTART WITH 1;

TRUNCATE TABLE current_autoruns;
ALTER SEQUENCE current_autoruns_id_seq RESTART WITH 1;

TRUNCATE TABLE instance;
ALTER SEQUENCE instance_id_seq RESTART WITH 1;

TRUNCATE TABLE previous_autoruns;
ALTER SEQUENCE previous_autoruns_id_seq RESTART WITH 1;