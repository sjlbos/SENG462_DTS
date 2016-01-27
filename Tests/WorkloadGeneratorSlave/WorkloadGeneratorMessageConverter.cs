
using System;
using Newtonsoft.Json;
using Newtonsoft.Json.Linq;

namespace WorkloadGeneratorSlave
{
    internal class WorkloadGeneratorMessageConverter : JsonConverter
    {
        public override void WriteJson(JsonWriter writer, object value, JsonSerializer serializer)
        {
            throw new NotImplementedException();
        }

        public override object ReadJson(JsonReader reader, Type objectType, object existingValue, JsonSerializer serializer)
        {
            JObject item = JObject.Load(reader);
            string messageType = item["MessageType"].Value<string>();
            switch (messageType)
            {
                case MessageType.ControlMessage:
                    return item.ToObject<ControlMessage>();
                case MessageType.BatchOrderMessage:
                    return item.ToObject<WorkloadBatchMessage>();
                default:
                    throw new UnrecognizedMessageTypeException(messageType);
            }
        }

        public override bool CanConvert(Type objectType)
        {
            return typeof (WorkloadGeneratorMessage).IsAssignableFrom(objectType);
        }
    }
}
