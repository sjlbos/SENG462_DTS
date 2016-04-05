package main

var getUserId string = "SELECT * FROM \"get_user_account_by_char_id\"($1)"
var addOrCreateUser string = "SELECT * FROM \"add_or_create_user_account\"($1::varchar, $2::money, $3::timestamptz)"
var updateBalance string = "SELECT * FROM \"update_user_account_balance\"($1, $2::money)"
var addPendingPurchase string = "SELECT * FROM \"add_pending_purchase\"($1,$2,$3::int,$4::money,$5, $6)"
var getLatestPendingPurchase string = "SELECT * FROM \"get_latest_pending_purchase_for_user\"($1)"
var commitPurchase string = "SELECT * FROM \"commit_pending_purchase\"($1,$2)"
var addPendingSale string = "SELECT * FROM \"add_pending_sale\"($1,$2,$3::int,$4::money,$5, $6)"
var getLatestPendingSale string = "SELECT * FROM \"get_latest_pending_sale_for_user\"($1::int)"
var commitSale string = "SELECT * FROM \"commit_pending_sale\"($1,$2)"
var cancelTransaction string = "SELECT * FROM \"cancel_pending_transaction\"($1)"
var addBuyTrigger string = "SELECT * FROM \"add_buy_trigger\"($1::int,$2::varchar,$3::money,$4::timestamptz)"
var addSellTrigger string = "SELECT * FROM \"add_sell_trigger\"($1::int,$2::varchar,$3::money,$4::timestamptz)"
var setBuyTrigger string = "SELECT * FROM \"commit_buy_trigger\"($1::int,$2::int,$3::money,$4::timestamptz)"
var setSellTrigger string = "SELECT * FROM \"commit_sell_trigger\"($1::int,$2::int,$3::money,$4::timestamptz)"
var cancelBuyTrigger string = "SELECT * FROM \"cancel_buy_trigger\"($1::int)"
var cancelSellTrigger string = "SELECT * FROM \"cancel_sell_trigger\"($1::int)"
var getBuyTriggerId string = "SELECT * FROM \"get_buy_trigger_id_for_user_and_stock\"($1::int, $2::varchar)"
var getSellTriggerId string = "SELECT * FROM \"get_sell_trigger_id_for_user_and_stock\"($1::int, $2::varchar)"
var getPendingTriggerId string = "SELECT * FROM \"get_pending_trigger_id_for_user_and_stock\"($1::int, $2::varchar, $3::trigger_type)"
var getTriggerById string = "SELECT * FROM \"get_trigger_by_id\"($1)"
var performBuyTrigger string = "SELECT * FROM \"perform_buy_trigger\"($1, $2::money)"
var performSellTrigger string = "SELECT * FROM \"perform_sell_trigger\"($1, $2::money)"
var getAllTriggers string = "SELECT * FROM \"get_all_triggers_for_uid\"($1::int)"
var getAllStocks string = "SELECT * FROM \"get_user_portfolio\"($1::int)"