﻿
namespace TransactionEvents
{
    public enum CommandType
    {
        ADD,
        QUOTE,
        BUY,
        COMMIT_BUY,
        CANCEL_BUY,
        SELL,
        COMMIT_SELL,
        CANCEL_SELL,
        SET_BUY_AMOUNT,
        CANCEL_SET_BUY,
        SET_BUY_TRIGGER,
        SET_SELL_AMOUNT,
        CANCEL_SET_SELL,
        SET_SELL_TRIGGER,
        DUMPLOG,
        DISPLAY_SUMMARY
    }
}
