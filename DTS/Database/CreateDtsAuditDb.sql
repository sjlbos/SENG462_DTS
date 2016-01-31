
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
	transaction_id uuid NOT NULL,
	user_id char[10] NOT NULL,
	service varchar[256] NOT NULL,
	server varchar[256] NOT NULL
);

CREATE FUNCTION add_base_event(
	_occured_at timestamptz,
	_type event_type,
	_transaction_id uuid,
	_user_id char[10],
	_service varchar[256],
	_server varchar[256]
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
	stock char[3],
	funds money,
	FOREIGN KEY (id) REFERENCES events (id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE FUNCTION log_user_command_event(
	_occured_at timestamptz,
	_transaction_id uuid,
	_user_id char[10],
	_service varchar[256],
	_server varchar[256],
	_command command_type,
	_stock char[3],
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
	stock char[3] NOT NULL,
	price money NOT NULL,
	quote_server_time timestamptz NOT NULL,
	cryptokey varchar[64],
	FOREIGN KEY (id) REFERENCES events (id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE FUNCTION log_quote_server_event(
	_occured_at timestamptz,
	_transaction_id uuid,
	_user_id char[10],
	_service varchar[256],
	_server varchar[256],
	_stock char[3],
	_price money,
	_quote_server_time timestamptz,
	_cryptokey varchar[64]
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

CREATE FUNCTION log_account_transaction_event(
	_occured_at timestamptz,
	_transaction_id uuid,
	_user_id char[10],
	_service varchar[256],
	_server varchar[256],
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
	stock char[3],
	funds money,
	filename varchar[1024],
	FOREIGN KEY (id) REFERENCES events (id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE FUNCTION log_system_event(
	_occured_at timestamptz,
	_transaction_id uuid,
	_user_id char[10],
	_service varchar[256],
	_server varchar[256],
	_command command_type,
	_stock char[3],
	_funds money,
	_filename varchar[1024]
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
	stock char[3],
	funds money,
	error_message varchar[1024],
	filename varchar[1024],
	FOREIGN KEY (id) REFERENCES events (id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE FUNCTION log_error_event(
	_occured_at timestamptz,
	_transaction_id uuid,
	_user_id char[10],
	_service varchar[256],
	_server varchar[256],
	_command command_type,
	_stock char[3],
	_funds money,
	_error_message varchar[1024],
	_filename varchar[1024]
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
	stock char[3],
	funds money,
	filename varchar[1024],
	debug_message varchar[2056],
	FOREIGN KEY (id) REFERENCES events (id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE FUNCTION log_debug_event(
	_occured_at timestamptz,
	_transaction_id uuid,
	_user_id char[10],
	_service varchar[256],
	_server varchar[256],
	_command command_type,
	_stock char[3],
	_funds money,
	_filename varchar[1024],
	_debug_message varchar[1024]
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



