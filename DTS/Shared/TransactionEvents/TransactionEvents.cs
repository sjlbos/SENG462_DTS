
using System;
using Newtonsoft.Json;
using Newtonsoft.Json.Converters;

namespace TransactionEvents
{
    public abstract class TransactionEvent
    {
        public Guid Id { get; set; }
        public DateTime LoggedAt { get; set; }
        public DateTime OccuredAt { get; set; }
        public Guid TransactionId { get; set; }
        public string UserId { get; set; }
        public string Service { get; set; }
        public string Server { get; set; }
    }

    public class UserCommandEvent : TransactionEvent
    {
        [JsonConverter(typeof(StringEnumConverter))]
        public CommandType CommandType { get; set; }
        public string StockSymbol { get; set; }
        public Decimal Funds { get; set; }
    }

    public class QuoteServerEvent : TransactionEvent
    {
        public string StockSymbol { get; set; }
        public Decimal Price { get; set; }
        public DateTime QuoteServerTime { get; set; }
        public string CryptoKey { get; set; }
    }

    public class AccountTransactionEvent : TransactionEvent
    {
        [JsonConverter(typeof(StringEnumConverter))]
        public AccountAction AccountAction { get; set; }
        public Decimal Funds { get; set; }
    }

    public class SystemEvent : TransactionEvent
    {
        [JsonConverter(typeof(StringEnumConverter))]
        public CommandType Command { get; set; }
        public string StockSymbol { get; set; }
        public Decimal Funds { get; set; }
        public string FileName { get; set; }
    }

    public class ErrorEvent : TransactionEvent
    {
        [JsonConverter(typeof(StringEnumConverter))]
        public CommandType Command { get; set; }
        public string StockSymbol { get; set; }
        public Decimal Funds { get; set; }
        public string ErrorMessage { get; set; }
        public string FileName { get; set; }
    }

    public class DebugEvent : TransactionEvent
    {
        [JsonConverter(typeof(StringEnumConverter))]
        public CommandType Command { get; set; }
        public string StockSymbol { get; set; }
        public Decimal Funds { get; set; }
        public string FileName { get; set; }
        public string DebugMessage { get; set; }
    }
}
