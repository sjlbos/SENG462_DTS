DROP TYPE IF EXISTS event_type CASCADE;
DROP TYPE IF EXISTS command_type CASCADE;
DROP TYPE IF EXISTS account_action CASCADE;
DROP TABLE IF EXISTS events CASCADE;
DROP TABLE IF EXISTS user_command_events CASCADE;
DROP TABLE IF EXISTS account_transaction_events CASCADE;
DROP TABLE IF EXISTS system_events CASCADE;
DROP TABLE IF EXISTS error_events CASCADE;
DROP TABLE IF EXISTS debug_events CASCADE;
DROP TABLE IF EXISTS quote_server_events CASCADE;



CREATE TYPE event_type AS ENUM(
	'UserCommandEvent',
	'QuoteServerEvent',
	'AccountTransactionEvent',
	'SystemEvent',
	'ErrorEvent',
	'DebugEvent'
);

CREATE TYPE command_type AS ENUM(
	'ADD',
	'QUOTE',
	'BUY',
	'COMMIT_BUY',
	'CANCEL_BUY',
	'SELL',
	'COMMIT_SELL',
	'CANCEL_SELL',
	'SET_BUY_AMOUNT',
	'CANCEL_SET_BUY',
	'SET_BUY_TRIGGER',
	'SET_SELL_AMOUNT',
	'SET_SELL_TRIGGER',
	'CANCEL_SET_SELL',
	'DUMPLOG',
	'DISPLAY_SUMMARY'
);

CREATE TYPE account_action AS ENUM(
	'add',
	'remove'
);

CREATE TABLE events(
	id SERIAL PRIMARY KEY,
	logged_at timestamptz NOT NULL DEFAULT current_timestamp,
	occured_at timestamptz NOT NULL,
	type event_type NOT NULL,
	transaction_id int NOT NULL,
	user_id char(10),
	service varchar(256) NOT NULL,
	server varchar(256) NOT NULL
);

CREATE OR REPLACE FUNCTION add_base_event(
	_occured_at timestamptz,
	_type event_type,
	_transaction_id int,
	_user_id char(10),
	_service varchar(256),
	_server varchar(256)
)
RETURNS int AS
$$
	INSERT INTO events(
		occured_at, 
		type,
		transaction_id,
		user_id,
		service,
		server)
	VALUES (
		_occured_at,
		_type,
		_transaction_id,
		_user_id,
		_service,
		_server
		)
	RETURNING id;
$$
LANGUAGE SQL VOLATILE;



CREATE TABLE user_command_events(
	id int NOT NULL PRIMARY KEY,
	command command_type NOT NULL,
	stock char(3),
	funds money,
	FOREIGN KEY (id) REFERENCES events (id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE OR REPLACE FUNCTION log_user_command_event(
	_occured_at timestamptz,
	_transaction_id int,
	_user_id char(10),
	_service varchar(256),
	_server varchar(256),
	_command command_type,
	_stock char(3),
	_funds money
)
RETURNS int AS
$$
DECLARE 
	event_id int;
BEGIN
	event_id = add_base_event(
		_occured_at,
		'UserCommandEvent',
		_transaction_id,
		_user_id,
		_service,
		_server
		);

	INSERT INTO user_command_events(
		id,
		command,
		stock,
		funds)
	VALUES(		
		event_id,
		_command,
		_stock,
		_funds);

	RETURN event_id;
END;
$$
LANGUAGE 'plpgsql' VOLATILE;


CREATE TABLE quote_server_events(
	id int NOT NULL PRIMARY KEY,
	stock char(3) NOT NULL,
	price money NOT NULL,
	quote_server_time timestamptz NOT NULL,
	cryptokey varchar(64),
	FOREIGN KEY (id) REFERENCES events (id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE OR REPLACE FUNCTION log_quote_server_event(
	_occured_at timestamptz,
	_transaction_id int,
	_user_id char(10),
	_service varchar(256),
	_server varchar(256),
	_stock char(3),
	_price money,
	_quote_server_time timestamptz,
	_cryptokey varchar(64)
)
RETURNS int AS
$$
DECLARE 
	event_id int;
BEGIN
	event_id = add_base_event(
		_occured_at,
		'QuoteServerEvent',
		_transaction_id,
		_user_id,
		_service,
		_server
		);

	INSERT INTO quote_server_events(
		id,
		stock,
		price,
		quote_server_time,
		cryptokey)
	VALUES(		
		event_id,
		_stock,
		_price,
		_quote_server_time,
		_cryptokey);

	RETURN event_id;
END;
$$
LANGUAGE 'plpgsql' VOLATILE;



CREATE TABLE account_transaction_events(
	id int NOT NULL PRIMARY KEY,
	action account_action NOT NULL,
	funds money NOT NULL,
	FOREIGN KEY (id) REFERENCES events (id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE OR REPLACE FUNCTION log_account_transaction_event(
	_occured_at timestamptz,
	_transaction_id int,
	_user_id char(10),
	_service varchar(256),
	_server varchar(256),
	_action account_action,
	_funds money
)
RETURNS int AS
$$
DECLARE 
	event_id int;
BEGIN
	event_id = add_base_event(
		_occured_at,
		'AccountTransactionEvent',
		_transaction_id,
		_user_id,
		_service,
		_server
		);

	INSERT INTO account_transaction_events(
		id,
		action,
		funds)
	VALUES(		
		event_id,
		_action,
		_funds);

	RETURN event_id;
END;
$$
LANGUAGE 'plpgsql' VOLATILE;



CREATE TABLE system_events(
	id int NOT NULL PRIMARY KEY,
	command command_type NOT NULL,
	stock char(3),
	funds money,
	filename varchar(1024),
	FOREIGN KEY (id) REFERENCES events (id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE OR REPLACE FUNCTION log_system_event(
	_occured_at timestamptz,
	_transaction_id int,
	_user_id char(10),
	_service varchar(256),
	_server varchar(256),
	_command command_type,
	_stock char(3),
	_funds money,
	_filename varchar(1024)
)
RETURNS int AS
$$
DECLARE 
	event_id int;
BEGIN
	event_id = add_base_event(
		_occured_at,
		'SystemEvent',
		_transaction_id,
		_user_id,
		_service,
		_server
		);

	INSERT INTO system_events(
		id,
		command,
		stock,
		funds,
		filename)
	VALUES(		
		event_id,
		_command,
		_stock,
		_funds,
		_filename);

	RETURN event_id;
END;
$$
LANGUAGE 'plpgsql' VOLATILE;



CREATE TABLE error_events(
	id int NOT NULL PRIMARY KEY,
	command command_type NOT NULL,
	stock char(3),
	funds money,
	error_message varchar(2048),
	filename varchar(1024),
	FOREIGN KEY (id) REFERENCES events (id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE OR REPLACE FUNCTION log_error_event(
	_occured_at timestamptz,
	_transaction_id int,
	_user_id char(10),
	_service varchar(256),
	_server varchar(256),
	_command command_type,
	_stock char(3),
	_funds money,
	_error_message varchar(2048),
	_filename varchar(1024)
)
RETURNS int AS
$$
DECLARE 
	event_id int;
BEGIN
	event_id = add_base_event(
		_occured_at,
		'ErrorEvent',
		_transaction_id,
		_user_id,
		_service,
		_server
		);

	INSERT INTO error_events(
		id,
		command,
		stock,
		funds,
		error_message,
		filename)
	VALUES(		
		event_id,
		_command,
		_stock,
		_funds,
		_error_message,
		_filename);

	RETURN event_id;
END;
$$
LANGUAGE 'plpgsql' VOLATILE;



CREATE TABLE debug_events(
	id int NOT NULL PRIMARY KEY,
	command command_type NOT NULL,
	stock char(3),
	funds money,
	filename varchar(1024),
	debug_message varchar(2048),
	FOREIGN KEY (id) REFERENCES events (id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE OR REPLACE FUNCTION log_debug_event(
	_occured_at timestamptz,
	_transaction_id int,
	_user_id char(10),
	_service varchar(256),
	_server varchar(256),
	_command command_type,
	_stock char(3),
	_funds money,
	_filename varchar(1024),
	_debug_message varchar(2048)
)
RETURNS int AS
$$
DECLARE 
	event_id int;
BEGIN
	event_id = add_base_event(
		_occured_at,
		'DebugEvent',
		_transaction_id,
		_user_id,
		_service,
		_server
		);

	INSERT INTO debug_events(
		id,
		command,
		stock,
		funds,
		filename,
		debug_message)
	VALUES(		
		event_id,
		_command,
		_stock,
		_funds,
		_filename,
		_debug_message);

	RETURN event_id;
END;
$$
LANGUAGE 'plpgsql' VOLATILE;


CREATE OR REPLACE FUNCTION get_all_events(_start timestamptz, _end timestamptz)
RETURNS TABLE(	id int, 
				logged_at timestamptz, 
				occured_at timestamptz,
				type event_type,
				transaction_id int,
				user_id char(10),
				service char(256),
				server char(256),
				command command_type,
				stock char(3),
				funds money,
				filename varchar(1024),
				message varchar(2048),
				action account_action,
				quote_server_time timestamptz,
				cryptokey varchar(64)
			) AS
$$
	WITH base_events AS(
		SELECT 	id,
				logged_at,
				occured_at,
				type,
				transaction_id,
				user_id,
				service,
				server
		FROM events
		WHERE occured_at BETWEEN _start AND _end
		)
	SELECT 	e.id, 
			e.logged_at, 
			e.occured_at, 
			e.type, 
			e.transaction_id, 
			e.user_id, 
			e.service, 
			e.server, 
			uc.command, 
			uc.stock, 
			uc.funds,
			CAST(NULL AS varchar(1024)),
			CAST(NULL AS varchar(2056)),
			CAST(NULL AS account_action),
			CAST(NULL AS timestamptz),
			CAST(NULL AS varchar(64))
	FROM base_events e
	INNER JOIN user_command_events uc ON e.id = uc.id

	UNION ALL 

		SELECT 	e.id, 
			e.logged_at, 
			e.occured_at, 
			e.type, 
			e.transaction_id, 
			e.user_id, 
			e.service, 
			e.server, 
			NULL,
			qs.stock,
			qs.price,
			NULL,
			NULL,
			NULL,
			qs.quote_server_time,
			qs.cryptokey
	FROM base_events e
	INNER JOIN quote_server_events qs ON e.id = qs.id

	UNION ALL

	SELECT 	e.id, 
			e.logged_at, 
			e.occured_at, 
			e.type, 
			e.transaction_id, 
			e.user_id, 
			e.service, 
			e.server, 
			NULL,
			NULL,
			at.funds,
			NULL,
			NULL,
			at.action,
			NULL,
			NULL
	FROM base_events e
	INNER JOIN account_transaction_events at ON e.id = at.id

	UNION ALL

	SELECT 	e.id, 
			e.logged_at, 
			e.occured_at, 
			e.type, 
			e.transaction_id, 
			e.user_id, 
			e.service, 
			e.server, 
			se.command, 
			se.stock, 
			se.funds,
			se.filename,
			NULL,
			NULL,
			NULL,
			NULL
	FROM base_events e
	INNER JOIN system_events se ON e.id = se.id

	UNION ALL

	SELECT 	e.id, 
			e.logged_at, 
			e.occured_at, 
			e.type, 
			e.transaction_id, 
			e.user_id, 
			e.service, 
			e.server, 
			er.command, 
			er.stock, 
			er.funds,
			er.filename,
			er.error_message,
			NULL,
			NULL,
			NULL
	FROM base_events e
	INNER JOIN error_events er ON e.id = er.id

	UNION ALL

	SELECT 	e.id, 
			e.logged_at, 
			e.occured_at, 
			e.type, 
			e.transaction_id, 
			e.user_id, 
			e.service, 
			e.server, 
			de.command, 
			de.stock, 
			de.funds,
			de.filename,
			de.debug_message,
			NULL,
			NULL,
			NULL
	FROM base_events e
	INNER JOIN debug_events de ON e.id = de.id

	ORDER BY occured_at DESC;
$$
LANGUAGE SQL VOLATILE;


CREATE OR REPLACE FUNCTION get_all_events_by_user(_user_id char(10), _start timestamptz, _end timestamptz)
RETURNS TABLE(	id int, 
				logged_at timestamptz, 
				occured_at timestamptz,
				type event_type,
				transaction_id int,
				user_id char(10),
				service char(256),
				server char(256),
				command command_type,
				stock char(3),
				funds money,
				filename varchar(1024),
				message varchar(2048),
				action account_action,
				quote_server_time timestamptz,
				cryptokey varchar(64)
			) AS
$$
	WITH base_events AS(
		SELECT 	id,
				logged_at,
				occured_at,
				type,
				transaction_id,
				user_id,
				service,
				server
		FROM events
		WHERE user_id = _user_id
			AND occured_at BETWEEN _start AND _end
		)
	SELECT 	e.id, 
			e.logged_at, 
			e.occured_at, 
			e.type, 
			e.transaction_id, 
			e.user_id, 
			e.service, 
			e.server, 
			uc.command, 
			uc.stock, 
			uc.funds,
			CAST(NULL AS varchar(1024)),
			CAST(NULL AS varchar(2056)),
			CAST(NULL AS account_action),
			CAST(NULL AS timestamptz),
			CAST(NULL AS varchar(64))
	FROM base_events e
	INNER JOIN user_command_events uc ON e.id = uc.id

	UNION ALL 

		SELECT 	e.id, 
			e.logged_at, 
			e.occured_at, 
			e.type, 
			e.transaction_id, 
			e.user_id, 
			e.service, 
			e.server, 
			NULL,
			qs.stock,
			qs.price,
			NULL,
			NULL,
			NULL,
			qs.quote_server_time,
			qs.cryptokey
	FROM base_events e
	INNER JOIN quote_server_events qs ON e.id = qs.id

	UNION ALL

	SELECT 	e.id, 
			e.logged_at, 
			e.occured_at, 
			e.type, 
			e.transaction_id, 
			e.user_id, 
			e.service, 
			e.server, 
			NULL,
			NULL,
			at.funds,
			NULL,
			NULL,
			at.action,
			NULL,
			NULL
	FROM base_events e
	INNER JOIN account_transaction_events at ON e.id = at.id

	UNION ALL

	SELECT 	e.id, 
			e.logged_at, 
			e.occured_at, 
			e.type, 
			e.transaction_id, 
			e.user_id, 
			e.service, 
			e.server, 
			se.command, 
			se.stock, 
			se.funds,
			se.filename,
			NULL,
			NULL,
			NULL,
			NULL
	FROM base_events e
	INNER JOIN system_events se ON e.id = se.id

	UNION ALL

	SELECT 	e.id, 
			e.logged_at, 
			e.occured_at, 
			e.type, 
			e.transaction_id, 
			e.user_id, 
			e.service, 
			e.server, 
			er.command, 
			er.stock, 
			er.funds,
			er.filename,
			er.error_message,
			NULL,
			NULL,
			NULL
	FROM base_events e
	INNER JOIN error_events er ON e.id = er.id

	UNION ALL

	SELECT 	e.id, 
			e.logged_at, 
			e.occured_at, 
			e.type, 
			e.transaction_id, 
			e.user_id, 
			e.service, 
			e.server, 
			de.command, 
			de.stock, 
			de.funds,
			de.filename,
			de.debug_message,
			NULL,
			NULL,
			NULL
	FROM base_events e
	INNER JOIN debug_events de ON e.id = de.id

	ORDER BY occured_at DESC;
$$
LANGUAGE SQL VOLATILE;

