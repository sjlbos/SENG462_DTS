using System;
using Newtonsoft.Json;
using Newtonsoft.Json.Linq;

namespace TransactionEvents
{
    public class TransactionEventConverter : JsonConverter
    {
        private static class EventType
        {
            public const string UserCommandEvent = "UserCommandEvent";
            public const string QuoteServerEvent = "QuoteServerEvent";
            public const string AcccountTransactionEvent = "AccountTransactionEvent";
            public const string SystemEvent = "SystemEvent";
            public const string ErrorEvent = "ErrorEvent";
            public const string DebugEvent = "DebugEvent";
        }

        public override void WriteJson(JsonWriter writer, object value, JsonSerializer serializer)
        {
            var jObject = JObject.FromObject(value);
            jObject.Add("EventType", GetTransactionEventTypeString(value as TransactionEvent));
            jObject.WriteTo(writer);
        }

        public override object ReadJson(JsonReader reader, Type objectType, object existingValue, JsonSerializer serializer)
        {
            var item = JObject.Load(reader);

            var typeToken = item["EventType"];
            if (typeToken == null)
                throw new UnrecognizedTransactionEventException();

            string eventType = typeToken.Value<string>();

            switch (eventType)
            {
                case EventType.UserCommandEvent:
                    return item.ToObject<UserCommandEvent>();
                case EventType.QuoteServerEvent:
                    return item.ToObject<QuoteServerEvent>();
                case EventType.AcccountTransactionEvent:
                    return item.ToObject<AccountTransactionEvent>();
                case EventType.SystemEvent:
                    return item.ToObject<SystemEvent>();
                case EventType.ErrorEvent:
                    return item.ToObject<ErrorEvent>();
                case EventType.DebugEvent:
                    return item.ToObject<DebugEvent>();
                default:
                    throw new UnrecognizedTransactionEventException(eventType);
            }
        }

        public override bool CanConvert(Type objectType)
        {
            return typeof (TransactionEvent).IsAssignableFrom(objectType);
        }

        private static string GetTransactionEventTypeString(TransactionEvent transactionEvent)
        {
            if (transactionEvent is UserCommandEvent)
                return EventType.UserCommandEvent;
            if (transactionEvent is QuoteServerEvent)
                return EventType.QuoteServerEvent;
            if (transactionEvent is AccountTransactionEvent)
                return EventType.AcccountTransactionEvent;
            if (transactionEvent is SystemEvent)
                return EventType.SystemEvent;
            if (transactionEvent is ErrorEvent)
                return EventType.ErrorEvent;
            if (transactionEvent is DebugEvent)
                return EventType.DebugEvent;

            throw new UnrecognizedTransactionEventException(transactionEvent.GetType().ToString());
        }
    }
}
