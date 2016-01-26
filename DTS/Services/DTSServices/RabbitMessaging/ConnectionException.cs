
using System;
using System.Runtime.Serialization;

namespace RabbitMessaging
{
    [Serializable]
    public class ConnectionException : Exception
    {
         public ConnectionException() { }

        public ConnectionException(string message) : base(message) { }

        public ConnectionException(string message, Exception innerException) : base(message, innerException) { }

        protected ConnectionException(SerializationInfo serializationInfo, StreamingContext streamingContext) : base(serializationInfo, streamingContext) { }
    }
}
