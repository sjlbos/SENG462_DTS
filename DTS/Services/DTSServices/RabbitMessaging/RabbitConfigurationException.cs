
using System;
using System.Runtime.Serialization;

namespace RabbitMessaging
{
    [Serializable]
    public class RabbitConfigurationException : Exception
    {
        public RabbitConfigurationException() { }

        public RabbitConfigurationException(string message) : base(message) { }

        public RabbitConfigurationException(string message, Exception innerException) : base(message, innerException) { }

        protected RabbitConfigurationException(SerializationInfo serializationInfo, StreamingContext streamingContext) : base(serializationInfo, streamingContext) { }
    }
}
