
using System;
using System.Globalization;
using System.Xml;
using System.Xml.Serialization;
using Newtonsoft.Json;
using Newtonsoft.Json.Converters;

namespace TransactionEvents
{
    public abstract class TransactionEvent
    {
        public Guid Id { get; set; }
        public DateTime LoggedAt { get; set; }
        public DateTime OccuredAt { get; set; }
        public int TransactionId { get; set; }
        public string UserId { get; set; }
        public string Service { get; set; }
        public string Server { get; set; }

        public abstract void WriteXml(XmlWriter w);

        protected void WriteCommonPropertyXml(XmlWriter w)
        {
            w.WriteElementString("timestamp", OccuredAt.Ticks.ToString(CultureInfo.InvariantCulture));
            w.WriteElementString("transactionNum", TransactionId.ToString(CultureInfo.InvariantCulture));
            w.WriteElementString("username", UserId);
            w.WriteElementString("server", Server);
        }
    }

    public class UserCommandEvent : TransactionEvent
    {
        [JsonConverter(typeof(StringEnumConverter))]
        public CommandType Command { get; set; }
        public string StockSymbol { get; set; }
        public Decimal? Funds { get; set; }

        public override void WriteXml(XmlWriter w)
        {
            w.WriteStartElement("userCommand");
            WriteCommonPropertyXml(w);
            w.WriteElementString("command", Command.ToString());
            if (StockSymbol != null)
            {
                w.WriteElementString("stockSymbol", StockSymbol);
            }
            if (Funds != null)
            {
                w.WriteElementString("funds", Funds.Value.ToString("#.##"));
            }
            w.WriteEndElement();
        }
    }

    public class QuoteServerEvent : TransactionEvent
    {
        public string StockSymbol { get; set; }
        public Decimal Price { get; set; }
        public DateTime QuoteServerTime { get; set; }
        public string CryptoKey { get; set; }

        public override void WriteXml(XmlWriter w)
        {
            w.WriteStartElement("quoteServer");
            WriteCommonPropertyXml(w);
            w.WriteElementString("price", Price.ToString("N2"));
            w.WriteElementString("stockSymbol", StockSymbol);
            w.WriteElementString("quoteSeverTime", QuoteServerTime.Ticks.ToString(CultureInfo.InvariantCulture));
            w.WriteElementString("cryptokey", CryptoKey);
            w.WriteEndElement();
        }
    }

    [XmlType(TypeName = "accountTransaction")]
    public class AccountTransactionEvent : TransactionEvent
    {
        [JsonConverter(typeof(StringEnumConverter))]
        public AccountAction AccountAction { get; set; }
        public Decimal Funds { get; set; }

        public override void WriteXml(XmlWriter w)
        {
            w.WriteStartElement("accountTransaction");
            WriteCommonPropertyXml(w);
            w.WriteElementString("action", AccountAction.ToString());
            w.WriteElementString("funds", Funds.ToString("N2"));
            w.WriteEndElement();
        }
    }

    public class SystemEvent : TransactionEvent
    {
        [JsonConverter(typeof(StringEnumConverter))]
        public CommandType Command { get; set; }
        public string StockSymbol { get; set; }
        public Decimal Funds { get; set; }
        public string FileName { get; set; }
        public override void WriteXml(XmlWriter w)
        {
            w.WriteStartElement("systemEvent");
            WriteCommonPropertyXml(w);
            w.WriteElementString("command", Command.ToString());
            w.WriteElementString("stockSymbol", StockSymbol);
            w.WriteElementString("funds", Funds.ToString("N2"));
            w.WriteElementString("filename", FileName);
            w.WriteEndElement();
        }
    }

    public class ErrorEvent : TransactionEvent
    {
        [JsonConverter(typeof(StringEnumConverter))]
        public CommandType Command { get; set; }
        public string StockSymbol { get; set; }
        public Decimal Funds { get; set; }
        public string ErrorMessage { get; set; }
        public string FileName { get; set; }

        public override void WriteXml(XmlWriter w)
        {
            w.WriteStartElement("errorEvent");
            WriteCommonPropertyXml(w);
            w.WriteElementString("command", Command.ToString());
            w.WriteElementString("stockSymbol", StockSymbol);
            w.WriteElementString("funds", Funds.ToString("N2"));
            w.WriteElementString("errorMessage", ErrorMessage);
            w.WriteElementString("filename", FileName);
            w.WriteEndElement();
        }
    }

    public class DebugEvent : TransactionEvent
    {
        [JsonConverter(typeof(StringEnumConverter))]
        public CommandType Command { get; set; }
        public string StockSymbol { get; set; }
        public Decimal Funds { get; set; }
        public string FileName { get; set; }
        public string DebugMessage { get; set; }

        public override void WriteXml(XmlWriter w)
        {
            w.WriteStartElement("debugEvent");
            WriteCommonPropertyXml(w);
            w.WriteElementString("command", Command.ToString());
            w.WriteElementString("stockSymbol", StockSymbol);
            w.WriteElementString("funds", Funds.ToString("N2"));
            w.WriteElementString("debugMessage", DebugMessage);
            w.WriteElementString("filename", FileName);
            w.WriteEndElement();
        }
    }
}
