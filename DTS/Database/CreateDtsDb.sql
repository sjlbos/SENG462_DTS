DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS portfolios CASCADE;
DROP TYPE IF EXISTS trigger_type CASCADE;
DROP TABLE IF EXISTS triggers CASCADE;
DROP TYPE IF EXISTS transaction_type CASCADE;
DROP TABLE IF EXISTS transactions CASCADE;
DROP TABLE IF EXISTS pending_transactions CASCADE;
DROP TABLE IF EXISTS pending_triggers CASCADE;


CREATE TABLE IF NOT EXISTS users(
	id SERIAL PRIMARY KEY,
	user_id varchar NOT NULL UNIQUE,
	balance money NOT NULL DEFAULT 0.00 CHECK(balance >= 0::money),
	created_at timestamptz default current_timestamp
);

CREATE TABLE IF NOT EXISTS portfolios(
	uid int NOT NULL,
	stock varchar NOT NULL,
	num_shares int NOT NULL CHECK(num_shares > 0),
	FOREIGN KEY (uid) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE,
	UNIQUE(uid, stock)
);

CREATE TYPE trigger_type AS ENUM ('buy', 'sell');

CREATE TABLE IF NOT EXISTS triggers(
	id SERIAL PRIMARY KEY,
	uid int NOT NULL,
	stock varchar NOT NULL,
	type trigger_type NOT NULL,
	trigger_price money CHECK(trigger_price >= 0::money),
	num_shares int NOT NULL CHECK(num_shares > 0),
	created_at timestamptz NOT NULL,
	FOREIGN KEY (uid) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE,
	UNIQUE (uid, stock, type)
);

CREATE TABLE IF NOT EXISTS pending_triggers(
	id SERIAL PRIMARY KEY,
	uid int NOT NULL,
	stock varchar NOT NULL,
	type trigger_type NOT NULL,
	dollar_amount money CHECK(dollar_amount >= 0::money),
	created_at timestamptz NOT NULL,
	FOREIGN KEY (uid) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE,
	UNIQUE (uid, stock, type)
);

CREATE TYPE transaction_type AS ENUM ('sale', 'purchase');

CREATE TABLE IF NOT EXISTS transactions(
	id SERIAL PRIMARY KEY,
	uid int NOT NULL,
	type transaction_type NOT NULL,
	stock varchar NOT NULL,
	num_shares int NOT NULL CHECK(num_shares > 0),
	share_price money NOT NULL CHECK(share_price >= 0::money),
	made_at timestamptz NOT NULL,
	FOREIGN KEY (uid) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS pending_transactions(
	id SERIAL PRIMARY KEY,
	uid int NOT NULL,
	type transaction_type NOT NULL,
	stock varchar NOT NULL,
	num_shares int NOT NULL CHECK(num_shares > 0),
	share_price money NOT NULL CHECK(share_price >= 0::money),
	requested_at timestamptz NOT NULL,
	expires_at timestamptz NOT NULL,
	FOREIGN KEY (uid) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE
);



CREATE OR REPLACE FUNCTION get_reserved_funds (_uid int) RETURNS money AS
$$
	SELECT SUM(value)
	FROM(
		SELECT num_shares * trigger_price AS value
		FROM triggers
		WHERE uid = _uid AND type = 'buy'
		) AS trigger_values;
$$
LANGUAGE SQL VOLATILE;



CREATE OR REPLACE FUNCTION add_user_account(_user_id varchar, _balance money, _created_at timestamptz)
RETURNS int AS
$$
	INSERT INTO users(user_id, balance)
	VALUES(_user_id, _balance)
	RETURNING id;
$$
LANGUAGE SQL VOLATILE;



CREATE OR REPLACE FUNCTION update_user_account_balance(_uid int, _balance money) 
RETURNS void AS
$$
	UPDATE users
	SET balance = _balance
	WHERE id = _uid;
$$
LANGUAGE SQL VOLATILE;

CREATE OR REPLACE FUNCTION get_user_by_id(_uid int)
RETURNS TABLE(id int, user_id varchar, balance money) AS
$$
	SELECT 	id, 
			user_id,
			balance
	FROM users
	WHERE id = _uid;
$$
LANGUAGE SQL VOLATILE;

CREATE OR REPLACE FUNCTION get_user_account_by_char_id(_user_id varchar)
RETURNS TABLE(id int, user_id varchar, balance money) AS
$$
	SELECT 	id,
			user_id,
			balance
	FROM users
	WHERE user_id = _user_id;
$$
LANGUAGE SQL VOLATILE;



CREATE OR REPLACE FUNCTION get_user_portfolio(_uid int)
RETURNS TABLE(stock varchar, num_shares int) AS
$$
	SELECT 	stock,
	num_shares
	FROM portfolios
	WHERE uid = _uid;
$$
LANGUAGE SQL VOLATILE;



CREATE OR REPLACE FUNCTION get_user_stock_amount(_uid int, _stock varchar)
RETURNS int AS
$$ 
	SELECT num_shares
	FROM portfolios
	WHERE uid = _uid AND stock = _stock;
$$
LANGUAGE SQL VOLATILE;



CREATE OR REPLACE FUNCTION perform_purchase_transaction(
	_uid int,
	_stock varchar,
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



CREATE OR REPLACE FUNCTION perform_sale_transaction(
	_uid int,
	_stock varchar,
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



CREATE OR REPLACE FUNCTION get_all_triggers()
RETURNS TABLE (	id int, 
				uid int,
				stock varchar, 
				type trigger_type, 
				trigger_price money, 
				num_shares int, 
				created_at timestamptz) AS
$$
	SELECT 	id,
			uid,
			stock,
			type,
			trigger_price,
			num_shares,
			created_at
	FROM triggers;
$$
LANGUAGE SQL VOLATILE;



CREATE OR REPLACE FUNCTION get_trigger_by_id(_id int)
RETURNS TABLE (	id int, 
		uid int,
		stock varchar, 
		type trigger_type, 
		trigger_price money, 
		num_shares int, 
		created_at timestamptz) AS
$$
	SELECT 	id,
			uid,
			stock,
			type,
			trigger_price,
			num_shares,
			created_at
	FROM triggers
	WHERE id = _id;
$$
LANGUAGE SQL VOLATILE;



CREATE OR REPLACE FUNCTION get_user_buy_triggers(_uid int)
RETURNS TABLE (	id int, 
				uid int,
				stock varchar, 
				type trigger_type, 
				trigger_price money, 
				num_shares int, 
				created_at timestamptz) AS
$$
	SELECT 	id,
			uid,
			stock,
			type,
			trigger_price,
			num_shares,
			created_at
	FROM triggers
	WHERE 	uid = _uid 
		AND type = 'buy';
$$
LANGUAGE SQL VOLATILE;



CREATE OR REPLACE FUNCTION get_user_sell_triggers(_uid int)
RETURNS TABLE (	id int, 
				uid int,
				stock varchar, 
				type trigger_type, 
				trigger_price money, 
				num_shares int, 
				created_at timestamptz) AS
$$
	SELECT 	id,
			uid,
			stock,
			type,
			trigger_price,
			num_shares,
			created_at
	FROM triggers
	WHERE 	uid = _uid 
		AND type = 'sell';
$$
LANGUAGE SQL VOLATILE;



CREATE OR REPLACE FUNCTION get_buy_trigger_for_user_and_stock(_uid int, _stock varchar)
RETURNS TABLE (	id int, 
				uid int,
				stock varchar, 
				type trigger_type, 
				trigger_price money, 
				num_shares int, 
				created_at timestamptz) AS
$$
	SELECT 	id,
			uid,
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

CREATE OR REPLACE FUNCTION get_buy_trigger_id_for_user_and_stock(_uid int, _stock varchar)
RETURNS TABLE (	id int ) AS
$$
	SELECT 	id
	FROM triggers
	WHERE uid = _uid
		AND stock = _stock
		AND type = 'buy';
$$
LANGUAGE SQL VOLATILE;



CREATE OR REPLACE FUNCTION get_sell_trigger_for_user_and_stock(_uid int, _stock varchar)
RETURNS TABLE (	id int, 
				uid int,
				stock varchar, 
				type trigger_type, 
				trigger_price money, 
				num_shares int, 
				created_at timestamptz) AS
$$
	SELECT 	id,
			uid,
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


CREATE OR REPLACE FUNCTION get_sell_trigger_id_for_user_and_stock(_uid int, _stock varchar)
RETURNS TABLE (	id int ) AS
$$
	SELECT 	id
	FROM triggers
	WHERE uid = _uid 
		AND stock = _stock
		AND type = 'sell';
$$
LANGUAGE SQL VOLATILE;



CREATE OR REPLACE FUNCTION get_pending_trigger_id_for_user_and_stock( _uid int, _stock varchar, _type trigger_type)
RETURNS TABLE ( id int ) AS
$$
	SELECT id
	FROM pending_triggers
	WHERE uid = _uid
		AND stock = _stock
		AND type = _type;
$$
LANGUAGE SQL VOLATILE;


CREATE OR REPLACE FUNCTION add_buy_trigger(
	_uid int,
	_stock varchar,
	_dollar_amount money,
	_created_at timestamptz
)
RETURNS int AS
$$
DECLARE
	_transaction_id int;
	_balance money;
BEGIN
	SELECT balance INTO _balance FROM users WHERE id = _uid;

	UPDATE users SET balance = _balance - _dollar_amount;

	INSERT INTO pending_triggers(
		uid,
		stock,
		type,
		dollar_amount,
		created_at
		)
	VALUES(
		_uid,
		_stock,
		'buy',
		_dollar_amount,
		_created_at
		)
	RETURNING id INTO _transaction_id;

	RETURN _transaction_id;
END;
$$
LANGUAGE 'plpgsql' VOLATILE;


CREATE OR REPLACE FUNCTION add_sell_trigger(
	_uid int,
	_stock varchar,
	_dollar_amount money,
	_created_at timestamptz
)
RETURNS int AS
$$
	INSERT INTO pending_triggers(
		uid,
		stock,
		type,
		dollar_amount,
		created_at
		)
	VALUES(
		_uid,
		_stock,
		'sell',
		_dollar_amount,
		_created_at
		)
	RETURNING id;
$$
LANGUAGE SQL VOLATILE;



CREATE OR REPLACE FUNCTION commit_buy_trigger(
	_id int,
	_uid int,
	_stock_price money,
	_made_at timestamptz
)
RETURNS int AS
$$
DECLARE
	num_shares int;
	rtnMoney money;
	moneySaved money;
BEGIN
	SELECT dollar_amount INTO moneySaved 
	FROM pending_triggers 
	WHERE uid = _uid;

	num_shares = moneySaved / _stock_price;

	INSERT INTO triggers(
					uid, 
					stock, 
					type, 
					created_at,
					num_shares,
					trigger_price)
	(	SELECT
			uid,
			stock,
			type,
			current_timestamp,
			_num_shares,
			_stock_price
		FROM pending_triggers
		WHERE id = _id);

	rtnMoney = moneySaved - (_stock_price * _num_shares);
	UPDATE users SET balance = balance + rtnMoney WHERE id = _user_id;

	DELETE FROM pending_triggers WHERE id = _id;
	RETURN _id;
END;
$$
LANGUAGE 'plpgsql' VOLATILE;



CREATE OR REPLACE FUNCTION commit_sell_trigger(
	_id int,
	_uid int,
	_stock_price money,
	_made_at timestamptz
)
RETURNS int AS
$$
DECLARE
	_transaction_id int;
	_num_shares int;
	_stock varchar;
	moneySaved money;
BEGIN
	SELECT dollar_amount INTO moneySaved 
	FROM pending_triggers 
	WHERE id = _id;

	_num_shares = moneySaved / _stock_price;

	UPDATE portfolios SET num_shares = numshares - _num_shares;

	INSERT INTO triggers(
				uid, 
				stock, 
				type, 
				created_at,
				num_shares,
				trigger_price)
	(	SELECT
			uid,
			stock,
			type,
			current_timestamp,
			_num_shares,
			_stock_price
		FROM pending_triggers
		WHERE id = _id)
		RETURNING stock INTO _stock;

	DELETE FROM pending_triggers WHERE id = _id RETURNING id INTO _transaction_id;

	RETURN _transaction_id;
END;
$$
LANGUAGE 'plpgsql' VOLATILE;


CREATE OR REPLACE FUNCTION cancel_buy_trigger(
	_id int
)
RETURNS void AS
$$
DECLARE
	_dollar_amount money;
	_user_id int;
BEGIN
	SELECT trigger_price * num_shares, uid INTO _dollar_amount, _user_id
	FROM triggers
	WHERE id = _id;	

	UPDATE users SET balance = balance + _dollar_amount WHERE id = _user_id;

	DELETE FROM triggers 
	WHERE id = _id;
END;
$$
LANGUAGE 'plpgsql' VOLATILE;



CREATE OR REPLACE FUNCTION cancel_sell_trigger(
	_id int
)
RETURNS void AS
$$
DECLARE
	_num_shares int;
	_uid int;
	_existing_shares int;
	_stock varchar;
BEGIN

	SELECT num_shares, uid, stock INTO _num_shares, _uid, _stock
	FROM triggers
	WHERE id = _id;

	_existing_shares = get_user_stock_amount(_uid, _stock);

	IF _existing_shares IS NULL THEN
		INSERT INTO portfolios (uid, stock, num_shares)
			VALUES (_uid, _stock, _num_shares);
	ELSE
		UPDATE portfolios SET num_shares = (_num_shares + _existing_shares) WHERE uid = _uid AND stock = _stock;
	END IF;

	DELETE FROM triggers 
	WHERE id = _id;
END;
$$
LANGUAGE 'plpgsql' VOLATILE;



CREATE OR REPLACE FUNCTION perform_sell_trigger(
	_id int,
	_quote_price money
)
RETURNS void AS
$$
DECLARE
	_num_shares int;
	_uid int;
	_sale_value money;
BEGIN
	SELECT num_shares,uid INTO _num_shares,_uid
	FROM triggers 
	WHERE id = _id;

	_sale_value = _num_shares * _quote_price;

	UPDATE users SET balance = balance + _sale_value WHERE id = _uid;

	DELETE FROM triggers WHERE id = _id;
END;
$$
LANGUAGE 'plpgsql' VOLATILE;

CREATE OR REPLACE FUNCTION perform_buy_trigger(
	_id int,
	_quote_price money
)
RETURNS void AS
$$
DECLARE
	_num_shares int;
	_stock varchar;
	_uid int;
	_sale_value money;
BEGIN
	SELECT num_shares,uid,stock INTO _num_shares,_uid,_stock
	FROM triggers 
	WHERE id = _id;

	INSERT INTO portfolios AS p (uid, stock, num_shares)
	VALUES (_uid, _stock, _num_shares)
	ON CONFLICT (uid, stock)
	DO UPDATE SET num_shares = p.num_shares + _num_shares
	WHERE p.uid = _uid and p.stock = _stock; 

	DELETE FROM triggers WHERE id = _id;
END;
$$
LANGUAGE 'plpgsql' VOLATILE;

	

CREATE OR REPLACE FUNCTION get_user_transaction_history(
	_uid int
)
RETURNS TABLE (id int, type transaction_type, stock varchar, num_shares int, share_price money, made_at timestamptz) AS
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



CREATE OR REPLACE FUNCTION add_pending_purchase(
	_uid int,
	_stock varchar,
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



CREATE OR REPLACE FUNCTION add_pending_sale(
	_uid int,
	_stock varchar,
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



CREATE OR REPLACE FUNCTION get_latest_pending_purchase_for_user(
	_uid int
)
RETURNS TABLE(id int, uid int, stock varchar, num_shares int, share_price money, requested_at timestamptz, expires_at timestamptz) AS
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



CREATE OR REPLACE FUNCTION get_all_pending_purchases_for_user(
	_uid int
)
RETURNS TABLE(id int, uid int, stock varchar, num_shares int, share_price money, requested_at timestamptz, expires_at timestamptz) AS
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


CREATE OR REPLACE FUNCTION commit_pending_purchase(
	_id int,
	_made_at timestamptz
)
RETURNS int AS
$$
DECLARE 
	_transaction_id int;
	_uid int;
	_stock varchar;
	_num_shares int;
	_share_price money;
BEGIN
	SELECT  uid, stock, num_shares, share_price 
	INTO    _uid, _stock,_num_shares,_share_price
	FROM    pending_transactions
	WHERE   id = _id;

	_transaction_id = perform_purchase_transaction(
		_uid,
		_stock,
		_num_shares,
		_share_price,
		_made_at
	);

	DELETE FROM pending_transactions WHERE id = _id;

	RETURN _transaction_id;
END;
$$
LANGUAGE 'plpgsql' VOLATILE;



CREATE OR REPLACE FUNCTION get_latest_pending_sale_for_user(
	_uid int
)
RETURNS TABLE(id int, uid int, stock varchar, num_shares int, share_price money, requested_at timestamptz, expires_at timestamptz) AS
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
	LIMIT 1;
$$
LANGUAGE SQL VOLATILE;



CREATE OR REPLACE FUNCTION get_all_pending_sales_for_user(
	_uid int
)
RETURNS TABLE(id int, uid int, stock varchar, num_shares int, share_price money, requested_at timestamptz, expires_at timestamptz) AS
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



CREATE OR REPLACE FUNCTION commit_pending_sale(
	_id int,
	_made_at timestamptz
)
RETURNS int AS
$$
DECLARE 
	_transaction_id int;
	_uid int;
	_stock varchar;
	_num_shares int;
	_share_price money;
BEGIN
	SELECT  uid, stock, num_shares, share_price 
	INTO    _uid, _stock,_num_shares,_share_price
	FROM    pending_transactions
	WHERE   id = _id;

	_transaction_id = perform_sale_transaction(
		_uid,
		_stock,
		_num_shares,
		_share_price,
		_made_at
	);

	DELETE FROM pending_transactions WHERE id = _id;

	RETURN _transaction_id;
END;
$$
LANGUAGE 'plpgsql' VOLATILE;



CREATE OR REPLACE FUNCTION cancel_pending_transaction(
	_id int
)
RETURNS void AS
$$
	DELETE FROM pending_transactions WHERE id = _id;
$$
LANGUAGE SQL VOLATILE;


