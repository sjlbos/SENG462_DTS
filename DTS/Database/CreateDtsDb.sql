CREATE TABLE users(
	id SERIAL PRIMARY KEY,
	user_id char[10] NOT NULL UNIQUE,
	balance money NOT NULL DEFAULT 0.00 CHECK(balance > 0::money),
	created_at timestamptz NOT NULL
);

CREATE TABLE portfolios(
	uid int NOT NULL,
	stock char[3] NOT NULL,
	num_shares int NOT NULL CHECK(num_shares > 0),
	FOREIGN KEY (uid) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE,
	UNIQUE(uid, stock)
);

CREATE TYPE trigger_type AS ENUM ('buy', 'sell');

CREATE TABLE triggers(
	id SERIAL PRIMARY KEY,
	uid int NOT NULL,
	stock char[3] NOT NULL,
	type trigger_type NOT NULL,
	trigger_price money CHECK(trigger_price >= 0::money),
	num_shares int NOT NULL CHECK(num_shares > 0),
	created_at timestamptz NOT NULL,
	FOREIGN KEY (uid) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE,
	UNIQUE (uid, stock, type)
);

CREATE TYPE transaction_type AS ENUM ('sale', 'purchase');

CREATE TABLE transactions(
	id SERIAL PRIMARY KEY,
	uid int NOT NULL,
	type transaction_type NOT NULL,
	stock char[3] NOT NULL,
	num_shares int NOT NULL CHECK(num_shares > 0),
	share_price money NOT NULL CHECK(share_price >= 0::money),
	made_at timestamptz NOT NULL,
	FOREIGN KEY (uid) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE pending_transactions(
	id SERIAL PRIMARY KEY,
	uid int NOT NULL,
	type transaction_type NOT NULL,
	stock char[3] NOT NULL,
	num_shares int NOT NULL CHECK(num_shares > 0),
	share_price money NOT NULL CHECK(share_price >= 0::money),
	requested_at timestamptz NOT NULL,
	expires_at timestamptz NOT NULL,
	FOREIGN KEY (uid) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE
);



CREATE FUNCTION get_reserved_funds (_uid int) RETURNS money AS
$$
	SELECT SUM(value)
	FROM(
		SELECT num_shares * trigger_price AS value
		FROM triggers
		WHERE uid = _uid AND type = 'buy'
		) AS trigger_values;
$$
LANGUAGE SQL VOLATILE;



CREATE FUNCTION add_user_account(_user_id char[10], _balance money)
RETURNS int AS
$$
	INSERT INTO users(user_id, balance)
	VALUES(_user_id, _balance)
	RETURNING id;
$$
LANGUAGE SQL VOLATILE;



CREATE FUNCTION update_user_account_balance(_uid int, _balance money) 
RETURNS void AS
$$
	UPDATE users
	SET balance = _balance
	WHERE id = _uid;
$$
LANGUAGE SQL VOLATILE;



CREATE FUNCTION get_user_account_by_char_id(_user_id char[10])
RETURNS TABLE(id int, user_id char[10], balance money) AS
$$
	SELECT 	id,
			user_id,
			balance
	FROM users
	WHERE user_id = _user_id;
$$
LANGUAGE SQL VOLATILE;



CREATE FUNCTION get_user_portfolio(_uid int)
RETURNS TABLE(stock char[3], num_shares int) AS
$$
	SELECT 	stock,
			num_shares
	FROM portfolios
	WHERE uid = _uid;
$$
LANGUAGE SQL VOLATILE;



CREATE FUNCTION get_user_stock_amount(_uid int, _stock char[3])
RETURNS int AS
$$ 
	SELECT num_shares
	FROM portfolios
	WHERE uid = _uid AND stock = _stock;
$$
LANGUAGE SQL VOLATILE;



CREATE FUNCTION perform_purchase_transaction(
	_uid int,
	_stock char[3],
	_num_shares int,
	_share_price money,
	_made_at timestamptz
)
RETURNS int AS
$$
DECLARE
	_transaction_id int;
	_existing_shares int;
BEGIN
	INSERT INTO transactions(
			uid,
			type, 
			stock,
			num_shares,
			share_price,
			made_at
		)
	VALUES(
		_uid,
		'purchase',
		_stock,
		_num_shares,
		_share_price,
		_made_at
		)
	RETURNING id INTO _transaction_id;

	UPDATE users
	SET balance = balance - (_num_shares * _share_price)
	WHERE id = _uid;

	_existing_shares = get_user_stock_amount(_uid, _stock);

	IF _existing_shares IS NULL THEN
		INSERT INTO portfolios (uid, stock, num_shares)
			VALUES (_uid, _stock, _num_shares);
	ELSE
		UPDATE portfolios SET num_shares = (_num_shares + _existing_shares) WHERE uid = _uid AND stock = _stock;
	END IF;

	RETURN _transaction_id;
END;
$$
LANGUAGE 'plpgsql' VOLATILE;



CREATE FUNCTION perform_sale_transaction(
	_uid int,
	_stock char[3],
	_num_shares int,
	_share_price money,
	_made_at timestamptz
)
RETURNS int AS
$$
DECLARE 
	_transaction_id int;
	_existing_shares int;
	_total_shares int;
BEGIN
	INSERT INTO transactions(
				uid,
				type, 
				stock,
				num_shares,
				share_price,
				made_at
			)
	VALUES(
		_uid,
		'sale',
		_stock,
		_num_shares,
		_share_price,
		_made_at
		)
	RETURNING id INTO _transaction_id;

	UPDATE users 
	SET balance = balance + (_share_price * _num_shares) 
	WHERE id = _uid;

	_existing_shares = get_user_stock_amount(_uid, _stock);
	_total_shares = _existing_shares - _num_shares;

	IF _total_shares = 0 THEN
		DELETE FROM portfolios WHERE uid = _uid AND stock = _stock;
	ELSE
		UPDATE portfolios SET num_shares = _total_shares WHERE uid = _uid AND stock = _stock;
	END IF;

	RETURN _transaction_id;
END;
$$
LANGUAGE 'plpgsql' VOLATILE;



CREATE FUNCTION get_user_triggers(_uid int)
RETURNS TABLE (	id int, 
				stock char[3], 
				type trigger_type, 
				trigger_price money, 
				num_shares int, 
				created_at timestamptz) AS
$$
	SELECT 	id,
			stock,
			type,
			trigger_price,
			num_shares,
			created_at
	FROM triggers
	WHERE uid = _uid;
$$
LANGUAGE SQL VOLATILE;



CREATE FUNCTION add_buy_trigger(
	_uid int,
	_stock char[3],
	_trigger_price money,
	_num_shares int,
	_created_at timestamptz
)
RETURNS int AS
$$
	INSERT INTO triggers(
		uid,
		stock,
		type,
		trigger_price,
		num_shares,
		created_at
		)
	VALUES(
		_uid,
		_stock,
		'buy',
		_trigger_price,
		_num_shares,
		_created_at
		)
	RETURNING id;
$$
LANGUAGE SQL VOLATILE;



CREATE FUNCTION commit_buy_trigger(
	_id int,
	_made_at timestamptz
)
RETURNS int AS
$$
DECLARE
	_transaction_id int;
BEGIN
	_transaction_id = perform_purchase_transaction((
		SELECT 	uid,
				stock,
				num_shares,
				share_price,
				_made_at
		FROM triggers
		WHERE id = _id
		));

	DELETE FROM triggers WHERE id = _id;

	RETURN _transaction_id;
END;
$$
LANGUAGE 'plpgsql' VOLATILE;



CREATE FUNCTION commit_sell_trigger(
	_id int,
	_made_at timestamptz
)
RETURNS int AS
$$
DECLARE
	_transaction_id int;
BEGIN
	_transaction_id = perform_sale_transaction((
		SELECT 	uid,
				stock,
				num_shares,
				share_price,
				_made_at
		FROM triggers
		WHERE id = _id
		));

	DELETE FROM triggers WHERE id = _id;

	RETURN _transaction_id;
END;
$$
LANGUAGE 'plpgsql' VOLATILE;



CREATE FUNCTION add_sell_trigger(
	_uid int,
	_stock char[3],
	_trigger_price money,
	_num_shares int,
	_created_at timestamptz
)
RETURNS int AS
$$
	INSERT INTO triggers(
		uid,
		stock,
		type,
		trigger_price,
		num_shares,
		created_at
		)
	VALUES(
		_uid,
		_stock,
		'sell',
		_trigger_price,
		_num_shares,
		_created_at
		)
	RETURNING id;
$$
LANGUAGE SQL VOLATILE;



CREATE FUNCTION set_trigger_price(
	_id int,
	_trigger_price money
)
RETURNS void AS 
$$
	UPDATE triggers 
	SET trigger_price = _trigger_price 
	WHERE id = _id;
$$
LANGUAGE SQL VOLATILE;



CREATE FUNCTION get_buy_trigger_for_user_and_stock(_uid int, _stock char[3])
RETURNS TABLE (id int, stock char[3], type trigger_type, trigger_price money, num_shares int, created_at timestamptz) AS
$$
	SELECT 	id,
			stock,
			type,
			trigger_price,
			num_shares,
			created_at
	FROM triggers
	WHERE uid = _uid
		AND stock = _stock
		AND type = 'buy';
$$
LANGUAGE SQL VOLATILE;



CREATE FUNCTION get_sell_trigger_for_user_and_stock(_uid int, _stock char[3])
RETURNS TABLE (id int, stock char[3], type trigger_type, trigger_price money, num_shares int, created_at timestamptz) AS
$$
	SELECT 	id,
			stock,
			type,
			trigger_price,
			num_shares,
			created_at
	FROM triggers
	WHERE uid = _uid 
		AND stock = _stock
		AND type = 'sell';
$$
LANGUAGE SQL VOLATILE;



CREATE FUNCTION cancel_trigger(
	_id int
)
RETURNS void AS
$$
	DELETE FROM triggers 
	WHERE id = _id;
$$
LANGUAGE SQL VOLATILE;



CREATE FUNCTION get_user_transaction_history(
	_uid int
)
RETURNS TABLE (id int, type transaction_type, stock char[3], num_shares int, share_price money, made_at timestamptz) AS
$$
	SELECT 	id,
			type,
			stock,
			num_shares,
			share_price,
			made_at
	FROM transactions
	WHERE uid = _uid;
$$
LANGUAGE SQL VOLATILE;



CREATE FUNCTION add_pending_purchase(
	_uid int,
	_stock char[3],
	_num_shares int,
	_share_price money,
	_requested_at timestamptz,
	_expires_at timestamptz
)
RETURNS int AS
$$
	INSERT INTO pending_transactions(
		uid,
		type,
		stock,
		num_shares,
		share_price,
		requested_at,
		expires_at
		)
	VALUES(
		_uid,
		'purchase',
		_stock,
		_num_shares,
		_share_price,
		_requested_at,
		_expires_at
		)
	RETURNING id;
$$
LANGUAGE SQL VOLATILE;



CREATE FUNCTION add_pending_sale(
	_uid int,
	_stock char[3],
	_num_shares int,
	_share_price money,
	_requested_at timestamptz,
	_expires_at timestamptz
)
RETURNS int AS
$$
	INSERT INTO pending_transactions(
		uid,
		type,
		stock,
		num_shares,
		share_price,
		requested_at,
		expires_at
		)
	VALUES(
		_uid,
		'sale',
		_stock,
		_num_shares,
		_share_price,
		_requested_at,
		_expires_at
		)
	RETURNING id;
$$
LANGUAGE SQL VOLATILE;



CREATE FUNCTION get_latest_pending_purchase_for_user(
	_uid int
)
RETURNS TABLE(id int, uid int, stock char[3], num_shares int, share_price money, requested_at timestamptz, expires_at timestamptz) AS
$$
	SELECT	id,
			uid,
			stock,
			num_shares,
			share_price,
			requested_at,
			expires_at
	FROM pending_transactions
	WHERE uid = _uid AND type = 'purchase'
	ORDER BY requested_at DESC
	LIMIT 1;
$$
LANGUAGE SQL VOLATILE;



CREATE FUNCTION get_all_pending_purchases_for_user(
	_uid int
)
RETURNS TABLE(id int, uid int, stock char[3], num_shares int, share_price money, requested_at timestamptz, expires_at timestamptz) AS
$$
	SELECT	id,
			uid,
			stock,
			num_shares,
			share_price,
			requested_at,
			expires_at
	FROM pending_transactions
	WHERE uid = _uid AND type = 'purchase'
	ORDER BY requested_at DESC
$$
LANGUAGE SQL VOLATILE;


CREATE FUNCTION commit_pending_purchase(
	_id int,
	_made_at timestamptz
)
RETURNS int AS
$$
DECLARE 
	_transaction_id int;
BEGIN
	_transaction_id = perform_purchase_transaction((
		SELECT 	uid,
				stock,
				num_shares,
				share_price,
				_made_at
		FROM pending_transactions
		WHERE id = _id
		));

	DELETE FROM pending_transactions WHERE id = _id;

	RETURN _transaction_id;
END;
$$
LANGUAGE 'plpgsql' VOLATILE;



CREATE FUNCTION get_latest_pending_sale(
	_user_id char[10]
)
RETURNS TABLE(id int, uid int, stock char[3], num_shares int, share_price money, requested_at timestamptz, expires_at timestamptz) AS
$$
DECLARE
	_uid int;
BEGIN
	SELECT id INTO _uid
	FROM users
	WHERE user_id = _user_id;

	RETURN QUERY
		SELECT	id,
				uid,
				stock,
				num_shares,
				share_price,
				requested_at,
				expires_at
		FROM pending_transactions
		WHERE uid = _uid AND type = 'sale'
		ORDER BY requested_at DESC
		LIMIT 1;
END;
$$
LANGUAGE 'plpgsql' VOLATILE;



CREATE FUNCTION get_all_pending_sales_for_user(
	_uid int
)
RETURNS TABLE(id int, uid int, stock char[3], num_shares int, share_price money, requested_at timestamptz, expires_at timestamptz) AS
$$
	SELECT	id,
			uid,
			stock,
			num_shares,
			share_price,
			requested_at,
			expires_at
	FROM pending_transactions
	WHERE uid = _uid AND type = 'sale'
	ORDER BY requested_at DESC
$$
LANGUAGE SQL VOLATILE;



CREATE FUNCTION commit_pending_sale(
	_id int,
	_made_at timestamptz
)
RETURNS int AS
$$
DECLARE 
	_transaction_id int;
BEGIN
	_transaction_id = perform_sale_transaction((
		SELECT 	uid,
				stock,
				num_shares,
				share_price,
				_made_at
		FROM pending_transactions
		WHERE id = _id
		));

	DELETE FROM pending_transactions WHERE id = _id;

	RETURN _transaction_id;
END;
$$
LANGUAGE 'plpgsql' VOLATILE;



CREATE FUNCTION cancel_pending_transaction(
	_id int
)
RETURNS void AS
$$
	DELETE FROM pending_transactions WHERE id = _id;
$$
LANGUAGE SQL VOLATILE;


